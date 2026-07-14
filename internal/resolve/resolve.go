// Package resolve provides functionality to resolve credentials based on a given context.
package resolve

import (
	"fmt"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
	"github.com/amiasea/terraform-credentials-tfcred/internal/tfcontext"
)

// Resolve retrieves the token associated with a context key from the credential store.
func Resolve(ctx tfcontext.Context) (string, error) {
	if ctx.Key == "" {
		return "", fmt.Errorf(
			"cannot resolve credential: target context key is empty",
		)
	}

	storeFile := store.Load()

	entry, exists := storeFile.Contexts[ctx.Key]
	if !exists {
		return "", fmt.Errorf(
			"cannot resolve credential: context '%s' does not exist",
			ctx.Key,
		)
	}

	token, err := store.GetToken(entry)
	if err != nil {
		return "", fmt.Errorf(
			"failed to retrieve token for context '%s': %w",
			ctx.Key,
			err,
		)
	}

	if token == "" {
		return "", fmt.Errorf(
			"retrieved token for context '%s' is empty",
			ctx.Key,
		)
	}

	return token, nil
}
