package services_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/services"
	_ "github.com/RHEnVision/provisioning-backend/internal/testing/initialization"
	"github.com/stretchr/testify/require"

	"github.com/RHEnVision/provisioning-backend/internal/dao/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/testing/factories"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/stretchr/testify/assert"
)

func TestListPubkeysHandler(t *testing.T) {
	ctx := stubs.WithAccountDaoOne(context.Background())
	ctx = identity.WithTenant(t, ctx)
	ctx = stubs.WithPubkeyDao(ctx)
	err := stubs.AddPubkey(ctx, &models.Pubkey{
		Name: factories.GetSequenceName("pubkey"),
		Body: factories.GenerateRSAPubKey(t),
	})
	require.NoError(t, err, "failed to add stubbed key")
	err = stubs.AddPubkey(ctx, &models.Pubkey{
		Name: factories.GetSequenceName("pubkey"),
		Body: factories.GenerateRSAPubKey(t),
	})
	require.NoError(t, err, "failed to add stubbed key")

	req, err := http.NewRequestWithContext(ctx, "GET", "/api/provisioning/pubkeys", nil)
	require.NoError(t, err, "failed to create request")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(services.ListPubkeys)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, "Wrong status code")

	var pubkeys []models.Pubkey

	err = json.NewDecoder(rr.Body).Decode(&pubkeys)
	require.NoError(t, err, "failed to decode response body")

	assert.Equal(t, 2, len(pubkeys), "expected two pubkeys in response json")
}

func TestCreatePubkeyHandler(t *testing.T) {
	var err error
	var json_data []byte
	ctx := stubs.WithAccountDaoOne(context.Background())
	ctx = identity.WithTenant(t, ctx)
	ctx = stubs.WithPubkeyDao(ctx)

	values := map[string]interface{}{
		"account_id": 1,
		"name":       "very cool key",
		"body":       "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap",
	}

	if json_data, err = json.Marshal(values); err != nil {
		t.Fatal("unable to marshal values to json")
	}
	req, err := http.NewRequestWithContext(ctx, "POST", "/api/provisioning/pubkeys", bytes.NewBuffer(json_data))
	require.NoError(t, err, "failed to create request")
	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(services.CreatePubkey)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

	stubCount := stubs.PubkeyStubCount(ctx)
	assert.Equal(t, 1, stubCount, "Pubkey has not been Created through DAO")
}
