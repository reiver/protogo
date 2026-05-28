package main

import (
	"codeberg.org/reiver/go-field"

	"protogo/cfg"
	"protogo/srv/log"
)

func showandtell() {
	log := logsrv.Begin()
	defer log.End()

	log.Highlight(
		field.S("CFG"),

		field.String("LOG_LEVEL", cfg.LogLevel()),

		field.String("DB_DIR_PATH", cfg.DBDirPath()),
		field.String("DB_FILE_NAME", cfg.DBFileName),
		field.String("DB_FILE_PATH", cfg.DBFilePath()),
	)
}
