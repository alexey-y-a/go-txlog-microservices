package apigateway_test

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/alexey-y-a/go-txlog-microservices/libs/logger"
	apiserver "github.com/alexey-y-a/go-txlog-microservices/services/api-gateway/internal/server"
)

type apiSetResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type apiGetResponse struct {
	Status string `json:"status"`
	Value  string `json:"value"`
}

type apiCommonResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func TestE2E_ApiGatewayAndKVService(t *testing.T) {
	t.Helper()

	logger.Init()

	cmdKv := exec.Command("go", "run", "./services/kv-service/cmd/kv")
	cmdKv.Stdout = os.Stdout
	cmdKv.Stderr = os.Stderr

	err := cmdKv.Start()
	require.NoError(t, err, "kv-service process should start")

	defer func() {
		_ = cmdKv.Process.Kill()
	}()

	time.Sleep(300 * time.Millisecond)

	apiSrv := apiserver.NewServer("http://localhost:8081")

	go func() {
		err = apiSrv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
		}
	}()

	time.Sleep(300 * time.Millisecond)

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	setBody := `{"key":"user42","value":"Alice"}`
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/set", strings.NewReader(setBody))
	require.NoError(t, err, "should create set request")

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err, "set request should not error")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "set should return 200")

	var setResp apiSetResponse
	err = json.NewDecoder(resp.Body).Decode(&setResp)
	require.NoError(t, err, "set response should be valid JSON")
	require.Equal(t, "ok", setResp.Status, "set response status should be ok")

	getReq, err := http.NewRequest(http.MethodGet, "http://localhost:8080/api/get?key=user42", nil)
	require.NoError(t, err, "should create get request")

	getResp, err := client.Do(getReq)
	require.NoError(t, err, "get request should not error")
	defer getResp.Body.Close()

	require.Equal(t, http.StatusOK, getResp.StatusCode, "get should return 200")

	var getBody apiGetResponse
	err = json.NewDecoder(getResp.Body).Decode(&getBody)
	require.NoError(t, err, "get response should be valid JSON")
	require.Equal(t, "ok", getBody.Status, "get response status should be ok")
	require.Equal(t, "Alice", getBody.Value, "get should return correct value")

	deleteReq, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/api/delete?key=user42", nil)
	require.NoError(t, err, "should create delete request")

	deleteResp, err := client.Do(deleteReq)
	require.NoError(t, err, "delete request should not error")
	defer deleteResp.Body.Close()

	require.Equal(t, http.StatusOK, deleteResp.StatusCode, "delete should return 200")

	var delBody apiCommonResponse
	err = json.NewDecoder(deleteResp.Body).Decode(&delBody)
	require.NoError(t, err, "delete response should be valid JSON")
	require.Equal(t, "ok", delBody.Status, "delete response status should be ok")

	getAfterDeleteReq, err := http.NewRequest(http.MethodGet, "http://localhost:8080/api/get?key=user42", nil)
	require.NoError(t, err, "should create get-after-delete request")

	getAfterDeleteResp, err := client.Do(getAfterDeleteReq)
	require.NoError(t, err, "get-after-delete request should not error")
	defer getAfterDeleteResp.Body.Close()

	require.Equal(t, http.StatusNotFound, getAfterDeleteResp.StatusCode, "get after delete should return 404")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = apiSrv.Shutdown(ctx)
	require.NoError(t, err, "api-gateway shutdown should not error")
}
