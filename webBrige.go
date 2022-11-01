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

type wsBrige struct {
	conn          *websocket.Conn
	answer        map[int]chan DataObject
	answerID      int
	answerMutex   sync.Mutex
	closed        bool
	buffer        strings.Builder
	updateScripts map[string]*strings.Builder
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 8096,
}

func CreateSocketBrige(w http.ResponseWriter, req *http.Request) webBrige {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		ErrorLog(err.Error())
		return nil
	}

	brige := new(wsBrige)
	brige.answerID = 1
	brige.answer = make(map[int]chan DataObject)
	brige.conn = conn
	brige.closed = false
	brige.updateScripts = map[string]*strings.Builder{}
	return brige
}

func (brige *wsBrige) close() {
	brige.closed = true
	brige.conn.Close()
}

func (brige *wsBrige) startUpdateScript(htmlID string) bool {
	if _, ok := brige.updateScripts[htmlID]; ok {
		return false
	}
	buffer := allocStringBuilder()
	brige.updateScripts[htmlID] = buffer
	buffer.WriteString("var element = document.getElementById('")
	buffer.WriteString(htmlID)
	buffer.WriteString("');\nif (element) {\n")
	return true
}

func (brige *wsBrige) finishUpdateScript(htmlID string) {
	if buffer, ok := brige.updateScripts[htmlID]; ok {
		buffer.WriteString("scanElementsSize();\n}\n")
		brige.writeMessage(buffer.String())
		freeStringBuilder(buffer)
		delete(brige.updateScripts, htmlID)
	}
}

func (brige *wsBrige) argToString(arg any) (string, bool) {
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

	default:
		if n, ok := isInt(arg); ok {
			return fmt.Sprintf("%d", n), true
		}
	}

	ErrorLog("Unsupported agument type")
	return "", false
}

func (brige *wsBrige) runFunc(funcName string, args ...any) bool {
	brige.buffer.Reset()
	brige.buffer.WriteString(funcName)
	brige.buffer.WriteRune('(')
	for i, arg := range args {
		argText, ok := brige.argToString(arg)
		if !ok {
			return false
		}

		if i > 0 {
			brige.buffer.WriteString(", ")
		}
		brige.buffer.WriteString(argText)
	}
	brige.buffer.WriteString(");")

	funcText := brige.buffer.String()
	if ProtocolInDebugLog {
		DebugLog("Run func: " + funcText)
	}
	if err := brige.conn.WriteMessage(websocket.TextMessage, []byte(funcText)); err != nil {
		ErrorLog(err.Error())
		return false
	}
	return true
}

func (brige *wsBrige) updateInnerHTML(htmlID, html string) {
	brige.runFunc("updateInnerHTML", htmlID, html)
}

func (brige *wsBrige) appendToInnerHTML(htmlID, html string) {
	brige.runFunc("appendToInnerHTML", htmlID, html)
}

func (brige *wsBrige) updateCSSProperty(htmlID, property, value string) {
	if buffer, ok := brige.updateScripts[htmlID]; ok {
		buffer.WriteString(`element.style['`)
		buffer.WriteString(property)
		buffer.WriteString(`'] = '`)
		buffer.WriteString(value)
		buffer.WriteString("';\n")
	} else {
		brige.runFunc("updateCSSProperty", htmlID, property, value)
	}
}

func (brige *wsBrige) updateProperty(htmlID, property string, value any) {
	if buffer, ok := brige.updateScripts[htmlID]; ok {
		if val, ok := brige.argToString(value); ok {
			buffer.WriteString(`element.setAttribute('`)
			buffer.WriteString(property)
			buffer.WriteString(`', `)
			buffer.WriteString(val)
			buffer.WriteString(");\n")
		}
	} else {
		brige.runFunc("updateProperty", htmlID, property, value)
	}
}

func (brige *wsBrige) removeProperty(htmlID, property string) {
	if buffer, ok := brige.updateScripts[htmlID]; ok {
		buffer.WriteString(`if (element.hasAttribute('`)
		buffer.WriteString(property)
		buffer.WriteString(`')) { element.removeAttribute('`)
		buffer.WriteString(property)
		buffer.WriteString("');}\n")
	} else {
		brige.runFunc("removeProperty", htmlID, property)
	}
}

func (brige *wsBrige) readMessage() (string, bool) {
	_, p, err := brige.conn.ReadMessage()
	if err != nil {
		if !brige.closed {
			ErrorLog(err.Error())
		}
		return "", false
	}

	return string(p), true
}

func (brige *wsBrige) writeMessage(script string) bool {
	if ProtocolInDebugLog {
		DebugLog("Run script:")
		DebugLog(script)
	}
	if brige.conn == nil {
		ErrorLog("No connection")
		return false
	}
	if err := brige.conn.WriteMessage(websocket.TextMessage, []byte(script)); err != nil {
		ErrorLog(err.Error())
		return false
	}
	return true
}

func (brige *wsBrige) canvasTextMetrics(htmlID, font, text string) TextMetrics {
	result := TextMetrics{}

	brige.answerMutex.Lock()
	answerID := brige.answerID
	brige.answerID++
	brige.answerMutex.Unlock()

	answer := make(chan DataObject)
	brige.answer[answerID] = answer

	if brige.runFunc("canvasTextMetrics", answerID, htmlID, font, text) {
		data := <-answer
		result.Width = dataFloatProperty(data, "width")
	}

	delete(brige.answer, answerID)
	return result
}

func (brige *wsBrige) htmlPropertyValue(htmlID, name string) string {
	brige.answerMutex.Lock()
	answerID := brige.answerID
	brige.answerID++
	brige.answerMutex.Unlock()

	answer := make(chan DataObject)
	brige.answer[answerID] = answer

	if brige.runFunc("getPropertyValue", answerID, htmlID, name) {
		data := <-answer
		if value, ok := data.PropertyValue("value"); ok {
			return value
		}
	}

	delete(brige.answer, answerID)
	return ""
}

func (brige *wsBrige) answerReceived(answer DataObject) {
	if text, ok := answer.PropertyValue("answerID"); ok {
		if id, err := strconv.Atoi(text); err == nil {
			if chanel, ok := brige.answer[id]; ok {
				chanel <- answer
				delete(brige.answer, id)
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

func (brige *wsBrige) remoteAddr() string {
	return brige.conn.RemoteAddr().String()
}
