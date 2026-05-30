package gui

import (
	"image/color"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"codeberg.org/reiver/go-fediverseid"
)

func (receiver *App) layoutOnboarding(gtx layout.Context) layout.Dimensions {
	var th *material.Theme = receiver.theme

	if !receiver.onboardingLoading {
		// Handle editor submit (Enter key).
		for {
			event, ok := receiver.fediIDEditor.Update(gtx)
			if !ok {
				break
			}
			if _, ok := event.(widget.SubmitEvent); ok {
				receiver.completeOnboarding()
			}
		}

		// Handle save button click.
		if receiver.onboardingSaveClick.Clicked(gtx) {
			receiver.completeOnboarding()
		}

		// Handle skip button click (only shown after error).
		if receiver.onboardingSkipClick.Clicked(gtx) {
			receiver.page = PageHome
		}
	}

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.X = gtx.Dp(unit.Dp(400))
		gtx.Constraints.Min.X = gtx.Constraints.Max.X

		return layout.Inset{Left: unit.Dp(32), Right: unit.Dp(32)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				// App name.
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.H4(th, "ProToGo")
					lbl.Color = color.NRGBA{R: 0x3F, G: 0x51, B: 0xB5, A: 0xFF}
					lbl.Alignment = 1 // center
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
				// Tagline.
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body2(th, "Your career on the Fediverse")
					lbl.Color = color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xFF}
					lbl.Alignment = 1 // center
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(48)}.Layout),
				// Prompt.
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					var msg string
					if receiver.onboardingLoading {
						msg = "Fetching your profile..."
					} else {
						msg = "Enter your Fediverse ID to get started"
					}
					lbl := material.Body1(th, msg)
					lbl.Alignment = 1 // center
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
				// Editor field.
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layoutCard(gtx, func(gtx layout.Context) layout.Dimensions {
						editor := material.Editor(th, &receiver.fediIDEditor, "@you@example.com")
						editor.Color = color.NRGBA{R: 0x3F, G: 0x51, B: 0xB5, A: 0xFF}
						editor.HintColor = color.NRGBA{R: 0x99, G: 0x99, B: 0x99, A: 0xFF}
						return editor.Layout(gtx)
					})
				}),
				// Error message.
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if "" == receiver.onboardingError {
						return layout.Dimensions{}
					}
					return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						lbl := material.Caption(th, receiver.onboardingError)
						lbl.Color = color.NRGBA{R: 0xD3, G: 0x2F, B: 0x2F, A: 0xFF}
						return lbl.Layout(gtx)
					})
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(24)}.Layout),
				// Buttons.
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if receiver.onboardingLoading {
						return layout.Dimensions{}
					}

					if "" != receiver.onboardingError {
						// After an error, show Retry + Continue anyway.
						return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(th, &receiver.onboardingSaveClick, "Retry")
								btn.Background = color.NRGBA{R: 0x3F, G: 0x51, B: 0xB5, A: 0xFF}
								return btn.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(th, &receiver.onboardingSkipClick, "Continue anyway")
								btn.Background = color.NRGBA{R: 0x75, G: 0x75, B: 0x75, A: 0xFF}
								return btn.Layout(gtx)
							}),
						)
					}

					btn := material.Button(th, &receiver.onboardingSaveClick, "Continue")
					btn.Background = color.NRGBA{R: 0x3F, G: 0x51, B: 0xB5, A: 0xFF}
					return btn.Layout(gtx)
				}),
			)
		})
	})
}

func (receiver *App) completeOnboarding() {
	var text string = strings.TrimSpace(receiver.fediIDEditor.Text())
	if "" == text {
		receiver.onboardingError = "Please enter your Fediverse ID."
		return
	}

	_, err := fediverseid.ParseFediverseIDString(text)
	if nil != err {
		receiver.onboardingError = "Not a valid Fediverse ID. Expected format: @name@host"
		return
	}

	receiver.onboardingError = ""
	receiver.me.FediID = text
	receiver.fediIDEditor.SetText(text)
	persistProfileFediID(text)
	receiver.startFetch(text)
}
