//go:build wasm

package rui

import (
	_ "embed"
	"syscall/js"
)

//go:embed app_wasm.js
var wasmScripts string

type wasmApp struct {
	params            AppParams
	createContentFunc func(Session) SessionContent
	session           Session
	brige             webBrige
}

func (app *wasmApp) Finish() {
	app.session.close()
}

func (app *wasmApp) startSession(this js.Value, args []js.Value) interface{} {
	if app.createContentFunc == nil || len(args) == 1 {
		return nil
	}

	params := ParseDataText(args[0].String())
	session := newSession(app, 0, "", params)
	session.setBrige(make(chan DataObject), app.brige)
	if !session.setContent(app.createContentFunc(session), session) {
		return nil
	}

	app.session = session

	answer := allocStringBuilder()
	defer freeStringBuilder(answer)

	session.writeInitScript(answer)
	answerText := answer.String()

	if ProtocolInDebugLog {
		DebugLog("Start session:")
		DebugLog(answerText)
	}
	return nil
}

func (app *wasmApp) removeSession(id int) {
}

// StartApp - create the new wasmApp and start it
func StartApp(addr string, createContentFunc func(Session) SessionContent, params AppParams) {
	app := new(wasmApp)
	app.params = params
	app.createContentFunc = createContentFunc

	if createContentFunc == nil {
		return
	}

	app.brige = createWasmBrige()
	js.Global().Set("startSession", js.FuncOf(app.startSession))

	/*
		script := defaultScripts + wasmScripts
		script = strings.ReplaceAll(script, "\\", `\\`)
		script = strings.ReplaceAll(script, "\n", `\n`)
		script = strings.ReplaceAll(script, "\t", `\t`)
		script = strings.ReplaceAll(script, "\"", `\"`)
		script = strings.ReplaceAll(script, "'", `\'`)

		js.Global().Call("execScript", `document.getElementById('ruiscript').text += "`+script+`"`)
	*/

	document := js.Global().Get("document")
	body := document.Call("querySelector", "body")
	body.Set("innerHTML", `<div class="ruiRoot" id="ruiRootView"></div>
<div class="ruiPopupLayer" id="ruiPopupLayer" style="visibility: hidden;" onclick="clickOutsidePopup(event)"></div>
<a id="ruiDownloader" download style="display: none;"></a>`)

	//js.Global().Call("execScript", "initSession()")
	js.Global().Call("initSession", "")
	//window.Call("execScript", "initSession()")

	for true {
		if message, ok := app.brige.readMessage(); ok && app.session != nil {
			if ProtocolInDebugLog {
				DebugLog(message)
			}

			if obj := ParseDataText(message); obj != nil {
				switch command := obj.Tag(); command {
				/*
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

										case "disconnect":
						session.onDisconnect()
						return

					case "session-close":
						session.onFinish()
						session.App().removeSession(session.ID())
						brige.close()

				*/
				case "answer":
					app.session.handleAnswer(obj)

				case "imageLoaded":
					app.session.imageManager().imageLoaded(obj, app.session)

				case "imageError":
					app.session.imageManager().imageLoadError(obj, app.session)

				default:
					app.session.handleEvent(command, obj)
				}
			}
		}
	}
}

func FinishApp() {
	//app.Finish()
}

func OpenBrowser(url string) bool {
	return false
}

/*
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
*/
