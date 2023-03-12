package util

import (
	"crypto/tls"
	"io"
	"net/http"
	"strconv"
	"strings"
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
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		Error(`http status code is invalid: ` + strconv.Itoa(resp.StatusCode) + ` resp message: ` + string(arr))
		return nil, StatusError
	}
	return arr, nil
}

func URLValid(path string) string {
	var build strings.Builder
	for i := 0; i < len(path); i++ {
		switch path[i] {
		case '+':
			build.WriteString(`%2B`)

		case ' ':
			build.WriteString(`%20`)

		case '/':
			build.WriteString("%2F")

		case '?':
			build.WriteString(`%3F`)

		case '%':
			build.WriteString(`%25`)

		case '#':
			build.WriteString(`%23`)

		case '&':
			build.WriteString(`%26`)

		case '=':
			build.WriteString(`%3D`)

		default:
			build.WriteByte(path[i])
		}
	}
	return build.String()
}
