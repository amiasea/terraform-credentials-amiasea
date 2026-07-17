// Package main
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

var version = "dev"

type creds struct {
	Token string `json:"token"`
}

func main() {
	helperLog(fmt.Sprintf(
		"starting credentials helper version=%s args=%v",
		version,
		os.Args,
	))

	if len(os.Args) < 2 {
		fmt.Println(version)
		fmt.Fprintln(os.Stderr, "terraform-credentials-tfcred is a Terraform credentials helper")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "get":
		handleGet()

	case "store", "forget":
		// tfcred manages credentials outside Terraform's helper lifecycle.
		// These commands are intentionally no-ops.
		helperLog(fmt.Sprintf("ignored command=%s", os.Args[1]))
		os.Exit(0)

	default:
		helperLog(fmt.Sprintf("unexpected command=%s", os.Args[1]))
		os.Exit(1)
	}
}

func handleGet() {
	if len(os.Args) < 3 {
		helperLog("get command missing hostname")
		os.Exit(1)
	}

	requestedDomain := os.Args[2]

	cwd, err := os.Getwd()
	if err != nil {
		helperLog(fmt.Sprintf("failed getting current directory: %v", err))
		os.Exit(1)
	}

	helperLog(fmt.Sprintf(
		"credential request hostname=%s cwd=%s",
		requestedDomain,
		cwd,
	))

	contextKey, found := store.ResolveContextByDir(cwd)
	if !found {
		helperLog("credential request failed: no matching context")
		writeCredentials("")
		return
	}

	storeFile := store.Load()

	entry, ok := storeFile.Contexts[contextKey]
	if !ok {
		helperLog(fmt.Sprintf(
			"credential request failed: unknown context=%s",
			contextKey,
		))
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
		helperLog(fmt.Sprintf(
			"credential request failed: domain mismatch requested=%s configured=%s",
			requestedDomain,
			contextDomain,
		))
		writeCredentials("")
		return
	}

	token, err := resolve.Resolve(tfcontext.Context{
		Key: contextKey,
	})
	if err != nil {
		helperLog(fmt.Sprintf(
			"credential resolution failed context=%s error=%v",
			contextKey,
			err,
		))
		os.Exit(1)
	}

	helperLog(fmt.Sprintf(
		"credential request succeeded hostname=%s context=%s",
		requestedDomain,
		contextKey,
	))

	writeCredentials(token)
}

func writeCredentials(token string) {
	if err := json.NewEncoder(os.Stdout).Encode(creds{
		Token: token,
	}); err != nil {
		helperLog(fmt.Sprintf(
			"failed writing credentials response: %v",
			err,
		))
		os.Exit(1)
	}
}

func helperLog(msg string) {
	logPath := os.Getenv("TFCRED_CREDENTIALS_HELPER_LOG")
	if logPath == "" {
		return
	}

	log.AppendFile(logPath, msg)
}
