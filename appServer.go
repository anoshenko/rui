//go:build !wasm

package rui

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//go:embed app_socket.js
var socketScripts string

//go:embed app_post.js
var httpPostScripts string

func debugLog(text string) {
	log.Println("\033[34m" + text)
}

func errorLog(text string) {
	log.Println("\033[31m" + text)
}

type sessionInfo struct {
	session  Session
	response chan string
}

type application struct {
	server            *http.Server
	params            AppParams
	createContentFunc func(Session) SessionContent
	sessions          map[int]sessionInfo
}

func (app *application) getStartPage() string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString("<!DOCTYPE html>\n<html>\n")
	getStartPage(buffer, app.params)
	buffer.WriteString("\n</html>")
	return buffer.String()
}

func (app *application) Finish() {
	for _, session := range app.sessions {
		session.session.close()
		if session.response != nil {
			close(session.response)
			session.response = nil
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		log.Println(err.Error())
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
	if info, ok := app.sessions[id]; ok {
		if info.response != nil {
			close(info.response)
		}
		delete(app.sessions, id)
	}
}

func (app *application) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if ProtocolInDebugLog {
		DebugLogF("%s %s", req.Method, req.URL.Path)
	}

	switch req.Method {
	case "POST":
		if req.URL.Path == "/" {
			app.postHandler(w, req)
		}

	case "GET":
		switch req.URL.Path {
		case "/":
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, app.getStartPage())

		case "/ws":
			if bridge := createSocketBridge(w, req); bridge != nil {
				go app.socketReader(bridge)
			}

		case "/script.js":
			w.WriteHeader(http.StatusOK)
			if app.params.NoSocket {
				io.WriteString(w, httpPostScripts)
			} else {
				io.WriteString(w, socketScripts)
			}
			io.WriteString(w, "\n")
			io.WriteString(w, defaultScripts)

		default:
			filename := req.URL.Path[1:]
			if size := len(filename); size > 0 && filename[size-1] == '/' {
				filename = filename[:size-1]
			}

			if !serveResourceFile(filename, w, req) &&
				!serveDownloadFile(filename, w, req) {
				w.WriteHeader(http.StatusNotFound)
			}
		}
	}
}

func setSessionIDCookie(w http.ResponseWriter, sessionID int) {
	cookie := http.Cookie{
		Name:     "session",
		Value:    strconv.Itoa(sessionID),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func (app *application) postHandler(w http.ResponseWriter, req *http.Request) {

	if reqBody, err := io.ReadAll(req.Body); err == nil {
		message := string(reqBody)

		if ProtocolInDebugLog {
			DebugLog(message)
		}

		if obj := ParseDataText(message); obj != nil {
			var session Session = nil
			var response chan string = nil

			if cookie, err := req.Cookie("session"); err == nil {
				sessionID, err := strconv.Atoi(cookie.Value)
				if err != nil {
					ErrorLog(err.Error())
				} else if info, ok := app.sessions[sessionID]; ok && info.response != nil {
					response = info.response
					session = info.session
				}
			}

			command := obj.Tag()
			startSession := false

			if session == nil || command == "startSession" {
				events := make(chan DataObject, 1024)
				bridge := createHttpBridge(req)
				response = bridge.response
				answer := ""
				session, answer = app.startSession(obj, events, bridge, response)

				bridge.writeMessage(answer)
				session.onStart()
				if command == "session-resume" {
					session.onResume()
				}
				bridge.sendResponse()

				setSessionIDCookie(w, session.ID())
				startSession = true

				go sessionEventHandler(session, events, bridge)
			}

			if !startSession {
				switch command {
				case "nop":
					session.sendResponse()

				case "session-close":
					session.onFinish()
					session.App().removeSession(session.ID())
					return

				default:
					if !session.handleAnswer(command, obj) {
						session.addToEventsQueue(obj)
					}
				}
			}

			io.WriteString(w, <-response)
			for len(response) > 0 {
				io.WriteString(w, <-response)
			}
		}
	}
}

func (app *application) socketReader(bridge *wsBridge) {
	var session Session
	events := make(chan DataObject, 1024)

	for {
		message, ok := bridge.readMessage()
		if !ok {
			events <- NewDataObject("disconnect")
			return
		}

		if ProtocolInDebugLog {
			DebugLog("🖥️ -> " + message)
		}

		if obj := ParseDataText(message); obj != nil {
			command := obj.Tag()
			switch command {
			case "startSession":
				answer := ""
				if session, answer = app.startSession(obj, events, bridge, nil); session != nil {
					if !bridge.writeMessage(answer) {
						return
					}
					session.onStart()
					go sessionEventHandler(session, events, bridge)
				}

			case "reconnect":
				if sessionText, ok := obj.PropertyValue("session"); ok {
					if sessionID, err := strconv.Atoi(sessionText); err == nil {
						if info, ok := app.sessions[sessionID]; ok {
							session := info.session
							session.setBridge(events, bridge)
							answer := allocStringBuilder()
							defer freeStringBuilder(answer)

							session.writeInitScript(answer)
							if !bridge.writeMessage(answer.String()) {
								return
							}
							session.onReconnect()
							go sessionEventHandler(session, events, bridge)
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
				if session, answer = app.startSession(obj, events, bridge, nil); session != nil {
					if !bridge.writeMessage(answer) {
						return
					}
					session.onStart()
					go sessionEventHandler(session, events, bridge)
				}

			default:
				if !session.handleAnswer(command, obj) {
					events <- obj
				}
			}
		}
	}
}

func sessionEventHandler(session Session, events chan DataObject, bridge bridge) {
	for {
		data := <-events

		switch command := data.Tag(); command {
		case "disconnect":
			session.onDisconnect()
			return

		case "session-close":
			session.onFinish()
			session.App().removeSession(session.ID())
			bridge.close()

		default:
			session.handleEvent(command, data)
		}
	}
}

func (app *application) startSession(params DataObject, events chan DataObject,
	bridge bridge, response chan string) (Session, string) {

	if app.createContentFunc == nil {
		return nil, ""
	}

	session := newSession(app, app.nextSessionID(), "", params)
	session.setBridge(events, bridge)
	if !session.setContent(app.createContentFunc(session)) {
		return nil, ""
	}

	app.sessions[session.ID()] = sessionInfo{
		session:  session,
		response: response,
	}

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

var apps = []*application{}

// StartApp - create the new application and start it
func StartApp(addr string, createContentFunc func(Session) SessionContent, params AppParams) {
	app := new(application)
	app.params = params
	app.sessions = map[int]sessionInfo{}
	app.createContentFunc = createContentFunc
	apps = append(apps, app)

	redirectAddr := ""
	if index := strings.IndexRune(addr, ':'); index >= 0 {
		redirectAddr = addr[:index] + ":80"
	} else {
		redirectAddr = addr + ":80"
		if params.CertFile != "" && params.KeyFile != "" {
			addr += ":443"
		} else {
			addr += ":80"
		}
	}

	app.server = &http.Server{Addr: addr}
	http.Handle("/", app)

	serverRun := func(err error) {
		if err != nil {
			if err == http.ErrServerClosed {
				log.Println(err)
			} else {
				log.Fatal(err)
			}
		}
	}

	if params.CertFile != "" && params.KeyFile != "" {
		if params.Redirect80 {
			redirectTLS := func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "https://"+addr+r.RequestURI, http.StatusMovedPermanently)
			}

			go func() {
				serverRun(http.ListenAndServe(redirectAddr, http.HandlerFunc(redirectTLS)))
			}()
		}
		serverRun(app.server.ListenAndServeTLS(params.CertFile, params.KeyFile))
	} else {
		serverRun(app.server.ListenAndServe())
	}
}

func FinishApp() {
	for _, app := range apps {
		app.Finish()
	}
	apps = []*application{}
}

func OpenBrowser(url string) bool {
	var err error

	switch runtime.GOOS {
	case "linux":
		for _, provider := range []string{"xdg-open", "x-www-browser", "www-browser"} {
			if _, err = exec.LookPath(provider); err == nil {
				if err = exec.Command(provider, url).Start(); err == nil {
					return true
				}
			}
		}

	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()

	case "darwin":
		err = exec.Command("open", url).Start()

	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err != nil
}
