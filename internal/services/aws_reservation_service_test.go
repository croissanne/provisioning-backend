package services_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	Clientstubs "github.com/RHEnVision/provisioning-backend/internal/clients/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/dao/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/services"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	_ "github.com/RHEnVision/provisioning-backend/internal/testing/initialization"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAWSReservationHandler(t *testing.T) {
	var err error
	var json_data []byte
	ctx := stubs.WithAccountDaoOne(context.Background())
	ctx = identity.WithTenant(t, ctx)
	ctx = Clientstubs.WithSourcesClient(ctx)
	ctx = Clientstubs.WithImageBuilderClient(ctx)
	ctx = stubs.WithReservationDao(ctx)
	ctx = stubs.WithPubkeyDao(ctx)
	pk := models.Pubkey{Name: "new", Body: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC8w6DONv1qn3IdgxSpkYOClq7oe7davWFqKVHPbLoS6+dFInru7gdEO5byhTih6+PwRhHv/b1I+Mtt5MDZ8Sv7XFYpX/3P/u5zQiy1PkMSFSz0brRRUfEQxhXLW97FJa7l+bej2HJDt7f9Gvcj+d/fNWC9Z58/GX11kWk4SIXaKotkN+kWn54xGGS7Zvtm86fP59Srt6wlklSsG8mZBF7jVUjyhAgm/V5gDFb2/6jfiwSb2HyJ9/NbhLkWNdwrvpdGZqQlYhnwTfEZdpwizW/Mj3MxP5O31HN45aE0wog0UeWY4gvTl4Ogb6kescizAM6pCff3RBslbFxLdOO7cR17 lzap+rsakey@redhat.com"}
	err = stubs.AddPubkey(ctx, &pk)
	require.NoError(t, err, "failed to generate pubkey")

	values := map[string]interface{}{
		"source_id":     "1",
		"image_id":      "2bc640f6-927a-404a-9594-5b2da7e06608",
		"amount":        1,
		"instance_type": "t1.micro",
		"pubkey_id":     pk.ID,
	}
	if json_data, err = json.Marshal(values); err != nil {
		t.Fatalf("unable to marshal values to json: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "/api/provisioning/reservations/aws", bytes.NewBuffer(json_data))
	require.NoError(t, err, "failed to create request")
	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(services.CreateAWSReservation)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

	stubCount := stubs.ReservationStubCount(ctx)
	assert.Equal(t, 1, stubCount, "Reservation has not been Created through DAO")
}
