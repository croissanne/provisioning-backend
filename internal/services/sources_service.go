package services

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/render"
)

func ListSources(w http.ResponseWriter, r *http.Request) {
	client, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}
	sourcesList, err := client.ListProvisioningSources(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}
	if err := render.RenderList(w, r, payloads.NewListSourcesResponse(sourcesList)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render sources list", err))
		return
	}
}
