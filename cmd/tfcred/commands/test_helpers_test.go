package commands

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/amiasea/terraform-credentials-tfcred/internal/store"
)

type fakeVault struct {
	tokens  map[string]string
	deleted []string
}

func newFakeVault() *fakeVault {
	return &fakeVault{
		tokens: make(map[string]string),
	}
}

func (v *fakeVault) Get(key string) (string, error) {
	token, ok := v.tokens[key]
	if !ok {
		return "", os.ErrNotExist
	}

	return token, nil
}

func (v *fakeVault) Set(key, token string) error {
	v.tokens[key] = token
	return nil
}

func (v *fakeVault) Delete(key string) error {
	delete(v.tokens, key)
	v.deleted = append(v.deleted, key)

	return nil
}

func setupCommandTest(t *testing.T) *fakeVault {
	t.Helper()

	vault := newFakeVault()

	store.SetVault(vault)

	store.PurgeAll()
	store.Init("app.terraform.io")

	return vault
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed creating pipe: %v", err)
	}

	os.Stdout = writer

	fn()

	_ = writer.Close()
	os.Stdout = old

	var buffer bytes.Buffer

	_, _ = io.Copy(
		&buffer,
		reader,
	)

	_ = reader.Close()

	return buffer.String()
}
