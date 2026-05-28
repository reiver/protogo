package dbsrv

import (
	"database/sql"
	"os"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"

	_ "modernc.org/sqlite"

	"protogo/cfg"
	"protogo/srv/log"
)

var (
	readDB  *sql.DB
	writeDB *sql.DB
)

func ReadDB() *sql.DB {
	return readDB
}

func WriteDB() *sql.DB {
	return writeDB
}

func init() {
	log := logsrv.Begin()
	defer log.End()

	{
		var dir string = cfg.DBDirPath()

	        err := os.MkdirAll(dir, 0700)
	        if nil != err {
			log.Error(
				field.S("failed to mkdir-all"),
				field.E(err),
				field.String("directory", dir),
			)
			panic(err)
	        }
	}

	var filename string = cfg.DBFilePath()

	{
		var err error

		writeDB, err = sql.Open("sqlite", filename)
		if nil != err {
			err = erorr.Wrap(err, "failed to open SQLite write database connection",
				field.String("file-name", filename),
			)
			panic(err)
		}
		if nil == writeDB {
			err = erorr.Wrap(err, "nil SQLite write database connection",
				field.String("file-name", filename),
			)
			panic(err)
		}

		writeDB.SetMaxOpenConns(1)

		_, err = writeDB.Exec(`PRAGMA journal_mode=WAL`)
		if nil != err {
			err = erorr.Wrap(err, "failed to enable WAL mode on SQLite database",
				field.String("file-name", filename),
			)
			panic(err)
		}

		err = writeDB.Ping()
		if nil != err {
			err = erorr.Wrap(err, "failed to ping SQLite write database connection",
				field.String("file-name", filename),
			)
			panic(err)
		}
	}

	{
		var err error

		readDB, err = sql.Open("sqlite", filename)
		if nil != err {
			err = erorr.Wrap(err, "failed to open SQLite read database connection",
				field.String("file-name", filename),
			)
			panic(err)
		}
		if nil == readDB {
			err = erorr.Wrap(err, "nil SQLite read database connection",
				field.String("file-name", filename),
			)
			panic(err)
		}

		err = readDB.Ping()
		if nil != err {
			err = erorr.Wrap(err, "failed to ping SQLite read database connection",
				field.String("file-name", filename),
			)
			panic(err)
		}
	}

	{
		err := CreateDatabase(logsrv.Logger(), writeDB)
		if nil != err {
			err = erorr.Wrap(err, "failed to create database table(s) in SQLite database",
				field.String("file-name", filename),
			)
			panic(err)
		}
	}
}

func Close() {
	log := logsrv.Begin()
	defer log.End()

	var filename string = cfg.DBFilePath()

	if nil != readDB {
		err := readDB.Close()
		if nil != err {
			log.Error(
				field.S("failed to close SQLite read database connection"),
				field.E(err),
				field.String("file-name", filename),
			)
		}
	}

	if nil != writeDB {
		err := writeDB.Close()
		if nil != err {
			log.Error(
				field.S("failed to close SQLite write database connection"),
				field.E(err),
				field.String("file-name", filename),
			)
		}
	}
}
