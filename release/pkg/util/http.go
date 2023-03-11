package util

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

func HttpDo(req *http.Request) ([]byte, error) {
	timeout := time.Duration(10 * time.Second) //超时时间50ms
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	arr, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil && err != io.EOF {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(`error http status code: ` + strconv.Itoa(resp.StatusCode) + ` resp message: ` + string(arr))
	}

	return arr, nil
}
