package handlers

import (
	"net/http"

	"git.klink.asia/paul/certman/views"
)

func NotFoundHandler(w http.ResponseWriter, req *http.Request) {
	view := views.New(req)
	view.RenderError(w, http.StatusNotFound)
}

func ErrorHandler(w http.ResponseWriter, req *http.Request) {
	view := views.New(req)
	view.RenderError(w, http.StatusInternalServerError)
}

func CSRFErrorHandler(w http.ResponseWriter, req *http.Request) {
	view := views.New(req)
	view.RenderError(w, http.StatusForbidden)
}
