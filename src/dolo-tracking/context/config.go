package context

// Configuration defines the configuration of the application
type Configuration struct {
	SparkPost SparkPostConfig
	Hubspot   HubspotConfig
}

// HubspotConfig defines the configuration of the Hubspot API
type HubspotConfig struct {
	APIKey string
}

// SparkPostConfig defines the configuration of the Sparkpost email service
type SparkPostConfig struct {
	APIKey string
}
