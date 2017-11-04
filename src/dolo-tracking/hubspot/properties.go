package hubspot

// Property good job hubspot Name and Property are actually the same, but used in different contexts
type Property struct {
	Name     string      `json:"name"`
	Property string      `json:"property"`
	Value    interface{} `json:"value"`
}
