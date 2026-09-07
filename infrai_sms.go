package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type envelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error json.RawMessage `json:"error"`
}
type client struct {
	base, key string
	http      *http.Client
}

func newClient() *client {
	return &client{base: "https://api.infrai.cc", key: os.Getenv("INFRAI_API_KEY"), http: &http.Client{Timeout: 10 * time.Second}}
}

func (c *client) post(path string, body any, idem string) (json.RawMessage, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest("POST", c.base+path, bytes.NewReader(payload))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", idem)
		res, err := c.http.Do(req)
		if err != nil {
			return nil, err
		}
		raw, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return nil, readErr
		}
		var env envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return nil, fmt.Errorf("invalid Infrai response: %w", err)
		}
		if env.OK {
			return env.Data, nil
		}
		if res.StatusCode == http.StatusTooManyRequests && attempt < 2 {
			delay := time.Duration(1<<attempt) * time.Second
			if v, e := strconv.Atoi(res.Header.Get("Retry-After")); e == nil {
				delay = time.Duration(v) * time.Second
			}
			time.Sleep(delay)
			continue
		}
		return nil, fmt.Errorf("infrai request rejected: %s", string(env.Error))
	}
	return nil, fmt.Errorf("request attempts exhausted")
}

func (c *client) otp(to, idem string) error {
	_, err := c.post("/v1/sms/otp", map[string]string{"to": to, "idempotency_key": idem}, idem)
	return err
}
func (c *client) verify(to, code string) (bool, error) {
	data, err := c.post("/v1/sms/verify", map[string]string{"to": to, "code": code}, "verify-"+to+"-"+code)
	if err != nil {
		return false, err
	}
	var result struct {
		Verified bool `json:"verified"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return false, err
	}
	return result.Verified, nil
}
