package main

import (
	"chi-users-project/config"
	"github.com/go-chi/httplog"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"strconv"
	"strings"
	"fmt"
	"os"
)


func main() {
	logger := httplog.NewLogger("chi-users-project", httplog.Options{
    JSON: true,
  })

	runmode := "dev"
	if len(os.Args) > 1 {
		if strings.ToLower(os.Args[1]) == "prod" || strings.ToLower(os.Args[1]) == "production" {
			runmode = "prod"
		} else if strings.ToLower(os.Args[1]) == "test" {
			runmode = "test"
		}
	}

	fmt.Printf("\n ----- Running app in *%s* mode -----\n\n", runmode)
	
	envErr := godotenv.Load(fmt.Sprintf(".env/%s.env", runmode))
	if envErr != nil {
		logger.Fatal().Err(envErr).Msg("")
		return
	}
	logger.Info().Msg("Successfully Loaded Env")

	dbConnection, connErr := initDB(logger)
	if connErr != nil {
		logger.Fatal().Err(connErr).Msg("")
		return
	}
	defer dbConnection.CoolDown()

	router := config.SetupRouter(logger, dbConnection.GetDB(), runmode)

	router.StartGracefulServer(os.Getenv("CHI_BASE_URL"), os.Getenv("CHI_PORT"))
}


func initDB(logger zerolog.Logger) (config.DBConn, error) {
	dbConnection := config.InitDBConn(logger, true)

	port, err := strconv.Atoi(os.Getenv("CHI_DBPORT"))
	if err != nil {
		return dbConnection, err
	}

	dbConnection.SetHost(os.Getenv("CHI_DBHOST"))
	dbConnection.SetUser(os.Getenv("CHI_DBUSER"))
	dbConnection.SetPassword(os.Getenv("CHI_DBPASSWORD"))
	dbConnection.SetName(os.Getenv("CHI_DBNAME"))
	dbConnection.SetPort(port)

	err = dbConnection.FireUp()
	if err != nil {
		return dbConnection, err
	}

	return dbConnection, nil
}
