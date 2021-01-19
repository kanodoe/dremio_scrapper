package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"

	"github.com/kanodoe/dremio_scrapper/internal/client"
)

var (
	timeout = time.Duration(5) * time.Second
	logger  log.Logger
)

func TestDremioClientTest(t *testing.T) {

	logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

	t.Run("Dremio Login Client ", func(t *testing.T) {

		response := `{
			"token": "6220073649bbeb2f374716f273",
			"userName": "mcfadden.dotson",
			"firstName": "Mcfadden",
			"lastName": "Dotson",
			"expires": 1510442336986,
			"email": "Mcfadden.Dotson@walmart.com",
			"userId": "d13930ce-e883-4450-9819-5bb6ad78899b",
			"admin": true,
			"clusterId": "715b05de-41ee-4bae-ac0a-ea7d4ffce86e",
			"clusterCreatedAt": 1571051858986,
			"showUserAndUserProperties": true,
			"version": "4.2.2-202004211133290458-b550b6fa",
			"permissions": {
				"canUploadProfiles": true,
				"canDownloadProfiles": true,
				"canEmailForSupport": true,
				"canChatForSupport": false
			},
			"userCreatedAt": 1406252836794
		}`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))
		}))
		c := client.NewDremioClient(server.URL, "username", "password", timeout)
		c = client.NewLoggingMiddleware(logger)(c)

		resp, err := c.Login(context.Background())
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "Mcfadden", resp.FirstName)

	})

	t.Run("Client Response error Decode response", func(t *testing.T) {
		responseNOK := ``
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(responseNOK))
		}))
		c := client.NewDremioClient(server.URL, "username", "passwordMala", timeout)
		c = client.NewLoggingMiddleware(logger)(c)
		resp, err := c.Login(context.Background())
		assert.Nil(t, resp)
		assert.NotNil(t, err)
	})

}
