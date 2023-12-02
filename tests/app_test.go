package tests

import (
	"youtube-downloader/tests/testhelper"
	"testing"
)


func TestThatIndexPageWorks2(t *testing.T) {
	testHelper := testhelper.InitTestDBAndService(t)
	defer testHelper.Cleanup()

	res := testHelper.SendAsNoOne("get", "/", nil)

	body := testHelper.GetResponseBody(res)

	testHelper.AssertStatus(res, 200)
	testHelper.Assert(string(body) == "Welcome to Teagan's chi app!", "")
}
