let socket;

function sendMessage(message) {
	if (!socket) {
		createSocket(function() {
			sendReconnectMessage();
			if (!windowFocus) {
				windowFocus = true;
				sendMessage( "session-resume{session=" + sessionID +"}" );
			}
			socket.send(message);
		});
	} else {
		socket.send(message);
	}
}

function createSocket(onopen) {
	let socketUrl = document.location.protocol == "https:" ? "wss://" : "ws://" ;
	socketUrl += document.location.hostname;
	const port = document.location.port;
	if (port) {
		socketUrl += ":" + port;
	}
	socketUrl += window.location.pathname + "ws";

	socket = new WebSocket(socketUrl);
	socket.onopen = onopen;
	socket.onclose = onSocketClose;
	socket.onerror = onSocketError;
	socket.onmessage = function(event) {
		window.execScript ? window.execScript(event.data) : window.eval(event.data);
	};
} 

function closeSocket() {
	if (socket) {
		socket.close();
	}
}

window.onload = createSocket(function() {
	sendMessage( sessionInfo("start-session") );
});

window.onfocus = function() {
	windowFocus = true;
	if (!socket) {
		createSocket(function() {
			sendReconnectMessage();
			sendMessage( "session-resume{session=" + sessionID +"}" );
		});
	} else {
		sendMessage( "session-resume{session=" + sessionID +"}" );
	}
}

function socketReconnect() {
	if (!socket) {
		createSocket(sendReconnectMessage);
	}
}

function onSocketClose(event) {
	console.log("socket closed");
	socket = null;
	if (!event.wasClean && windowFocus) {
		window.setTimeout(socketReconnect, 10000);
	}
}

function onSocketError(error) {
	console.log(error);
}
