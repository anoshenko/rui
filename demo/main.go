package main

import (
	"embed"
	"fmt"

	"github.com/anoshenko/rui"
)

//go:embed resources
var resources embed.FS

const rootViewText = `
GridLayout {
	id = rootLayout, width = 100%, height = 100%, cell-height = "auto, 1fr",
	content = [
		GridLayout {
			id = rootTitle, width = 100%, cell-width = "auto, 1fr", 
			cell-vertical-align = center, background-color = #ffc0ded9,
			content = [
				ImageView { 
					id = rootTitleButton, padding = 8px, src = menu_icon.svg,
				},
				TextView { 
					id = rootTitleText, column = 1, padding-left = 8px, text = "Title",
				}
			],
		},
		StackLayout {
			id = rootViews, row = 1,
		}
	]
}
`

type demoPage struct {
	title   string
	creator func(session rui.Session) rui.View
	view    rui.View
}

type demoSession struct {
	rootView rui.View
	pages    []demoPage
}

func (demo *demoSession) OnStart(session rui.Session) {
	rui.DebugLog("Session start")
}

func (demo *demoSession) OnFinish(session rui.Session) {
	rui.DebugLog("Session finish")
}

func (demo *demoSession) OnResume(session rui.Session) {
	rui.DebugLog("Session resume")
}

func (demo *demoSession) OnPause(session rui.Session) {
	rui.DebugLog("Session pause")
}

func (demo *demoSession) OnDisconnect(session rui.Session) {
	rui.DebugLog("Session disconnect")
}

func (demo *demoSession) OnReconnect(session rui.Session) {
	rui.DebugLog("Session reconnect")
}

func createDemo(session rui.Session) rui.SessionContent {
	sessionContent := new(demoSession)
	sessionContent.pages = []demoPage{
		{"Text style", createTextStyleDemo, nil},
		{"View border", viewDemo, nil},
		{"Background image", createBackgroundDemo, nil},
		{"ListLayout", createListLayoutDemo, nil},
		{"GridLayout", createGridLayoutDemo, nil},
		{"ColumnLayout", createColumnLayoutDemo, nil},
		{"StackLayout", createStackLayoutDemo, nil},
		{"AbsoluteLayout", createAbsoluteLayoutDemo, nil},
		{"Resizable", createResizableDemo, nil},
		{"ListView", createListViewDemo, nil},
		{"Checkbox", createCheckboxDemo, nil},
		{"Controls", createControlsDemo, nil},
		{"TableView", createTableViewDemo, nil},
		{"EditView", createEditDemo, nil},
		{"ImageView", createImageViewDemo, nil},
		{"Canvas", createCanvasDemo, nil},
		{"VideoPlayer", createVideoPlayerDemo, nil},
		{"AudioPlayer", createAudioPlayerDemo, nil},
		{"Popups", createPopupDemo, nil},
		{"Filter", createFilterDemo, nil},
		{"Clip", createClipDemo, nil},
		{"Transform", transformDemo, nil},
		{"Transition", createTransitionDemo, nil},
		{"Key events", createKeyEventsDemo, nil},
		{"Mouse events", createMouseEventsDemo, nil},
		{"Pointer events", createPointerEventsDemo, nil},
		{"Touch events", createTouchEventsDemo, nil},
		//{"Tabs", createTabsDemo, nil},
	}

	return sessionContent
}

func (demo *demoSession) CreateRootView(session rui.Session) rui.View {
	demo.rootView = rui.CreateViewFromText(session, rootViewText)
	if demo.rootView == nil {
		return nil
	}

	rui.Set(demo.rootView, "rootTitleButton", rui.ClickEvent, demo.clickMenuButton)
	demo.showPage(0)

	return demo.rootView
}

func (demo *demoSession) clickMenuButton() {
	items := make([]string, len(demo.pages))
	for i, page := range demo.pages {
		items[i] = page.title
	}

	rui.ShowMenu(demo.rootView.Session(), rui.Params{
		rui.Items:           items,
		rui.OutsideClose:    true,
		rui.VerticalAlign:   rui.TopAlign,
		rui.HorizontalAlign: rui.LeftAlign,
		rui.PopupMenuResult: func(n int) {
			demo.showPage(n)
		},
	})
}

func (demo *demoSession) showPage(index int) {
	if index < 0 || index >= len(demo.pages) {
		return
	}

	if stackLayout := rui.StackLayoutByID(demo.rootView, "rootViews"); stackLayout != nil {
		if demo.pages[index].view == nil {
			demo.pages[index].view = demo.pages[index].creator(demo.rootView.Session())
			stackLayout.Append(demo.pages[index].view)
		} else {
			stackLayout.MoveToFront(demo.pages[index].view)
		}
		rui.Set(demo.rootView, "rootTitleText", rui.Text, demo.pages[index].title)
	}
	// TODO
}

func main() {
	rui.ProtocolInDebugLog = true
	rui.AddEmbedResources(&resources)
	app := rui.NewApplication("RUI demo", "icon.svg", createDemo)

	//addr := rui.GetLocalIP() + ":8080"
	addr := "localhost:8000"
	fmt.Print(addr)
	rui.OpenBrowser("http://" + addr)
	app.Start(addr)
}
