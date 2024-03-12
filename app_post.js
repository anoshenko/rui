
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

/*
window.onload = function() {
	socketUrl = document.location.protocol == "https:" ? "wss://" : "ws://" 
	socketUrl += document.location.hostname
	const port = document.location.port
	if (port) {
		socketUrl += ":" + port
	}
	socketUrl += window.location.pathname + "ws"

	socket = new WebSocket(socketUrl);
	socket.onopen = socketOpen;
	socket.onclose = socketClose;
	socket.onerror = socketError;
	socket.onmessage = function(event) {
		window.execScript ? window.execScript(event.data) : window.eval(event.data);
	};
};

function socketOpen() {
	sendMessage( sessionInfo() );
}

function socketReopen() {
	sendMessage( "reconnect{session=" + sessionID + "}" );
}

function socketReconnect() {
	if (!socket) {
		socket = new WebSocket(socketUrl);
		socket.onopen = socketReopen;
		socket.onclose = socketClose;
		socket.onerror = socketError;
		socket.onmessage = function(event) {
			window.execScript ? window.execScript(event.data) : window.eval(event.data);
		};
	}
}

function socketClose(event) {
	console.log("socket closed")
	socket = null;
	if (!event.wasClean && windowFocus) {
		window.setTimeout(socketReconnect, 10000);
	}
}
*/

function socketError(error) {
	console.log(error);
}

window.onfocus = function(event) {
	windowFocus = true
	sendMessage( "session-resume{session=" + sessionID +"}" );
}
