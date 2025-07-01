//go:build !wasm

package rui

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type webBridge struct {
	answer              map[int]chan DataObject
	answerID            int
	answerMutex         sync.Mutex
	writeMutex          sync.Mutex
	closed              bool
	canvasBuffer        strings.Builder
	canvasVarNumber     int
	updateScripts       map[string]*strings.Builder
	writeMessage        func(string) bool
	callFuncImmediately func(funcName string, args ...any) bool
}

type wsBridge struct {
	webBridge
	conn *websocket.Conn
}

type httpBridge struct {
	webBridge
	responseBuffer strings.Builder
	response       chan string
	remoteAddress  string
	//conn *websocket.Conn
}

type canvasVar struct {
	name string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 8096,
}

func createSocketBridge(w http.ResponseWriter, req *http.Request) *wsBridge {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		ErrorLog(err.Error())
		return nil
	}

	bridge := new(wsBridge)
	bridge.initBridge()
	bridge.conn = conn
	bridge.writeMessage = func(script string) bool {
		if ProtocolInDebugLog {
			DebugLog("🖥️ <- " + script)
		}

		if bridge.conn == nil {
			ErrorLog("No connection")
			return false
		}

		bridge.writeMutex.Lock()
		err := bridge.conn.WriteMessage(websocket.TextMessage, []byte(script))
		bridge.writeMutex.Unlock()

		if err != nil {
			ErrorLog(err.Error())
			return false
		}
		return true
	}
	bridge.callFuncImmediately = bridge.callFunc
	return bridge
}

func createHttpBridge(req *http.Request) *httpBridge {
	bridge := new(httpBridge)
	bridge.initBridge()
	bridge.response = make(chan string, 10)
	bridge.writeMessage = func(script string) bool {
		if script != "" {
			if ProtocolInDebugLog {
				DebugLog(script)
			}

			if bridge.responseBuffer.Len() > 0 {
				bridge.responseBuffer.WriteRune('\n')
			}
			bridge.responseBuffer.WriteString(script)
		}
		return true
	}
	bridge.callFuncImmediately = bridge.callImmediately
	bridge.remoteAddress = req.RemoteAddr
	return bridge
}

func (bridge *webBridge) initBridge() {
	bridge.answerID = 1
	bridge.answer = make(map[int]chan DataObject)
	bridge.closed = false
	bridge.updateScripts = map[string]*strings.Builder{}
}

func (bridge *webBridge) startUpdateScript(htmlID string) bool {
	if _, ok := bridge.updateScripts[htmlID]; ok {
		return false
	}
	buffer := allocStringBuilder()
	bridge.updateScripts[htmlID] = buffer
	buffer.WriteString("{\nlet element = document.getElementById('")
	buffer.WriteString(htmlID)
	buffer.WriteString("');\nif (element) {\n")
	return true
}

func (bridge *webBridge) finishUpdateScript(htmlID string) {
	if buffer, ok := bridge.updateScripts[htmlID]; ok {
		buffer.WriteString("scanElementsSize();\n}\n}\n")
		bridge.writeMessage(buffer.String())

		freeStringBuilder(buffer)
		delete(bridge.updateScripts, htmlID)
	}
}

func (bridge *webBridge) argToString(arg any) (string, bool) {
	switch arg := arg.(type) {
	case string:
		arg = strings.ReplaceAll(arg, "\\", `\\`)
		arg = strings.ReplaceAll(arg, "'", `\'`)
		arg = strings.ReplaceAll(arg, "\n", `\n`)
		arg = strings.ReplaceAll(arg, "\r", `\r`)
		arg = strings.ReplaceAll(arg, "\t", `\t`)
		arg = strings.ReplaceAll(arg, "\b", `\b`)
		arg = strings.ReplaceAll(arg, "\f", `\f`)
		arg = strings.ReplaceAll(arg, "\v", `\v`)
		return `'` + arg + `'`, true

	case rune:
		switch arg {
		case '\t':
			return `'\t'`, true
		case '\r':
			return `'\r'`, true
		case '\n':
			return `'\n'`, true
		case '\b':
			return `'\b'`, true
		case '\f':
			return `'\f'`, true
		case '\v':
			return `'\v'`, true
		case '\'':
			return `'\''`, true
		case '\\':
			return `'\\'`, true
		}
		if arg < ' ' {
			return fmt.Sprintf(`'\x%02d'`, int(arg)), true
		}
		return `'` + string(arg) + `'`, true

	case bool:
		if arg {
			return "true", true
		} else {
			return "false", true
		}

	case float32:
		return fmt.Sprintf("%g", float64(arg)), true

	case float64:
		return fmt.Sprintf("%g", arg), true

	case []float64:
		buffer := allocStringBuilder()
		defer freeStringBuilder(buffer)
		lead := '['
		for _, val := range arg {
			buffer.WriteRune(lead)
			lead = ','
			fmt.Fprintf(buffer, "%g", val)
		}
		buffer.WriteRune(']')
		return buffer.String(), true

	case canvasVar:
		return arg.name, true

	default:
		if n, ok := isInt(arg); ok {
			return fmt.Sprintf("%d", n), true
		}
	}

	ErrorLog("Unsupported argument type")
	return "", false
}

func (bridge *webBridge) callFuncScript(funcName string, args ...any) (string, bool) {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString(funcName)
	buffer.WriteRune('(')
	for i, arg := range args {
		argText, ok := bridge.argToString(arg)
		if !ok {
			return "", false
		}

		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(argText)
	}
	buffer.WriteString(");")

	return buffer.String(), true
}

func (bridge *webBridge) callFunc(funcName string, args ...any) bool {
	if funcText, ok := bridge.callFuncScript(funcName, args...); ok {
		return bridge.writeMessage(funcText)
	}
	return false
}

func (bridge *webBridge) updateInnerHTML(htmlID, html string) {
	bridge.callFunc("updateInnerHTML", htmlID, html)
}

func (bridge *webBridge) appendToInnerHTML(htmlID, html string) {
	bridge.callFunc("appendToInnerHTML", htmlID, html)
}

func (bridge *webBridge) updateCSSProperty(htmlID, property, value string) {
	if buffer, ok := bridge.updateScripts[htmlID]; ok {
		buffer.WriteString(`element.style['`)
		buffer.WriteString(property)
		buffer.WriteString(`'] = '`)
		buffer.WriteString(value)
		buffer.WriteString("';\n")
	} else {
		bridge.callFunc("updateCSSProperty", htmlID, property, value)
	}
}

func (bridge *webBridge) updateProperty(htmlID, property string, value any) {
	if buffer, ok := bridge.updateScripts[htmlID]; ok {
		if val, ok := bridge.argToString(value); ok {
			buffer.WriteString(`element.setAttribute('`)
			buffer.WriteString(property)
			buffer.WriteString(`', `)
			buffer.WriteString(val)
			buffer.WriteString(");\n")
		}
	} else {
		bridge.callFunc("updateProperty", htmlID, property, value)
	}
}

func (bridge *webBridge) removeProperty(htmlID, property string) {
	if buffer, ok := bridge.updateScripts[htmlID]; ok {
		buffer.WriteString(`if (element.hasAttribute('`)
		buffer.WriteString(property)
		buffer.WriteString(`')) { element.removeAttribute('`)
		buffer.WriteString(property)
		buffer.WriteString("');}\n")
	} else {
		bridge.callFunc("removeProperty", htmlID, property)
	}
}

func (bridge *webBridge) appendAnimationCSS(css string) {
	bridge.writeMessage(`{
	let styles = document.getElementById('ruiAnimations');
	if (styles) {
		styles.textContent += '` + css + `';
	}
}`)
}

func (bridge *webBridge) setAnimationCSS(css string) {
	bridge.writeMessage(`{
	let styles = document.getElementById('ruiAnimations');
	if (styles) {
		styles.textContent = '` + css + `';
	}
}`)
}

func (bridge *webBridge) canvasStart(htmlID string) {
	bridge.canvasBuffer.Reset()
	bridge.canvasBuffer.WriteString("{\nconst ctx = getCanvasContext('")
	bridge.canvasBuffer.WriteString(htmlID)
	bridge.canvasBuffer.WriteString(`');`)
}

func (bridge *webBridge) callCanvasFunc(funcName string, args ...any) {
	bridge.canvasBuffer.WriteString("\nctx.")
	bridge.canvasBuffer.WriteString(funcName)
	bridge.canvasBuffer.WriteRune('(')
	for i, arg := range args {
		if i > 0 {
			bridge.canvasBuffer.WriteString(", ")
		}
		argText, _ := bridge.argToString(arg)
		bridge.canvasBuffer.WriteString(argText)
	}
	bridge.canvasBuffer.WriteString(");")
}

func (bridge *webBridge) updateCanvasProperty(property string, value any) {
	bridge.canvasBuffer.WriteString("\nctx.")
	bridge.canvasBuffer.WriteString(property)
	bridge.canvasBuffer.WriteString(" = ")
	argText, _ := bridge.argToString(value)
	bridge.canvasBuffer.WriteString(argText)
	bridge.canvasBuffer.WriteString(";")
}

func (bridge *webBridge) createCanvasVar(funcName string, args ...any) any {
	bridge.canvasVarNumber++
	result := canvasVar{name: fmt.Sprintf("v%d", bridge.canvasVarNumber)}
	bridge.canvasBuffer.WriteString("\nlet ")
	bridge.canvasBuffer.WriteString(result.name)
	bridge.canvasBuffer.WriteString(" = ctx.")
	bridge.canvasBuffer.WriteString(funcName)
	bridge.canvasBuffer.WriteRune('(')
	for i, arg := range args {
		if i > 0 {
			bridge.canvasBuffer.WriteString(", ")
		}
		argText, _ := bridge.argToString(arg)
		bridge.canvasBuffer.WriteString(argText)
	}
	bridge.canvasBuffer.WriteString(");")
	return result
}

func (bridge *webBridge) createPath2D(arg string) any {
	bridge.canvasVarNumber++
	result := canvasVar{name: fmt.Sprintf("v%d", bridge.canvasVarNumber)}
	bridge.canvasBuffer.WriteString("\nlet ")
	bridge.canvasBuffer.WriteString(result.name)
	bridge.canvasBuffer.WriteString(` = new Path2D(`)
	if arg != "" {
		argText, _ := bridge.argToString(arg)
		bridge.canvasBuffer.WriteString(argText)
	}
	bridge.canvasBuffer.WriteString(`);`)
	return result
}

func (bridge *webBridge) callCanvasVarFunc(v any, funcName string, args ...any) {
	varName, ok := v.(canvasVar)
	if !ok {
		return
	}
	bridge.canvasBuffer.WriteString("\n")
	bridge.canvasBuffer.WriteString(varName.name)
	bridge.canvasBuffer.WriteRune('.')
	bridge.canvasBuffer.WriteString(funcName)
	bridge.canvasBuffer.WriteRune('(')
	for i, arg := range args {
		if i > 0 {
			bridge.canvasBuffer.WriteString(", ")
		}
		argText, _ := bridge.argToString(arg)
		bridge.canvasBuffer.WriteString(argText)
	}
	bridge.canvasBuffer.WriteString(");")
}

func (bridge *webBridge) callCanvasImageFunc(url string, property string, funcName string, args ...any) {

	bridge.canvasBuffer.WriteString("\nimg = images.get('")
	bridge.canvasBuffer.WriteString(url)
	bridge.canvasBuffer.WriteString("');\nif (img) {\n")
	if property != "" {
		bridge.canvasBuffer.WriteString("ctx.")
		bridge.canvasBuffer.WriteString(property)
		bridge.canvasBuffer.WriteString(" = ")
	}
	bridge.canvasBuffer.WriteString("ctx.")
	bridge.canvasBuffer.WriteString(funcName)
	bridge.canvasBuffer.WriteString("(img")
	for _, arg := range args {
		bridge.canvasBuffer.WriteString(", ")
		argText, _ := bridge.argToString(arg)
		bridge.canvasBuffer.WriteString(argText)
	}
	bridge.canvasBuffer.WriteString(");\n}")
}

func (bridge *webBridge) canvasFinish() {
	bridge.canvasBuffer.WriteString("\n}\n")
	bridge.writeMessage(bridge.canvasBuffer.String())
}

func (bridge *webBridge) remoteValue(funcName string, args ...any) (DataObject, bool) {
	bridge.answerMutex.Lock()
	answerID := bridge.answerID
	bridge.answerID++
	bridge.answerMutex.Unlock()

	answer := make(chan DataObject)
	bridge.answer[answerID] = answer

	funcArgs := append([]any{answerID}, args...)
	var result DataObject = nil

	if bridge.callFuncImmediately(funcName, funcArgs...) {
		result = <-answer
	}

	close(answer)
	delete(bridge.answer, answerID)
	return result, true
}

func (bridge *webBridge) canvasTextMetrics(htmlID, font, text string) TextMetrics {
	result := TextMetrics{}
	if data, ok := bridge.remoteValue("canvasTextMetrics", htmlID, font, text); ok {
		result.Width = dataFloatProperty(data, "width")
	}
	return result
}

func (bridge *webBridge) htmlPropertyValue(htmlID, name string) string {
	if data, ok := bridge.remoteValue("getPropertyValue", htmlID, name); ok {
		if value, ok := data.PropertyValue("value"); ok {
			return value
		}
	}
	return ""
}

func (bridge *webBridge) answerReceived(answer DataObject) {
	if text, ok := answer.PropertyValue("answerID"); ok {
		if id, err := strconv.Atoi(text); err == nil {
			if chanel, ok := bridge.answer[id]; ok {
				chanel <- answer
				delete(bridge.answer, id)
			} else {
				ErrorLog("Bad answerID = " + text + " (chan not found)")
			}
		} else {
			ErrorLog("Invalid answerID = " + text)
		}
	} else {
		ErrorLog("answerID not found")
	}
}

func (bridge *wsBridge) close() {
	bridge.closed = true
	defer bridge.conn.Close()
	bridge.conn = nil
}

func (bridge *wsBridge) readMessage() (string, bool) {
	_, p, err := bridge.conn.ReadMessage()
	if err != nil {
		if !bridge.closed {
			ErrorLog(err.Error())
		}
		return "", false
	}

	return string(p), true
}

func (bridge *wsBridge) sendResponse() {
}

func (bridge *wsBridge) remoteAddr() string {
	return bridge.conn.RemoteAddr().String()
}

func (bridge *httpBridge) close() {
	bridge.closed = true
	// TODO
}

func (bridge *httpBridge) callImmediately(funcName string, args ...any) bool {
	if funcText, ok := bridge.callFuncScript(funcName, args...); ok {
		if ProtocolInDebugLog {
			DebugLog("Run func: " + funcText)
		}
		bridge.response <- funcText
		return true
	}
	return false
}

func (bridge *httpBridge) sendResponse() {
	bridge.writeMutex.Lock()
	text := bridge.responseBuffer.String()
	bridge.responseBuffer.Reset()
	bridge.writeMutex.Unlock()
	bridge.response <- text
}

func (bridge *httpBridge) remoteAddr() string {
	return bridge.remoteAddress
}
