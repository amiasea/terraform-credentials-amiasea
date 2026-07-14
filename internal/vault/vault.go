package vault

type Vault interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}
