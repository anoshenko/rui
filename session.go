package rui

import (
	"fmt"
	"strconv"
	"strings"
)

// SessionContent is the interface of a session content
type SessionContent interface {
	CreateRootView(session Session) View
}

// Session provide interface to session parameters assess
type Session interface {
	// App return the current application interface
	App() Application
	// ID return the id of the session
	ID() int

	// DarkTheme returns "true" if the dark theme is used
	DarkTheme() bool
	// Mobile returns "true" if current session is displayed on a touch screen device
	TouchScreen() bool
	// PixelRatio returns the ratio of the resolution in physical pixels to the resolution
	// in logical pixels for the current display device.
	PixelRatio() float64
	// TextDirection returns the default text direction (LeftToRightDirection (1) or RightToLeftDirection (2))
	TextDirection() int
	// Constant returns the constant with "tag" name or "" if it is not exists
	Constant(tag string) (string, bool)
	// Color returns the color with "tag" name or 0 if it is not exists
	Color(tag string) (Color, bool)
	// SetCustomTheme set the custom theme
	SetCustomTheme(name string) bool
	// Language returns the current session language
	Language() string
	// SetLanguage set the current session language
	SetLanguage(lang string)
	// GetString returns the text for the current language
	GetString(tag string) (string, bool)

	Content() SessionContent
	setContent(content SessionContent, self Session) bool

	// RootView returns the root view of the session
	RootView() View
	// Get returns a value of the view (with id defined by the first argument) property with name defined by the second argument.
	// The type of return value depends on the property. If the property is not set then nil is returned.
	Get(viewID, tag string) interface{}
	// Set sets the value (third argument) of the property (second argument) of the view with id defined by the first argument.
	// Return "true" if the value has been set, in the opposite case "false" are returned and
	// a description of the error is written to the log
	Set(viewID, tag string, value interface{}) bool

	DownloadFile(path string)
	DownloadFileData(filename string, data []byte)

	registerAnimation(props []AnimatedProperty) string

	resolveConstants(value string) (string, bool)
	checkboxOffImage() string
	checkboxOnImage() string
	radiobuttonOffImage() string
	radiobuttonOnImage() string

	viewByHTMLID(id string) View
	nextViewID() string
	styleProperty(styleTag, property string) (string, bool)
	stylePropertyNode(styleTag, propertyTag string) DataNode

	setBrige(events chan DataObject, brige WebBrige)
	writeInitScript(writer *strings.Builder)
	runScript(script string)
	runGetterScript(script string) DataObject //, answer chan DataObject)
	handleAnswer(data DataObject)
	handleResize(data DataObject)
	handleViewEvent(command string, data DataObject)
	close()

	onStart()
	onFinish()
	onPause()
	onResume()
	onDisconnect()
	onReconnect()

	ignoreViewUpdates() bool
	setIgnoreViewUpdates(ignore bool)

	popupManager() *popupManager
	imageManager() *imageManager
}

type sessionData struct {
	customTheme      *theme
	darkTheme        bool
	touchScreen      bool
	textDirection    int
	pixelRatio       float64
	language         string
	languages        []string
	checkboxOff      string
	checkboxOn       string
	radiobuttonOff   string
	radiobuttonOn    string
	app              Application
	sessionID        int
	viewCounter      int
	content          SessionContent
	rootView         View
	ignoreUpdates    bool
	popups           *popupManager
	images           *imageManager
	brige            WebBrige
	events           chan DataObject
	animationCounter int
}

func newSession(app Application, id int, customTheme string, params DataObject) Session {
	session := new(sessionData)
	session.app = app
	session.sessionID = id
	session.darkTheme = false
	session.touchScreen = false
	session.pixelRatio = 1
	session.textDirection = LeftToRightDirection
	session.languages = []string{}
	session.viewCounter = 0
	session.ignoreUpdates = false
	session.animationCounter = 0

	if customTheme != "" {
		if theme, ok := newTheme(customTheme); ok {
			session.customTheme = theme
		}
	}

	if value, ok := params.PropertyValue("touch"); ok {
		session.touchScreen = (value == "1" || value == "true")
	}

	if value, ok := params.PropertyValue("direction"); ok {
		if value == "rtl" {
			session.textDirection = RightToLeftDirection
		}
	}

	if value, ok := params.PropertyValue("languages"); ok {
		session.languages = strings.Split(value, ",")
	}

	if value, ok := params.PropertyValue("dark"); ok {
		session.darkTheme = (value == "1" || value == "true")
	}

	if value, ok := params.PropertyValue("pixel-ratio"); ok {
		if f, err := strconv.ParseFloat(value, 64); err != nil {
			ErrorLog(err.Error())
		} else {
			session.pixelRatio = f
		}
	}

	return session
}

func (session *sessionData) App() Application {
	return session.app
}

func (session *sessionData) ID() int {
	return session.sessionID
}

func (session *sessionData) setBrige(events chan DataObject, brige WebBrige) {
	session.events = events
	session.brige = brige
}

func (session *sessionData) close() {
	if session.events != nil {
		session.events <- ParseDataText(`session-close{session="` + strconv.Itoa(session.sessionID) + `"}`)
	}
}

func (session *sessionData) styleProperty(styleTag, propertyTag string) (string, bool) {
	var style DataObject
	ok := false
	if session.customTheme != nil {
		style, ok = session.customTheme.styles[styleTag]
	}
	if !ok {
		style, ok = defaultTheme.styles[styleTag]
	}

	if ok {
		if node := style.PropertyWithTag(propertyTag); node != nil && node.Type() == TextNode {
			return session.resolveConstants(node.Text())
		}
	}

	//errorLogF(`property "%v" not found`, propertyTag)
	return "", false
}

func (session *sessionData) stylePropertyNode(styleTag, propertyTag string) DataNode {
	var style DataObject
	ok := false
	if session.customTheme != nil {
		style, ok = session.customTheme.styles[styleTag]
	}
	if !ok {
		style, ok = defaultTheme.styles[styleTag]
	}

	if ok {
		return style.PropertyWithTag(propertyTag)
	}

	return nil
}

func (session *sessionData) nextViewID() string {
	session.viewCounter++
	return fmt.Sprintf("id%06d", session.viewCounter)
}

func (session *sessionData) viewByHTMLID(id string) View {
	if session.rootView == nil {
		return nil
	}
	popupManager := session.popupManager()
	for _, popup := range popupManager.popups {
		if view := popup.viewByHTMLID(id); view != nil {
			return view
		}
	}
	return viewByHTMLID(id, session.rootView)
}

func (session *sessionData) Content() SessionContent {
	return session.content
}

func (session *sessionData) setContent(content SessionContent, self Session) bool {
	if content != nil {
		session.content = content
		session.rootView = content.CreateRootView(self)
		if session.rootView != nil {
			return true
		}
	}
	return false
}

func (session *sessionData) RootView() View {
	return session.rootView
}

func (session *sessionData) writeInitScript(writer *strings.Builder) {
	var workTheme *theme
	if session.customTheme == nil {
		workTheme = defaultTheme
	} else {
		workTheme = new(theme)
		workTheme.init()
		workTheme.concat(defaultTheme)
		workTheme.concat(session.customTheme)
	}

	if css := workTheme.cssText(session); css != "" {
		writer.WriteString(`document.querySelector('style').textContent += "`)
		writer.WriteString(css)
		writer.WriteString("\";\n")
	}

	if session.rootView != nil {
		writer.WriteString(`document.getElementById('ruiRootView').innerHTML = '`)
		viewHTML(session.rootView, writer)
		writer.WriteString("';\nscanElementsSize();")
	}
}

func (session *sessionData) reload() {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	session.writeInitScript(buffer)
	session.runScript(buffer.String())
}

func (session *sessionData) ignoreViewUpdates() bool {
	return session.brige == nil || session.ignoreUpdates
}

func (session *sessionData) setIgnoreViewUpdates(ignore bool) {
	session.ignoreUpdates = ignore
}

func (session *sessionData) Get(viewID, tag string) interface{} {
	if view := ViewByID(session.RootView(), viewID); view != nil {
		return view.Get(tag)
	}
	return nil
}

func (session *sessionData) Set(viewID, tag string, value interface{}) bool {
	if view := ViewByID(session.RootView(), viewID); view != nil {
		return view.Set(tag, value)
	}
	return false
}

func (session *sessionData) popupManager() *popupManager {
	if session.popups == nil {
		session.popups = new(popupManager)
		session.popups.popups = []Popup{}
	}
	return session.popups
}

func (session *sessionData) imageManager() *imageManager {
	if session.images == nil {
		session.images = new(imageManager)
		session.images.images = make(map[string]*imageData)
	}
	return session.images
}

func (session *sessionData) runScript(script string) {
	if session.brige != nil {
		session.brige.WriteMessage(script)
	} else {
		ErrorLog("No connection")
	}
}

func (session *sessionData) runGetterScript(script string) DataObject { //}, answer chan DataObject) {
	if session.brige != nil {
		return session.brige.RunGetterScript(script)
	}

	ErrorLog("No connection")
	result := NewDataObject("error")
	result.SetPropertyValue("text", "No connection")
	return result
}

func (session *sessionData) handleAnswer(data DataObject) {
	session.brige.AnswerReceived(data)
}

func (session *sessionData) handleResize(data DataObject) {
	if node := data.PropertyWithTag("views"); node != nil && node.Type() == ArrayNode {
		for _, el := range node.ArrayElements() {
			if el.IsObject() {
				obj := el.Object()
				getFloat := func(tag string) float64 {
					if value, ok := obj.PropertyValue(tag); ok {
						f, err := strconv.ParseFloat(value, 64)
						if err == nil {
							return f
						}
						ErrorLog(`Resize event error: ` + err.Error())
					} else {
						ErrorLogF(`Resize event error: the property "%s" not found`, tag)
					}
					return 0
				}
				if viewID, ok := obj.PropertyValue("id"); ok {
					if n := strings.IndexRune(viewID, '-'); n > 0 {
						if index, err := strconv.Atoi(viewID[n+1:]); err == nil {
							if view := session.viewByHTMLID(viewID[:n]); view != nil {
								view.onItemResize(view, index, getFloat("x"), getFloat("y"), getFloat("width"), getFloat("height"))
							} else {
								ErrorLogF(`View with id == %s not found`, viewID[:n])
							}
						} else {
							ErrorLogF(`Invalid view id == %s not found`, viewID)
						}
					} else if view := session.viewByHTMLID(viewID); view != nil {
						view.onResize(view, getFloat("x"), getFloat("y"), getFloat("width"), getFloat("height"))
						view.setScroll(getFloat("scroll-x"), getFloat("scroll-y"), getFloat("scroll-width"), getFloat("scroll-height"))
					} else {
						ErrorLogF(`View with id == %s not found`, viewID)
					}
				} else {
					ErrorLog(`"id" property not found`)
				}
			} else {
				ErrorLog(`Resize event error: views element is not object`)
			}
		}
	} else {
		ErrorLog(`Resize event error: invalid "views" property`)
	}
}

func (session *sessionData) handleViewEvent(command string, data DataObject) {
	if viewID, ok := data.PropertyValue("id"); ok {
		if view := session.viewByHTMLID(viewID); view != nil {
			view.handleCommand(view, command, data)
		}
	} else {
		ErrorLog(`"id" property not found. Event: ` + command)
	}
}
