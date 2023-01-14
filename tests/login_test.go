package tests

import (
	"chi-users-project/app/services/dtos"
	"chi-users-project/tests/testhelper"
	"encoding/json"
	"testing"
	"strings"
)


func TestValidLogin(t *testing.T) {
	// first we need to initialize the test db or something, and add it to r as context...
	// would be cool if we could get this into some kind of before all hook
	testHelper := testhelper.InitTestDBAndService(t)

	// now init data needed in the DB
	testHelper.CreateTestAuthUsers()
	defer testHelper.CleanupAuth()

	reqData, jsonErr := json.Marshal(
		map[string]interface{}{
		"email": "regular@test.com",
		"password": "testpassword9",
	})
	if jsonErr != nil {
		t.Errorf("expected error to be nil got %v", jsonErr)
	}

	res := testHelper.SendAsNoOne("post", "/auth/login", reqData)
	
	data := testHelper.GetResponseBody(res)
	testHelper.AssertStatus(res, 200)

	var tokenDTO dtos.LoginTokenDTO
	jsonErr = json.Unmarshal(data, &tokenDTO)
	if jsonErr != nil {
		t.Errorf("expected error to be nil got %v", jsonErr)
	}

	splitJWT := strings.Split(tokenDTO.Token, ".")
	testHelper.Assert(len(splitJWT) == 3, "")
}


func TestInvalidLogin(t *testing.T) {
	testHelper := testhelper.InitTestDBAndService(t)
	testHelper.CreateTestAuthUsers()
	defer testHelper.CleanupAuth()
	
	regularData, jsonErr := json.Marshal(
		map[string]interface{}{
		"email": "fakeemail1@test.com",
		"password": "password",
	})
	if jsonErr != nil {
		t.Errorf("expected error to be nil got %v", jsonErr)
	}

	res := testHelper.SendAsNoOne("post", "/auth/login", regularData)

	testHelper.AssertStatus(res, 401)
	testHelper.ValidateErrDTOPresent(res)
}
