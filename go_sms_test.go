package goes_ms

import (
	"context"
	"io"
	"net/http"
	"testing"
)

func TestSendSMS(t *testing.T) {
	const (
		sender    = ""
		apiKey    = ""
		secretKey = ""
	)

	s := NewSender(&Credentials{
		ApiKey:    apiKey,
		SecretKey: secretKey,
	})

	resp, err := s.Send(context.Background(), &Body{
		Message: "Hello World!!!",
		Phone:   "09999999999",
		Sender:  sender,
	})
	if err != nil {
		t.Errorf("failed to send sms: %v", err)
		return
	}
	if resp.StatusCode == http.StatusOK {
		return
	}

	t.Errorf("unexpected http status code: %v", resp.StatusCode)
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Logf("unable to read response body: %v", err)
	}
	t.Logf("response body: %v", string(buf))
}
