package hubspot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Contact contains the properties as represented in the Hubspot API
type Contact struct {
	ContactID int
	Email     string
}

// AddContactRequest is a request to add contacts
type AddContactRequest struct {
	Properties []Property `json:"properties"`
}

// ContactResponseResult is an actual search result
type ContactResponseResult struct {
	ContactID  int                             `json:"vid"`
	Properties ContactResponseResultProperties `json:"properties"`
}

// ContactResponseResultProperties properties are wrapped... thanks hubspot
type ContactResponseResultProperties struct {
	Email SearchProperty `json:"email"`
}

// AddContact will add the contact in Hubspot
// http://developers.hubspot.com/docs/methods/contacts/v2/get_contacts_properties
func AddContact(apiKey string, email string) (*Contact, error) {
	const (
		hubspotURL = "https://api.hubapi.com/contacts/v1/contact/?hapikey=%s"
	)
	var (
		err        error
		url        string
		req        AddContactRequest
		reqBytes   []byte
		httpClient http.Client
		httpReq    *http.Request
		respRaw    *http.Response
		bodyBytes  []byte
		resp       ContactResponseResult
	)

	req = AddContactRequest{
		Properties: []Property{
			Property{
				Property: "email",
				Value:    email,
			},
		},
	}

	if reqBytes, err = json.Marshal(req); err != nil {
		return nil, err
	}
	url = fmt.Sprintf(hubspotURL, apiKey)

	if httpReq, err = http.NewRequest("POST", url, bytes.NewBuffer(reqBytes)); err != nil {
		return nil, err
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")

	httpClient = http.Client{}
	if respRaw, err = httpClient.Do(httpReq); err != nil {
		return nil, err
	}
	defer respRaw.Body.Close()

	if bodyBytes, err = ioutil.ReadAll(respRaw.Body); err != nil {
		return nil, err
	}

	if respRaw.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: %s \n%s", respRaw.Status, string(bodyBytes))
	}

	if err = json.Unmarshal(bodyBytes, &resp); err != nil {
		return nil, err
	}

	return &Contact{
		ContactID: resp.ContactID,
		Email:     email,
	}, nil
}

// GetContact returns a contact with this email
func GetContact(apiKey string, id int) (*Contact, error) {
	const (
		hubspotURL = "https://api.hubapi.com/contacts/v1/contact/vid/%d/profile?hapikey=%s"
	)

	var (
		err        error
		url        string
		httpClient http.Client
		req        *http.Request
		resp       ContactResponseResult
		respRaw    *http.Response
		bodyBytes  []byte
	)

	url = fmt.Sprintf(hubspotURL, id, apiKey)

	if req, err = http.NewRequest("GET", url, bytes.NewBuffer([]byte{})); err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	httpClient = http.Client{}
	if respRaw, err = httpClient.Do(req); err != nil {
		return nil, err
	}
	defer respRaw.Body.Close()

	if respRaw.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if respRaw.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: %s", respRaw.Status)
	}

	if bodyBytes, err = ioutil.ReadAll(respRaw.Body); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(bodyBytes, &resp); err != nil {
		return nil, err
	}

	return &Contact{
		ContactID: resp.ContactID,
		Email:     resp.Properties.Email.Value,
	}, nil
}
