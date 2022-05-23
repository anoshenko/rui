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
	// ImageConstant returns the image constant with "tag" name or "" if it is not exists
	ImageConstant(tag string) (string, bool)
	// SetCustomTheme set the custom theme
	SetCustomTheme(name string) bool
	// UserAgent returns the "user-agent" text of the client browser
	UserAgent() string
	// RemoteAddr returns the client address.
	RemoteAddr() string
	// Language returns the current session language
	Language() string
	// SetLanguage set the current session language
	SetLanguage(lang string)
	// GetString returns the text for the current language
	GetString(tag string) (string, bool)

	// Content returns the SessionContent of session
	Content() SessionContent
	setContent(content SessionContent, self Session) bool

	// SetTitle sets the text of the browser title/tab
	SetTitle(title string)
	// SetTitleColor sets the color of the browser navigation bar. Supported only in Safari and Chrome for android
	SetTitleColor(color Color)

	// RootView returns the root view of the session
	RootView() View
	// Get returns a value of the view (with id defined by the first argument) property with name defined by the second argument.
	// The type of return value depends on the property. If the property is not set then nil is returned.
	Get(viewID, tag string) interface{}
	// Set sets the value (third argument) of the property (second argument) of the view with id defined by the first argument.
	// Return "true" if the value has been set, in the opposite case "false" are returned and
	// a description of the error is written to the log
	Set(viewID, tag string, value interface{}) bool

	// DownloadFile downloads (saves) on the client side the file located at the specified path on the server.
	DownloadFile(path string)
	//DownloadFileData downloads (saves) on the client side a file with a specified name and specified content.
	DownloadFileData(filename string, data []byte)

	registerAnimation(props []AnimatedProperty) string

	resolveConstants(value string) (string, bool)
	checkboxOffImage() string
	checkboxOnImage() string
	radiobuttonOffImage() string
	radiobuttonOnImage() string

	viewByHTMLID(id string) View
	nextViewID() string
	styleProperty(styleTag, property string) interface{}

	setBrige(events chan DataObject, brige WebBrige)
	writeInitScript(writer *strings.Builder)
	runScript(script string)
	runGetterScript(script string) DataObject //, answer chan DataObject)
	handleAnswer(data DataObject)
	handleRootSize(data DataObject)
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
	customTheme      Theme
	currentTheme     Theme
	darkTheme        bool
	touchScreen      bool
	screenWidth      int
	screenHeight     int
	textDirection    int
	pixelRatio       float64
	userAgent        string
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
	animationCSS     string
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
	session.animationCSS = ""

	if customTheme != "" {
		if theme, ok := CreateThemeFromText(customTheme); ok {
			session.customTheme = theme
			session.currentTheme = nil
		}
	}

	if value, ok := params.PropertyValue("touch"); ok {
		session.touchScreen = (value == "1" || value == "true")
	}

	if value, ok := params.PropertyValue("user-agent"); ok {
		session.userAgent = value
	}

	if value, ok := params.PropertyValue("direction"); ok {
		if value == "rtl" {
			session.textDirection = RightToLeftDirection
		}
	}

	if value, ok := params.PropertyValue("language"); ok {
		session.language = value
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

func (session *sessionData) styleProperty(styleTag, propertyTag string) interface{} {
	if style := session.getCurrentTheme().style(styleTag); style != nil {
		return style.getRaw(propertyTag)
	}
	//errorLogF(`property "%v" not found`, propertyTag)
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
			session.rootView.setParentID("ruiRootView")
			return true
		}
	}
	return false
}

func (session *sessionData) RootView() View {
	return session.rootView
}

func (session *sessionData) writeInitScript(writer *strings.Builder) {
	if css := session.getCurrentTheme().cssText(session); css != "" {
		css = strings.ReplaceAll(css, "\n", `\n`)
		css = strings.ReplaceAll(css, "\t", `\t`)
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

	css := appStyles + session.getCurrentTheme().cssText(session) + session.animationCSS
	css = strings.ReplaceAll(css, "\n", `\n`)
	css = strings.ReplaceAll(css, "\t", `\t`)
	buffer.WriteString(`document.querySelector('style').textContent = "`)
	buffer.WriteString(css)
	buffer.WriteString("\";\n")

	if session.rootView != nil {
		buffer.WriteString(`document.getElementById('ruiRootView').innerHTML = '`)
		viewHTML(session.rootView, buffer)
		buffer.WriteString("';\nscanElementsSize();")
	}

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

func (session *sessionData) handleRootSize(data DataObject) {
	getValue := func(tag string) int {
		if value, ok := data.PropertyValue(tag); ok {
			float, err := strconv.ParseFloat(value, 64)
			if err == nil {
				return int(float)
			}
			ErrorLog(`Resize event error: ` + err.Error())
		} else {
			ErrorLogF(`Resize event error: the property "%s" not found`, tag)
		}
		return 0
	}

	if w := getValue("width"); w > 0 {
		session.screenWidth = w
	}
	if h := getValue("height"); h > 0 {
		session.screenHeight = h
	}
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
						if view := session.viewByHTMLID(viewID[:n]); view != nil {
							view.onItemResize(view, viewID[n+1:], getFloat("x"), getFloat("y"), getFloat("width"), getFloat("height"))
						} else {
							ErrorLogF(`View with id == %s not found`, viewID[:n])
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

func (session *sessionData) SetTitle(title string) {
	title, _ = session.GetString(title)
	session.runScript(`document.title = "` + title + `";`)
}

func (session *sessionData) SetTitleColor(color Color) {
	session.runScript(`setTitleColor("` + color.cssString() + `");`)
}

func (session *sessionData) RemoteAddr() string {
	return session.brige.remoteAddr()
}
