package main

import (
	"dolo-tracking/context"
	"dolo-tracking/email"
	"dolo-tracking/hubspot"
	"dolo-tracking/logger"
	"flag"
	"fmt"
	"os"
	"time"
)

// FIXME all

const emailBody = `Bonjour,<br /><br />

Je me demandais si vous aviez eu quelques minutes pour écouter le nouveau single de Doloréanne : Comme une actrice. L'extrait été lancé sur 45tours.ca le 19 septembre dernier. Je vous laisse l'extrait sur le Bandcamp du groupe https://doloreanne.bandcamp.com/track/comme-une-actrice.  Sachez que votre critique par rapport à votre programmation nous ferait le plus grand des plaisirs.<br /><br />

Merci beaucoup et n'hésitez pas à m'appeler si vous voulez en discuter de vive voix. <br /><br />

---<br />
Alexandre Vézina<br />
(581) 982-5190`

const emailSubject = "Doloréanne: Comme une actrice"

const fromEmail = "alex@doloreanne.com"
const fromFirstname = "Alexandre"
const fromLastname = "Vézina"

const hubspotPipeline = "3682c689-4605-4437-abde-e0604828bf06"

// https://app.hubspot.com/property-settings/2213414/deal/dealstage
const hubspotDealstage = "aecf09f6-1e9a-4a4f-b9fa-384218cf6214"

func newConfiguration(hubspotKey string, sparkpostKey string) (*context.Configuration, error) {
	return &context.Configuration{
		SparkPost: context.SparkPostConfig{
			APIKey: sparkpostKey,
		},
		Hubspot: context.HubspotConfig{
			APIKey: hubspotKey,
		},
	}, nil
}

func buildContext(hubspotKey string, sparkpostKey string) (*context.App, error) {
	var (
		err    error
		config *context.Configuration
		ctx    *context.App
	)

	if config, err = newConfiguration(hubspotKey, sparkpostKey); err != nil {
		return ctx, err
	}

	ctx = &context.App{
		Config: *config,
	}

	return ctx, nil
}

func sendEmail(ctx *context.App, to string) error {
	sp := email.NewMapperSparkpost(ctx.Config.SparkPost.APIKey, fromEmail)

	return sp.SendHTML(email.HTMLEmailOpts{
		HTML:     emailBody,
		To:       to,
		FromName: fromFirstname + " " + fromLastname,
		Subject:  emailSubject,
	})
}

func main() {
	var (
		err           error
		ctx           *context.App
		companyList   []hubspot.Company
		contactIDList []int
		contactList   []hubspot.Contact
		contact       *hubspot.Contact
		deal          *hubspot.Deal
	)

	hubspotKey := flag.String("hubspot", "", "Hubspot API key")
	sparkpostAPIKey := flag.String("sparkpost", "", "SparkPost API key")
	flag.Parse()

	if *hubspotKey == "" || *sparkpostAPIKey == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Build the app context
	if ctx, err = buildContext(*hubspotKey, *sparkpostAPIKey); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Find the companies
	if companyList, err = hubspot.FindCompanies(ctx.Config.Hubspot.APIKey); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	i := 0
	sent := 0
	for _, company := range companyList {
		i++
		logger.Debug(fmt.Sprintf("Processing company '%s'", company.Name))

		// Find the contacts
		if contactIDList, err = hubspot.GetCompanyContactList(ctx.Config.Hubspot.APIKey, company.CompanyID); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		// Find the contacts
		contactList = []hubspot.Contact{}
		for _, contactID := range contactIDList {
			if contact, err = hubspot.GetContact(ctx.Config.Hubspot.APIKey, contactID); err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
			contactList = append(contactList, *contact)

			time.Sleep(time.Millisecond * 133)
		}
		if len(contactList) == 0 {
			continue
		}

		// Find the deal or create it if missing
		if deal, err = hubspot.FindDeal(ctx.Config.Hubspot.APIKey, company.CompanyID, hubspotPipeline); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		if deal == nil {
			sent++
			// Send the email
			if err = sendEmail(ctx, contactList[0].Email); err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			// Add the deal
			if deal, err = hubspot.AddDeal(ctx.Config.Hubspot.APIKey, company, contactList[0], hubspotPipeline, hubspotDealstage); err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
			if deal == nil {
				logger.Error("No deal returned")
				os.Exit(1)
			}

			// Log the email
			emailMeta := hubspot.MetadataEmail{
				From: hubspot.MetadataEmailFrom{
					Email:     fromEmail,
					Firstname: fromFirstname,
					Lastname:  fromLastname,
				},
				To: []hubspot.MetadataEmailTo{
					hubspot.MetadataEmailTo{
						Email: contact.Email,
					},
				},
				Subject: emailSubject,
				HTML:    emailBody,
			}
			if err = hubspot.AddEngagementEmail(ctx.Config.Hubspot.APIKey, company, contactList[0], *deal, emailMeta); err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

		}

		time.Sleep(time.Millisecond * 1000)
	}
	fmt.Printf("\nSent %d emails out of %d companies\n", sent, i)
}
