package core

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Engine struct {
	status string
}

func (e Engine) Get(url string, params map[string]string, headers map[string]string, ssl bool) (response string, err error) {
	if (&url == nil) || (url == "") {
		fmt.Printf("Need a valid URL\n")
	}

	if !strings.HasPrefix(strings.ToLower(url), "http") {
		if ssl {
			url = fmt.Sprintf("%s://", "https") + url
		} else {
			url = fmt.Sprintf("%s://", "http") + url
		}
	}

	fullUrl := fmt.Sprintf("%s?%s", url, joinParams(params))

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		panic(err.Error())
	}

	for header_key, header_value := range headers {
		req.Header.Add(header_key, header_value)
	}

	response, err = sendRequest(req)
	return response, err
}

func sendRequest(req *http.Request) (response string, err error) {
	client := http.Client{}

	resp, err := client.Do(req)
    if err != nil {
        response = ""
    } else {
        defer resp.Body.Close()
        body, readErr := ioutil.ReadAll(resp.Body)
        if readErr != nil {
            err = readErr
        } else {
            response = string(body)
        }
    }

	return response, err
}

func joinParams(params map[string]string) (joinedParams string) {

	for k, v := range params {
		joinedParams += fmt.Sprintf("%s=%s&", k, v)
	}
	return joinedParams
}
