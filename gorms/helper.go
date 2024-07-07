package gorms

import (
	"strings"
)

func getTableNameFromQuery(query string) string {
	query = strings.Replace(query, "`", "", -1)
	query = strings.Replace(query, `"`, "", -1)

	match := regex.FindStringSubmatch(query)
	if len(match) > 1 {
		if match[1] == "" {
			return match[2]
		}
		return match[1]
	}
	return ErrCantFindTableName
}
