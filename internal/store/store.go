// Package store provides access to credential manager and contexts.json.
package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	vaultpkg "github.com/amiasea/terraform-credentials-tfcred/internal/vault"
)

var credentialVault vaultpkg.Vault = vaultpkg.NewWindowsVault()

// Entry represents a single context entry in the contexts.json file.
type Entry struct {
	Org       string `json:"org"`
	TokenType string `json:"tokenType"`
	Domain    string `json:"domain"`
}

// File represents the contents of tfcred_contexts.json.
type File struct {
	Contexts      map[string]Entry  `json:"contexts"`
	Directories   map[string]string `json:"directories"`
	DefaultDomain string            `json:"default_domain"`
}

func SetVault(v vaultpkg.Vault) {
	credentialVault = v
}

func getStoragePath() string {
	if contextDir := os.Getenv("TF_CRED_CONTEXT_DIR"); contextDir != "" {
		return filepath.Join(contextDir, "tfcred_contexts.json")
	}

	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		return "tfcred_contexts.json"
	}

	targetDir := filepath.Join(localAppData, "tfcred")
	_ = os.MkdirAll(targetDir, 0o755)

	return filepath.Join(targetDir, "tfcred_contexts.json")
}

// BindDirectory binds a working directory to a context.
func BindDirectory(
	dir,
	contextKey string,
) error {
	f := load()

	if f.Directories == nil {
		f.Directories = map[string]string{}
	}

	f.Directories[filepath.Clean(dir)] = contextKey

	return write(f)
}

// ResolveContextByDir resolves the context associated with a directory.
func ResolveContextByDir(
	startDir string,
) (string, bool) {
	f := load()

	if f.Directories == nil {
		return "", false
	}

	contextKey, exists := f.Directories[filepath.Clean(startDir)]
	return contextKey, exists
}

// Init initializes tfcred storage.
func Init(defaultDomain string) {
	fmt.Println("TF_CRED_CONTEXT_DIR:", os.Getenv("TF_CRED_CONTEXT_DIR"))
	fmt.Println("storage path:", getStoragePath())

	storageFile := getStoragePath()

	initial := File{
		DefaultDomain: defaultDomain,
		Contexts:      map[string]Entry{},
		Directories:   map[string]string{},
	}

	if _, err := os.Stat(storageFile); err == nil {
		current := load()

		if current.DefaultDomain == "" {
			current.DefaultDomain = defaultDomain
		}

		if current.Contexts == nil {
			current.Contexts = map[string]Entry{}
		}

		if current.Directories == nil {
			current.Directories = map[string]string{}
		}

		_ = write(current)

		fmt.Println("[tfcred] already initialized globally at:", storageFile)
		return
	}

	_ = write(initial)
}

// Add creates or updates a context and securely stores its token.
func Add(
	ctx,
	org,
	tokenType,
	domain,
	token string,
) {
	f := load()

	if f.Contexts == nil {
		f.Contexts = map[string]Entry{}
	}

	entry := Entry{
		Org:       org,
		TokenType: tokenType,
		Domain:    domain,
	}

	f.Contexts[ctx] = entry

	if token != "" {
		if err := saveToken(entry, token); err != nil {
			fmt.Printf(
				"[tfcred][error] Failed to store token securely: %v\n",
				err,
			)
			os.Exit(1)
		}

		fmt.Println("[tfcred] Token securely vaulted in Windows Credential Manager.")
	}

	_ = write(f)
}

// Remove deletes a context and its associated credential.
func Remove(
	ctx string,
) (Entry, bool) {
	f := load()

	entry, exists := f.Contexts[ctx]
	if !exists {
		return Entry{}, false
	}

	delete(f.Contexts, ctx)

	for dir, boundKey := range f.Directories {
		if boundKey == ctx {
			delete(f.Directories, dir)
		}
	}

	_ = deleteToken(entry)

	_ = write(f)

	return entry, true
}

// SetDefaultDomain updates the default Terraform domain.
func SetDefaultDomain(
	domain string,
) {
	f := load()

	f.DefaultDomain = domain

	if f.Contexts == nil {
		f.Contexts = map[string]Entry{}
	}

	_ = write(f)
}

// PurgeDomain removes every context associated with a domain.
func PurgeDomain(
	domain string,
) []string {
	f := load()

	if f.Contexts == nil {
		f.Contexts = map[string]Entry{}
	}

	var removed []string

	for name, entry := range f.Contexts {
		if entry.Domain != domain {
			continue
		}

		_ = deleteToken(entry)

		delete(f.Contexts, name)
		removed = append(removed, name)
	}

	for dir, boundKey := range f.Directories {
		if _, exists := f.Contexts[boundKey]; !exists {
			delete(f.Directories, dir)
		}
	}

	_ = write(f)

	return removed
}

// PurgeAll removes every configured context.
func PurgeAll() []string {
	f := load()

	if f.Contexts == nil {
		f.Contexts = map[string]Entry{}
	}

	var removed []string

	for name, entry := range f.Contexts {
		_ = deleteToken(entry)

		delete(f.Contexts, name)
		removed = append(removed, name)
	}

	f.Directories = map[string]string{}

	_ = write(f)

	return removed
}

// Load returns the current configuration.
func Load() File {
	return load()
}

// List prints all configured contexts.
func List() File {
	f := load()

	if len(f.Contexts) == 0 {
		fmt.Println("[tfcred] no contexts configured")
		return f
	}

	fmt.Println("[tfcred] configured contexts:")

	w := tabwriter.NewWriter(
		os.Stdout,
		0,
		0,
		3,
		' ',
		0,
	)

	_, _ = fmt.Fprintln(
		w,
		" CONTEXT\tTYPE\tORGANIZATION\tDOMAIN",
	)

	_, _ = fmt.Fprintln(
		w,
		" -------\t----\t------------\t------",
	)

	for name, entry := range f.Contexts {
		_, _ = fmt.Fprintf(
			w,
			" %s\t%s\t%s\t%s\n",
			name,
			entry.TokenType,
			entry.Org,
			entry.Domain,
		)
	}

	_ = w.Flush()

	return f
}

// GetToken retrieves the credential associated with an entry.
func GetToken(
	entry Entry,
) (string, error) {
	return credentialVault.Get(
		TokenVaultKey(
			entry.Domain,
			entry.TokenType,
			entry.Org,
		),
	)
}

func saveToken(
	entry Entry,
	token string,
) error {
	return credentialVault.Set(
		TokenVaultKey(
			entry.Domain,
			entry.TokenType,
			entry.Org,
		),
		token,
	)
}

func deleteToken(
	entry Entry,
) error {
	return credentialVault.Delete(
		TokenVaultKey(
			entry.Domain,
			entry.TokenType,
			entry.Org,
		),
	)
}

// TokenVaultBase converts a domain name into the base Windows Credential
// Manager namespace used by tfcred.
func TokenVaultBase(
	domain string,
) string {
	return fmt.Sprintf(
		"tfcred:domain:%s",
		sanitizeDomain(domain),
	)
}

// TokenVaultKey builds the unique credential manager key.
func TokenVaultKey(
	domain,
	tokenType,
	org string,
) string {
	base := TokenVaultBase(domain)

	if tokenType == "" {
		return base
	}

	return fmt.Sprintf(
		"%s:%s:%s",
		base,
		sanitizeTokenComponent(tokenType),
		sanitizeTokenComponent(org),
	)
}

func sanitizeDomain(
	domain string,
) string {
	domain = strings.TrimSpace(
		strings.ToLower(domain),
	)

	domain = strings.ReplaceAll(
		domain,
		".",
		"_",
	)

	domain = strings.ReplaceAll(
		domain,
		"-",
		"_",
	)

	return domain
}

func sanitizeTokenComponent(
	input string,
) string {
	component := strings.TrimSpace(
		strings.ToLower(input),
	)

	component = strings.ReplaceAll(
		component,
		"-",
		"_",
	)

	return component
}

// CleanOrphanedDirectories scans bindings and removes invalid entries.
func CleanOrphanedDirectories() ([]string, []string) {
	f := load()

	if len(f.Directories) == 0 {
		return nil, nil
	}

	var (
		deadPaths     []string
		deadContexts  []string
		pathsToDelete []string
	)

	if f.Contexts == nil {
		f.Contexts = map[string]Entry{}
	}

	for dir, contextKey := range f.Directories {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			deadPaths = append(deadPaths, dir)
			pathsToDelete = append(pathsToDelete, dir)
			continue
		}

		if _, exists := f.Contexts[contextKey]; !exists {
			deadContexts = append(
				deadContexts,
				fmt.Sprintf("%s -> [%s]", dir, contextKey),
			)

			pathsToDelete = append(pathsToDelete, dir)
		}
	}

	for _, dir := range pathsToDelete {
		delete(f.Directories, dir)
	}

	if len(pathsToDelete) > 0 {
		_ = write(f)
	}

	return deadPaths, deadContexts
}

// CheckOrphanedDirectories performs a read-only scan of bindings.
func CheckOrphanedDirectories() ([]string, []string) {
	f := load()

	if len(f.Directories) == 0 {
		return nil, nil
	}

	var (
		missingPaths    []string
		missingContexts []string
	)

	if f.Contexts == nil {
		f.Contexts = map[string]Entry{}
	}

	for dir, contextKey := range f.Directories {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			missingPaths = append(
				missingPaths,
				dir,
			)

			continue
		}

		if _, exists := f.Contexts[contextKey]; !exists {
			missingContexts = append(
				missingContexts,
				fmt.Sprintf("%s -> [%s]", dir, contextKey),
			)
		}
	}

	return missingPaths, missingContexts
}

func load() File {
	storageFile := getStoragePath()

	b, err := os.ReadFile(storageFile)
	if err != nil {
		return File{
			Contexts:    map[string]Entry{},
			Directories: map[string]string{},
		}
	}

	var f File

	if err := json.Unmarshal(b, &f); err != nil {
		return File{
			Contexts:    map[string]Entry{},
			Directories: map[string]string{},
		}
	}

	if f.Contexts == nil {
		f.Contexts = map[string]Entry{}
	}

	if f.Directories == nil {
		f.Directories = map[string]string{}
	}

	return f
}

func write(
	f File,
) error {
	storageFile := getStoragePath()

	b, err := json.MarshalIndent(
		f,
		"",
		" ",
	)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(
		filepath.Dir(storageFile),
		0o755,
	); err != nil {
		return err
	}

	return os.WriteFile(
		storageFile,
		b,
		0o644,
	)
}
