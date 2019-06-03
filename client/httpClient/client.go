package httpClient

import (
	"io"
	"net/http"
	"time"
)

type HttpClient struct {
	*http.Client
	addr     string
	protocol string
}

func NewHttpClient(addr, protocol string, maxIdleConns int, idleConnTimeoutSec int, disableCompress bool) *HttpClient {
	var tr = http.Transport{
		MaxIdleConns:       maxIdleConns,
		IdleConnTimeout:    time.Duration(idleConnTimeoutSec) * time.Second,
		DisableCompression: disableCompress,
	}

	return &HttpClient{
		&http.Client{
			Transport: &tr,
		},
		addr,
		protocol,
	}
}

func (cli *HttpClient) Do(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, cli.protocol+"://"+cli.addr+url, body)
	if err != nil {
		return nil, err
	}

	resp, err := cli.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
