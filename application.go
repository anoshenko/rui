package rui

import (
	_ "embed"
	"strings"
)

//go:embed app_scripts.js
var defaultScripts string

//go:embed app_styles.css
var appStyles string

//go:embed defaultTheme.rui
var defaultThemeText string

// Application represent generic application interface, see also [Session]
type Application interface {
	// Finish finishes the application
	Finish()

	// Params returns application parameters set by StartApp function
	Params() AppParams

	removeSession(id int)
}

// AppParams defines parameters of the app
type AppParams struct {
	// Title - title of the app window/tab
	Title string

	// TitleColor - background color of the app window/tab (applied only for Safari and Chrome for Android)
	TitleColor Color

	// Icon - the icon file name
	Icon string

	// CertFile - path of a certificate for the server must be provided
	// if neither the Server's TLSConfig.Certificates nor TLSConfig.GetCertificate are populated.
	// If the certificate is signed by a certificate authority, the certFile should be the concatenation
	// of the server's certificate, any intermediates, and the CA's certificate.
	CertFile string

	AutoCertDomain string

	// KeyFile - path of a private key for the server must be provided
	// if neither the Server's TLSConfig.Certificates nor TLSConfig.GetCertificate are populated.
	KeyFile string

	// Redirect80 - if true then the function of redirect from port 80 to 443 is created
	Redirect80 bool

	// NoSocket - if true then WebSockets will not be used and information exchange
	// between the client and the server will be carried out only via http.
	NoSocket bool

	// SocketAutoClose - time in seconds after which the socket is automatically closed for an inactive session.
	// The countdown begins after the OnPause event arrives.
	// If the value of this property is less than or equal to 0 then the socket is not closed.
	SocketAutoClose int
}

func getStartPage(buffer *strings.Builder, params AppParams) {
	buffer.WriteString(`<head>
		<meta charset="utf-8">
		<title>`)
	buffer.WriteString(params.Title)
	buffer.WriteString("</title>")
	if params.Icon != "" {
		buffer.WriteString(`
		<link rel="icon" href="`)
		buffer.WriteString(params.Icon)
		buffer.WriteString(`">`)
	}

	if params.TitleColor != 0 {
		buffer.WriteString(`
		<meta name="theme-color" content="`)
		buffer.WriteString(params.TitleColor.cssString())
		buffer.WriteString(`">`)
	}

	buffer.WriteString(`
		<base target="_blank" rel="noopener">
		<meta name="viewport" content="width=device-width">
		<style>`)
	buffer.WriteString(appStyles)
	buffer.WriteString(`</style>
		<style id="ruiAnimations"></style>
		<script src="/script.js"></script>
	</head>
	<body id="body" onkeydown="keyDownEvent(this, event)">
		<div class="ruiRoot" id="ruiRootView"></div>
		<div class="ruiPopupLayer" id="ruiPopupLayer" style="visibility: hidden; isolation: isolate;"></div>
		<div class="ruiTooltipLayer" id="ruiTooltipLayer" style="visibility: hidden; opacity: 0;">
		<div id="ruiTooltipText" class="ruiTooltipText"></div>
		<div id="ruiTooltipTopArrow" class="ruiTooltipTopArrow"></div>
			<div id="ruiTooltipBottomArrow" class="ruiTooltipBottomArrow"></div>
		</div>
		
		<a id="ruiDownloader" download style="display: none;"></a>
	</body>`)
}
