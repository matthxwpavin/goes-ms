package goes_ms

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type SmsSender struct {
	client *http.Client
	cred   *Credentials
}

func NewSender(cred *Credentials) *SmsSender {
	return NewSenderWithProxy(cred, nil)
}

func NewSenderWithProxy(cred *Credentials, proxyURL *url.URL) *SmsSender {
	transport := http.DefaultTransport
	if proxyURL != nil {
		transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}
	return &SmsSender{
		client: &http.Client{Transport: transport},
		cred:   cred,
	}
}

func (s *SmsSender) Send(ctx context.Context, body *Body) (*http.Response, error) {
	if body == nil {
		return nil, errors.New("the body must not be nil")
	}

	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://portal-otp.smsmkt.com/api/send-message", bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api_key", s.cred.ApiKey)
	req.Header.Add("secret_key", s.cred.SecretKey)

	return s.client.Do(req)
}

type Credentials struct {
	ApiKey, SecretKey string
}

type Body struct {

	// Message is a message to recipients.
	Message string `json:"message"`

	// Phone is a destination's phone number, and, or can be multiple numbers seperated by comma.
	Phone string `json:"phone"`

	// sender name that approved by the APIs provider.
	Sender string `json:"sender"`

	// SendDate is a scheduled date in YYYY-MM-DD HH:mm:ss format.
	SendDate string `json:"send_date"`

	// URL is an url to provide to receive a result of sending the message.
	URL string `json:"url"`

	// Expire is the message's expiry in HH:mm format. valid range: 5m < expire < 11h
	// When expired, recipients will not receive the message.
	Expire string `json:"expire"`
}
