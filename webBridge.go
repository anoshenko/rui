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

type wsBridge struct {
	conn            *websocket.Conn
	answer          map[int]chan DataObject
	answerID        int
	answerMutex     sync.Mutex
	closed          bool
	buffer          strings.Builder
	canvasBuffer    strings.Builder
	canvasVarNumber int
	updateScripts   map[string]*strings.Builder
}

type canvasVar struct {
	name string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 8096,
}

func CreateSocketBridge(w http.ResponseWriter, req *http.Request) webBridge {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		ErrorLog(err.Error())
		return nil
	}

	bridge := new(wsBridge)
	bridge.answerID = 1
	bridge.answer = make(map[int]chan DataObject)
	bridge.conn = conn
	bridge.closed = false
	bridge.updateScripts = map[string]*strings.Builder{}
	return bridge
}

func (bridge *wsBridge) close() {
	bridge.closed = true
	bridge.conn.Close()
}

func (bridge *wsBridge) startUpdateScript(htmlID string) bool {
	if _, ok := bridge.updateScripts[htmlID]; ok {
		return false
	}
	buffer := allocStringBuilder()
	bridge.updateScripts[htmlID] = buffer
	buffer.WriteString("var element = document.getElementById('")
	buffer.WriteString(htmlID)
	buffer.WriteString("');\nif (element) {\n")
	return true
}

func (bridge *wsBridge) finishUpdateScript(htmlID string) {
	if buffer, ok := bridge.updateScripts[htmlID]; ok {
		buffer.WriteString("scanElementsSize();\n}\n")
		bridge.writeMessage(buffer.String())
		freeStringBuilder(buffer)
		delete(bridge.updateScripts, htmlID)
	}
}

func (bridge *wsBridge) argToString(arg any) (string, bool) {
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
			buffer.WriteString(fmt.Sprintf("%g", val))
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

func (bridge *wsBridge) callFunc(funcName string, args ...any) bool {
	bridge.buffer.Reset()
	bridge.buffer.WriteString(funcName)
	bridge.buffer.WriteRune('(')
	for i, arg := range args {
		argText, ok := bridge.argToString(arg)
		if !ok {
			return false
		}

		if i > 0 {
			bridge.buffer.WriteString(", ")
		}
		bridge.buffer.WriteString(argText)
	}
	bridge.buffer.WriteString(");")

	funcText := bridge.buffer.String()
	if ProtocolInDebugLog {
		DebugLog("Run func: " + funcText)
	}
	if err := bridge.conn.WriteMessage(websocket.TextMessage, []byte(funcText)); err != nil {
		ErrorLog(err.Error())
		return false
	}
	return true
}

func (bridge *wsBridge) updateInnerHTML(htmlID, html string) {
	bridge.callFunc("updateInnerHTML", htmlID, html)
}

func (bridge *wsBridge) appendToInnerHTML(htmlID, html string) {
	bridge.callFunc("appendToInnerHTML", htmlID, html)
}

func (bridge *wsBridge) updateCSSProperty(htmlID, property, value string) {
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

func (bridge *wsBridge) updateProperty(htmlID, property string, value any) {
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

func (bridge *wsBridge) removeProperty(htmlID, property string) {
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

func (bridge *wsBridge) addAnimationCSS(css string) {
	bridge.writeMessage(`var styles = document.getElementById('ruiAnimations');
if (styles) {
	styles.textContent += '` + css + `';
}`)
}

func (bridge *wsBridge) clearAnimation() {
	bridge.writeMessage(`var styles = document.getElementById('ruiAnimations');
if (styles) {
	styles.textContent = '';
}`)
}

func (bridge *wsBridge) canvasStart(htmlID string) {
	bridge.canvasBuffer.Reset()
	bridge.canvasBuffer.WriteString(`const ctx = getCanvasContext('`)
	bridge.canvasBuffer.WriteString(htmlID)
	bridge.canvasBuffer.WriteString(`');`)
}

func (bridge *wsBridge) callCanvasFunc(funcName string, args ...any) {
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

func (bridge *wsBridge) updateCanvasProperty(property string, value any) {
	bridge.canvasBuffer.WriteString("\nctx.")
	bridge.canvasBuffer.WriteString(property)
	bridge.canvasBuffer.WriteString(" = ")
	argText, _ := bridge.argToString(value)
	bridge.canvasBuffer.WriteString(argText)
	bridge.canvasBuffer.WriteString(";")
}

func (bridge *wsBridge) createCanvasVar(funcName string, args ...any) any {
	bridge.canvasVarNumber++
	result := canvasVar{name: fmt.Sprintf("v%d", bridge.canvasVarNumber)}
	bridge.canvasBuffer.WriteString("\nvar ")
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

func (bridge *wsBridge) callCanvasVarFunc(v any, funcName string, args ...any) {
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

func (bridge *wsBridge) callCanvasImageFunc(url string, property string, funcName string, args ...any) {

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

func (bridge *wsBridge) canvasFinish() {
	bridge.canvasBuffer.WriteString("\n")
	script := bridge.canvasBuffer.String()
	if ProtocolInDebugLog {
		DebugLog("Run script:")
		DebugLog(script)
	}
	if bridge.conn == nil {
		ErrorLog("No connection")
	} else if err := bridge.conn.WriteMessage(websocket.TextMessage, []byte(script)); err != nil {
		ErrorLog(err.Error())
	}
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

func (bridge *wsBridge) writeMessage(script string) bool {
	if ProtocolInDebugLog {
		DebugLog("Run script:")
		DebugLog(script)
	}
	if bridge.conn == nil {
		ErrorLog("No connection")
		return false
	}
	if err := bridge.conn.WriteMessage(websocket.TextMessage, []byte(script)); err != nil {
		ErrorLog(err.Error())
		return false
	}
	return true
}

func (bridge *wsBridge) canvasTextMetrics(htmlID, font, text string) TextMetrics {
	result := TextMetrics{}

	bridge.answerMutex.Lock()
	answerID := bridge.answerID
	bridge.answerID++
	bridge.answerMutex.Unlock()

	answer := make(chan DataObject)
	bridge.answer[answerID] = answer

	if bridge.callFunc("canvasTextMetrics", answerID, htmlID, font, text) {
		data := <-answer
		result.Width = dataFloatProperty(data, "width")
	}

	delete(bridge.answer, answerID)
	return result
}

func (bridge *wsBridge) htmlPropertyValue(htmlID, name string) string {
	bridge.answerMutex.Lock()
	answerID := bridge.answerID
	bridge.answerID++
	bridge.answerMutex.Unlock()

	answer := make(chan DataObject)
	bridge.answer[answerID] = answer

	if bridge.callFunc("getPropertyValue", answerID, htmlID, name) {
		data := <-answer
		if value, ok := data.PropertyValue("value"); ok {
			return value
		}
	}

	delete(bridge.answer, answerID)
	return ""
}

func (bridge *wsBridge) answerReceived(answer DataObject) {
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

func (bridge *wsBridge) remoteAddr() string {
	return bridge.conn.RemoteAddr().String()
}
