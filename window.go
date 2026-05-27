package main

import (
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"

	"protogo/srv/log"
)

func window() {
	log := logsrv.Begin()
	defer log.End()

	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("ProToGo"),
			app.Size(unit.Dp(800), unit.Dp(600)),
		)

		var ops op.Ops
		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				if e.Err != nil {
					log.Errorf("window error: %s", e.Err)
				}
				os.Exit(0)
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				e.Frame(gtx.Ops)
			}
		}
	}()

	app.Main()
}
