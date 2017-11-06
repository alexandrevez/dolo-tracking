package hubspot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// EngagementTypeEmail is an email type
const EngagementTypeEmail = "EMAIL"

// AddEngagementRequest is the elements we add to the deal
type AddEngagementRequest struct {
	Engagement   Engagement             `json:"engagement"`
	Associations EngagementAssociations `json:"associations"`
	Metadata     interface{}            `json:"metadata"`
}

// EngagementAssociations is associations to deals
type EngagementAssociations struct {
	DealIDList []int `json:"dealIds"`
}

// Engagement defines what type of engagement
type Engagement struct {
	Type string `json:"type"`
}

// AddEngagementResult is the result of add deal
type AddEngagementResult struct {
	DealID int `json:"dealId"`
}

// MetadataEmail is what is listed in an email engagement
type MetadataEmail struct {
	From    MetadataEmailFrom `json:"from"`
	To      []MetadataEmailTo `json:"to"`
	Subject string            `json:"subject"`
	HTML    string            `json:"html"`
}

// MetadataEmailFrom defines who sent the email
type MetadataEmailFrom struct {
	Email     string `json:"email"`
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
}

// MetadataEmailTo defines a recipient
type MetadataEmailTo struct {
	Email string `json:"email"`
}

// AddEngagementEmail will add a deal in hubspot with associations to company and contact
func AddEngagementEmail(apiKey string, company Company, contact Contact, deal Deal, meta MetadataEmail) error {
	const (
		hubspotURL = "https://api.hubapi.com/engagements/v1/engagements?hapikey=%s"
	)
	var (
		err        error
		url        string
		req        AddEngagementRequest
		reqBytes   []byte
		httpClient http.Client
		httpReq    *http.Request
		respRaw    *http.Response
		bodyBytes  []byte
	)

	req = AddEngagementRequest{
		Engagement: Engagement{
			Type: EngagementTypeEmail,
		},
		Associations: EngagementAssociations{
			DealIDList: []int{deal.DealID},
		},
		Metadata: meta,
	}

	if reqBytes, err = json.Marshal(req); err != nil {
		return err
	}

	url = fmt.Sprintf(hubspotURL, apiKey)

	if httpReq, err = http.NewRequest("POST", url, bytes.NewBuffer(reqBytes)); err != nil {
		return err
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")

	httpClient = http.Client{}
	if respRaw, err = httpClient.Do(httpReq); err != nil {
		return err
	}
	defer respRaw.Body.Close()

	if bodyBytes, err = ioutil.ReadAll(respRaw.Body); err != nil {
		return err
	}

	if respRaw.StatusCode != http.StatusOK {
		return fmt.Errorf("Error: %s \n%s", respRaw.Status, string(bodyBytes))
	}
	return nil
}
