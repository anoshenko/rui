
window.onfocus = function() {
	windowFocus = true
	sendMessage( "session-resume{session=" + sessionID +"}" );
}

function closeSocket() {
}
