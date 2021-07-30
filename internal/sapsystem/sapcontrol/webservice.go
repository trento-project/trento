package sapcontrol

import (
	"context"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"path"

	"github.com/hooklift/gowsdl/soap"
)

//go:generate mockery --all

type WebService interface {
	GetInstanceProperties() (*GetInstancePropertiesResponse, error)
	GetProcessList() (*GetProcessListResponse, error)
	GetSystemInstanceList() (*GetSystemInstanceListResponse, error)
}

type STATECOLOR string
type STATECOLOR_CODE int

const (
	STATECOLOR_GRAY   STATECOLOR = "SAPControl-GRAY"
	STATECOLOR_GREEN  STATECOLOR = "SAPControl-GREEN"
	STATECOLOR_YELLOW STATECOLOR = "SAPControl-YELLOW"
	STATECOLOR_RED    STATECOLOR = "SAPControl-RED"

	// NOTE: This was just copy-pasted from sap_host_exporter, not used right now
	// see: https://github.com/SUSE/sap_host_exporter/blob/68bbf2f1b490ab0efaa2dd7b878b778f07fba2ab/lib/sapcontrol/webservice.go#L42
	STATECOLOR_CODE_GRAY   STATECOLOR_CODE = 1
	STATECOLOR_CODE_GREEN  STATECOLOR_CODE = 2
	STATECOLOR_CODE_YELLOW STATECOLOR_CODE = 3
	STATECOLOR_CODE_RED    STATECOLOR_CODE = 4
)

type GetInstanceProperties struct {
	XMLName xml.Name `xml:"urn:SAPControl GetInstanceProperties"`
}

type GetProcessList struct {
	XMLName xml.Name `xml:"urn:SAPControl GetProcessList"`
}

type GetProcessListResponse struct {
	XMLName   xml.Name     `xml:"urn:SAPControl GetProcessListResponse"`
	Processes []*OSProcess `xml:"process>item,omitempty" json:"process>item,omitempty"`
}
type GetInstancePropertiesResponse struct {
	XMLName    xml.Name            `xml:"urn:SAPControl GetInstancePropertiesResponse"`
	Properties []*InstanceProperty `xml:"properties>item,omitempty" json:"properties>item,omitempty"`
}

type GetSystemInstanceList struct {
	XMLName xml.Name `xml:"urn:SAPControl GetSystemInstanceList"`
	Timeout int32    `xml:"timeout,omitempty" json:"timeout,omitempty"`
}

type GetSystemInstanceListResponse struct {
	XMLName   xml.Name       `xml:"urn:SAPControl GetSystemInstanceListResponse"`
	Instances []*SAPInstance `xml:"instance>item,omitempty" json:"instance>item,omitempty"`
}

type OSProcess struct {
	Name        string     `xml:"name,omitempty" json:"name,omitempty" mapstructure:"name,omitempty"`
	Description string     `xml:"description,omitempty" json:"description,omitempty" mapstructure:"description,omitempty"`
	Dispstatus  STATECOLOR `xml:"dispstatus,omitempty" json:"dispstatus,omitempty" mapstructure:"dispstatus,omitempty"`
	Textstatus  string     `xml:"textstatus,omitempty" json:"textstatus,omitempty" mapstructure:"textstatus,omitempty"`
	Starttime   string     `xml:"starttime,omitempty" json:"starttime,omitempty" mapstructure:"starttime,omitempty"`
	Elapsedtime string     `xml:"elapsedtime,omitempty" json:"elapsedtime,omitempty" mapstructure:"elapsedtime,omitempty"`
	Pid         int32      `xml:"pid,omitempty" json:"pid,omitempty" mapstructure:"pid,omitempty"`
}

type InstanceProperty struct {
	Property     string `xml:"property,omitempty" json:"property,omitempty" mapstructure:"property,omitempty"`
	Propertytype string `xml:"propertytype,omitempty" json:"propertytype,omitempty" mapstructure:"propertytype,omitempty"`
	Value        string `xml:"value,omitempty" json:"value,omitempty" mapstructure:"value,omitempty"`
}

type SAPInstance struct {
	Hostname      string     `xml:"hostname,omitempty" json:"hostname,omitempty" mapstructure:"hostname,omitempty"`
	InstanceNr    int32      `xml:"instanceNr,omitempty" json:"instanceNr" mapstructure:"instancenr"`
	HttpPort      int32      `xml:"httpPort,omitempty" json:"httpPort,omitempty" mapstructure:"httpport,omitempty"`
	HttpsPort     int32      `xml:"httpsPort,omitempty" json:"httpsPort,omitempty" mapstructure:"httpsport,omitempty"`
	StartPriority string     `xml:"startPriority,omitempty" json:"startPriority,omitempty" mapstructure:"startpriority,omitempty"`
	Features      string     `xml:"features,omitempty" json:"features,omitempty" mapstructure:"features,omitempty"`
	Dispstatus    STATECOLOR `xml:"dispstatus,omitempty" json:"dispstatus,omitempty" mapstructure:"dispstatus,omitempty"`
}

type webService struct {
	client *soap.Client
}

func NewWebService(instNumber string) WebService {
	socket := path.Join("/tmp", fmt.Sprintf(".sapstream5%s13", instNumber))

	udsClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				d := net.Dialer{}
				return d.DialContext(ctx, "unix", socket)
			},
		},
	}

	// The url used here is just phony:
	// we need a well formed url to create the instance but the above DialContext function won't actually use it.
	client := soap.NewClient("http://unix", soap.WithHTTPClient(udsClient))

	return &webService{
		client: client,
	}
}

// GetInstanceProperties returns a list of available instance features and information how to get it.
func (s *webService) GetInstanceProperties() (*GetInstancePropertiesResponse, error) {
	request := &GetInstanceProperties{}
	response := &GetInstancePropertiesResponse{}
	err := s.client.Call("''", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetProcessList returns a list of all processes directly started by the webservice
// according to the SAP start profile.
func (s *webService) GetProcessList() (*GetProcessListResponse, error) {
	request := &GetProcessList{}
	response := &GetProcessListResponse{}
	err := s.client.Call("''", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetSystemInstanceList returns a list of all processes directly started by the webservice
// according to the SAP start profile.
func (s *webService) GetSystemInstanceList() (*GetSystemInstanceListResponse, error) {
	request := &GetSystemInstanceList{}
	response := &GetSystemInstanceListResponse{}
	err := s.client.Call("''", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
