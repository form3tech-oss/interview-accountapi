package env

func GetDbHost() string {
	return getEnvVar("DATABASE-HOST", "localhost").(string)
}

func GetDbPort() string {
	return getEnvVar("DATABASE-PORT", "localhost").(string)
}

func GetDbUser() string {
	return getEnvVar("DATABASE-USERNAME", "postgres").(string)
}

func GetDbPassword() string {
	return getEnvVar("DATABASE-PASSWORD", "postgres").(string)
}

func GetDbName() string {
	return getEnvVar("DATABASE-NAME", "interview_accountapi").(string)
}
