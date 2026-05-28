package gui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"protogo/cfg"
)

func (receiver *App) layoutGigsPage(gtx layout.Context) layout.Dimensions {
	var gigs []Gig = receiver.gigs

	// Ensure we have enough clickables.
	for len(receiver.gigClicks) < len(gigs) {
		receiver.gigClicks = append(receiver.gigClicks, widget.Clickable{})
	}

	// Handle gig clicks — navigate to the poster's profile.
	for i, gig := range gigs {
		if receiver.gigClicks[i].Clicked(gtx) {
			var personIndex int = receiver.findPersonIndex(gig.PostedBy)
			if 0 <= personIndex {
				receiver.selectedPerson = personIndex
				receiver.personFrom = PageGigs
				receiver.page = PagePersonDetail
			}
		}
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutTopBar(gtx, cfg.Name, false)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return receiver.gigsList.Layout(gtx, len(gigs), func(gtx layout.Context, index int) layout.Dimensions {
				return receiver.layoutGigItem(gtx, gigs[index], index)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutBottomNav(gtx)
		}),
	)
}

func (receiver *App) layoutGigItem(gtx layout.Context, gig Gig, clickIndex int) layout.Dimensions {
	return material.Clickable(gtx, &receiver.gigClicks[clickIndex], func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layoutCard(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Subtitle1(receiver.theme, gig.Title).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						var line string
						if "" != gig.Company {
							line = gig.Company
						}
						if "" != gig.Location {
							if "" != line {
								line = line + " · " + gig.Location
							} else {
								line = gig.Location
							}
						}
						if "" == line {
							return layout.Dimensions{}
						}
						lbl := material.Body2(receiver.theme, line)
						lbl.Color = color.NRGBA{R: 0x44, G: 0x44, B: 0x44, A: 0xFF}
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									lbl := material.Caption(receiver.theme, gig.Type)
									lbl.Color = color.NRGBA{R: 0x3F, G: 0x51, B: 0xB5, A: 0xFF}
									return lbl.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									lbl := material.Caption(receiver.theme, " · "+gig.Timestamp)
									lbl.Color = color.NRGBA{R: 0x99, G: 0x99, B: 0x99, A: 0xFF}
									return lbl.Layout(gtx)
								}),
							)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(6)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							lbl := material.Caption(receiver.theme, gig.Description)
							lbl.Color = color.NRGBA{R: 0x55, G: 0x55, B: 0x55, A: 0xFF}
							lbl.MaxLines = 2
							return lbl.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if "" == gig.PostedBy {
							return layout.Dimensions{}
						}
						return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return layoutAvatar(gtx, receiver.theme, gig.PostedBy, unit.Dp(20))
								}),
								layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									lbl := material.Caption(receiver.theme, gig.PostedBy)
									lbl.Color = color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xFF}
									return lbl.Layout(gtx)
								}),
							)
						})
					}),
				)
			})
		})
	})
}
