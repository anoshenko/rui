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

type application struct {
	server            *http.Server
	params            AppParams
	createContentFunc func(Session) SessionContent
	sessions          map[int]Session
}

func (app *application) getStartPage() string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString("<!DOCTYPE html>\n<html>\n")
	getStartPage(buffer, app.params, socketScripts)
	buffer.WriteString("\n</html>")
	return buffer.String()
}

func (app *application) Finish() {
	for _, session := range app.sessions {
		session.close()
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

			if !serveResourceFile(filename, w, req) &&
				!serveDownloadFile(filename, w, req) {
				w.WriteHeader(http.StatusNotFound)
			}
		}
	}
}

func (app *application) socketReader(brige webBrige) {
	var session Session
	events := make(chan DataObject, 1024)

	for {
		message, ok := brige.readMessage()
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
					if !brige.writeMessage(answer) {
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
							if !brige.writeMessage(answer.String()) {
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
					if !brige.writeMessage(answer) {
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

func sessionEventHandler(session Session, events chan DataObject, brige webBrige) {
	for {
		data := <-events

		switch command := data.Tag(); command {
		case "disconnect":
			session.onDisconnect()
			return

		case "session-close":
			session.onFinish()
			session.App().removeSession(session.ID())
			brige.close()

		default:
			session.handleEvent(command, data)
		}
	}
}

func (app *application) startSession(params DataObject, events chan DataObject, brige webBrige) (Session, string) {
	if app.createContentFunc == nil {
		return nil, ""
	}

	session := newSession(app, app.nextSessionID(), "", params)
	session.setBrige(events, brige)
	if !session.setContent(app.createContentFunc(session)) {
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

var apps = []*application{}

// StartApp - create the new application and start it
func StartApp(addr string, createContentFunc func(Session) SessionContent, params AppParams) {
	app := new(application)
	app.params = params
	app.sessions = map[int]Session{}
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
				if exec.Command(provider, url).Start(); err == nil {
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
