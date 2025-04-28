package support

func IsDev() bool {
	return Config("app.env") == "dev"
}

func IsProd() bool {
	return Config("app.env") == "prod"
}
