package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)


type KVClient struct {
    baseURL string
    client *http.Client
}

func NewKVClient(baseURL string, timeout time.Duration) *KVClient {
    return &KVClient{
        baseURL: baseURL,
        client: &http.Client {
            Timeout: timeout,
        },
    }
}

type setRequest struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

type commonResponse struct {
    Status string `json:"status"`
    Message string `json:message,omitempty"`
}

func (c *KVClient) Set(key, value string) error {
    requestBody := setRequest {
        Key: key,
        Value: value,
    }

    bodyBytes, err := json.Marshal(requestBody)
    if err != nil {
        return fmt.Errorf("kvclient: marshal set request: %w", err)
    }

    url := c.baseURL + "kv/set"

    req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
    if err != nil {
        return fmt.Errorf("kvclient: new POST request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")

    resp, err := c.client.Do(req)
    if err != nil {
        return fmt.Errorf("kvclient: do POST request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("kvclient: set failed with status %d", resp.StatusCode)
    }

    return nil
}

type getResponse struct {
    Status string `json:"status"`
    Value string `json:"value,omitempty"`
}

func (c *KVClient) Get(key string) (string, bool, error) {
    url := c.baseURL + "kv/get?key=" + key

    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        return "", false, fmt.Errorf("kvclient: new GET request: %w", err)
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return "", false, fmt.Errorf("kvclient: do GET request: %w, err")
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusNotFound {
        return "", false, nil
    }

    if resp.StatusCode != http.StatusOK {
        return "", false, fmt.Errorf("kvclient: get failed with status %d", resp.StatusCode)
    }

    var response getResponse
    decoder := json.NewDecoder(resp.Body)
    err = decoder.Decode(&response)
    if err != nil {
        return "", false, fmt.Errorf("kvclient: decode get response: %w", err)
    }

    return response.Value, true, nil
}

func (c *KVClient) Delete(key string) error {
    url := c.baseURL + "kv/delete?key=" + key

    req, err := http.NewRequest(http.MethodDelete, url, nil)
    if err != nil {
        return fmt.Errorf("kvclient: new DELETE request: %w", err)
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return fmt.Errorf("kvclient: do DELETE request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("kvclient: delete failed with status %d", resp.StatusCode)
    }

    return nil
}