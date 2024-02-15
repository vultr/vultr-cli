package utils

const (
	// APIKeyError is the error used when missing the API key on commands which
	// require it
	//nolint:gosec
	APIKeyError string = `
Please export your VULTR API key as an environment variable or add 'api-key' to your config file, eg:
export VULTR_API_KEY='<api_key_from_vultr_account>'
	`
)
