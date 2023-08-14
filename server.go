package main

import (
	"chatproxy/middlewares"
	"chatproxy/routes"
	"chatproxy/util"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"net/http"
	"os"
	"strings"
)

// init() function
// - Called at startup to initialize the application
// - Loads config from .env file using viper
// - Configures logging using logrus
// - Sets log level from LOG_LEVEL environment var
// - Log rotation with lumberjack
// - Outputs logs to file and stdout
func init() {
	viper.SetConfigFile(util.EnvironmentFile)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error when read config file")
		return
	}
	// setup logrus
	logLevel, err := log.ParseLevel(viper.GetString(util.LogLevel))
	if err != nil {
		logLevel = log.InfoLevel
	}
	log.SetLevel(logLevel)
	log.SetFormatter(&log.JSONFormatter{})
	// Log to file&console
	//file, _ := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// Set the Lumberjack logger
	lumberjackLogger := &lumberjack.Logger{
		Filename:  viper.GetString(util.LogFilePath),
		MaxSize:   viper.GetInt(util.LogSize),
		LocalTime: true,
	}
	log.SetOutput(io.MultiWriter(os.Stdout, lumberjackLogger))

}

// main() function
// - Creates new Echo server instance
// - Registers middleware for logging, errors, etc.
// - Attaches route handlers from routes package
// - Defines handler for default "/" route
// - Starts server on configured port
func main() {
	// Echo instance
	e := echo.New()
	e.HideBanner = true
	// Middleware
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(100)))
	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	e.Use(middlewares.LoggingMiddleware)
	e.Use(middlewares.ServerHeader)
	//Domain White List Functionality
	if viper.GetBool(util.DomainWhitelistEnable) {
		configCORS(e)
	}
	//Routes
	routes := routes.Routes{}
	routes.ActivateRoutes(e)

	e.GET("/", defaultRoute)

	// Start server
	port := viper.GetInt(util.ServerPort)
	log.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

func configCORS(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: parseStringToStringList(util.DomainWhitelistValues, util.DomainWhiteListAllowAll),
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost,
			http.MethodDelete},
	}))
}

// defaultRoute() handler
// - Returns simple response for default route
// - Business logic will be in separate handlers
func defaultRoute(c echo.Context) error {
	r := "server running at port:" + viper.GetString(util.ServerPort)
	return c.String(http.StatusOK, r)
}

func parseStringToStringList(variable string, preset string) []string {
	var stringArrays []string
	if value, found := os.LookupEnv(variable); found {
		stringArrays = strings.Split(value, ",")
	}

	if stringArrays == nil {
		stringArrays = []string{preset}
	}

	return stringArrays
}
