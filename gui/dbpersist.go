package gui

import (
	"time"

	dbsrv "protogo/srv/db"
	"protogo/srv/log"
)

func persistPersonFavorite(person Person) {
	if 0 == person.DBID {
		return
	}
	go func() {
		logger := logsrv.Logger()
		db := dbsrv.WriteDB()
		dbsrv.UpdatePersonFavorite(logger, db, person.DBID, person.Favorite)
	}()
}

func persistGroupFavorite(group Group) {
	if 0 == group.DBID {
		return
	}
	go func() {
		logger := logsrv.Logger()
		db := dbsrv.WriteDB()
		dbsrv.UpdateGroupFavorite(logger, db, group.DBID, group.Favorite)
	}()
}

func persistPersonMessage(person Person, text string, timestamp time.Time) {
	if 0 == person.DBID {
		return
	}
	go func() {
		logger := logsrv.Logger()
		db := dbsrv.WriteDB()
		pid := person.DBID
		dbsrv.InsertMessage(logger, db, dbsrv.MessageRow{
			PersonID:  &pid,
			FromMe:    true,
			Text:      text,
			Timestamp: timestamp.Format("2006-01-02 15:04"),
		})
	}()
}

func persistGroupMessage(group Group, text string, timestamp time.Time) {
	if 0 == group.DBID {
		return
	}
	go func() {
		logger := logsrv.Logger()
		db := dbsrv.WriteDB()
		gid := group.DBID
		dbsrv.InsertMessage(logger, db, dbsrv.MessageRow{
			GroupID:   &gid,
			FromMe:    true,
			Text:      text,
			Timestamp: timestamp.Format("2006-01-02 15:04"),
		})
	}()
}
