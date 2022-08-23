package env

func GetServerHost() string {
	v := getEnvVar("SERVER_HOST", "http://localhost:8080")
	return v.(string)
}
