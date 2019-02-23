package env

// Initialize prepares environment
func Initialize() error {
	if err := initHome(); err != nil {
		return err
	}
	if err := initConfig(); err != nil {
		return err
	}

	return nil
}
