package subscription

import (
	"encoding/json"
	"os/exec"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Subscriptions []*Subscription

type Subscription struct {
	Identifier string `json:"identifier,omitempty" mapstructure:"identifier,omitempty"`
	Version    string `json:"version,omitempty" mapstructure:"version,omitempty"`
	Arch       string `json:"arch,omitempty" mapstructure:"arch,omitempty"`
	Status     string `json:"status,omitempty" mapstructure:"status,omitempty"`
	//RegCode string `json:"regcode,omitempty" mapstructure:"regcode,omitempty"`
	StartsAt           string `json:"starts_at,omitempty" mapstructure:"starts_at,omitempty"`
	ExpiresAt          string `json:"expires_at,omitempty" mapstructure:"expires_at,omitempty"`
	SubscriptionStatus string `json:"subscription_status,omitempty" mapstructure:"subscription_status,omitempty"`
	Type               string `json:"type,omitempty" mapstructure:"type,omitempty"`
}

type CustomCommand func(name string, arg ...string) *exec.Cmd

var customExecCommand CustomCommand = exec.Command

func NewSubscriptions() (Subscriptions, error) {
	var subs Subscriptions

	log.Info("Identifying the SUSE subscription details...")
	output, err := customExecCommand("SUSEConnect", "-s").Output()
	if err != nil {
		return nil, err
	}

	log.Debugf("SUSEConnect output: %s", string(output))

	err = json.Unmarshal(output, &subs)
	if err != nil {
		return nil, errors.Wrap(err, "error while decoding the subscription details")
	}
	log.Infof("Subscription (%d entries) discovered", len(subs))

	return subs, nil
}
