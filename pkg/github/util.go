package github

import (
	"net/http"
	"strings"

	"github.com/guguducken/auto-release/pkg/util"
)

func basicSet(req *http.Request, token string) {
	req.Header.Set(`Accept`, `*/*`)
	req.Header.Set(`Authorization`, `Bearer `+token)
	req.Header.Set(`X-GitHub-Api-Version`, `2022-11-28`)
}

func get(url string, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	basicSet(req, token)
	return util.HttpDo(req)
}

func post(url string, token string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(`POST`, url, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	basicSet(req, token)
	return util.HttpDo(req)
}

func delete(url string, token string) ([]byte, error) {
	req, err := http.NewRequest(`DELETE`, url, nil)
	if err != nil {
		return nil, err
	}

	basicSet(req, token)
	return util.HttpDo(req)
}

func put(url, token, body string) ([]byte, error) {
	req, err := http.NewRequest(`PUT`, url, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	basicSet(req, token)
	return util.HttpDo(req)
}
