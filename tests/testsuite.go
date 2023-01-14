package tests

// import (
// 	// "io/ioutil"
// 	"net/http"
// )

// type TestRequest struct {
// 	*http.Request
// 	testSuite				*TestSuite
// }

// type TestSuite struct {
// 	Response      *http.Response
// 	ResponseBody  []byte

// 	baseURL				string
// }


// func(this *TestSuite) GetTestRequest(req *http.Request) {
// 	return &TestRequest {
// 		Request: req,
// 		testSuite: this,
// 	}
// }


// func(this *TestRequest) SendRequest() {
// 	var err error
// 	this.testSuite.Response, err = testSuite.Client.Do(r.Request)
// 	if err != nil {
// 		panic(err)
// 	}

// 	r.testSuite.ResponseBody, err = ioutil.ReadAll(this.testSuite.Response.Body)
// 	if err != nil {
// 		panic(err)
// 	}
// }
