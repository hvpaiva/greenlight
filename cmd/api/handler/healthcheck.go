package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/hvpaiva/greenlight/pkg/ujson"
)

func (h *Handler) healthcheckHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	h.App.Logger.Debug("healthcheck endpoint hit")

	data := map[string]any{
		"status": "available",
		"system_info": map[string]string{
			"version":     h.App.Version,
			"environment": h.App.Env,
		},
	}

	if err := ujson.Write(w, http.StatusOK, data, nil); err != nil {
		h.App.HandleError(w, r, err)
	}
}
