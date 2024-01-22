package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/template/html"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"html/template"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func setupUnprotectedRoutes(app *fiber.App) {
	app.Get("/login", getLogin)
	app.Post("/login", loginHandler)
	app.Get("/logout", logout)
	app.Get("/register", getRegister)
	app.Post("/register", register)
}

func setupProtectedRoutes(app *fiber.App) {
	app.Post("/beeload/add/home", addDataForBeeLoad)
	app.Post("/beeload/compare/data", compareData)
	//app.Post("/beeload/create/bucket", checkUserPermission, createBucket)
	app.Post("/beeload/create/bucket", createBucket)
	app.Post("/beeload/compare/release", compareRelease)
	app.Post("/beeload/set/project", setActiveUserProject)
	app.Post("/beeload/get/version", getVersion)
	app.Post("/beeload/add/methodic", addMethodic) // TODO: добавить обработку методики
	app.Post("/beeload/add/version", addVersion)
	app.Post("/beeload/add/project", addProject)
	app.Post("/beeload/add/confl_page", addConflPage)
	app.Post("/beeload/add/user_to_project", addUserToProject)
	app.Post("/beeload/test/create", testCreate) //TODO: transfer request
	app.Post("/beeload/make/report", makeReport) //TODO: transfer request
	app.Get("/", startPage)
	app.Get("/main_page", getMainPage)
	app.Get("/compare", getCompare)
	app.Get("/current_tests", getCurrentTests)
	app.Get("/report_history", getReportHistory)
	app.Get("/test_history", getTestHistory)
	app.Get("/settings", getSettings)
	app.Get("/create_bucket", getCreateBucket)
	app.Get("/set_report_homepage", getReportHomePage)
	app.Get("/set_methodic", getSetMethodic)
	app.Get("/make_report", getMakeReport)
	app.Get("/start_test", getStartTest)
	app.Get("/adminPanel", getAdminPanel)
	app.Get("/adminPanel/subscription", getAdminSubscription)
	app.Get("/bucket/sync", getSyncBuckets)
	app.Post("/get_project_buckets", GetProjectBuckets)
	app.Post("/get_list_of_tests", GetListOfTests)
	app.Post("/get_bucket_projects", GetBucketProjects)
	app.Get("/compare_release", getCompareRelease)
	app.Post("/get_version_list", GetVersionsList)
	app.Post("/get_host_list", GetHostList)
	app.Post("/assignProjects", assignProjects)
	app.Post("/parse/env", startTestParseEnv)
}

var db *sql.DB

func main() {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
		return
	} else {
		fmt.Println("Config file: readed sucessfully.")
	}

	// часть отвечающая за логи //
	InitLogger(viper.GetBool("server.debug"), viper.GetString("server.log_level"), viper.GetString("server.log_filename"))

	engine := html.New("build/views", ".html")
	engine.AddFunc(
		// add unescape function
		"unescape", func(s string) template.HTML {
			return template.HTML(s)
		},
	)
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static(
		"/",            // mount static
		"build/public", // path to the static file folder
	)
	app.Use(favicon.New(favicon.Config{File: "build/public/image/512x512.png"}))
	setupUnprotectedRoutes(app)
	app.Use(jwtMiddleware())
	setupProtectedRoutes(app)

	DBInit()
	err := createSQLiteDB()
	if err != nil {
		logrus.Error("Error creating db:", err)
		return
	}

	// Запуск переодической синхронизации
	c := cron.New()
	_, _ = c.AddFunc(viper.GetString("server.sync_period"), func() {
		logrus.Debug("CRON STARTED!")
		syncBuckets() // синхронизация бакетов с беком
		logrus.Debug("CRON FINISHED!")
		fmt.Println("Выполнена периодическая синхронизация.")
	})
	c.Start()

	// Запуск сервера в горутине
	go func() {
		port := viper.GetInt("server.port")
		if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
			logrus.Error("Error starting server:", err)
			fmt.Println("Error starting server:", err)
		}
	}()

	// Обработка сигналов для завершения программы
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	logrus.Error("Shutting down...")
	fmt.Println("Shutting down...")
	db.Close()
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Shutdown(); err != nil {
		logrus.Error(err)
		fmt.Println("Error shutting down server:", err)
	}
}

//TODO: доделать проверку полей при создании отчёта
//TODO: необходимо добавить синхронизацию подписок с беком assignProjects
//TODO: сделать проверку заполнения полей на странице настроек
//TODO:
//TODO:
