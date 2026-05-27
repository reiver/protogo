package main

import (
	"protogo/srv/log"
)

func main() {
	shout()

	log := logsrv.Begin()
	defer log.End()

	log.Highlightf("ProToGo ⚡")
	defer log.Highlightf("ProToGo 👻")

	log.Informf("Here we go…")
	gui()
}
