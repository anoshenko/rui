package rui

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type WebBrige interface {
	ReadMessage() (string, bool)
	WriteMessage(text string) bool
	RunGetterScript(script string) DataObject
	AnswerReceived(answer DataObject)
	Close()
	remoteAddr() string
}

type wsBrige struct {
	conn        *websocket.Conn
	answer      map[int]chan DataObject
	answerID    int
	answerMutex sync.Mutex
	closed      bool
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 8096,
}

func CreateSocketBrige(w http.ResponseWriter, req *http.Request) WebBrige {
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
	return brige
}

func (brige *wsBrige) Close() {
	brige.closed = true
	brige.conn.Close()
}

func (brige *wsBrige) ReadMessage() (string, bool) {
	//messageType, p, err := brige.conn.ReadMessage()
	_, p, err := brige.conn.ReadMessage()
	if err != nil {
		if !brige.closed {
			ErrorLog(err.Error())
		}
		return "", false
	}

	return string(p), true
}

func (brige *wsBrige) WriteMessage(script string) bool {
	if ProtocolInDebugLog {
		DebugLog("Run script:")
		DebugLog(script)
	}
	if err := brige.conn.WriteMessage(websocket.TextMessage, []byte(script)); err != nil {
		ErrorLog(err.Error())
		return false
	}
	return true
}

func (brige *wsBrige) RunGetterScript(script string) DataObject {
	brige.answerMutex.Lock()
	answerID := brige.answerID
	brige.answerID++
	brige.answerMutex.Unlock()

	answer := make(chan DataObject)
	brige.answer[answerID] = answer
	errorText := ""
	if brige.conn != nil {
		script = "var answerID = " + strconv.Itoa(answerID) + ";\n" + script
		if ProtocolInDebugLog {
			DebugLog("\n" + script)
		}
		err := brige.conn.WriteMessage(websocket.TextMessage, []byte(script))
		if err == nil {
			return <-answer
		}
		errorText = err.Error()
	} else {
		if ProtocolInDebugLog {
			DebugLog("\n" + script)
		}
		errorText = "No connection"
	}

	result := NewDataObject("error")
	result.SetPropertyValue("text", errorText)
	delete(brige.answer, answerID)
	return result
}

func (brige *wsBrige) AnswerReceived(answer DataObject) {
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
