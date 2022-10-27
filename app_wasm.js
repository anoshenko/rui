

function initSession() {

	const touch_screen = (('ontouchstart' in document.documentElement) || (navigator.maxTouchPoints > 0) || (navigator.msMaxTouchPoints > 0)) ? "1" : "0";
	var message = "sessionInfo{touch=" + touch_screen 
	
	const style = window.getComputedStyle(document.body);
	if (style) {
		var direction = style.getPropertyValue('direction');
		if (direction) {
			message += ",direction=" + direction
		}
	}

	const lang = window.navigator.language;
	if (lang) {
		message += ",language=\"" + lang + "\"";
	}

	const langs = window.navigator.languages;
	if (langs) {
		message += ",languages=\"" + langs + "\"";
	}

	const userAgent = window.navigator.userAgent
	if (userAgent) {
		message += ",user-agent=\"" + userAgent + "\"";
	}

	const darkThemeMq = window.matchMedia("(prefers-color-scheme: dark)");
	if (darkThemeMq.matches) {
		message += ",dark=1";
	} 

	const pixelRatio = window.devicePixelRatio;
	if (pixelRatio) {
		message += ",pixel-ratio=" + pixelRatio;
	}

	startSession( message + "}" );
}


window.onfocus = function(event) {
	windowFocus = true
	sendMessage( "session-resume{session=" + sessionID +"}" );
}
