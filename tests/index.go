package tests

// import (
// 	"chi-users-project/app/services/dtos"
// 	"chi-users-project/app/services"
// 	"chi-users-project/app/models"
// 	"chi-users-project/app"
// 	"testing"
// 	"errors"
// )

// func urlBase() string {
// 	return "http://localhost:8000"
// }

// type BaseTest struct {
// 	*TestSuite

// 	regularUser				models.User
// 	regularToken			string
// 	adminUser					models.User
// 	adminToken				string
// 	superAdminUser		models.User
// 	superAdminToken		string
// }

// func(this *BaseTest) SendAsNoOne(typ string, urlEnd string, body []byte) {
// 	req := getRequest(typ, urlEnd, body)

// 	tr := this.NewTestRequest(req)

// 	tr.Send()
// }


// func(this *BaseTest) SendAsRegularUser(typ string, urlEnd string, body []byte) {
// 	if this.regularToken == "" {
// 		panic(errors.New("Must call InitAuth in Before hook to send as Regular User"))
// 	}

// 	req := getRequest(typ, urlEnd, body)

// 	req.Header.Set("Authorization", this.regularToken)

// 	tr := this.NewTestRequest(req)

// 	tr.Send()
// }


// func(this *BaseTest) SendAsAdmin(typ string, urlEnd string, body []byte) {
// 	if this.adminToken == "" {
// 		panic(errors.New("Must call InitAuth in Before hook to send as Admin"))
// 	}

// 	req := getRequest(typ, urlEnd, body)

// 	req.Header.Set("Authorization", this.adminToken)

// 	tr := this.NewTestRequest(req)

// 	tr.Send()
// }


// func(this *BaseTest) SendAsSuperAdmin(typ string, urlEnd string, body []byte) {
// 	if this.superAdminToken == "" {
// 		panic(errors.New("Must call InitAuth in Before hook to send as Super Admin"))
// 	}

// 	req := getRequest(typ, urlEnd, body)

// 	req.Header.Set("Authorization", this.superAdminToken)

// 	tr := this.NewTestRequest(req)

// 	tr.Send()
// }


// func(this *BaseTest) InitAuth() {
// 	this.createTestAuthUsers()

// 	superDTO := dtos.LoginDTO {
// 		Email: "superadmin@test.com",
// 		Password: "testpassword7",
// 	}

// 	this.superAdminToken = this.getTokenViaService(superDTO) //this.getToken(superData)

// 	adminDTO := dtos.LoginDTO {
// 		Email: "admin@test.com",
// 		Password: "testpassword8",
// 	}

// 	this.adminToken = this.getTokenViaService(adminDTO)

// 	regularDTO := dtos.LoginDTO {
// 		Email: "regular@test.com",
// 		Password: "testpassword9",
// 	}

// 	this.regularToken = this.getTokenViaService(regularDTO)
// }


// func(this *BaseTest) createTestAuthUsers() {
// 	// create 3 users (1 of each auth level)
// 	// bypass controllers/services save straight to DB
// 	db := app.DBConnection.GetDB()

// 	superAdmin := models.User {
// 		FirstName: "Super",
// 		LastName: "Admin",
// 		Email: "superadmin@test.com",
// 		Password: "testpassword7",
// 		Role: 3,
// 	}
// 	if createErr := db.Create(&superAdmin).Error; createErr != nil {
// 		panic(createErr)
// 	}

// 	this.superAdminUser = superAdmin

// 	admin := models.User {
// 		FirstName: "Admin",
// 		LastName: "Tester",
// 		Email: "admin@test.com",
// 		Password: "testpassword8",
// 		Role: 2,
// 	}
// 	if createErr := db.Create(&admin).Error; createErr != nil {
// 		panic(createErr)
// 	}
// 	this.adminUser = admin

// 	regular := models.User {
// 		FirstName: "Regular",
// 		LastName: "User",
// 		Email: "regular@test.com",
// 		Password: "testpassword9",
// 	}
// 	if createErr := db.Create(&regular).Error; createErr != nil {
// 		panic(createErr)
// 	}
// 	this.regularUser = regular
// }


// func(this *BaseTest) CleanupAuth() {
// 	db := app.DBConnection.GetDB()

// 	if deleteErr := db.Unscoped().Delete(&this.superAdminUser).Error; deleteErr != nil {
// 		panic(deleteErr)
// 	}
// 	if deleteErr := db.Unscoped().Delete(&this.adminUser).Error; deleteErr != nil {
// 		panic(deleteErr)
// 	}
// 	if deleteErr := db.Unscoped().Delete(&this.regularUser).Error; deleteErr != nil {
// 		panic(deleteErr)
// 	}
// }


// func(this BaseTest) getTokenViaService(creds dtos.LoginDTO) string {
// 	bs := services.InitService(revel.AppLog)
// 	service := services.LoginService{BaseService: &bs}

// 	tokenDTO, _ := service.LoginUser(creds, false)

// 	if tokenDTO.Token == "" {
// 		panic(errors.New("No Token!!"))
// 	}

// 	return tokenDTO.Token
// }
