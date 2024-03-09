package config

func InitConfig() error {
	err := InitApplicationConfig()
	if err != nil {
		return err
	}
	err = InitStorageConfig()
	if err != nil {
		return err
	}
	return nil
}

func MustInit() {
	err := InitConfig()
	if err != nil {
		panic(err)
	}
}
