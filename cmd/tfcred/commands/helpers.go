// Package commands contains all CLI subcommands for the tfcred tool.
package commands

import (
	"fmt"
	"regexp"
	"strings"
)

// supportedDomains lists all officially supported Terraform domains.
var supportedDomains = []string{
	"app.terraform.io",
	"app.eu.terraform.io",
}

// isSupportedDomain returns true if the given domain is supported.
func isSupportedDomain(domain string) bool {
	for _, d := range supportedDomains {
		if d == domain {
			return true
		}
	}
	return false
}

// promptDefaultDomain prompts the user for a default domain.
func promptDefaultDomain() string {
	fmt.Print("Enter default domain [app.terraform.io]: ")
	var domain string
	_, _ = fmt.Scanln(&domain)
	if domain == "" {
		return "app.terraform.io"
	}
	return domain
}

var tokenFormatPattern = regexp.MustCompile(`^[a-zA-Z0-9]{14}\.atlasv1\.[a-zA-Z0-9]{30,70}$`)

func isValidTokenFormat(token string) bool {
	return tokenFormatPattern.MatchString(strings.TrimSpace(token))
}
