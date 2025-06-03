package support

func IsDev() bool {
	isDev, err := Config("app.env")
	if err != nil {
		return false
	}
	return isDev == "dev"
}

func IsProd() bool {
	isProd, err := Config("app.env")
	if err != nil {
		return false
	}
	return isProd == "prod"
}
