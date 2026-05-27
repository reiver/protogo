package main

import (
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"

	"protogo/srv/log"
)

func main() {
	shout()

	log := logsrv.Begin()
	defer log.End()

	log.Highlightf("ProToGo ⚡")
	defer log.Highlightf("ProToGo 👻")

	log.Informf("Here we go…")
	window()
}
