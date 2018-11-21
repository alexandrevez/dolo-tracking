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

const emailBody = `Bonjour,<br /><br />

Je me demandais si vous aviez vu passer le nouveau single de Doloréanne : Jusqu'au bout. 
L'extrait a été lancé sur postedecoute.ca le 21 octobre dernier et est en recommandation 
radios amplifiées sur <a href="https://palmaresadisq.ca/fr/sur-les-ondes/radio-amplifiees/" target="_BLANK">Palmares ADISQ</a> 
en ce moment. 
<br /><br />
Je vous laisse ausi l'extrait sur le Bandcamp pour une écoute rapide <a href="https://doloreanne.bandcamp.com/track/jusquau-bout" target="_BLANK">https://doloreanne.bandcamp.com/track/jusquau-bout</a>. 
<br/><br />
Sachez que votre critique par rapport à votre programmation nous fait, comme toujours, 
le plus grand des plaisirs.<br /><br />

Merci beaucoup et n'hésitez pas à m'appeler si vous voulez en discuter de vive voix. <br /><br />

---<br />
Alexandre Vézina<br />
(581) 982-5190`

const emailSubject = "Doloréanne : Jusqu'au bout"

const fromEmail = "alex@doloreanne.com"
const fromFirstname = "Alexandre"
const fromLastname = "Vézina"

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
	)

	hubspotKey := flag.String("hubspot", "", "Hubspot API key")
	sparkpostAPIKey := flag.String("sparkpost", "", "SparkPost API key")
	flag.Parse()

	if *hubspotKey == "" || *sparkpostAPIKey == "" {
		flag.Usage()
		os.Exit(1)
	}

	if ctx, err = buildContext(*hubspotKey, *sparkpostAPIKey); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	if companyList, err = hubspot.FindCompanies(ctx.Config.Hubspot.APIKey); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	i := 0
	sent := 0
	for _, company := range companyList {
		i++
		logger.Debug(fmt.Sprintf("Processing company '%s'", company.Name))

		if contactIDList, err = hubspot.GetCompanyContactList(ctx.Config.Hubspot.APIKey, company.CompanyID); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

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
		if contactList[0].Email == "" {
			logger.Warn(fmt.Sprintf("No email found for deal: %s", company.Name))
			continue
		}

		logger.Debug(fmt.Sprintf("Sending email to %s", contactList[0].Email))
		if err = sendEmail(ctx, contactList[0].Email); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		sent++

		time.Sleep(time.Millisecond * 2000)
	}
	logger.Debug(fmt.Sprintf("Sent %d emails in %d companies\n", sent, i))
}
