package utils

import (
	"regexp"
	"strings"
)

func SimpleSlug(s string) string {
    s = strings.ToLower(s)
    s = strings.TrimSpace(s)
    s = strings.ReplaceAll(s, " ", "-")
    return regexp.MustCompile(`[^a-z0-9\-]`).ReplaceAllString(s, "")
}