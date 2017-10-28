package context

// Configuration defines the configuration of the application
type Configuration struct {
	SparkPost SparkPostConfig `json:"sparkpost"`
}

// SparkPostConfig defines the configuration of the Sparkpost email service
type SparkPostConfig struct {
	APIKey      string `json:"api_key"`
	FromAddress string `json:"from_address"`
}
