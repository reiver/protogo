package gui

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"protogo/cfg"
)

func (receiver *App) layoutProfilePage(gtx layout.Context) layout.Dimensions {
	var me Person = receiver.me

	var widgets []layout.Widget

	// Profile header with banner + avatar.
	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layoutProfileHeader(gtx, receiver.theme, me.Name, me.Title, me.Company, me.FediID)
	})

	// Detail sections.
	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layoutDetailSection(gtx, receiver.theme, "Title", me.Title)
		})
	})

	if "" != me.Company {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layoutDetailSection(gtx, receiver.theme, "Company", me.Company)
			})
		})
	}

	if "" != me.FediID {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layoutDetailSection(gtx, receiver.theme, "Fediverse ID", me.FediID)
			})
		})
	}

	if "" != me.Note {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layoutDetailSection(gtx, receiver.theme, "Notes", me.Note)
			})
		})
	}

	// Resumes section.
	if 0 < len(me.Resumes) {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return receiver.layoutResumesSection(gtx, me.Resumes, PageProfile)
			})
		})
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutTopBar(gtx, cfg.Name, false)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return receiver.profileList.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
				return widgets[index](gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutBottomNav(gtx)
		}),
	)
}
