package email

import (
	sp "github.com/SparkPost/gosparkpost"
)

// HTMLEmailOpts represents the options needed to send an HTML email
type HTMLEmailOpts struct {
	HTML     string
	To       string
	FromName string
	Subject  string
}

// SparkPostURL is the URL to hit the API
const SparkPostURL = "https://api.sparkpost.com"

// MapperSparkpost is an implementation of email mapper for SparkPost
type MapperSparkpost struct {
	Config      sp.Config
	FromAddress string
}

// NewMapperSparkpost create a sparkpost mapper
func NewMapperSparkpost(apiKey string, fromAddress string) (m MapperSparkpost) {
	return MapperSparkpost{
		Config: sp.Config{
			BaseUrl:    SparkPostURL,
			ApiKey:     apiKey,
			ApiVersion: 1,
		},
		FromAddress: fromAddress,
	}
}

// SendHTML sends an html email using Sparkpost service
func (m MapperSparkpost) SendHTML(opts HTMLEmailOpts) error {
	var (
		err    error
		client sp.Client
		tx     sp.Transmission
	)

	if err = client.Init(&m.Config); err != nil {
		return err
	}

	tx = sp.Transmission{
		Recipients: []string{opts.To},
		Content: sp.Content{
			HTML:    opts.HTML,
			Subject: opts.Subject,
			From: sp.From{
				Email: m.FromAddress,
				Name:  opts.FromName,
			},
		},
	}

	if _, _, err = client.Send(&tx); err != nil {
		return err
	}

	return nil
}
