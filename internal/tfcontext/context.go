// Package tfcontext provides a standardized way to represent and validate context keys for Terraform credentials management.
package tfcontext

import (
	"errors"
	"strings"
)

// Context represents a standardized context key used for managing Terraform credentials.
type Context struct {
	Key string
}

// Parse standardizes and validates a given context key string.
func Parse(key string) (Context, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return Context{}, errors.New("context key cannot be empty or blank")
	}

	// The string is the direct map key. No more splitting or complex string sanitizations.
	return Context{
		Key: key,
	}, nil
}
