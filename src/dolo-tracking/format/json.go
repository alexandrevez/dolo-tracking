package format

import (
	"dolo-tracking/logger"
	"encoding/json"
)

// NewJSONString converts an interface into a well formatted JSON string
func NewJSONString(i interface{}) string {
	data, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		logger.Warn(err.Error())
		return ""
	}
	return string(data)
}
