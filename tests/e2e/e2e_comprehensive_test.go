package e2e

import (
	"net/http"
	"testing"

	"free2free/tests/testutils"
	"github.com/stretchr/testify/assert"
)

func TestFullApplicationE2E(t *testing.T) {
	testServer := testutils.NewTestServer()
	defer testServer.Close()

	t.Run("ApplicationStartup", func(t *testing.T) {
		resp, err := testServer.DoRequest("GET", "/swagger/index.html", nil, nil)
		assert.NoError(t, err)
		if resp.StatusCode == http.StatusNotFound {
			t.Skip("Swagger endpoint not configured, skipping")
			return
		}
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	})

	t.Run("AuthenticationFlow", func(t *testing.T) {
		user, err := testServer.CreateTestUser()
		assert.NoError(t, err)

		token, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		_, err = testutils.ValidateJWTToken(token)
		assert.NoError(t, err)

		resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", "/profile", token, nil)
		assert.NoError(t, err)
		assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
		resp.Body.Close()
	})

	t.Run("AdminFeatures", func(t *testing.T) {
		if testServer.DB == nil {
			t.Skip("Database not available, skipping admin features test")
			return
		}

		adminUser, err := testServer.CreateTestUser()
		assert.NoError(t, err)

		testServer.DB.Model(&adminUser).Update("is_admin", true)

		adminToken, err := testutils.CreateMockJWTToken(adminUser.ID, adminUser.Name, true)
		assert.NoError(t, err)

		resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", "/admin/activities", adminToken, nil)
		assert.NoError(t, err)
		assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
		resp.Body.Close()
	})

	t.Run("UserFeatures", func(t *testing.T) {
		user, err := testServer.CreateTestUser()
		assert.NoError(t, err)

		token, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
		assert.NoError(t, err)

		resp, err := testutils.MakeAuthenticatedRequest(testServer, "GET", "/user/matches", token, nil)
		assert.NoError(t, err)
		assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
		resp.Body.Close()
	})

	t.Run("OrganizerFeatures", func(t *testing.T) {
		user, err := testServer.CreateTestUser()
		assert.NoError(t, err)

		token, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
		assert.NoError(t, err)

		resp, err := testutils.MakeAuthenticatedRequest(testServer, "PUT", "/organizer/approve-participant/1", token, nil)
		assert.NoError(t, err)
		assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
		resp.Body.Close()
	})

	t.Run("ReviewFeatures", func(t *testing.T) {
		user, err := testServer.CreateTestUser()
		assert.NoError(t, err)

		token, err := testutils.CreateMockJWTToken(user.ID, user.Name, user.IsAdmin)
		assert.NoError(t, err)

		reviewData := map[string]interface{}{
			"reviewee_id": 2,
			"score":       5,
			"comment":     "Great experience!",
		}

		resp, err := testutils.MakeAuthenticatedRequest(testServer, "POST", "/review/matches/1", token, reviewData)
		assert.NoError(t, err)
		assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
		resp.Body.Close()
	})
}
