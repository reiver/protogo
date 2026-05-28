package gui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

var navProfileIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.SocialPerson)
	return icon
}()

var navContactsIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.SocialPeople)
	return icon
}()

var navChatsIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.CommunicationChat)
	return icon
}()

func (receiver *App) layoutBottomNav(gtx layout.Context) layout.Dimensions {
	if receiver.navProfileClick.Clicked(gtx) {
		receiver.page = PageProfile
	}
	if receiver.navContactsClick.Clicked(gtx) {
		receiver.page = PageHome
	}
	if receiver.navChatsClick.Clicked(gtx) {
		receiver.page = PageChats
	}

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			size := gtx.Constraints.Min
			defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, color.NRGBA{R: 0xFA, G: 0xFA, B: 0xFA, A: 0xFF})
			return layout.Dimensions{Size: size}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// Top border line.
			var lineHeight int = gtx.Dp(unit.Dp(1))
			defer clip.Rect{Max: image.Point{X: gtx.Constraints.Max.X, Y: lineHeight}}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, color.NRGBA{R: 0xDD, G: 0xDD, B: 0xDD, A: 0xFF})
			return layout.Dimensions{Size: image.Point{X: gtx.Constraints.Max.X, Y: lineHeight}}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(1)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layoutNavTab(gtx, receiver.theme, &receiver.navProfileClick, navProfileIcon, "Profile", receiver.page == PageProfile)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layoutNavTab(gtx, receiver.theme, &receiver.navContactsClick, navContactsIcon, "Contacts", receiver.page == PageHome)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layoutNavTab(gtx, receiver.theme, &receiver.navChatsClick, navChatsIcon, "Chats", receiver.page == PageChats)
					}),
				)
			})
		}),
	)
}

func layoutNavTab(gtx layout.Context, th *material.Theme, click *widget.Clickable, icon *widget.Icon, label string, active bool) layout.Dimensions {
	var tabColor color.NRGBA
	if active {
		tabColor = color.NRGBA{R: 0x3F, G: 0x51, B: 0xB5, A: 0xFF} // indigo
	} else {
		tabColor = color.NRGBA{R: 0x99, G: 0x99, B: 0x99, A: 0xFF} // gray
	}

	return material.Clickable(gtx, click, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						var sizePx int = gtx.Dp(unit.Dp(24))
						gtx.Constraints.Min = image.Point{X: sizePx, Y: sizePx}
						gtx.Constraints.Max = image.Point{X: sizePx, Y: sizePx}
						return icon.Layout(gtx, tabColor)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Caption(th, label)
						lbl.Color = tabColor
						return lbl.Layout(gtx)
					}),
				)
			})
		})
	})
}
