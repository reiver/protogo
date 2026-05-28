package gui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type App struct {
	theme *material.Theme

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
	onboardingSaveClick widget.Clickable
	onboardingError     string
	chatEditor          widget.Editor
	searchEditor        widget.Editor
	fediIDEditor        widget.Editor

	homeList       layout.List
	gigsList       layout.List
	chatsList      layout.List
	profileList    layout.List
	personList     layout.List
	groupList      layout.List
	resumeList     layout.List
	chatList       layout.List
}

func newApp() *App {
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

	return &App{
		theme: material.NewTheme(),

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
	}
}

func (receiver *App) Layout(gtx layout.Context) layout.Dimensions {
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
