package qywx

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

func HttpRequest(urlStr string, method string, headers map[string]string,
	params map[string]string, data any) (*http.Response, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	u.RawQuery = query.Encode()
	buf := new(bytes.Buffer)
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
