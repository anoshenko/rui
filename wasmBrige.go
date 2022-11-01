//go:build wasm

package rui

import (
	"fmt"
	"strconv"
	"sync"
	"syscall/js"
)

type wasmBrige struct {
	queue       chan string
	answer      map[int]chan DataObject
	answerID    int
	answerMutex sync.Mutex
	closed      bool
}

func createWasmBrige() webBrige {
	brige := new(wasmBrige)
	brige.queue = make(chan string, 1000)
	brige.answerID = 1
	brige.answer = make(map[int]chan DataObject)
	brige.closed = false

	return brige
}

func (brige *wasmBrige) startUpdateScript(htmlID string) bool {
	return false
}

func (brige *wasmBrige) finishUpdateScript(htmlID string) {
}

func (brige *wasmBrige) runFunc(funcName string, args ...any) bool {
	if ProtocolInDebugLog {
		text := funcName + "("
		for i, arg := range args {
			if i > 0 {
				text += fmt.Sprintf(", `%v`", arg)
			} else {
				text += fmt.Sprintf("`%v`", arg)
			}
		}
		DebugLog(text + ")")
	}

	js.Global().Call(funcName, args...)
	return true
}

func (brige *wasmBrige) updateInnerHTML(htmlID, html string) {
	brige.runFunc("updateInnerHTML", htmlID, html)
}

func (brige *wasmBrige) appendToInnerHTML(htmlID, html string) {
	brige.runFunc("appendToInnerHTML", htmlID, html)
}

func (brige *wasmBrige) updateCSSProperty(htmlID, property, value string) {
	brige.runFunc("updateCSSProperty", htmlID, property, value)
}

func (brige *wasmBrige) updateProperty(htmlID, property string, value any) {
	brige.runFunc("updateProperty", htmlID, property, value)
}

func (brige *wasmBrige) removeProperty(htmlID, property string) {
	brige.runFunc("removeProperty", htmlID, property)
}

func (brige *wasmBrige) close() {
}

func (brige *wasmBrige) readMessage() (string, bool) {
	return "", false
}

func (brige *wasmBrige) writeMessage(script string) bool {
	if ProtocolInDebugLog {
		DebugLog("Run script:")
		DebugLog(script)
	}

	window := js.Global().Get("window")
	window.Call("execScript", script)

	return true
}

/*
	func (brige *wasmBrige) runGetterScript(script string) DataObject {
		brige.answerMutex.Lock()
		answerID := brige.answerID
		brige.answerID++
		brige.answerMutex.Unlock()

		answer := make(chan DataObject)
		brige.answer[answerID] = answer
		errorText := ""

		js.Global().Set("answerID", strconv.Itoa(answerID))

		window := js.Global().Get("window")
		window.Call("execScript", script)

		result := NewDataObject("error")
		result.SetPropertyValue("text", errorText)
		delete(brige.answer, answerID)
		return result
	}
*/

func (brige *wasmBrige) canvasTextMetrics(htmlID, font, text string) TextMetrics {
	// TODO
	return TextMetrics{}
}

func (brige *wasmBrige) htmlPropertyValue(htmlID, name string) string {
	return ""
}

func (brige *wasmBrige) answerReceived(answer DataObject) {
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

func (brige *wasmBrige) remoteAddr() string {
	return "localhost"
}
