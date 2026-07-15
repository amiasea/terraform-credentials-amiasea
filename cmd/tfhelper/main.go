package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/amiasea/terraform-credentials-tfcred/internal/log"
	"github.com/amiasea/terraform-credentials-tfcred/internal/resolve"
	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
	"github.com/amiasea/terraform-credentials-tfcred/internal/tfcontext"
)

type creds struct {
	Token string `json:"token"`
}

func main() {
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "get":
			handleGet()
		case "store", "forget":
			// Gracefully ignore
			os.Exit(0)
		default:
			fmt.Fprintln(os.Stderr, "This program should only be invoked by Terraform.")
			os.Exit(1)
		}
	} else {
		fmt.Fprintln(os.Stderr, "This program should only be invoked by Terraform.")
		os.Exit(1)
	}
}

func handleGet() {
	if len(os.Args) < 3 {
		log.Err("get command requires a hostname/domain")
		os.Exit(1)
	}

	requestedDomain := os.Args[2]

	// === All your existing get logic goes here (exactly as you posted) ===
	cwd, err := os.Getwd()
	if err != nil {
		log.Err(fmt.Sprintf("failed to get current working directory: %v", err))
		os.Exit(1)
	}

	contextKey, found := store.ResolveContextByDir(cwd)
	if !found {
		_ = json.NewEncoder(os.Stdout).Encode(creds{Token: ""})
		os.Exit(0)
	}

	storeFile := store.Load()
	entry, ok := storeFile.Contexts[contextKey]
	if !ok {
		log.Err(fmt.Sprintf("unknown context key '%s'", contextKey))
		os.Exit(1)
	}

	contextDomain := entry.Domain
	if contextDomain == "" {
		contextDomain = storeFile.DefaultDomain
	}
	if contextDomain == "" {
		contextDomain = "app.terraform.io"
	}

	if contextDomain != requestedDomain {
		_ = json.NewEncoder(os.Stdout).Encode(creds{Token: ""})
		os.Exit(0)
	}

	activeCtx := tfcontext.Context{Key: contextKey}
	token, err := resolve.Resolve(activeCtx)

	fmt.Fprintln(os.Stderr, "resolved token:", token)

	if err != nil {
		log.Err(err.Error())
		os.Exit(1)
	}

	_ = json.NewEncoder(os.Stdout).Encode(creds{Token: token})
	os.Exit(0)
}
