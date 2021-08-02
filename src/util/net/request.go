package net

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func Get(c *http.Client, url string, header map[string]string) (string, error) {
	return Do(c, "GET", url, header, bytes.NewBuffer([]byte("")))
}

func POST(c *http.Client, url string, header map[string]string, body io.Reader) (string, error) {
	return Do(c, "POST", url, header, body)
}

func Do(c *http.Client, method, url string, header map[string]string, body io.Reader) (string, error) {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", err
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	defer req.Body.Close()

	c.Timeout = 5 * time.Second

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		data, er := ioutil.ReadAll(resp.Body)
		if er == nil && len(data) > 0 {
			return "", errors.New(string(data))
		}

		return "", er
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBytes), nil
}
