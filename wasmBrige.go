//go:build wasm

package rui

import (
	"strconv"
	"sync"
	"syscall/js"

	"github.com/gorilla/websocket"
)

type wasmBrige struct {
	queue       chan string
	answer      map[int]chan DataObject
	answerID    int
	answerMutex sync.Mutex
	closed      bool
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 8096,
}

func createWasmBrige() webBrige {
	brige := new(wasmBrige)
	brige.queue = make(chan string, 1000)
	brige.answerID = 1
	brige.answer = make(map[int]chan DataObject)
	brige.closed = false

	js.Global().Set("sendMessage", js.FuncOf(brige.sendMessage))

	return brige
}

func (brige *wasmBrige) sendMessage(this js.Value, args []js.Value) interface{} {
	if len(args) > 0 {
		brige.queue <- args[0].String()
	}
	return nil
}

func (brige *wasmBrige) close() {
}

func (brige *wasmBrige) readMessage() (string, bool) {
	return <-brige.queue, true
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
