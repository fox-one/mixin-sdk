package utils

import (
	"bytes"
	"compress/gzip"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// BuildURL build url
func BuildURL(path string, paras ...string) (string, error) {
	u, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for i := 0; i < len(paras)-1; i += 2 {
		q.Set(paras[i], paras[i+1])
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

// SetHeaders set http headers
func SetHeaders(r *http.Request, headers ...string) {
	// r.Header.Add("Connection", "Keep-Alive")
	for _i := 0; _i < len(headers)-1; _i += 2 {
		r.Header.Add(headers[_i], headers[_i+1])
	}
}

// NewRequest new request
func NewRequest(url, method, payload string, headers ...string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBufferString(payload))
	if err == nil {
		SetHeaders(req, headers...)
	}

	return req, err
}

// GetRequest new get request
func GetRequest(url string, headers ...string) (*http.Response, error) {
	return SendRequest(url, "GET", "", headers...)
}

// PostRequest new post request
func PostRequest(url string, payload string, headers ...string) (*http.Response, error) {
	return SendRequest(url, "POST", payload, headers...)
}

// PutRequest new put request
func PutRequest(url string, payload string, headers ...string) (*http.Response, error) {
	return SendRequest(url, "PUT", payload, headers...)
}

// DeleteRequest new delete request
func DeleteRequest(url string, payload string, headers ...string) (*http.Response, error) {
	return SendRequest(url, "DELETE", payload, headers...)
}

var (
	httpClient *http.Client
	once       sync.Once
)

// SendRequest send request
func SendRequest(url string, method string, payload string, headers ...string) (*http.Response, error) {
	req, _ := http.NewRequest(method, url, bytes.NewBufferString(payload))
	SetHeaders(req, headers...)
	return DoRequest(req)
}

// DoRequest do request
func DoRequest(request *http.Request) (*http.Response, error) {
	once.Do(func() {
		httpClient = &http.Client{}

		tr := &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				dialer := net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}
				return dialer.Dial(network, addr)
			},
		}

		httpClient.Transport = tr
	})

	resp, err := httpClient.Do(request)
	if err == nil && resp != nil && resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
	}

	return resp, err
}

// ReadResponse read response
func ReadResponse(resp *http.Response) ([]byte, error) {
	if resp == nil || resp.Body == nil {
		return nil, nil
	}

	defer resp.Body.Close()
	body := &bytes.Buffer{}

	var reader = resp.Body
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
	}

	if _, err := body.ReadFrom(reader); err != nil {
		return nil, err
	}

	return body.Bytes(), nil
}
