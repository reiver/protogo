package cfg

import (
	"os"
	"path/filepath"
)

func DBDirPath() string {
	configDir, err := os.UserConfigDir()
	if nil != err {
		configDir = "."
	}

	return filepath.Join(configDir, "protogo")
}

const DBFileName string = "protogo.db"

func DBFilePath() string {
	return filepath.Join(DBDirPath(), "protogo.db")
}
