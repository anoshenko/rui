package rui

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
			//session.callFunc("updateCSSStyle", view.htmlID(), builder.finish())
			session.updateProperty(view.htmlID(), "style", builder.finish())
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
