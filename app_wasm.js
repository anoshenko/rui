
function log(s) {
	console.log(s);
}

window.onfocus = function(event) {
	windowFocus = true
	sendMessage( "session-resume{session=" + sessionID +"}" );
}
