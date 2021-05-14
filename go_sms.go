package goes_ms

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

const (
	keyConfig = "Clicknic-sender"
	user      = "user"
	password  = "password"
	sender    = "sender"
)

type SmsSender struct {
	client    *http.Client
	smsApiUrl string
	user      string
	password  string
	sender    string
}

func NewDefaultClicknicSender(url string) *SmsSender {
	clicknicSender := viper.GetStringMapString(keyConfig)
	return &SmsSender{
		client:    http.DefaultClient,
		smsApiUrl: url,
		user:      clicknicSender[user],
		password:  clicknicSender[password],
		sender:    clicknicSender[sender],
	}
}

func NewClicknicSenderToProxy(url string, proxyURL *url.URL) *SmsSender {
	clicknicSender := viper.GetStringMapString(keyConfig)
	return &SmsSender{
		client:    &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}},
		smsApiUrl: url,
		user:      clicknicSender[user],
		password:  clicknicSender[password],
		sender:    clicknicSender[sender],
	}
}

func (sd *SmsSender) ChangeSender(name string) {
	sd.sender = name
}

func (sd *SmsSender) SendSms(mobileNo, message string) *SmsStatus {
	var mobileNumberForSending string
	if strings.HasPrefix(mobileNo, "0") {
		mobileNumberForSending = "+66" + mobileNo[1:]
	} else {
		mobileNumberForSending = mobileNo
	}

	return sd.send(mobileNumberForSending, message)
}

const (
	headerUserAgent   = "User-Agent"
	headerContentType = "Content-Type"
)

func (sd *SmsSender) send(mobileNo, message string) *SmsStatus {
	body := make(url.Values)
	sd.addBasicParams(body)
	body.Add("Msnlist", mobileNo)
	body.Add("Msg", message)

	httpReq, err := http.NewRequest(http.MethodPost, sd.smsApiUrl, bytes.NewReader([]byte(body.Encode())))
	if err != nil {
		panic(err)
	}
	httpReq.Header.Set(headerUserAgent, "Mozilla/5.0")
	httpReq.Header.Set(headerContentType, "application/x-www-form-urlencoded")

	statusSms := new(SmsStatus)
	resp, err := sd.client.Do(httpReq)
	if err != nil {
		statusSms.HttpStatus = resp.StatusCode
		statusSms.Status = STATUS_ERROR
		statusSms.Reason = err.Error()
	} else if resp.StatusCode != http.StatusOK {
		statusSms.HttpStatus = resp.StatusCode
		statusSms.Status = STATUS_ERROR
	} else {
		statusSms.Status = STATUS_SUCCESS
	}
	return statusSms
}

func (sd *SmsSender) addBasicParams(body url.Values) {
	body.Add("User", sd.user)
	body.Add("Password", sd.password)
	body.Add("Sender", sd.sender)
}

const (
	STATUS_SUCCESS = iota + 1
	STATUS_ERROR
)

type SmsStatus struct {
	Status     int
	HttpStatus int
	Reason     string
}
