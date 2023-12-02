package testhelper

import (
	"youtube-downloader/app/services/dtos"
	"net/http/httptest"
	"encoding/json"
	"path/filepath"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"bytes"
	"fmt"
)


// ----- helpers/datahandlers for the tests to use -----


func(this TestHelper) GetResponseBody(res *http.Response) []byte {
	data, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		this.t.Errorf("expected error to be nil got %v", err)
	}

	return data
}


// Assertions

func(this TestHelper) Assert(value bool, msg string) {
	this.assert_it(value, msg, 2)
}


func(this TestHelper) AssertStatus(res *http.Response, expected int) {
	this.assert_it(res.StatusCode == expected, fmt.Sprintf("Expected status: %d; actual status: %d", expected, res.StatusCode), 2)
}


func(this TestHelper) AssertErrDTOPresent(res *http.Response) {
	errDTO := dtos.ErrorDTO{}
	jsonErr := json.Unmarshal(this.GetResponseBody(res), &errDTO)
	if jsonErr != nil {
		this.t.Fatal(jsonErr.Error())
	}

	this.assert_it(errDTO.Exists(), "ErrorDTO response expected", 2)
}


func(this TestHelper) assert_it(value bool, msg string, skip int) {
	if skip == 0 {
		skip = 1
	}
	if !value {
		_, file, line, _ := runtime.Caller(skip)
		this.t.Errorf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line})...)
	}
}


// ----- send data -----

func(this TestHelper) SendAsNoOne(typ string, urlEnd string, body []byte) *http.Response {
	req := this.getRequest(typ, urlEnd, body)
	w := httptest.NewRecorder()

	this.server.Router.ServeHTTP(w, req)

	return w.Result()
}


func(this TestHelper) SendAsRegularUser(typ string, urlEnd string, body []byte) *http.Response {
	if this.RegularToken == "" {
		this.t.Fatal("Must call InitAuth in Before hook to send as Regular User")
	}

	req := this.getRequest(typ, urlEnd, body)
	w := httptest.NewRecorder()

	req.Header.Set("Authorization", this.RegularToken)

	this.server.Router.ServeHTTP(w, req)

	return w.Result()
}


func(this TestHelper) SendAsAdmin(typ string, urlEnd string, body []byte) *http.Response {
	if this.AdminToken == "" {
		this.t.Fatal("Must call InitAuth in Before hook to send as Admin")
	}

	req := this.getRequest(typ, urlEnd, body)
	w := httptest.NewRecorder()

	req.Header.Set("Authorization", this.AdminToken)

	this.server.Router.ServeHTTP(w, req)

	return w.Result()
}


func(this TestHelper) SendAsSuperAdmin(typ string, urlEnd string, body []byte) *http.Response {
	if this.SuperAdminToken == "" {
		this.t.Fatal("Must call InitAuth in Before hook to send as Super Admin")
	}

	req := this.getRequest(typ, urlEnd, body)
	w := httptest.NewRecorder()

	req.Header.Set("Authorization", this.SuperAdminToken)

	this.server.Router.ServeHTTP(w, req)

	return w.Result()
}


func(this TestHelper) getRequest(typ string, url string, body []byte) (*http.Request) {
	var req *http.Request
	var err error

	if body == nil {
		req, err = http.NewRequest(strings.ToUpper(typ), url, nil)
	} else {
		buf := bytes.NewBuffer(body)
		req, err = http.NewRequest(strings.ToUpper(typ), url, buf)
	}
	if err != nil {
		this.t.Errorf("request error: %s", err.Error())
	}

	return req
}
