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
	"slices"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/acme/autocert"
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
	usedSessionIds    []int
}

func (app *application) getStartPage() string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString("<!DOCTYPE html>\n<html>\n")
	getStartPage(buffer, app.nextSessionID(), app.params)
	buffer.WriteString("\n</html>")
	return buffer.String()
}

func (app *application) Params() AppParams {
	params := app.params
	if params.NoSocket {
		params.SocketAutoClose = 0
	}
	return params
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

func (app *application) getCreateContentFunc() func(Session) SessionContent {
	return app.createContentFunc
}

func (app *application) nextSessionID() int {
	if app.usedSessionIds == nil {
		app.usedSessionIds = []int{}
	}

	for {
		n := rand.Intn(0x7FFFFFFE) + 1
		if !slices.Contains(app.usedSessionIds, n) {
			app.usedSessionIds = append(app.usedSessionIds, n)
			slices.Sort(app.usedSessionIds)
			return n
		}
	}
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
	case http.MethodPost:
		if req.URL.Path == "/" {
			app.postHandler(w, req)
		}

	case http.MethodGet:
		switch req.URL.Path {
		case "/":
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, app.getStartPage())

		case "/e":
			app.sseHandler(w, req)

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

/*
func setSessionIDCookie(w http.ResponseWriter, sessionID int) {
	cookie := http.Cookie{
		Name:     "session",
		Value:    strconv.Itoa(sessionID),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func getSessionIDCookie(req *http.Request) (int, error) {
	cookie, err := req.Cookie("session")
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(cookie.Value)
}
*/

func (app *application) postHandler(w http.ResponseWriter, req *http.Request) {

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		ErrorLog(err.Error())
		return
	}

	message := string(reqBody)
	if ProtocolInDebugLog {
		DebugLog(message)
	}

	obj, err := ParseDataText(message)
	if err != nil {
		ErrorLog(err.Error())
		return
	}

	sessionID, ok := getSessionID(obj)
	if !ok {
		return
	}

	var session Session = nil
	var response chan string = nil
	if info, ok := app.sessions[sessionID]; ok && info.response != nil {
		response = info.response
		session = info.session
	}

	command := obj.Tag()

	if session == nil {
		if command != "start-session" {
			io.WriteString(w, "reloadPage();")
			return
		}

		events := make(chan DataObject, 1024)
		bridge := createHttpBridge(req)
		response = bridge.response

		session = app.createSession(obj, events, bridge, response)
		if session == nil {
			return
		}

		go sessionEventHandler(session, events, bridge)
	}

	switch command {

	case "session-close":
		session.onFinish()
		session.App().removeSession(session.ID())
		return

	case "nop":
		if len(response) == 0 {
			session.addToEventsQueue(obj)
		}

	default:
		if !session.handleAnswer(command, obj) {
			session.addToEventsQueue(obj)
		} else {
			io.WriteString(w, "sendNop();")
		}
	}

	io.WriteString(w, <-response)
	for len(response) > 0 {
		io.WriteString(w, <-response)
	}
}

func (app *application) sseHandler(w http.ResponseWriter, req *http.Request) {
	/*
		sessionID, err := strconv.Atoi(req.URL.Query().Get("id"))
		if err != nil {
			ErrorLog("SessionID error: " + err.Error())
			return
		}

		if sessionInfo, ok := app.sessions[sessionID]; ok {

		}
	*/
}

func getSessionID(obj DataObject) (int, bool) {
	sessionText, ok := obj.PropertyValue("session")
	if !ok {
		ErrorLog(`"session" key not found`)
		return 0, false
	}

	sessionID, err := strconv.Atoi(sessionText)
	if err != nil {
		ErrorLog(`"session" key text strconv.Atoi error: ` + err.Error())
		return 0, false
	}

	return sessionID, true
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
			DebugLog("🖥️ → " + message)
		}

		obj, err := ParseDataText(message)
		if err != nil {
			ErrorLog(err.Error())
			continue
		}

		switch command := obj.Tag(); command {
		case "start-session":
			if session = app.createSession(obj, events, bridge, nil); session != nil {
				go sessionEventHandler(session, events, bridge)
				events <- obj
			}

		case "reconnect":
			session = nil
			if sessionID, ok := getSessionID(obj); ok {
				if info, ok := app.sessions[sessionID]; ok {
					session = info.session
					session.setBridge(events, bridge)

					go sessionEventHandler(session, events, bridge)
					session.onReconnect()
				} else {
					DebugLogF("Session #%d not exists", sessionID)
				}
			}

			if session == nil {
				bridge.writeMessage("reloadPage();")
				return
			}

		default:
			if !session.handleAnswer(command, obj) {
				events <- obj
			}
		}
	}
}

func sessionEventHandler(session Session, events chan DataObject, bridge bridge) {
	for {
		data := <-events

		switch command := data.Tag(); command {
		case "disconnect":
			session.setBridge(nil, nil)
			session.onDisconnect()
			return

		case "session-close":
			session.onFinish()
			session.App().removeSession(session.ID())
			bridge.close()
			return

		case "nop":
			session.sendResponse()

		default:
			session.handleEvent(command, data)
		}
	}
}

func (app *application) createSession(params DataObject, events chan DataObject,
	bridge bridge, response chan string) Session {

	sessionID, ok := getSessionID(params)
	if !ok || app.createContentFunc == nil {
		return nil
	}

	session := newSession(app, sessionID, "", params)
	session.setBridge(events, bridge)

	app.sessions[sessionID] = sessionInfo{
		session:  session,
		response: response,
	}

	return session
}

var apps = []*application{}

// StartApp - create the new application and start it
func StartApp(addr string, createContentFunc func(Session) SessionContent, params AppParams) {
	resources.scanDefaultResourcePath()

	app := new(application)
	app.params = params
	app.sessions = map[int]sessionInfo{}
	app.createContentFunc = createContentFunc
	apps = append(apps, app)

	redirectAddr := ""
	https := params.AutoCertDomain != "" || (params.CertFile != "" && params.KeyFile != "")

	if index := strings.IndexRune(addr, ':'); index >= 0 {
		redirectAddr = addr[:index] + ":80"
	} else {
		redirectAddr = addr + ":80"
		if https {
			addr += ":443"
		} else {
			addr += ":80"
		}
	}

	serverRun := func(err error) {
		if err != nil {
			if err == http.ErrServerClosed {
				log.Println(err)
			} else {
				log.Fatal(err)
			}
		}
	}

	if https {
		if params.Redirect80 {
			redirectTLS := func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "https://"+addr+r.RequestURI, http.StatusMovedPermanently)
			}

			go func() {
				serverRun(http.ListenAndServe(redirectAddr, http.HandlerFunc(redirectTLS)))
			}()
		}

		if params.AutoCertDomain != "" {
			mux := http.NewServeMux()
			mux.Handle("/", app)
			serverRun(http.Serve(autocert.NewListener(params.AutoCertDomain), mux))
		} else {
			app.server = &http.Server{Addr: addr}
			http.Handle("/", app)
			serverRun(app.server.ListenAndServeTLS(params.CertFile, params.KeyFile))
		}
	} else {
		app.server = &http.Server{Addr: addr}
		http.Handle("/", app)
		serverRun(app.server.ListenAndServe())
	}
}

// FinishApp finishes application
func FinishApp() {
	for _, app := range apps {
		app.Finish()
	}
	apps = []*application{}
}

// OpenBrowser open browser with specific URL locally. Useful for applications which run on local machine
// or for debug purposes.
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
