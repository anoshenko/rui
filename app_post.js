
async function sendMessage(message) {
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

window.onload = function() {
    sendMessage( sessionInfo() );
}

window.onfocus = function() {
	windowFocus = true
	sendMessage( "session-resume{}" );
}
