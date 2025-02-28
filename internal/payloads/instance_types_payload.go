package payloads

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/go-chi/render"
)

type InstanceTypeResponse struct {
	*clients.InstanceType
}

func (s *InstanceTypeResponse) Bind(_ *http.Request) error {
	return nil
}

func (s *InstanceTypeResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewListInstanceTypeResponse(sl []*clients.InstanceType) []render.Renderer {
	list := make([]render.Renderer, len(sl))
	for i, instanceType := range sl {
		list[i] = &InstanceTypeResponse{instanceType}
	}
	return list
}
