package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type Payload struct {
	Event      string                 `json:"event"`
	Repository string                 `json:"repository"`
	Commit     string                 `json:"commit"`
	Ref        string                 `json:"ref"`
	Head       string                 `json:"head"`
	Workflow   string                 `json:"workflow"`
	Extras     map[string]interface{} `json:"extras,omitempty"`
}

var (
	webhookURL, webhookSecret string
)

func main() {
	hookUrl, ok := os.LookupEnv("WEBHOOK_URL")
	if !ok {
		fmt.Println("::error title=Error::WEBHOOK_URL env not set")
		return
	}
	secret, ok := os.LookupEnv("WEBHOOK_SECRET")
	if !ok {
		fmt.Println("::error title=Error::WEBHOOK_SECRET env not set")
		return
	}
	u, err := url.Parse(hookUrl)
	if err != nil {
		fmt.Println("::error title=Error::Not a valid url")
		return
	}
	if u.Scheme != "https" {
		fmt.Println("::error title=Error::WEBHOOK_URL is not https")
		return
	}

	webhookURL = u.String()
	webhookSecret = secret

	payload := Payload{
		Event:      os.Getenv("GITHUB_EVENT_NAME"),
		Repository: os.Getenv("GITHUB_REPOSITORY"),
		Commit:     os.Getenv("GITHUB_SHA"),
		Ref:        os.Getenv("GITHUB_REF"),
		Head:       os.Getenv("GITHUB_HEAD_REF"),
		Workflow:   os.Getenv("GITHUB_WORKFLOW"),
	}

	json.Unmarshal([]byte(os.Getenv("EXTRAS")), &payload.Extras)

	if err := SendWebhook(payload); err != nil {
		fmt.Printf("::error title=Webhook error::Error sending webhook %s", err)
		return
	}

}

func SendWebhook(payload Payload) error {

	p, _ := json.Marshal(payload)

	hmac := ComputeHmac256(p, webhookSecret)

	client := &http.Client{}
	req, err := http.NewRequest("POST", webhookURL, bytes.NewReader(p))
	if err != nil {
		return err
	}

	fmt.Println(hmac)

	req.Header.Set("X-Hub-Signature", fmt.Sprintf("sha256=%s", hmac))
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("Invalid status code %d", res.StatusCode)
}

func ComputeHmac256(message []byte, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write(message)
	return hex.EncodeToString(h.Sum(nil))
}
