let eventSource

async function sendMessage(message) {
	if (!eventSource) {
		createEventSource();
	}

	const response = await fetch('/', {
		method			: 'POST',
		body			: message,
		"Content-Type"	: "text/plain",
	  });

	const text = await response.text();
	if (text != "") {
		window.eval(text)
	}
}

function createEventSource() {
	/*
	eventSource = new EventSource("e?id="+sessionID);
	eventSource.onmessage = onEventMessage;
	eventSource.onerror = onEventError;
	*/
}

function onEventMessage(event) {
	let script = base64ToString(event.data);
	if (script != "") {
		window.eval(script)
	}
}

function onEventError(err) {
	console.log(err);
	eventSource = null;
}


window.onload = function() {
    sendMessage( sessionInfo("start-session") );
}

window.onfocus = function() {
	windowFocus = true;
	sendMessage( "session-resume{session=" + sessionID +"}" );
}

function closeSocket() {
	if (eventSource) {
		eventSource.close();
		eventSource = null;
	}
}

function sendNop() {
	sendMessage( "nop{session=" + sessionID +"}" );
}