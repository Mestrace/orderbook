package resources

func Init() {
	if err := InitDB(); err != nil {
		panic(err)
	}

	InitRateLimit()
}
