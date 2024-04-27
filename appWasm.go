//go:build wasm

package rui

import (
	_ "embed"
	"encoding/base64"
	"path/filepath"
	"strings"
	"syscall/js"
)

//go:embed app_wasm.js
var wasmScripts string

type wasmApp struct {
	params            AppParams
	createContentFunc func(Session) SessionContent
	session           Session
	bridge            bridge
	close             chan DataObject
}

func (app *wasmApp) Finish() {
	app.session.close()
}

func (app *wasmApp) Params() AppParams {
	params := app.params
	params.SocketAutoClose = 0
	return params
}

func debugLog(text string) {
	js.Global().Get("console").Call("log", text)
}

func errorLog(text string) {
	js.Global().Get("console").Call("log", "%c"+text, "color: #F00;")
}

func (app *wasmApp) handleMessage(this js.Value, args []js.Value) any {
	if len(args) > 0 {
		text := args[0].String()
		if ProtocolInDebugLog {
			DebugLog(text)
		}
		if obj := ParseDataText(text); obj != nil {
			switch command := obj.Tag(); command {
			case "session-close":
				app.close <- obj

			default:
				if !app.session.handleAnswer(command, obj) {
					app.session.handleEvent(command, obj)
				}
			}
		}
	}
	return nil
}

func (app *wasmApp) removeSession(id int) {
}

func (app *wasmApp) createSession() Session {
	session := newSession(app, 0, "", ParseDataText(js.Global().Call("sessionInfo", "").String()))
	session.setBridge(app.close, app.bridge)
	session.setContent(app.createContentFunc(session))
	return session
}

func (app *wasmApp) init(params AppParams) {

	app.params = params

	document := js.Global().Get("document")
	body := document.Call("querySelector", "body")
	head := document.Call("querySelector", "head")

	meta := document.Call("createElement", "meta")
	meta.Set("name", "viewport")
	meta.Set("content", "width=device-width")
	head.Call("appendChild", meta)

	meta = document.Call("createElement", "base")
	meta.Set("target", "_blank")
	meta.Set("rel", "noopener")
	head.Call("appendChild", meta)

	if params.Icon != "" {
		url := params.Icon
		if image, ok := resources.images[params.Icon]; ok && image.fs != nil {
			dataType := map[string]string{
				".svg":  "data:image/svg+xml",
				".png":  "data:image/png",
				".jpg":  "data:image/jpg",
				".jpeg": "data:image/jpg",
				".gif":  "data:image/gif",
			}
			ext := strings.ToLower(filepath.Ext(params.Icon))
			if prefix, ok := dataType[ext]; ok {
				if data, err := image.fs.ReadFile(image.path); err == nil {
					url = prefix + ";base64," + base64.StdEncoding.EncodeToString(data)
				} else {
					DebugLog(err.Error())
				}
			}
		}

		meta = document.Call("createElement", "link")
		meta.Set("rel", "icon")
		meta.Set("href", url)
		head.Call("appendChild", meta)
	}

	script := document.Call("createElement", "script")
	script.Set("type", "text/javascript")
	script.Set("textContent", defaultScripts+wasmScripts)
	body.Call("appendChild", script)

	js.Global().Set("sendMessage", js.FuncOf(app.handleMessage))

	app.close = make(chan DataObject)
	app.session = app.createSession()

	style := document.Call("createElement", "style")
	css := appStyles + app.session.getCurrentTheme().cssText(app.session)
	css = strings.ReplaceAll(css, `\n`, "\n")
	css = strings.ReplaceAll(css, `\t`, "\t")
	style.Set("textContent", css)
	document.Call("querySelector", "head").Call("appendChild", style)

	style = document.Call("createElement", "style")
	style.Set("id", "ruiAnimations")
	document.Call("querySelector", "head").Call("appendChild", style)

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	div := document.Call("createElement", "div")
	div.Set("className", "ruiRoot")
	div.Set("id", "ruiRootView")
	viewHTML(app.session.RootView(), buffer)
	div.Set("innerHTML", buffer.String())
	body.Call("appendChild", div)

	div = document.Call("createElement", "div")
	div.Set("className", "ruiPopupLayer")
	div.Set("id", "ruiPopupLayer")
	div.Set("style", "visibility: hidden;")
	body.Call("appendChild", div)

	div = document.Call("createElement", "a")
	div.Set("id", "ruiDownloader")
	div.Set("download", "")
	div.Set("style", "display: none;")
	body.Call("appendChild", div)

	if params.TitleColor != 0 {
		app.bridge.callFunc("setTitleColor", params.TitleColor.cssString())
	}

}

// StartApp - create the new wasmApp and start it
func StartApp(addr string, createContentFunc func(Session) SessionContent, params AppParams) {
	if createContentFunc == nil {
		return
	}

	app := new(wasmApp)
	app.createContentFunc = createContentFunc
	app.close = make(chan DataObject)
	app.bridge = createWasmBridge(app.close)

	app.init(params)
	<-app.close
}

func FinishApp() {
}

func OpenBrowser(url string) bool {
	return false
}
