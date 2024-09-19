package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/hvpaiva/greenlight/pkg/ujson"
)

func (a Application) healthcheckHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	a.Logger.Debug("healthcheck endpoint hit")

	data := map[string]any{
		"status": "available",
		"system_info": map[string]string{
			"version":     a.Config.Version,
			"environment": a.Config.Env,
		},
	}

	if err := ujson.Write(w, http.StatusOK, data, nil); err != nil {
		a.HandleError(w, r, err)
	}
}
