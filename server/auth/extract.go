package auth

import (
	"strings"
)

const BearerPrefix = "Bearer "

func ExtractBearerToken(authHeader string) string {
	if strings.HasPrefix(authHeader, BearerPrefix) {
		return strings.TrimPrefix(authHeader, BearerPrefix)
	}
	return ""
}
