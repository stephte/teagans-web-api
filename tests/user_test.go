package tests

import (
	"youtube-downloader/tests/testhelper"
	"encoding/json"
	"testing"
	"fmt"
)


// ----- User Index -----

func TestThatUserIndexRequiresAuth(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	res := helper.SendAsNoOne("get", "/users", nil)

	helper.AssertStatus(res, 401)
}


func TestThatUserIndexRequiresAdminAuth(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	res := helper.SendAsRegularUser("get", "/users", nil)

	helper.AssertStatus(res, 401)
	helper.AssertErrDTOPresent(res)
}


func TestThatUserIndexAcceptsAdminAndSuperAdminAuth(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()
	
	res := helper.SendAsAdmin("get", "/users", nil)
	helper.AssertStatus(res, 200)
}


// ----- User find -----

func TestUserFindWorks(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	res := helper.SendAsRegularUser("get", fmt.Sprintf("/users/%s", helper.RegularUser.ID), nil)
	helper.AssertStatus(res, 200)
}


func TestUserFindDoesntWorkForOtherUser(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	res := helper.SendAsRegularUser("get", fmt.Sprintf("/users/%s", helper.AdminUser.ID), nil)

	helper.AssertStatus(res, 401)
	helper.AssertErrDTOPresent(res)
}


func TestUserFindWorksForAdminUsers(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	res := helper.SendAsSuperAdmin("get", fmt.Sprintf("/users/%s", helper.AdminUser.ID), nil)
	helper.AssertStatus(res, 200)

	res = helper.SendAsSuperAdmin("get", fmt.Sprintf("/users/%s", helper.RegularUser.ID), nil)
	helper.AssertStatus(res, 200)
}


// ----- User create -----

func TestUserCreate(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	data := map[string]interface{}{
		"firstName": "Testy",
		"lastName": "McTest",
		"email": "testymctest@test.com",
		"password": "testpassword7",
	}
	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	res := helper.SendAsNoOne("post", "/users", reqData)

	helper.AssertStatus(res, 200)

	body := helper.GetUserDTO(res)
	helper.Assert(body.FirstName == "Testy", "Firstname value incorrect")
	helper.Assert(body.LastName == "McTest", "Lastname value incorrect")
	helper.Assert(body.Email == "testymctest@test.com", "Email incorrect")
	helper.Assert(body.Role == 1, "Role incorrect")

	// get id from ResponseBody
	id := body.ID

	dres := helper.SendAsSuperAdmin("delete", fmt.Sprintf("/users/%s", id), nil)
	helper.AssertStatus(dres, 204)
}


func TestInvalidUserCreate(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	data := map[string]interface{}{
		"firstName": "Testy",
		"lastName": "McTest",
		"email": "testymctest2@test.com",
		"password": "testpassword12",
		"role": 2,
	}
	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	res := helper.SendAsNoOne("post", "/users", reqData)

	helper.AssertStatus(res, 400)
	helper.AssertErrDTOPresent(res)
}


func TestValidAdminUserCreate(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	data := map[string]interface{}{
		"firstName": "Testy",
		"lastName": "McTest",
		"email": "testymctest@test.com",
		"password": "testpassword7",
		"role": 2,
	}
	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	res := helper.SendAsSuperAdmin("post", "/users", reqData)

	helper.AssertStatus(res, 200)

	body := helper.GetUserDTO(res)
	helper.Assert(body.FirstName == "Testy", "Firstname value incorrect")
	helper.Assert(body.LastName == "McTest", "Lastname value incorrect")
	helper.Assert(body.Email == "testymctest@test.com", "Email incorrect")
	helper.Assert(body.Role == 2, "Role incorrect")

	// get id from ResponseBody
	id := body.ID

	dres := helper.SendAsSuperAdmin("delete", fmt.Sprintf("/users/%s", id), nil)
	helper.AssertStatus(dres, 204)
}


func  TestUserCreateInvalidEmail(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	data := map[string]interface{}{
		"firstName": "Testy",
		"lastName": "McTest",
		"email": "testymctest@test",
		"password": "testpassword7",
	}
	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	res := helper.SendAsNoOne("post", "/users", reqData)

	helper.AssertStatus(res, 400)
	helper.AssertErrDTOPresent(res)
}


func  TestUserCreateInvalidPassword(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	data := map[string]interface{}{
		"firstName": "Testy",
		"lastName": "McTest",
		"email": "testymctest@test",
		"password": "testpassword7",
	}
	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	res := helper.SendAsNoOne("post", "/users", reqData)

	helper.AssertStatus(res, 400)
	helper.AssertErrDTOPresent(res)
}


// ----- User update (PATCH) -----

func  TestUserUpdateWorksForUser(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	data := map[string]interface{}{
		"firstName": "Testie",
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	res := helper.SendAsRegularUser("patch", fmt.Sprintf("/users/%s", helper.RegularUser.ID), reqData)
	helper.AssertStatus(res, 200)

	body := helper.GetUserDTO(res)
	helper.Assert(body.FirstName == "Testie", "Firstname should be updated")
	helper.Assert(body.LastName == helper.RegularUser.LastName, "Lastname should be the same")
	helper.Assert(body.Email == helper.RegularUser.Email, "Email should stay the same")
	helper.Assert(body.Role == helper.RegularUser.Role, "Role should stay the same")

	helper.Assert(body.FirstName == "Testie", "Firstname should be updated")
}


func  TestUserUpdateFailsForOtherUser(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	data := map[string]interface{}{
		"firstName": "Testie",
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	res := helper.SendAsRegularUser("patch", fmt.Sprintf("/users/%s", helper.AdminUser.ID), reqData)

	helper.AssertStatus(res, 401)
	helper.AssertErrDTOPresent(res)
}


func  TestUserUpdateWorksForSuperAdmin(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	data := map[string]interface{}{
		"lastName": "Test",
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	res := helper.SendAsSuperAdmin("patch", fmt.Sprintf("/users/%s", helper.RegularUser.ID), reqData)
	helper.AssertStatus(res, 200)

	body := helper.GetUserDTO(res)
	helper.Assert(body.FirstName == helper.RegularUser.FirstName, "Firstname should stay the same")
	helper.Assert(body.LastName == "Test", "Lastname should be updated")
	helper.Assert(body.Email == helper.RegularUser.Email, "Email should stay the same")
	helper.Assert(body.Role == helper.RegularUser.Role, "Role should stay the same")
}


func  TestSuperAdminCanUpdateRole(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	data := map[string]interface{}{
		"email": "TEST@test.test",
		"role": 2,
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	res := helper.SendAsSuperAdmin("patch", fmt.Sprintf("/users/%s", helper.RegularUser.ID), reqData)
	helper.AssertStatus(res, 200)

	body := helper.GetUserDTO(res)
	helper.Assert(body.Email == "test@test.test", "Email should be updated and should be lowercase")
	helper.Assert(body.Role == 2, "Role should be updated")
	helper.Assert(body.FirstName == helper.RegularUser.FirstName, "Firstname should stay the same")
	helper.Assert(body.LastName == helper.RegularUser.LastName, "Lastname should stay the same")
}


func TestAdminCanUpdateRegularUserToAdmin(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	data := map[string]interface{}{
		"firstName": "Test",
		"role": 2,
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	res := helper.SendAsAdmin("patch", fmt.Sprintf("/users/%s", helper.RegularUser.ID), reqData)
	helper.AssertStatus(res, 200)

	body := helper.GetUserDTO(res)
	helper.Assert(body.Email == helper.RegularUser.Email, "Email should stay the same")
	helper.Assert(body.Role == 2, "Role should be updated")
	helper.Assert(body.FirstName == "Test", "Firstname should be updated")
	helper.Assert(body.LastName == helper.RegularUser.LastName, "Lastname should stay the same")
}


func  TestAdminCanNotUpdateRoleAboveThemselves(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	data := map[string]interface{}{
		"firstName": "Test",
		"role": 3,
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	res := helper.SendAsAdmin("patch", fmt.Sprintf("/users/%s", helper.RegularUser.ID), reqData)

	helper.AssertStatus(res, 401)
	helper.AssertErrDTOPresent(res)
}


func  TestAdminCanNotUpdateSuperAdmin(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	data := map[string]interface{}{
		"firstName": "Test",
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	res := helper.SendAsAdmin("patch", fmt.Sprintf("/users/%s", helper.SuperAdminUser.ID), reqData)

	helper.AssertStatus(res, 401)
	helper.AssertErrDTOPresent(res)
}


func  TestAdminCanNotRaiseOwnRole(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	data := map[string]interface{}{
		"role": 3,
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	res := helper.SendAsAdmin("patch", fmt.Sprintf("/users/%s", helper.AdminUser.ID), reqData)

	helper.AssertStatus(res, 401)
	helper.AssertErrDTOPresent(res)
}


func TestUserUpdateInvalidEmail(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	data := map[string]interface{}{
		"email": "test.com",
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	res := helper.SendAsRegularUser("patch", fmt.Sprintf("/users/%s", helper.RegularUser.ID), reqData)

	helper.AssertStatus(res, 400)
	helper.AssertErrDTOPresent(res)
}


// ----- User update OG (PUT) -----

func TestUserUpdateOGWorksForUser(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	data := map[string]interface{}{
		"firstName": "Testie",
		"lastName": helper.RegularUser.LastName,
		"email": helper.RegularUser.Email,
		"role": helper.RegularUser.Role,
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	res := helper.SendAsRegularUser("put", fmt.Sprintf("/users/%s", helper.RegularUser.ID), reqData)
	helper.AssertStatus(res, 200)

	body := helper.GetUserDTO(res)
	helper.Assert(body.FirstName == "Testie", "Firstname should be updated")
	helper.Assert(body.LastName == helper.RegularUser.LastName, "Lastname should be the same")
	helper.Assert(body.Email == helper.RegularUser.Email, "Email should be the same")
	helper.Assert(body.Role == helper.RegularUser.Role, "Role should be the same")
}


func TestUserUpdateOGFailsForOtherUser(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	data := map[string]interface{}{
		"firstName": "Testie",
		"lastName": helper.AdminUser.LastName,
		"email": helper.AdminUser.Email,
		"role": helper.AdminUser.Role,
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	res := helper.SendAsRegularUser("put", fmt.Sprintf("/users/%s", helper.AdminUser.ID), reqData)

	helper.AssertStatus(res, 401)
	helper.AssertErrDTOPresent(res)
}


func TestUserUpdateOGWorksForSuperAdmin(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	data := map[string]interface{}{
		"firstName": helper.RegularUser.FirstName,
		"lastName": "Test",
		"email": helper.RegularUser.Email,
		"role": helper.RegularUser.Role,
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	res := helper.SendAsSuperAdmin("put", fmt.Sprintf("/users/%s", helper.RegularUser.ID), reqData)
	helper.AssertStatus(res, 200)

	body := helper.GetUserDTO(res)
	helper.Assert(body.FirstName == helper.RegularUser.FirstName, "Firstname should be the same")
	helper.Assert(body.LastName == "Test", "Lastname should be updated")
	helper.Assert(body.Email == helper.RegularUser.Email, "Lastname should be the same")
	helper.Assert(body.Role == helper.RegularUser.Role, "Lastname should be the same")
}


func TestSuperAdminCanUpdateRoleWithOG(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	data := map[string]interface{}{
		"firstName": helper.RegularUser.FirstName,
		"lastName": helper.RegularUser.LastName,
		"email": "testing@mail.test",
		"role": 2,
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	res := helper.SendAsSuperAdmin("put", fmt.Sprintf("/users/%s", helper.RegularUser.ID), reqData)
	helper.AssertStatus(res, 200)

	body := helper.GetUserDTO(res)
	helper.Assert(body.FirstName == helper.RegularUser.FirstName, "Firstname should be the same")
	helper.Assert(body.Role == 2, "Role should be updated")
	helper.Assert(body.LastName == helper.RegularUser.LastName, "Lastname should be the same")
	helper.Assert(body.Email == "testing@mail.test", "Email should be updated")
}


func TestUserUpdateOGInvalidEmail(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	data := map[string]interface{}{
		"firstName": helper.RegularUser.FirstName,
		"lastName": helper.RegularUser.LastName,
		"email": "fake@email.",
		"role": helper.RegularUser.Role,
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	res := helper.SendAsRegularUser("put", fmt.Sprintf("/users/%s", helper.RegularUser.ID), reqData)

	helper.AssertStatus(res, 400)
	helper.AssertErrDTOPresent(res)
}


func TestAdminCanNotUpdateRoleWithOG(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	data := map[string]interface{}{
		"firstName": "Test",
		"lastName": helper.RegularUser.LastName,
		"email": helper.RegularUser.Email,
		"role": 2,
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	res := helper.SendAsAdmin("put", fmt.Sprintf("/users/%s", helper.RegularUser.ID), reqData)

	helper.AssertStatus(res, 401)
	helper.AssertErrDTOPresent(res)
}


// ----- User delete -----

func TestDeleteWorksForCurrentUser(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	res := helper.SendAsRegularUser("delete", fmt.Sprintf("/users/%s", helper.RegularUser.ID), nil)
	helper.AssertStatus(res, 204)
}


func TestDeleteDoesNotWorksForOtherUser(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	res := helper.SendAsRegularUser("delete", fmt.Sprintf("/users/%s", helper.AdminUser.ID), nil)

	helper.AssertStatus(res, 401)
	helper.AssertErrDTOPresent(res)
}


func TestDeleteWorksForSuperUser(t *testing.T) {
	helper := testhelper.InitTestDBAndService(t)
	helper.InitAuth()
	defer helper.CleanupAuth()

	res := helper.SendAsSuperAdmin("delete", fmt.Sprintf("/users/%s", helper.AdminUser.ID), nil)
	helper.AssertStatus(res, 204)
}
