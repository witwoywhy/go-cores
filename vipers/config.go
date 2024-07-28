package vipers

var (
	fileNameConfig      = "config"
	fileExtensionConfig = "yaml"
	pathConfig          = "./configs"
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
