package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const baseURL = "http://0.0.0.0:8080/api/v1"

func startServerFromRoot(t *testing.T) *exec.Cmd {
	os.Chdir("../")
	mainPath := filepath.Join("cmd", "orchestrator", "main.go")

	cmd := exec.Command("go", "run", mainPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	err := cmd.Start()
	require.NoError(t, err)

	time.Sleep(13 * time.Second)
	return cmd
}

func stopServer(cmd *exec.Cmd) {
	_ = cmd.Process.Kill()
}

func TestIntegrationFlow(t *testing.T) {
	cmd := startServerFromRoot(t)
	defer stopServer(cmd)

	client := &http.Client{}

	registerPayload := map[string]string{
		"login":    "testuser",
		"password": "strongpassword",
	}
	registerBody, _ := json.Marshal(registerPayload)
	resp, err := client.Post(baseURL+"/register", "application/json", bytes.NewBuffer(registerBody))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	resp.Body.Close()

	resp, err = client.Post(baseURL+"/login", "application/json", bytes.NewBuffer(registerBody))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	var loginResp struct {
		Token string `json:"token"`
	}
	body, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &loginResp)
	resp.Body.Close()
	require.NotEmpty(t, loginResp.Token)

	exprPayload := map[string]string{
		"expression": "2 + 3 * 4",
	}
	exprBody, _ := json.Marshal(exprPayload)

	req, _ := http.NewRequest("POST", baseURL+"/calculate", bytes.NewBuffer(exprBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)

	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	var exprResp struct {
		Id string `json:"id"`
	}
	body, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &exprResp)
	resp.Body.Close()
	require.NotEmpty(t, exprResp.Id)

	req, _ = http.NewRequest("GET", baseURL+"/expressions", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)

	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	var listResp struct {
		List []struct {
			ID     string `json:"ID"`
			Input  string `json:"Input"`
			Status string `json:"Status"`
			Result string `json:"Result"`
		} `json:"list"`
	}
	body, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &listResp)
	resp.Body.Close()

	require.NotEmpty(t, listResp.List)
	found := false
	for _, item := range listResp.List {
		if item.ID == exprResp.Id {
			found = true
			break
		}
	}
	require.True(t, found, "expression should be in list")
}
