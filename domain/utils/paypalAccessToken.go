package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GeneratePayPalAccessToken() (string, error) {
	clientID := "AcAvim6aHXQxhW3XryTaxzQf-PmqnKnY9Mlq_bliOCoD55clHs5O-7hDFshYDQ6TRPmvMBgMFLx-4FYq"
	clientSecret := "EOJGhmFQymy8TajUpbHk2bcyzVvzBb3w-sjhzdt8_CgwF9zhs_uNLDAqZgVYBADCOqzOzgl6po0bNYAV"
	baseURL := "https://api-m.sandbox.paypal.com/v1/oauth2/token"

	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", baseURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return "", errors.New("Failed to create request:")
	}

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("Failed to send request:")
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Failed to read response:")
	}

	var accessTokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal(data, &accessTokenResponse)
	if err != nil {
		return "", errors.New("Failed to unmarshal response:")
	}

	return accessTokenResponse.AccessToken, nil
}
