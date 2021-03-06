package rui

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//go:embed app_scripts.js
var defaultScripts string

//go:embed app_styles.css
var appStyles string

//go:embed defaultTheme.rui
var defaultThemeText string

// Application - app interface
type Application interface {
	Finish()
	nextSessionID() int
	removeSession(id int)
}

type application struct {
	server            *http.Server
	params            AppParams
	createContentFunc func(Session) SessionContent
	sessions          map[int]Session
}

// AppParams defines parameters of the app
type AppParams struct {
	// Title - title of the app window/tab
	Title string
	// TitleColor - background color of the app window/tab (applied only for Safari and Chrome for Android)
	TitleColor Color
	// Icon - the icon file name
	Icon string
	// CertFile - path of a certificate for the server must be provided
	// if neither the Server's TLSConfig.Certificates nor TLSConfig.GetCertificate are populated.
	// If the certificate is signed by a certificate authority, the certFile should be the concatenation
	// of the server's certificate, any intermediates, and the CA's certificate.
	CertFile string
	// KeyFile - path of a private key for the server must be provided
	// if neither the Server's TLSConfig.Certificates nor TLSConfig.GetCertificate are populated.
	KeyFile string
	// Redirect80 - if true then the function of redirect from port 80 to 443 is created
	Redirect80 bool
}

func (app *application) getStartPage() string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString(`<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>`)
	buffer.WriteString(app.params.Title)
	buffer.WriteString("</title>")
	if app.params.Icon != "" {
		buffer.WriteString(`
		<link rel="icon" href="`)
		buffer.WriteString(app.params.Icon)
		buffer.WriteString(`">`)
	}

	if app.params.TitleColor != 0 {
		buffer.WriteString(`
		<meta name="theme-color" content="`)
		buffer.WriteString(app.params.TitleColor.cssString())
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
		<a id="ruiDownloader" download style="display: none;"></a>
	</body>
</html>`)

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

		case "root-size":
			session.handleRootSize(data)

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

type downloadFile struct {
	filename string
	path     string
	data     []byte
}

var currentDownloadId = int(rand.Int31())
var downloadFiles = map[string]downloadFile{}

func (session *sessionData) startDownload(file downloadFile) {
	currentDownloadId++
	id := strconv.Itoa(currentDownloadId)
	downloadFiles[id] = file
	session.runScript(fmt.Sprintf(`startDowndload("%s", "%s")`, id, file.filename))
}

func serveDownloadFile(id string, w http.ResponseWriter, r *http.Request) bool {
	if file, ok := downloadFiles[id]; ok {
		delete(downloadFiles, id)
		if file.data != nil {
			http.ServeContent(w, r, file.filename, time.Now(), bytes.NewReader(file.data))
			return true
		} else if _, err := os.Stat(file.path); err == nil {
			http.ServeFile(w, r, file.path)
			return true
		}
	}
	return false
}

// DownloadFile starts downloading the file on the client side.
func (session *sessionData) DownloadFile(path string) {
	if _, err := os.Stat(path); err != nil {
		ErrorLog(err.Error())
		return
	}

	_, filename := filepath.Split(path)
	session.startDownload(downloadFile{
		filename: filename,
		path:     path,
		data:     nil,
	})
}

// DownloadFileData starts downloading the file on the client side. Arguments specify the name of the downloaded file and its contents
func (session *sessionData) DownloadFileData(filename string, data []byte) {
	if data == nil {
		ErrorLog("Invalid download data. Must be not nil.")
		return
	}

	session.startDownload(downloadFile{
		filename: filename,
		path:     "",
		data:     data,
	})
}
