package ewc

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

type ApiService struct {
	client http.Client
}

func NewApiService() *ApiService {
	srv := new(ApiService)
	srv.client = http.Client{
		Timeout: time.Second * 10,
	}

	return srv
}

func (srv *ApiService) get(url string, closure func(r *http.Response)) {
	request := ApiRequest{
		Body:       nil,
		Method:     http.MethodGet,
		RequestUrl: url,
	}
	srv.createRequest(request, closure)
}

func (srv *ApiService) post(url string, data interface{}, closure func(r *http.Response)) {
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
	srv.createRequest(request, closure)
}

func (srv *ApiService) put(url string, data interface{}, closure func(r *http.Response)) {
	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Println("marshall POST data error:", err)
		return
	}

	request := ApiRequest{
		Body:       jsonData,
		Method:     http.MethodPut,
		RequestUrl: url,
	}
	srv.createRequest(request, closure)
}

func (srv *ApiService) delete(url string, closure func(r *http.Response)) {
	request := ApiRequest{
		Body:       nil,
		Method:     http.MethodDelete,
		RequestUrl: url,
	}
	srv.createRequest(request, closure)
}

func (srv *ApiService) createRequest(data ApiRequest, closure func(r *http.Response)) {
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

	response, err := srv.client.Do(request)

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
