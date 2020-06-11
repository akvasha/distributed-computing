package sms

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type SMSClient struct {
	host   string
	apiID  string
	client *http.Client
}

func BuildClient(host string, apiID string) (smsClient SMSClient) {
	smsClient.host = host
	smsClient.apiID = apiID
	smsClient.client = http.DefaultClient
	return
}

func InitClient() (smsClient SMSClient) {
	smsClient = BuildClient(os.Getenv("SMS_HOST"),
		os.Getenv("SMS_API_ID"))
	return
}

type SMSInfo struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	SmsID      string `json:"sms_id"`
	StatusText string `json:"status_text"`
}

type SendResponse struct {
	Status     string             `json:"status"`
	StatusCode int                `json:"status_code"`
	SMS        map[string]SMSInfo `json:"sms"`
	Balance    float64            `json:"balance"`
}

func (c *SMSClient) Send(receiver string, text string) (err error) {
	req := make(url.Values)
	req.Add("to", receiver)
	req.Add("msg", text)
	req.Add("api_id", c.apiID)
	req.Add("json", "1")
	urlQuery := c.host + "/sms/send/?" + req.Encode()

	var postResp *http.Response
	if postResp, err = c.client.Post(urlQuery, "application/json", nil); err != nil {
		return
	}
	defer postResp.Body.Close()

	var bytes []byte
	if bytes, err = ioutil.ReadAll(postResp.Body); err != nil {
		return
	}
	var resp SendResponse
	err = json.Unmarshal(bytes, &resp)
	if err == nil && resp.Status != "OK" {
		err = errors.New("Response status is not ok")
	}
	return
}
