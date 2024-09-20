package handler

import (
	"net/http"

	"github.com/hvpaiva/greenlight/cmd/api/erro"
	"github.com/hvpaiva/greenlight/pkg/ujson"
)

func (h *Handler) healthcheckHandler(w http.ResponseWriter, _ *http.Request) error {
	h.App.Logger.Debug("healthcheck endpoint hit")

	data := map[string]any{
		"status": "available",
		"system_info": map[string]string{
			"version":     h.App.Version,
			"environment": h.App.Env,
		},
	}

	if err := ujson.Write(w, http.StatusOK, data, nil); err != nil {
		return erro.ThrowInternalServer("error writing response", err)
	}

	return nil
}
