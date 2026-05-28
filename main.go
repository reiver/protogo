package main

import (
	"protogo/srv/db"
	"protogo/srv/log"
)

func main() {
	shout()

	log := logsrv.Begin()
	defer log.End()

	log.Highlightf("ProToGo ⚡")
	defer log.Highlightf("ProToGo 👻")

	defer dbsrv.Close()

	seeddb()

	showandtell()

	log.Informf("Here we go…")
	gui()
}
