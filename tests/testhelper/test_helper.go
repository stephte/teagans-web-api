package testhelper

import (
	"youtube-downloader/app/services/dtos"
	"youtube-downloader/app/services"
	"youtube-downloader/app/models"
	"github.com/go-chi/httplog"
	"youtube-downloader/config"
	"github.com/joho/godotenv"
	"testing"
	"strconv"
	"errors"
	"fmt"
	"os"
)


type TestHelper struct {
	RegularUser			models.User
	RegularToken		string
	RegularCsrf			string
	AdminUser			models.User
	AdminToken			string
	AdminCsrf			string
	SuperAdminUser		models.User
	SuperAdminToken		string
	SuperAdminCsrf		string

	dbConn				config.DBConn
	service				services.BaseService
	server				config.Router
	t					*testing.T
}


// ----- Initialization/cleanup methods -----


func InitTestDBAndService(t *testing.T) (TestHelper) {
	fmt.Print("\n")
	logger := httplog.NewLogger("youtube-downloader-tests", httplog.Options{ JSON: true })

	envErr := godotenv.Load("../.testenv")
	if envErr != nil {
		t.Fatalf("Test environment failed to load: %s", envErr.Error())
	}

	dbConnection := config.InitDBConn(logger, false)

	port, err := strconv.Atoi(os.Getenv("CHI_YT_DBPORT"))
	if err != nil {
		t.Fatalf("failed to convert port to int: %s", err.Error())
	}

	dbConnection.SetHost(os.Getenv("CHI_YT_DBHOST"))
	dbConnection.SetUser(os.Getenv("CHI_YT_DBUSER"))
	dbConnection.SetPassword(os.Getenv("CHI_YT_DBPASSWORD"))
	dbConnection.SetName(os.Getenv("CHI_YT_DBNAME"))
	dbConnection.SetPort(port)

	err = dbConnection.FireUp()
	if err != nil {
		t.Fatalf("Database failed to initialize: %s", err.Error())
	}

	server := config.SetupRouter(logger, dbConnection.GetDB(), "test")

	service := services.InitService(logger, dbConnection.GetDB())

	return TestHelper {
		t: t,
		dbConn: dbConnection,
		server: server,
		service: service,
	}
}


func(this *TestHelper) InitAuth() {
	this.CreateTestAuthUsers()

	superDTO := dtos.LoginDTO {
		Email: "superadmin@test.com",
		Password: "testpassword7",
	}

	this.SuperAdminToken, this.SuperAdminCsrf = this.getTokenViaService(superDTO)

	adminDTO := dtos.LoginDTO {
		Email: "admin@test.com",
		Password: "testpassword8",
	}

	this.AdminToken, this.AdminCsrf = this.getTokenViaService(adminDTO)

	regularDTO := dtos.LoginDTO {
		Email: "regular@test.com",
		Password: "testpassword9",
	}

	this.RegularToken, this.RegularCsrf = this.getTokenViaService(regularDTO)
}


func(this *TestHelper) CreateTestAuthUsers() {
	db := this.dbConn.GetDB()

	superAdmin := models.User {
		FirstName: "Super",
		LastName: "Admin",
		Email: "superadmin@test.com",
		Password: "testpassword7",
		Role: 3,
	}
	if createErr := db.Create(&superAdmin).Error; createErr != nil {
		this.t.Log("SuperAdmin already exists in DB")
	}

	this.SuperAdminUser = superAdmin

	admin := models.User {
		FirstName: "Admin",
		LastName: "Tester",
		Email: "admin@test.com",
		Password: "testpassword8",
		Role: 2,
	}
	if createErr := db.Create(&admin).Error; createErr != nil {
		this.t.Log("Admin already exists in DB")
	}
	this.AdminUser = admin

	regular := models.User {
		FirstName: "Regular",
		LastName: "User",
		Email: "regular@test.com",
		Password: "testpassword9",
	}
	if createErr := db.Create(&regular).Error; createErr != nil {
		this.t.Log("RegularUser already exists in DB")
	}
	this.RegularUser = regular
}


func(this TestHelper) CleanupAuth() {
	db := this.dbConn.GetDB()

	if deleteErr := db.Unscoped().Delete(&this.SuperAdminUser).Error; deleteErr != nil {
		panic(deleteErr)
	}
	if deleteErr := db.Unscoped().Delete(&this.AdminUser).Error; deleteErr != nil {
		panic(deleteErr)
	}
	if deleteErr := db.Unscoped().Delete(&this.RegularUser).Error; deleteErr != nil {
		panic(deleteErr)
	}

	this.Cleanup()
}


func(this TestHelper) Cleanup() {
	this.dbConn.CoolDown()
	fmt.Print("\n")
}


func(this TestHelper) getTokenViaService(creds dtos.LoginDTO) (string, string) {
	service := services.LoginService{BaseService: &this.service}

	tokenDTO, _, _ := service.LoginUser(creds, false)

	if tokenDTO.Token == "" {
		panic(errors.New("No Token!!"))
	}

	return tokenDTO.Token, tokenDTO.CSRF
}
