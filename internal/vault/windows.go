package vault

import (
	"github.com/zalando/go-keyring"
)

type WindowsVault struct {
	service string
}

func NewWindowsVault() *WindowsVault {
	return &WindowsVault{
		service: "tfcred",
	}
}

func (v *WindowsVault) Set(
	key string,
	value string,
) error {
	return keyring.Set(
		v.service,
		key,
		value,
	)
}

func (v *WindowsVault) Get(
	key string,
) (string, error) {
	return keyring.Get(
		v.service,
		key,
	)
}

func (v *WindowsVault) Delete(
	key string,
) error {
	return keyring.Delete(
		v.service,
		key,
	)
}
