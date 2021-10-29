package collector

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type Client interface {
	Publish(discoveryType string, payload interface{}) error
}

type client struct {
	config     *Config
	machineID  string
	httpClient *http.Client
}

type Config struct {
	CollectorHost string
	CollectorPort int
	EnablemTLS    bool
	Cert          string
	Key           string
	CA            string
}

const machineIdPath = "/etc/machine-id"

var fileSystem = afero.NewOsFs()

func NewCollectorClient(config *Config) (*client, error) {
	var tlsConfig *tls.Config
	var err error

	if config.EnablemTLS {
		tlsConfig, err = getTLSConfig(config.Cert, config.Key, config.CA)
		if err != nil {
			return nil, err
		}
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	machineIDBytes, err := afero.ReadFile(fileSystem, machineIdPath)

	if err != nil {
		return nil, err
	}

	machineID := strings.TrimSpace(string(machineIDBytes))

	return &client{
		config:     config,
		httpClient: httpClient,
		machineID:  machineID,
	}, nil
}

func (c *client) Publish(discoveryType string, payload interface{}) error {
	// TODO: remove this when we want to start collecting
	if !viper.GetBool("collector-enabled") {
		return nil
	}

	log.Debugf("Sending %s to data collector", discoveryType)

	requestBody, err := json.Marshal(map[string]interface{}{
		"agent_id":       c.machineID,
		"discovery_type": discoveryType,
		"payload":        payload,
	})
	if err != nil {
		return err
	}

	protocol := "http"
	if c.config.EnablemTLS {
		protocol = "https"
	}

	endpoint := fmt.Sprintf("%s://%s:%d/api/collect", protocol, c.config.CollectorHost, c.config.CollectorPort)
	resp, err := c.httpClient.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf(
			"something wrong happened while publishing data to the collector. Status: %d, Agent: %s, discovery: %s",
			resp.StatusCode, c.machineID, discoveryType)
	}

	return nil
}

func getTLSConfig(cert, key, ca string) (*tls.Config, error) {
	caCert, err := ioutil.ReadFile(ca)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	certificate, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            caCertPool,
		Certificates:       []tls.Certificate{certificate},
	}, nil
}
