package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jefflinse/githubsecret"
)

type github struct {
	username string
	token    string
}

func (gh github) getPublicKey(owner string, repo string) (publicKey, error) {
	path := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/secrets/public-key", owner, repo)
	key := publicKey{}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return key, err
	}

	// https://developer.github.com/v3/#current-version
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.SetBasicAuth(gh.username, gh.token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return key, err
	}

	if res.StatusCode != 200 {
		return key, fmt.Errorf("HTTP status %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return key, err
	}

	if err := json.Unmarshal(body, &key); err != nil {
		return key, err
	}

	return key, nil
}

func (gh github) storeSecret(owner string, repo string, key publicKey, secretID string, secretValue string) error {
	path := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/secrets/%s", owner, repo, secretID)

	encryptedValue, err := githubsecret.Encrypt(key.Key, secretValue)
	if err != nil {
		return err
	}

	body := struct {
		KeyID          string `json:"key_id"`
		EncryptedValue string `json:"encrypted_value"`
	}{
		KeyID:          key.KeyID,
		EncryptedValue: encryptedValue,
	}

	params, err := json.Marshal(&body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", path, ioutil.NopCloser(bytes.NewBuffer(params)))
	if err != nil {
		return err
	}

	// https://developer.github.com/v3/#current-version
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.SetBasicAuth(gh.username, gh.token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 201 && res.StatusCode != 204 {
		return fmt.Errorf("HTTP status %d", res.StatusCode)
	}

	return nil
}
