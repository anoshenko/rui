
window.onfocus = function(event) {
	windowFocus = true
	sendMessage( "session-resume{session=" + sessionID +"}" );
}
