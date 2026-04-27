package errs

import (
	"encoding/json"

	"github.com/witwoywhy/go-cores/enum/language"
	"github.com/witwoywhy/go-cores/logger"
)

type ErrorCodeMapping map[string]map[language.Language]ErrorCodeMappingMessage

type ErrorCodeMappingMessage struct {
	Message string `mapstructure:"message"`
}

func ParseToErrorCodeMapping(v string, l logger.Logger) ErrorCodeMapping {
	var mapping ErrorCodeMapping
	err := json.Unmarshal([]byte(v), &mapping)
	if err != nil {
		l.Warnf("failed to json.Unmarshal parse to error code mapping: %v", err)
		return nil
	}

	return mapping
}
