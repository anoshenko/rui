package rui

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
)

//go:embed app_scripts.js
var defaultScripts string

//go:embed app_styles.css
var appStyles string

//go:embed defaultTheme.rui
var defaultThemeText string

// Application - app interface
type Application interface {
	// Start - start the application life cycle
	Start(addr string)
	Finish()
	nextSessionID() int
	removeSession(id int)
}

type application struct {
	name, icon        string
	createContentFunc func(Session) SessionContent
	sessions          map[int]Session
}

func (app *application) getStartPage() string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString(`<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>`)
	buffer.WriteString(app.name)
	buffer.WriteString("</title>")
	if app.icon != "" {
		buffer.WriteString(`
		<link rel="icon" href="`)
		buffer.WriteString(app.icon)
		buffer.WriteString(`">`)
	}

	buffer.WriteString(`
		<base target="_blank" rel="noopener">
		<meta name="viewport" content="width=device-width">
		<style>`)
	buffer.WriteString(appStyles)
	buffer.WriteString(`</style>
		<script>`)
	buffer.WriteString(defaultScripts)
	buffer.WriteString(`</script>
	</head>
	<body>
		<div class="ruiRoot" id="ruiRootView"></div>
		<div class="ruiPopupLayer" id="ruiPopupLayer" style="visibility: hidden;" onclick="clickOutsidePopup(event)"></div>
	</body>
</html>`)

	return buffer.String()
}

func (app *application) init(name, icon string) {
	app.name = name
	app.icon = icon
	app.sessions = map[int]Session{}
}

func (app *application) Start(addr string) {
	http.Handle("/", app)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (app *application) Finish() {
	for _, session := range app.sessions {
		session.close()
	}

}

func (app *application) nextSessionID() int {
	n := rand.Intn(0x7FFFFFFE) + 1
	_, ok := app.sessions[n]
	for ok {
		n = rand.Intn(0x7FFFFFFE) + 1
		_, ok = app.sessions[n]
	}
	return n
}

func (app *application) removeSession(id int) {
	delete(app.sessions, id)
}

func (app *application) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if ProtocolInDebugLog {
		DebugLogF("%s %s", req.Method, req.URL.Path)
	}

	switch req.Method {
	case "GET":
		switch req.URL.Path {
		case "/":
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, app.getStartPage())

		case "/ws":
			if brige := CreateSocketBrige(w, req); brige != nil {
				go app.socketReader(brige)
			}

		default:
			filename := req.URL.Path[1:]
			if size := len(filename); size > 0 && filename[size-1] == '/' {
				filename = filename[:size-1]
			}

			if !serveResourceFile(filename, w, req) {
				w.WriteHeader(http.StatusNotFound)
			}
		}
	}
}

func (app *application) socketReader(brige WebBrige) {
	var session Session
	events := make(chan DataObject, 1024)

	for {
		message, ok := brige.ReadMessage()
		if !ok {
			events <- NewDataObject("disconnect")
			return
		}

		if ProtocolInDebugLog {
			DebugLog(message)
		}

		if obj := ParseDataText(message); obj != nil {
			command := obj.Tag()
			switch command {
			case "startSession":
				answer := ""
				if session, answer = app.startSession(obj, events, brige); session != nil {
					if !brige.WriteMessage(answer) {
						return
					}
					session.onStart()
					go sessionEventHandler(session, events, brige)
				}

			case "reconnect":
				if sessionText, ok := obj.PropertyValue("session"); ok {
					if sessionID, err := strconv.Atoi(sessionText); err == nil {
						if session = app.sessions[sessionID]; session != nil {
							session.setBrige(events, brige)
							answer := allocStringBuilder()
							defer freeStringBuilder(answer)

							session.writeInitScript(answer)
							if !brige.WriteMessage(answer.String()) {
								return
							}
							session.onReconnect()
							go sessionEventHandler(session, events, brige)
							return
						}
						DebugLogF("Session #%d not exists", sessionID)
					} else {
						ErrorLog(`strconv.Atoi(sessionText) error: ` + err.Error())
					}
				} else {
					ErrorLog(`"session" key not found`)
				}

				answer := ""
				if session, answer = app.startSession(obj, events, brige); session != nil {
					if !brige.WriteMessage(answer) {
						return
					}
					session.onStart()
					go sessionEventHandler(session, events, brige)
				}

			case "answer":
				session.handleAnswer(obj)

			case "imageLoaded":
				session.imageManager().imageLoaded(obj, session)

			case "imageError":
				session.imageManager().imageLoadError(obj, session)

			default:
				events <- obj
			}
		}
	}
}

func sessionEventHandler(session Session, events chan DataObject, brige WebBrige) {
	for {
		data := <-events

		switch command := data.Tag(); command {
		case "disconnect":
			session.onDisconnect()
			return

		case "session-close":
			session.onFinish()
			session.App().removeSession(session.ID())
			brige.Close()

		case "session-pause":
			session.onPause()

		case "session-resume":
			session.onResume()

		case "resize":
			session.handleResize(data)

		default:
			session.handleViewEvent(command, data)
		}
	}
}

func (app *application) startSession(params DataObject, events chan DataObject, brige WebBrige) (Session, string) {
	if app.createContentFunc == nil {
		return nil, ""
	}

	session := newSession(app, app.nextSessionID(), "", params)
	session.setBrige(events, brige)
	if !session.setContent(app.createContentFunc(session), session) {
		return nil, ""
	}

	app.sessions[session.ID()] = session

	answer := allocStringBuilder()
	defer freeStringBuilder(answer)

	answer.WriteString("sessionID = '")
	answer.WriteString(strconv.Itoa(session.ID()))
	answer.WriteString("';\n")
	session.writeInitScript(answer)
	answerText := answer.String()

	if ProtocolInDebugLog {
		DebugLog("Start session:")
		DebugLog(answerText)
	}
	return session, answerText
}

// NewApplication - create the new application of the single view type.
func NewApplication(name, icon string, createContentFunc func(Session) SessionContent) Application {
	app := new(application)
	app.init(name, icon)
	app.createContentFunc = createContentFunc
	return app
}

func OpenBrowser(url string) bool {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err != nil
}
