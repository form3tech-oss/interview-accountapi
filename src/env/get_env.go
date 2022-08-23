package env

import "os"

// Convenience function to retrieve a value otherwise
// rely on a default value. This allows me to run the code
// in a docker container in a production like environment
// whils also allowing me to run locally for debugging purposes.
func getEnvVar(key string, def any) any {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return def
}
