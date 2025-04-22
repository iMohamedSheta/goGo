package support

func IsDev() bool {
	return Config("app.environment") == "dev"
}

func IsProd() bool {
	return Config("app.environment") == "prod"
}
