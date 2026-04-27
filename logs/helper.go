package logs

import (
	"github.com/witwoywhy/go-cores/apps"
	"github.com/witwoywhy/go-cores/utils"
)

func masking(m map[string]any) {
	for k, v := range m {
		if v == nil {
			continue
		}

		switch val := v.(type) {
		case map[string]any:
			masking(val)
		case string:
			if _, sensitive := maskingList[k]; sensitive {
				m[k] = utils.MaskString(val, apps.MaskingChar)
			}
		case []any:
			for _, item := range val {
				if mapVal, ok := item.(map[string]any); ok {
					masking(mapVal)
				}
			}
		}
	}
}
