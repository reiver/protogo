package gui

import (
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/richtext"

	"github.com/reiver/go-giotoast"
)

type App struct {
	theme  *material.Theme
	window *app.Window

	page           Page
	chatFrom       Page
	resumeFrom     Page
	personFrom     Page
	selectedPerson int
	selectedGroup  int
	selectedResume int

	me     Person
	people []Person
	groups []Group
	gigs   []Gig

	personClicks   []widget.Clickable
	groupClicks    []widget.Clickable
	personFavClicks []widget.Clickable
	groupFavClicks  []widget.Clickable
	memberClicks   []widget.Clickable
	resumeClicks   []widget.Clickable
	backClick      widget.Clickable
	chatClick      widget.Clickable
	groupChatClick widget.Clickable
	favClick       widget.Clickable
	sendClick      widget.Clickable
	navProfileClick  widget.Clickable
	navGigsClick     widget.Clickable
	navContactsClick widget.Clickable
	navChatsClick    widget.Clickable
	gigClicks        []widget.Clickable
	chatItemClicks   []widget.Clickable
	onboardingSaveClick    widget.Clickable
	onboardingSkipClick    widget.Clickable
	onboardingError        string
	onboardingLoading      bool
	profileFediIDError     string
	chatEditor             widget.Editor
	searchEditor           widget.Editor
	fediIDEditor           widget.Editor

	homeList       layout.List
	gigsList       layout.List
	chatsList      layout.List
	profileList    layout.List
	personList     layout.List
	groupList      layout.List
	resumeList     layout.List
	chatList       layout.List

	fetchChan chan fetchResult
	toasts    giotoast.Queue
	bioState  richtext.InteractiveText
	bioSpans  []richtext.SpanStyle
}

func newApp(w *app.Window) *App {
	var people []Person = loadPeopleFromDB()
	var groups []Group = loadGroupsFromDB()
	var me Person = loadMeFromDB()

	var fediIDEditor widget.Editor
	fediIDEditor.SingleLine = true
	fediIDEditor.Submit = true
	fediIDEditor.SetText(me.FediID)

	var startPage Page = PageHome
	if "" == me.FediID {
		startPage = PageOnboarding
	}

	var bioSpans []richtext.SpanStyle
	if "" != me.SummaryHTML {
		bioSpans = htmlToSpans(me.SummaryHTML)
	}

	return &App{
		theme:  material.NewTheme(),
		window: w,

		page: startPage,

		me:     me,
		people: people,
		groups: groups,
		gigs:   loadGigsFromDB(),

		personClicks:    make([]widget.Clickable, len(people)),
		groupClicks:     make([]widget.Clickable, len(groups)),
		personFavClicks: make([]widget.Clickable, len(people)),
		groupFavClicks:  make([]widget.Clickable, len(groups)),

		homeList: layout.List{
			Axis: layout.Vertical,
		},
		gigsList: layout.List{
			Axis: layout.Vertical,
		},
		profileList: layout.List{
			Axis: layout.Vertical,
		},
		personList: layout.List{
			Axis: layout.Vertical,
		},
		groupList: layout.List{
			Axis: layout.Vertical,
		},
		chatsList: layout.List{
			Axis: layout.Vertical,
		},
		resumeList: layout.List{
			Axis: layout.Vertical,
		},
		chatEditor: widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
		fediIDEditor: fediIDEditor,
		chatList: layout.List{
			Axis:      layout.Vertical,
			Alignment: layout.End,
		},

		fetchChan: make(chan fetchResult, 1),
		bioSpans:  bioSpans,
	}
}

func (receiver *App) Layout(gtx layout.Context) layout.Dimensions {
	// Drain fetch results from goroutines.
	receiver.drainFetchResults(gtx)

	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return receiver.layoutPage(gtx)
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return receiver.toasts.Layout(gtx, receiver.theme)
		}),
	)
}

func (receiver *App) layoutPage(gtx layout.Context) layout.Dimensions {
	switch receiver.page {
	case PageOnboarding:
		return receiver.layoutOnboarding(gtx)
	case PageProfile:
		return receiver.layoutProfilePage(gtx)
	case PageGigs:
		return receiver.layoutGigsPage(gtx)
	case PageChats:
		return receiver.layoutChatsPage(gtx)
	case PagePersonDetail:
		return receiver.layoutPersonDetail(gtx)
	case PageGroupDetail:
		return receiver.layoutGroupDetail(gtx)
	case PageResumeDetail:
		return receiver.layoutResumeDetail(gtx)
	case PageChat:
		return receiver.layoutChat(gtx)
	case PageGroupChat:
		return receiver.layoutGroupChat(gtx)
	default:
		return receiver.layoutHome(gtx)
	}
}

func (receiver *App) drainFetchResults(gtx layout.Context) {
	select {
	case result := <-receiver.fetchChan:
		receiver.onboardingLoading = false

		if nil != result.err {
			receiver.toasts.ShowType(giotoast.TypeError, result.err.Error(), 5*time.Second, gtx.Now)
			if receiver.page == PageOnboarding {
				receiver.onboardingError = result.err.Error()
			} else {
				receiver.profileFediIDError = result.err.Error()
			}
			return
		}

		if "" != result.name {
			receiver.me.Name = result.name
		}
		receiver.me.SummaryHTML = result.summaryHTML
		receiver.me.IconURL = result.iconURL
		receiver.me.BannerURL = result.bannerURL
		receiver.me.ProfileURL = result.profileURL

		if "" != result.summaryHTML {
			receiver.bioSpans = htmlToSpans(result.summaryHTML)
		} else {
			receiver.bioSpans = nil
		}

		err := persistProfileFromActor(result.name, result.summaryHTML, result.iconURL, result.bannerURL, result.profileURL)
		if nil != err {
			receiver.toasts.ShowType(giotoast.TypeError, "Profile fetched but failed to save", 5*time.Second, gtx.Now)
		} else {
			receiver.toasts.ShowType(giotoast.TypeSuccess, "Profile updated", 3*time.Second, gtx.Now)
		}

		if receiver.page == PageOnboarding {
			receiver.page = PageHome
		}
	default:
	}
}

func (receiver *App) startFetch(fediID string) {
	if receiver.onboardingLoading {
		return
	}

	receiver.onboardingLoading = true
	receiver.onboardingError = ""
	receiver.profileFediIDError = ""

	go func() {
		result := fetchActor(fediID)
		receiver.fetchChan <- result
		receiver.window.Invalidate()
	}()
}
