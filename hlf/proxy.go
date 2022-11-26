package hlf

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	ProxyBaseUrl   = "http://localhost:9001"
	ProxyInvokeUrl = ProxyBaseUrl + "/invoke"
	ProxyQueryUrl  = ProxyBaseUrl + "/query"
)

func newPostRequest(url string, data []byte) (*http.Request, error) {

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Basic test")

	return req, nil
}

func sendInvokeRequest(data []byte) error {
	return sendHttpRequest(ProxyInvokeUrl, data)
}

func sendQueryRequest(data []byte) error {
	return sendHttpRequest(ProxyQueryUrl, data)
}

func sendHttpRequest(url string, data []byte) error {

	req, err := newPostRequest(url, data)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		log.Println(string(body))
		if err = resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return errors.New("invocation failed")
	}

	return nil
}
