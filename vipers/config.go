package vipers

import "os"

var (
	pathConfig          = getENVWithDefault("VIPER_PATH_CONFIG", "./config")
	fileNameConfig      = getENVWithDefault("VIPER_FILENAME_CONFIG", "config")
	fileExtensionConfig = getENVWithDefault("VIPER_FILE_EXTENSION_CONFIG", "yaml")
)

func SetFileNameConfig(config string) {
	fileNameConfig = config
}

func SetFileExtensionConfig(config string) {
	fileExtensionConfig = config
}

func SetPathConfig(config string) {
	pathConfig = config
}

func getENVWithDefault(key, def string) string {
	env := os.Getenv(key)
	if env == "" {
		return def
	}

	return env
}
