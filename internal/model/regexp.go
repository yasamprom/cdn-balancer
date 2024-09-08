package usecases

import "regexp"

var (
	RegexpLink = regexp.MustCompile(`^(https?://)([^/]+)(/.*)?(\\?.*)?$`)
)
