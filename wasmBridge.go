//go:build wasm

package rui

import (
	"fmt"
	"strconv"
	"sync"
	"syscall/js"
)

type wasmBridge struct {
	queue       chan string
	answer      map[int]chan DataObject
	answerID    int
	answerMutex sync.Mutex
	closed      bool
	canvas      js.Value
}

func createWasmBridge() webBridge {
	bridge := new(wasmBridge)
	bridge.queue = make(chan string, 1000)
	bridge.answerID = 1
	bridge.answer = make(map[int]chan DataObject)
	bridge.closed = false

	return bridge
}

func (bridge *wasmBridge) startUpdateScript(htmlID string) bool {
	return false
}

func (bridge *wasmBridge) finishUpdateScript(htmlID string) {
}

func (bridge *wasmBridge) callFunc(funcName string, args ...any) bool {
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

func (bridge *wasmBridge) updateInnerHTML(htmlID, html string) {
	bridge.callFunc("updateInnerHTML", htmlID, html)
}

func (bridge *wasmBridge) appendToInnerHTML(htmlID, html string) {
	bridge.callFunc("appendToInnerHTML", htmlID, html)
}

func (bridge *wasmBridge) updateCSSProperty(htmlID, property, value string) {
	bridge.callFunc("updateCSSProperty", htmlID, property, value)
}

func (bridge *wasmBridge) updateProperty(htmlID, property string, value any) {
	bridge.callFunc("updateProperty", htmlID, property, value)
}

func (bridge *wasmBridge) removeProperty(htmlID, property string) {
	bridge.callFunc("removeProperty", htmlID, property)
}

func (bridge *wasmBridge) close() {
}

func (bridge *wasmBridge) readMessage() (string, bool) {
	return "", false
}

func (bridge *wasmBridge) writeMessage(script string) bool {
	if ProtocolInDebugLog {
		DebugLog("Run script:")
		DebugLog(script)
	}

	window := js.Global().Get("window")
	window.Call("execScript", script)

	return true
}

func (bridge *wasmBridge) cavnasStart(htmlID string) {
	bridge.canvas = js.Global().Call("getCanvasContext", htmlID)
}

func (bridge *wasmBridge) callCanvasFunc(funcName string, args ...any) {
	if !bridge.canvas.IsNull() {
		for i, arg := range args {
			if array, ok := arg.([]float64); ok {
				arr := make([]any, len(array))
				for k, x := range array {
					arr[k] = x
				}
				args[i] = js.ValueOf(array)
			}
		}

		bridge.canvas.Call(funcName, args...)
	}
}

func (bridge *wasmBridge) callCanvasVarFunc(v any, funcName string, args ...any) {
	if jsVar, ok := v.(js.Value); ok && !jsVar.IsNull() {
		jsVar.Call(funcName, args...)
	}
}

func (bridge *wasmBridge) callCanvasImageFunc(url string, property string, funcName string, args ...any) {
	image := js.Global().Get("images").Call("get", url)
	if !image.IsUndefined() && !image.IsNull() {
		result := image.Call(funcName, args...)
		if property != "" {
			bridge.canvas.Set(property, result)
		}
	}
}

func (bridge *wasmBridge) createCanvasVar(funcName string, args ...any) any {
	if bridge.canvas.IsNull() {
		return bridge.canvas
	}
	return bridge.canvas.Call(funcName, args...)
}

func (bridge *wasmBridge) updateCanvasProperty(property string, value any) {
	if !bridge.canvas.IsNull() {
		bridge.canvas.Set(property, value)
	}
}

func (bridge *wasmBridge) cavnasFinish() {
}

func (bridge *wasmBridge) canvasTextMetrics(htmlID, font, text string) TextMetrics {
	result := js.Global().Call("canvasTextMetrics", 0, htmlID, font, text).String()
	DebugLog(result)

	// TODO
	return TextMetrics{}
}

func (bridge *wasmBridge) htmlPropertyValue(htmlID, name string) string {
	return ""
}

func (bridge *wasmBridge) answerReceived(answer DataObject) {
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

func (bridge *wasmBridge) remoteAddr() string {
	return "localhost"
}
