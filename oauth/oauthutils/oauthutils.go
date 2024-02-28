package oauthutils

import (
	"app/files/filemodes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
)

type credentials struct {
	Installed struct {
		ClientID                string   `json:"client_id"`
		ProjectID               string   `json:"project_id"`
		AuthURI                 string   `json:"auth_uri"`
		TokenURI                string   `json:"token_uri"`
		AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url"`
		ClientSecret            string   `json:"client_secret"`
		RedirectURIs            []string `json:"redirect_uris"`
	} `json:"installed"`
}

// this function basically edits the google's credentials.json
// to add the right redirect_uris while keeping the rest of
// the file intact
func FixCredentialsJSON(path string, port string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var c credentials
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return err
	}

	// we check if any of all the fields inside c.Installed
	// are empty strings; if they are, chances are the
	// credentials.json is not a valid file
	if hasEmptyStringFields(c.Installed) {
		return errors.New("wrong credentials file")
	}

	c.Installed.RedirectURIs = []string{fmt.Sprintf("http://localhost%v/", port)}

	newJSON, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, newJSON, filemodes.RW_R__R__)
	if err != nil {
		return err
	}

	return nil
}

func hasEmptyStringFields(s interface{}) bool {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if field.Kind() == reflect.String {
			if field.String() == "" {
				return true
			}
		}
	}

	return false
}
