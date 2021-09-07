package rui

// SessionStartListener is the listener interface of a session start event
type SessionStartListener interface {
	OnStart(session Session)
}

// SessionFinishListener is the listener interface of a session start event
type SessionFinishListener interface {
	OnFinish(session Session)
}

// SessionResumeListener is the listener interface of a session resume event
type SessionResumeListener interface {
	OnResume(session Session)
}

// SessionPauseListener is the listener interface of a session pause event
type SessionPauseListener interface {
	OnPause(session Session)
}

// SessionPauseListener is the listener interface of a session disconnect event
type SessionDisconnectListener interface {
	OnDisconnect(session Session)
}

// SessionPauseListener is the listener interface of a session reconnect event
type SessionReconnectListener interface {
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
		if listener, ok := session.content.(SessionPauseListener); ok {
			listener.OnPause(session)
		}
	}
}

func (session *sessionData) onResume() {
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
