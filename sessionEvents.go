package rui

import "time"

// SessionStartListener is the listener interface of a session start event
type SessionStartListener interface {
	// OnStart is a function that is called by the library after the creation of the root view of the application
	OnStart(session Session)
}

// SessionFinishListener is the listener interface of a session start event
type SessionFinishListener interface {
	// OnFinish is a function that is called by the library when the user closes the application page in the browser
	OnFinish(session Session)
}

// SessionResumeListener is the listener interface of a session resume event
type SessionResumeListener interface {
	// OnResume is a function that is called by the library when the application page in the client's browser becomes
	// active and is also called immediately after OnStart
	OnResume(session Session)
}

// SessionPauseListener is the listener interface of a session pause event
type SessionPauseListener interface {
	// OnPause is a function that is called by the library when the application page in the client's browser becomes
	// inactive and is also called when the user switches to a different browser tab/window, minimizes the browser,
	// or switches to another application
	OnPause(session Session)
}

// SessionPauseListener is the listener interface of a session disconnect event
type SessionDisconnectListener interface {
	// OnDisconnect is a function that is called by the library if the server loses connection with the client and
	// this happens when the connection is broken
	OnDisconnect(session Session)
}

// SessionPauseListener is the listener interface of a session reconnect event
type SessionReconnectListener interface {
	// OnReconnect is a function that is called by the library after the server reconnects with the client
	// and this happens when the connection is restored
	OnReconnect(session Session)
}

func (session *sessionData) onStart() {
	if session.content != nil {
		if listener, ok := session.content.(SessionStartListener); ok {
			listener.OnStart(session)
		}
		session.onResume()
	}
}

func (session *sessionData) onFinish() {
	if session.content != nil {
		session.onPause()
		if listener, ok := session.content.(SessionFinishListener); ok {
			listener.OnFinish(session)
		}
	}
}

func (session *sessionData) onPause() {
	if session.content != nil {
		session.pauseTime = time.Now().Unix()
		if listener, ok := session.content.(SessionPauseListener); ok {
			listener.OnPause(session)
		}
		if timeout := session.app.Params().SocketAutoClose; timeout > 0 {
			go session.autoClose(session.pauseTime, timeout)
		}
	}
}

func (session *sessionData) autoClose(start int64, timeout int) {
	time.Sleep(time.Second * time.Duration(timeout))
	if session.pauseTime == start {
		session.bridge.callFunc("closeSocket")
	}
}

func (session *sessionData) onResume() {
	session.pauseTime = 0
	if session.content != nil {
		if listener, ok := session.content.(SessionResumeListener); ok {
			listener.OnResume(session)
		}
	}
}

func (session *sessionData) onDisconnect() {
	if session.content != nil {
		if listener, ok := session.content.(SessionDisconnectListener); ok {
			listener.OnDisconnect(session)
		}
	}
}

func (session *sessionData) onReconnect() {
	if session.content != nil {
		if listener, ok := session.content.(SessionReconnectListener); ok {
			listener.OnReconnect(session)
		}
	}
}
