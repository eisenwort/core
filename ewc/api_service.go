package ewc

import (
	"bytes"
	"log"
	"net/http"
	"net/url"
	"time"
)

type ApiService struct {
	RequestErrorChan chan bool
}

func NewApiService() *ApiService {
	srv := new(ApiService)
	srv.RequestErrorChan = make(chan bool, chanSize)

	return srv
}

func (srv *ApiService) get(url string, closure func(r *http.Response)) {
	request := ApiRequest{
		Body:       nil,
		Method:     http.MethodGet,
		RequestUrl: baseUrl + url,
	}
	srv.createRequest(request, closure)
}

func (srv *ApiService) post(url string, data []byte, closure func(r *http.Response)) {
	request := ApiRequest{
		Body:       bytes.NewReader(data),
		Method:     http.MethodPost,
		RequestUrl: baseUrl + url,
	}
	srv.createRequest(request, closure)
}

func (srv *ApiService) put(url string, data []byte, closure func(r *http.Response)) {
	request := ApiRequest{
		Body:       bytes.NewReader(data),
		Method:     http.MethodPut,
		RequestUrl: baseUrl + url,
	}
	srv.createRequest(request, closure)
}

func (srv *ApiService) delete(url string, closure func(r *http.Response)) {
	request := ApiRequest{
		Body:       nil,
		Method:     http.MethodDelete,
		RequestUrl: baseUrl + url,
	}
	srv.createRequest(request, closure)
}

func (srv *ApiService) createRequest(data ApiRequest, closure func(r *http.Response)) {
	client := &http.Client{}
	client.Timeout = time.Second * 10
	request, err := http.NewRequest(data.Method, data.RequestUrl, data.Body)

	if err != nil {
		log.Printf("create %s request on %s error: %s", data.Method, data.RequestUrl, err.Error())
		return
	}

	request.Header.Set("X-API", "true")
	request.Header.Set("X-Auth-Id", userIDHeader)
	request.Header.Set("X-Auth-Token", jwt)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	if err != nil {
		log.Printf("send %s request on %s error: %s", data.Method, data.RequestUrl, err.Error())
		return
	}

	closure(response)
	_ = response.Body.Close()
}

func (srv *ApiService) createUrl(hostUrl string, params map[string]string) string {
	if len(params) == 0 {
		return hostUrl
	}

	values := url.Values{}

	for key, value := range params {
		values.Add(key, value)
	}

	return hostUrl + "?" + values.Encode()
}
