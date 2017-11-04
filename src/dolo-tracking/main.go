package main

import (
	"dolo-tracking/context"
	"dolo-tracking/email"
	"dolo-tracking/format"
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

const subject = "Doloréanne: Comme une actrice"

const fromEmail = "alex@doloreanne.com"
const fromName = "Alexandre Vézina"

const hubspotPipeline = "3682c689-4605-4437-abde-e0604828bf06"

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
		FromName: fromName,
		Subject:  subject,
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

	for _, company := range companyList {
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
			fmt.Println(format.NewJSONString(*contact))
			contactList = append(contactList, *contact)

			time.Sleep(time.Millisecond * 133)
		}
		if len(contactList) == 0 {
			continue
		}

		// // Find the deal or create it if missing
		// if deal, err = hubspot.FindDeal(ctx.Config.Hubspot.APIKey, company.CompanyID, hubspotPipeline); err != nil {
		// 	logger.Error(err.Error())
		// 	os.Exit(1)
		// }
		// if deal == nil {
		// 	if err = sendEmail(ctx, []string{"avezina@ubikvoip.com"}); err != nil {
		// 		logger.Error(err.Error())
		// 		os.Exit(1)
		// 	}
		// 	if err = hubspot.AddDeal(ctx.Config.Hubspot.APIKey, company, hubspotPipeline); err != nil {
		// 		logger.Error(err.Error())
		// 		os.Exit(1)
		// 	}
		// }

		time.Sleep(time.Millisecond * 1000)
	}

}
