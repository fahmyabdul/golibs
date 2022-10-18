package golibs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

var encoder = schema.NewEncoder()

// ResponseJSON :
type ResponseJSON struct {
	Success *bool       `json:"success,omitempty"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Count   *int        `json:"count,omitempty"`
	Data    interface{} `json:"data"`
}

// BoolPointer :
func BoolPointer(input bool) *bool {
	return &input
}

// NumberPointer :
func NumberPointer(input int) *int {
	return &input
}

// RespondWithJSON : Gives JSON Response on Request
func RespondWithJSON(w http.ResponseWriter, code int, message string, payload interface{}, count int) {
	responseStruct := ResponseJSON{
		Code:    code,
		Message: message,
		Data:    payload,
		Count:   &count,
	}
	response, _ := json.Marshal(responseStruct)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Respond : New JSON Response
func (p *ResponseJSON) Respond(w http.ResponseWriter) {
	response, _ := json.Marshal(p)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(p.Code)
	w.Write(response)
}

// RespondNoStruct : Gives JSON Response Without Using ResponseJSON struct
func RespondNoStruct(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// GetRequest : Make Get Request to Service/API
func GetRequest(endpoint string) (*http.Response, []byte, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}
	io.Copy(ioutil.Discard, resp.Body)

	return resp, body, nil
}

// PostRequest :
func PostRequest(endpoint string, payload interface{}) (*http.Response, []byte, error) {
	jsonValue, _ := json.Marshal(payload)
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}
	io.Copy(ioutil.Discard, resp.Body)

	return resp, body, nil
}

// SendRequestWithHeader :
func SendRequestWithHeader(method string, endpoint string, header map[string]string, payload interface{}) (*http.Response, []byte, error) {
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}

	httpRequest, err := http.NewRequest(method, endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, nil, err
	}

	// Set Request Headers
	for key, value := range header {
		httpRequest.Header.Set(key, value)
	}

	client := &http.Client{}
	response, err := client.Do(httpRequest)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}

	return response, responseBody, nil
}

// PostFormRequest :
func PostFormRequest(endpoint string, payload interface{}) (*http.Response, []byte, error) {
	form := url.Values{}

	err := encoder.Encode(payload, form)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return nil, nil, err
	}

	resp, err := http.PostForm(endpoint, form)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return resp, nil, err
	}
	io.Copy(ioutil.Discard, resp.Body)

	return resp, body, nil
}

// PostMultipartFormData :
func PostMultipartFormData(endpoint string, payload interface{}) (*http.Response, []byte, error) {
	form := url.Values{}

	err := encoder.Encode(payload, form)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return nil, nil, err
	}

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	for k, v := range form {
		for _, value := range v {
			bodyWriter.WriteField(k, value)
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(endpoint, contentType, bodyBuf)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}

	return resp, body, nil
}
