package grafana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	_ "embed"
	"net/http"
	"net/url"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

//go:embed node_exporter.json
var nodeDashboard []byte

type Config struct {
	PublicURL string
	ApiURL    string
	User      string
	Password  string
}

func (config Config) BaseUrl() string {
	if config.PublicURL != "" {
		return config.PublicURL
	}
	return config.ApiURL
}

func InitGrafana(ctx context.Context, config *Config) error {
	return retry.Do(
		func() error {
			log.Info("Initializing Grafana Dashboards")
			token, err := createToken(config)
			if err != nil {
				log.Error("Failed to create Grafana token")
				return err
			}

			dashboardsURL := fmt.Sprintf("%s/%s", config.ApiURL, "api/dashboards/db")
			req, err := http.NewRequest("POST", dashboardsURL, bytes.NewBuffer(nodeDashboard))
			if err != nil {
				return err
			}

			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
			req.Header.Add("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("failed to create Grafana Dashboard: %s", resp.Status)
			}

			return err
		},
		retry.OnRetry(func(_ uint, err error) {
			log.Info("Grafana Dashboards initialization failed")
			log.Error(err)
		}),
		retry.Delay(2*time.Second),
		retry.MaxJitter(3*time.Second),
		retry.Attempts(8),
		retry.LastErrorOnly(true),
		retry.Context(ctx),
	)
}

func createToken(config *Config) (string, error) {
	tokenName, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	tokenRequest, _ := json.Marshal(map[string]interface{}{
		"role":          "Admin",
		"name":          tokenName.String(),
		"secondsToLive": 60,
	})

	u, err := url.Parse(config.ApiURL)
	if err != nil {
		log.Fatal("Invalid Grafana URL provided")
	}
	authenticatedURL := fmt.Sprintf("%s://%s:%s@%s%s/api/auth/keys", u.Scheme, config.User, config.Password, u.Host, u.Path)
	req, err := http.NewRequest("POST", authenticatedURL, bytes.NewBuffer(tokenRequest))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create Grafana token: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println(string(body))

	parsedBody := struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Key  string `json:"key"`
	}{}
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		return "", err
	}

	return parsedBody.Key, nil
}
