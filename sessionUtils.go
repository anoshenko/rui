package rui

import (
	"fmt"
)

func sizeConstant(session Session, tag string) (SizeUnit, bool) {
	if text, ok := session.Constant(tag); ok {
		return StringToSizeUnit(text)
	}
	return AutoSize(), false
}

func updateCSSStyle(htmlID string, session Session) {
	if !session.ignoreViewUpdates() {
		if view := session.viewByHTMLID(htmlID); view != nil {
			var builder viewCSSBuilder

			builder.buffer = allocStringBuilder()
			builder.buffer.WriteString(`updateCSSStyle('`)
			builder.buffer.WriteString(view.htmlID())
			builder.buffer.WriteString(`', '`)
			view.cssStyle(view, &builder)
			builder.buffer.WriteString(`');`)
			view.Session().runScript(builder.finish())
		}
	}
}

func updateInnerHTML(htmlID string, session Session) {
	if !session.ignoreViewUpdates() {
		if view := session.viewByHTMLID(htmlID); view != nil {
			script := allocStringBuilder()
			defer freeStringBuilder(script)

			script.Grow(32 * 1024)
			view.htmlSubviews(view, script)
			view.Session().runScript(fmt.Sprintf(`updateInnerHTML('%v', '%v');`, view.htmlID(), script.String()))
			//view.updateEventHandlers()
		}
	}
}

func appendToInnerHTML(htmlID, content string, session Session) {
	if !session.ignoreViewUpdates() {
		if view := session.viewByHTMLID(htmlID); view != nil {
			view.Session().runScript(fmt.Sprintf(`appendToInnerHTML('%v', '%v');`, view.htmlID(), content))
			//view.updateEventHandlers()
		}
	}
}

func updateProperty(htmlID, property, value string, session Session) {
	if !session.ignoreViewUpdates() {
		session.runScript(fmt.Sprintf(`updateProperty('%v', '%v', '%v');`, htmlID, property, value))
	}
}

func updateCSSProperty(htmlID, property, value string, session Session) {
	if !session.ignoreViewUpdates() {
		session.runScript(fmt.Sprintf(`updateCSSProperty('%v', '%v', '%v');`, htmlID, property, value))
	}
}

func updateBoolProperty(htmlID, property string, value bool, session Session) {
	if !session.ignoreViewUpdates() {
		if value {
			session.runScript(fmt.Sprintf(`updateProperty('%v', '%v', true);`, htmlID, property))
		} else {
			session.runScript(fmt.Sprintf(`updateProperty('%v', '%v', false);`, htmlID, property))
		}
	}
}

func removeProperty(htmlID, property string, session Session) {
	if !session.ignoreViewUpdates() {
		session.runScript(fmt.Sprintf(`removeProperty('%v', '%v');`, htmlID, property))
	}
}

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
