package ewc

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

var httpClient = http.Client{
	Timeout: time.Second * 10,
}

func httpGet(url string, closure func(r *http.Response)) {
	request := ApiRequest{
		Body:       nil,
		Method:     http.MethodGet,
		RequestUrl: url,
	}
	createRequest(request, closure)
}

func httpPost(url string, data interface{}, closure func(r *http.Response)) {
	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Println("marshall POST data error:", err)
		return
	}

	request := ApiRequest{
		Body:       jsonData,
		Method:     http.MethodPost,
		RequestUrl: url,
	}
	createRequest(request, closure)
}

func httpPut(url string, data interface{}, closure func(r *http.Response)) {
	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Println("marshall PUT data error:", err)
		return
	}

	request := ApiRequest{
		Body:       jsonData,
		Method:     http.MethodPut,
		RequestUrl: url,
	}
	createRequest(request, closure)
}

func httpDelete(url string, closure func(r *http.Response)) {
	request := ApiRequest{
		Body:       nil,
		Method:     http.MethodDelete,
		RequestUrl: url,
	}
	createRequest(request, closure)
}

func createRequest(data ApiRequest, closure func(r *http.Response)) {
	if currentUser.Reseted {
		return
	}
	request, err := http.NewRequest(
		data.Method,
		baseUrl+data.RequestUrl,
		bytes.NewReader(data.Body),
	)

	if err != nil {
		log.Printf("create %s request on %s error: %s", data.Method, data.RequestUrl, err.Error())
		return
	}

	request.Header.Set("X-API", "true")
	request.Header.Set(IdHeader, userIDHeader)
	request.Header.Set(TokenHeader, jwtToken)
	request.Header.Set("Content-Type", "application/json")

	response, err := httpClient.Do(request)

	if err != nil {
		log.Printf("send %s request on %s error: %s", data.Method, data.RequestUrl, err.Error())
		return
	}

	closure(response)

	if err := response.Body.Close(); err != nil {
		log.Println("close response error:", err)
	}
}

func createUrl(hostUrl string, params map[string]string) string {
	if len(params) == 0 {
		return hostUrl
	}

	values := url.Values{}

	for key, value := range params {
		values.Add(key, value)
	}

	return hostUrl + "?" + values.Encode()
}
