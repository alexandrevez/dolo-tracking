package hubspot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Deal represents a deal yo
type Deal struct {
	DealID   int
	Name     string
	Pipeline string
}

// FindDealResponse is what is returned from FindDeal
type FindDealResponse struct {
	Results []FindDealResponseElem `json:"deals"`
	HasMore bool                   `json:"hasMore"`
	Offset  int                    `json:"offset"`
}

// FindDealResponseElem is an element that contains what we need in a FindDealResponses
type FindDealResponseElem struct {
	ID           int                        `json:"dealId"`
	Associations DealAssociations           `json:"associations"`
	Properties   FindDealResponseProperties `json:"properties"`
}

// DealAssociations is associations to deals
type DealAssociations struct {
	CompanyIDList []int `json:"associatedCompanyIds"`
	ContactIDList []int `json:"associatedVids"`
}

// FindDealResponseProperties is the properties needed to find a deal
type FindDealResponseProperties struct {
	Name     SearchProperty `json:"dealname"`
	Pipeline SearchProperty `json:"pipeline"`
}

// FindDeal will try to find if a company has a deal in pipeline
func FindDeal(apiKey string, companyID int, pipeline string) (*Deal, error) {
	const (
		hubspotURL = "https://api.hubapi.com/deals/v1/deal/paged?hapikey=%s&includeAssociations=true&limit=250&properties=dealname&properties=pipeline&offset=%d"
	)
	var (
		err           error
		url           string
		httpClient    http.Client
		req           *http.Request
		resp          FindDealResponse
		respRaw       *http.Response
		respBodyBytes []byte
		deal          *Deal
	)

	httpClient = http.Client{}
	url = fmt.Sprintf(hubspotURL, apiKey, 0)

	for {
		if req, err = http.NewRequest("GET", url, bytes.NewBuffer([]byte{})); err != nil {
			return deal, err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		if respRaw, err = httpClient.Do(req); err != nil {
			return deal, err
		}
		defer respRaw.Body.Close()

		if respBodyBytes, err = ioutil.ReadAll(respRaw.Body); err != nil {
			return deal, err
		}

		if respRaw.StatusCode != http.StatusOK {
			return deal, fmt.Errorf("Error: %s \n%s", respRaw.Status, string(respBodyBytes))
		}

		if err = json.Unmarshal(respBodyBytes, &resp); err != nil {
			return deal, err
		}

		for _, dealResult := range resp.Results {
			if dealResult.Properties.Pipeline.Value != pipeline {
				continue
			}

			for _, assocCompanyID := range dealResult.Associations.CompanyIDList {
				if assocCompanyID == companyID {
					return &Deal{
						DealID:   dealResult.ID,
						Name:     dealResult.Properties.Name.Value,
						Pipeline: dealResult.Properties.Pipeline.Value,
					}, nil
				}
			}
		}

		if !resp.HasMore {
			break
		}

		url = fmt.Sprintf(hubspotURL, apiKey, resp.Offset)
	}
	return nil, nil
}

// AddDealRequest is the elements we add to the deal
type AddDealRequest struct {
	Associations DealAssociations `json:"associations"`
	Properties   []Property       `json:"properties"`
}

// AddDealResult is the result of add deal
type AddDealResult struct {
	DealID int `json:"dealId"`
}

// AddDeal will add a deal in hubspot with associations to company and contact
func AddDeal(apiKey string, company Company, contact Contact, pipeline string, dealstage string) (*Deal, error) {
	const (
		hubspotURL = "https://api.hubapi.com/deals/v1/deal?hapikey=%s"
	)
	var (
		err        error
		url        string
		req        AddDealRequest
		reqBytes   []byte
		httpClient http.Client
		httpReq    *http.Request
		respRaw    *http.Response
		bodyBytes  []byte
		resp       AddDealResult
	)

	req = AddDealRequest{
		Associations: DealAssociations{
			CompanyIDList: []int{company.CompanyID},
			ContactIDList: []int{contact.ContactID},
		},
		Properties: []Property{
			Property{
				Name:  "dealname",
				Value: company.Name,
			},
			Property{
				Name:  "pipeline",
				Value: pipeline,
			},
			Property{
				Name:  "dealstage",
				Value: dealstage,
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

	return &Deal{
		DealID:   resp.DealID,
		Name:     company.Name,
		Pipeline: pipeline,
	}, nil
}
