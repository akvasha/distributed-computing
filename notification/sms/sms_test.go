package sms

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const DefaultHost = "https://sms.ru"

func TestSMSClient_Send(t *testing.T) {
	apiId := os.Getenv("TEST_API_ID")
	if apiId == "" {
		t.Skip("TEST_API_ID is not set")
	}
	client := BuildClient(DefaultHost, apiId)
	err := client.Send(os.Getenv("TEST_SMS_RECEIVER"), "Test Message")
	assert.Nil(t, err)
}
