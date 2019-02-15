package utils

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Result response result
type Result struct {
	status     string
	statusCode int
	data       []byte
	err        error
}

// Err response error
func (r *Result) Err() error {
	return r.err
}

// Status response status
func (r *Result) Status() (int, string) {
	return r.statusCode, r.status
}

// Bytes response body
func (r *Result) Bytes() ([]byte, error) {
	return r.data, r.err
}

// Reader response body reader
func (r *Result) Reader() (io.Reader, error) {
	return bytes.NewReader(r.data), r.err
}

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

var (
	httpClient *http.Client
	once       sync.Once
)

// SendRequest send request
func SendRequest(ctx context.Context, url string, method string, payload string, headers ...string) *Result {
	req, err := http.NewRequest(method, url, bytes.NewBufferString(payload))
	if err != nil {
		return &Result{err: err}
	}
	req = req.WithContext(ctx)
	SetHeaders(req, headers...)
	return DoRequest(req)
}

// DoRequest do request
func DoRequest(request *http.Request) *Result {
	once.Do(func() {
		httpClient = &http.Client{}

		tr := &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				dialer := net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 90 * time.Second,
				}
				return dialer.Dial(network, addr)
			},
			ResponseHeaderTimeout: time.Second * 10,
		}

		httpClient.Transport = tr
	})

	resp, err := httpClient.Do(request)
	if err != nil {
		return &Result{err: err}
	}

	result := &Result{}
	result.statusCode, result.status = resp.StatusCode, resp.Status
	result.data, result.err = ReadResponse(resp)
	return result
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
