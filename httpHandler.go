package rui

import (
	"net/http"
	"strings"
)

type httpHandler struct {
	app    *application
	prefix string
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		path := `/` + strings.TrimPrefix(req.URL.Path, `/`)
		req.URL.Path = `/` + strings.TrimPrefix(strings.TrimPrefix(path, h.prefix), `/`)

		h.app.ServeHTTP(w, req)
	}
}

/*
NewHandler is used to embed the rui application in third-party web frameworks (net/http, gin, echo...).
Example for echo:

	e := echo.New()
	e.Any(`/ui/*`, func()echo.HandlerFunc{
		rui.AddEmbedResources(&resources)

		h := rui.NewHandler("/ui", CreateSessionContent, rui.AppParams{
			Title: `Awesome app`,
			Icon: `favicon.png`,
		})

		return func(c echo.Context) error {
			h.ServeHTTP(c.Response(), c.Request())
			return nil
		}
	})
*/
func NewHandler(urlPrefix string, createContentFunc func(Session) SessionContent, params AppParams) *httpHandler {
	app := new(application)
	app.params = params
	app.sessions = map[int]Session{}
	app.createContentFunc = createContentFunc
	apps = append(apps, app)

	h := &httpHandler{
		app:    app,
		prefix: `/` + strings.Trim(urlPrefix, `/`),
	}

	return h
}
