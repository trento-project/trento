package subscription

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/mock"
	"github.com/trento-project/trento/internal/subscription/mocks"
)

func mockSUSEConnect() *exec.Cmd {
	return exec.Command("echo", `[{"identifier":"SLES_SAP","version":"15.2","arch":"x86_64",
    "status":"Registered","name":"SUSE Employee subscription for SUSE Linux Enterprise Server for SAP Applications",
    "regcode":"my-code","starts_at":"2019-03-20 09:55:32 UTC",
    "expires_at":"2024-03-20 09:55:32 UTC","subscription_status":"ACTIVE","type":"internal"},
    {"identifier":"sle-module-public-cloud","version":"15.2",
    "arch":"x86_64","status":"Registered"}]`)
}

func mockSUSEConnectErr() *exec.Cmd {
	return exec.Command("error")
}

func TestNewSubscriptions(t *testing.T) {
	mockCommand := new(mocks.CustomCommand)

	customExecCommand = mockCommand.Execute

	mockCommand.On("Execute", "SUSEConnect", "-s").Return(
		mockSUSEConnect(),
	)

	subs, err := NewSubscriptions()

	expectedSubs := Subscriptions{
		&Subscription{
			Identifier:         "SLES_SAP",
			Version:            "15.2",
			Arch:               "x86_64",
			Status:             "Registered",
			StartsAt:           "2019-03-20 09:55:32 UTC",
			ExpiresAt:          "2024-03-20 09:55:32 UTC",
			SubscriptionStatus: "ACTIVE",
			Type:               "internal",
		},
		&Subscription{
			Identifier: "sle-module-public-cloud",
			Version:    "15.2",
			Arch:       "x86_64",
			Status:     "Registered",
		},
	}

	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedSubs, subs)
}

func TestNewSubscriptionsErr(t *testing.T) {
	mockCommand := new(mocks.CustomCommand)

	customExecCommand = mockCommand.Execute

	mockCommand.On("Execute", "SUSEConnect", "-s").Return(
		mockSUSEConnectErr(),
	)

	subs, err := NewSubscriptions()

	assert.Equal(t, Subscriptions(nil), subs)
	assert.EqualError(t, err, "exec: \"error\": executable file not found in $PATH")
}
