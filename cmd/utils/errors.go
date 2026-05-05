package utils

const (
	// APIKeyError is the error used when missing the API key on commands which
	// require it
	//nolint:gosec
	APIKeyError string = `set VULTR_API_KEY as an environment variable or add 'api-key' to your config file`
)
