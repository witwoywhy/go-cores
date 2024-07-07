package gorms

import "regexp"

const (
	StartOutbound = "START OUTBOUND"
	EndOutbound   = "END OUTBOUND | %v | %s"
)

const (
	ErrCantFindTableName = "can't find table name"
)

var (
	regex = regexp.MustCompile(`(?i)^(?:INSERT INTO|UPDATE|DELETE FROM)\s+(\S+)|SELECT\s+.+?\bFROM\s+(\S+)`)
)
