package rui

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type webBridge interface {
	startUpdateScript(htmlID string) bool
	finishUpdateScript(htmlID string)
	callFunc(funcName string, args ...any) bool
	updateInnerHTML(htmlID, html string)
	appendToInnerHTML(htmlID, html string)
	updateCSSProperty(htmlID, property, value string)
	updateProperty(htmlID, property string, value any)
	removeProperty(htmlID, property string)
	readMessage() (string, bool)
	writeMessage(text string) bool
	addAnimationCSS(css string)
	clearAnimation()
	cavnasStart(htmlID string)
	callCanvasFunc(funcName string, args ...any)
	callCanvasVarFunc(v any, funcName string, args ...any)
	callCanvasImageFunc(url string, property string, funcName string, args ...any)
	createCanvasVar(funcName string, args ...any) any
	updateCanvasProperty(property string, value any)
	cavnasFinish()
	canvasTextMetrics(htmlID, font, text string) TextMetrics
	htmlPropertyValue(htmlID, name string) string
	answerReceived(answer DataObject)
	close()
	remoteAddr() string
}

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
	setContent(content SessionContent) bool

	// SetTitle sets the text of the browser title/tab
	SetTitle(title string)
	// SetTitleColor sets the color of the browser navigation bar. Supported only in Safari and Chrome for android
	SetTitleColor(color Color)

	// RootView returns the root view of the session
	RootView() View
	// Get returns a value of the view (with id defined by the first argument) property with name defined by the second argument.
	// The type of return value depends on the property. If the property is not set then nil is returned.
	Get(viewID, tag string) any
	// Set sets the value (third argument) of the property (second argument) of the view with id defined by the first argument.
	// Return "true" if the value has been set, in the opposite case "false" are returned and
	// a description of the error is written to the log
	Set(viewID, tag string, value any) bool

	// DownloadFile downloads (saves) on the client side the file located at the specified path on the server.
	DownloadFile(path string)
	//DownloadFileData downloads (saves) on the client side a file with a specified name and specified content.
	DownloadFileData(filename string, data []byte)
	// OpenURL opens the url in the new browser tab
	OpenURL(url string)

	getCurrentTheme() Theme
	registerAnimation(props []AnimatedProperty) string

	resolveConstants(value string) (string, bool)
	checkboxOffImage() string
	checkboxOnImage() string
	radiobuttonOffImage() string
	radiobuttonOnImage() string

	viewByHTMLID(id string) View
	nextViewID() string
	styleProperty(styleTag, property string) any

	setBridge(events chan DataObject, bridge webBridge)
	writeInitScript(writer *strings.Builder)
	callFunc(funcName string, args ...any)
	updateInnerHTML(htmlID, html string)
	appendToInnerHTML(htmlID, html string)
	updateCSSProperty(htmlID, property, value string)
	updateProperty(htmlID, property string, value any)
	removeProperty(htmlID, property string)
	startUpdateScript(htmlID string) bool
	finishUpdateScript(htmlID string)
	addAnimationCSS(css string)
	clearAnimation()
	cavnasStart(htmlID string)
	callCanvasFunc(funcName string, args ...any)
	createCanvasVar(funcName string, args ...any) any
	callCanvasVarFunc(v any, funcName string, args ...any)
	callCanvasImageFunc(url string, property string, funcName string, args ...any)
	updateCanvasProperty(property string, value any)
	cavnasFinish()
	canvasTextMetrics(htmlID, font, text string) TextMetrics
	htmlPropertyValue(htmlID, name string) string
	handleAnswer(data DataObject)
	handleRootSize(data DataObject)
	handleResize(data DataObject)
	handleEvent(command string, data DataObject)
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
	bridge           webBridge
	events           chan DataObject
	animationCounter int
	animationCSS     string
	updateScripts    map[string]*strings.Builder
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
	session.updateScripts = map[string]*strings.Builder{}

	if customTheme != "" {
		if theme, ok := CreateThemeFromText(customTheme); ok {
			session.customTheme = theme
			session.currentTheme = nil
		}
	}

	if params != nil {
		session.handleSessionInfo(params)
	}

	return session
}

func (session *sessionData) App() Application {
	return session.app
}

func (session *sessionData) ID() int {
	return session.sessionID
}

func (session *sessionData) setBridge(events chan DataObject, bridge webBridge) {
	session.events = events
	session.bridge = bridge
}

func (session *sessionData) close() {
	if session.events != nil {
		session.events <- ParseDataText(`session-close{session="` + strconv.Itoa(session.sessionID) + `"}`)
	}
}

func (session *sessionData) styleProperty(styleTag, propertyTag string) any {
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

func (session *sessionData) setContent(content SessionContent) bool {
	if content != nil {
		session.content = content
		session.rootView = content.CreateRootView(session)
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

	session.bridge.writeMessage(buffer.String())
}

func (session *sessionData) ignoreViewUpdates() bool {
	return session.bridge == nil || session.ignoreUpdates
}

func (session *sessionData) setIgnoreViewUpdates(ignore bool) {
	session.ignoreUpdates = ignore
}

func (session *sessionData) Get(viewID, tag string) any {
	if view := ViewByID(session.RootView(), viewID); view != nil {
		return view.Get(tag)
	}
	return nil
}

func (session *sessionData) Set(viewID, tag string, value any) bool {
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

func (session *sessionData) callFunc(funcName string, args ...any) {
	if session.bridge != nil {
		session.bridge.callFunc(funcName, args...)
	} else {
		ErrorLog("No connection")
	}
}

func (session *sessionData) updateInnerHTML(htmlID, html string) {
	if !session.ignoreViewUpdates() {
		if session.bridge != nil {
			session.bridge.updateInnerHTML(htmlID, html)
		} else {
			ErrorLog("No connection")
		}
	}
}

func (session *sessionData) appendToInnerHTML(htmlID, html string) {
	if !session.ignoreViewUpdates() {
		if session.bridge != nil {
			session.bridge.appendToInnerHTML(htmlID, html)
		} else {
			ErrorLog("No connection")
		}
	}
}

func (session *sessionData) updateCSSProperty(htmlID, property, value string) {
	if !session.ignoreViewUpdates() && session.bridge != nil {
		session.bridge.updateCSSProperty(htmlID, property, value)
	}
}

func (session *sessionData) updateProperty(htmlID, property string, value any) {
	if !session.ignoreViewUpdates() && session.bridge != nil {
		session.bridge.updateProperty(htmlID, property, value)
	}
}

func (session *sessionData) removeProperty(htmlID, property string) {
	if !session.ignoreViewUpdates() && session.bridge != nil {
		session.bridge.removeProperty(htmlID, property)
	}
}

func (session *sessionData) startUpdateScript(htmlID string) bool {
	if session.bridge != nil {
		return session.bridge.startUpdateScript(htmlID)
	}
	return false
}

func (session *sessionData) finishUpdateScript(htmlID string) {
	if session.bridge != nil {
		session.bridge.finishUpdateScript(htmlID)
	}
}

func (session *sessionData) addAnimationCSS(css string) {
	if session.bridge != nil {
		session.bridge.addAnimationCSS(css)
	}
}

func (session *sessionData) clearAnimation() {
	if session.bridge != nil {
		session.bridge.clearAnimation()
	}
}

func (session *sessionData) cavnasStart(htmlID string) {
	if session.bridge != nil {
		session.bridge.cavnasStart(htmlID)
	}
}

func (session *sessionData) callCanvasFunc(funcName string, args ...any) {
	if session.bridge != nil {
		session.bridge.callCanvasFunc(funcName, args...)
	}
}

func (session *sessionData) updateCanvasProperty(property string, value any) {
	if session.bridge != nil {
		session.bridge.updateCanvasProperty(property, value)
	}
}

func (session *sessionData) createCanvasVar(funcName string, args ...any) any {
	if session.bridge != nil {
		return session.bridge.createCanvasVar(funcName, args...)
	}
	return nil
}

func (session *sessionData) callCanvasVarFunc(v any, funcName string, args ...any) {
	if session.bridge != nil && v != nil {
		session.bridge.callCanvasVarFunc(v, funcName, args...)
	}
}

func (session *sessionData) callCanvasImageFunc(url string, property string, funcName string, args ...any) {
	if session.bridge != nil {
		session.bridge.callCanvasImageFunc(url, property, funcName, args...)
	}
}

func (session *sessionData) cavnasFinish() {
	if session.bridge != nil {
		session.bridge.cavnasFinish()
	}
}

func (session *sessionData) canvasTextMetrics(htmlID, font, text string) TextMetrics {
	if session.bridge != nil {
		return session.bridge.canvasTextMetrics(htmlID, font, text)
	}

	ErrorLog("No connection")
	return TextMetrics{Width: 0}
}

func (session *sessionData) htmlPropertyValue(htmlID, name string) string {
	if session.bridge != nil {
		return session.bridge.htmlPropertyValue(htmlID, name)
	}

	ErrorLog("No connection")
	return ""
}

func (session *sessionData) handleAnswer(data DataObject) {
	session.bridge.answerReceived(data)
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
							DebugLogF(`View with id == %s not found`, viewID[:n])
						}
					} else if view := session.viewByHTMLID(viewID); view != nil {
						view.onResize(view, getFloat("x"), getFloat("y"), getFloat("width"), getFloat("height"))
						view.setScroll(getFloat("scroll-x"), getFloat("scroll-y"), getFloat("scroll-width"), getFloat("scroll-height"))
					} else {
						DebugLogF(`View with id == %s not found`, viewID)
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

func (session *sessionData) handleSessionInfo(params DataObject) {
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
}

func (session *sessionData) handleEvent(command string, data DataObject) {
	switch command {
	case "session-pause":
		session.onPause()

	case "session-resume":
		session.onResume()

	case "root-size":
		session.handleRootSize(data)

	case "resize":
		session.handleResize(data)

	case "sessionInfo":
		session.handleSessionInfo(data)

	default:
		if viewID, ok := data.PropertyValue("id"); ok {
			if view := session.viewByHTMLID(viewID); view != nil {
				view.handleCommand(view, command, data)
			}
		} else if command != "clickOutsidePopup" {
			ErrorLog(`"id" property not found. Event: ` + command)
		}
	}
}

func (session *sessionData) SetTitle(title string) {
	title, _ = session.GetString(title)
	session.callFunc("setTitle", title)
}

func (session *sessionData) SetTitleColor(color Color) {
	session.callFunc("setTitleColor", color.cssString())
}

func (session *sessionData) RemoteAddr() string {
	return session.bridge.remoteAddr()
}

func (session *sessionData) OpenURL(urlStr string) {
	if _, err := url.ParseRequestURI(urlStr); err != nil {
		ErrorLog(err.Error())
		return
	}
	session.callFunc("openURL", urlStr)
}
