package gui

import (
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"

	"protogo/cfg"
	"protogo/srv/log"
)

func Run() {
	log := logsrv.Begin()
	defer log.End()

	application := newApp()

	go func() {
		w := new(app.Window)
		w.Option(
			app.Title(cfg.Name),
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
				application.Layout(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}()

	app.Main()
}
