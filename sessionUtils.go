package rui

import (
	"strings"
)

func (session *sessionData) startUpdateScript(htmlID string) {
	buffer := allocStringBuilder()
	session.updateScripts[htmlID] = buffer
	buffer.WriteString("var element = document.getElementById('")
	buffer.WriteString(htmlID)
	buffer.WriteString("');\nif (element) {\n")
}

func (session *sessionData) updateScript(htmlID string) *strings.Builder {
	if buffer, ok := session.updateScripts[htmlID]; ok {
		return buffer
	}
	return nil
}

func (session *sessionData) finishUpdateScript(htmlID string) {
	if buffer, ok := session.updateScripts[htmlID]; ok {
		buffer.WriteString("scanElementsSize();\n}\n")
		session.runScript(buffer.String())
		freeStringBuilder(buffer)
		delete(session.updateScripts, htmlID)
	}
}

func sizeConstant(session Session, tag string) (SizeUnit, bool) {
	if text, ok := session.Constant(tag); ok {
		return StringToSizeUnit(text)
	}
	return AutoSize(), false
}

func updateCSSStyle(htmlID string, session Session) {
	if !session.ignoreViewUpdates() {
		if view := session.viewByHTMLID(htmlID); view != nil {
			builder := viewCSSBuilder{buffer: allocStringBuilder()}
			view.cssStyle(view, &builder)
			session.runFunc("updateCSSStyle", view.htmlID(), builder.finish())
		}
	}
}

func updateInnerHTML(htmlID string, session Session) {
	if !session.ignoreViewUpdates() {
		var view View
		if htmlID == "ruiRootView" {
			view = session.RootView()
		} else {
			view = session.viewByHTMLID(htmlID)
		}
		if view != nil {
			script := allocStringBuilder()
			defer freeStringBuilder(script)

			script.Grow(32 * 1024)
			view.htmlSubviews(view, script)
			session.updateInnerHTML(view.htmlID(), script.String())
		}
	}
}

/*
func updateProperty(htmlID, property, value string, session Session) {
	if !session.ignoreViewUpdates() {
		if buffer := session.updateScript(htmlID); buffer != nil {
			buffer.WriteString(fmt.Sprintf(`element.setAttribute('%v', '%v');`, property, value))
			buffer.WriteRune('\n')
		} else {
			session.runFunc("updateProperty", htmlID, property, value)
		}
	}
}

func updateCSSProperty(htmlID, property, value string, session Session) {
	if !session.ignoreViewUpdates() {
		if buffer := session.updateScript(htmlID); buffer != nil {
			buffer.WriteString(fmt.Sprintf(`element.style['%v'] = '%v';`, property, value))
			buffer.WriteRune('\n')
		} else {
			session.runFunc("updateCSSProperty", htmlID, property, value)
		}
	}
}

func updateBoolProperty(htmlID, property string, value bool, session Session) {
	if !session.ignoreViewUpdates() {
		if buffer := session.updateScript(htmlID); buffer != nil {
			if value {
				buffer.WriteString(fmt.Sprintf(`element.setAttribute('%v', true);`, property))
			} else {
				buffer.WriteString(fmt.Sprintf(`element.setAttribute('%v', false);`, property))
			}
			buffer.WriteRune('\n')
		} else {
			session.runFunc("updateProperty", htmlID, property, value)
		}
	}
}

func removeProperty(htmlID, property string, session Session) {
	if !session.ignoreViewUpdates() {
		if buffer := session.updateScript(htmlID); buffer != nil {
			buffer.WriteString(fmt.Sprintf(`if (element.hasAttribute('%v')) { element.removeAttribute('%v');}`, property, property))
			buffer.WriteRune('\n')
		} else {
			session.runFunc("removeProperty", htmlID, property)
		}
	}
}
*/
/*
func setDisabled(htmlID string, disabled bool, session Session) {
	if !session.ignoreViewUpdates() {
		if disabled {
			session.runScript(fmt.Sprintf(`setDisabled('%v', true);`, htmlID))
		} else {
			session.runScript(fmt.Sprintf(`setDisabled('%v', false);`, htmlID))
		}
	}
}
*/

func viewByHTMLID(id string, startView View) View {
	if startView != nil {
		if startView.htmlID() == id {
			return startView
		}
		if container, ok := startView.(ParanetView); ok {
			for _, view := range container.Views() {
				if view != nil {
					if v := viewByHTMLID(id, view); v != nil {
						return v
					}
				}
			}
		}
	}
	return nil
}
