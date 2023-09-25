package main

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/template/html"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"html/template"
	"os"
	"os/signal"
	"strings"
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
	app.Post("/beeload/create/bucket", checkUserPermission, createBucket)
	app.Post("/beeload/compare/release", compareRelease)
	app.Post("/beeload/set/project", setActiveUserProject)
	app.Post("/beeload/add/methodic", addMethodic) // TODO: добавить обработку методики
	app.Post("beeload/add/version", addVersion)    // TODO: добавить обработку дополнения версии
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
	app.Post("/get_project_buckets", GetProjectBuckets)
	app.Post("/get_bucket_projects", GetBucketProjects)
	app.Get("/compare_release", getCompareRelease)
	app.Post("/get_version_list", GetVersionsList)
}

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
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Shutdown(); err != nil {
		logrus.Error(err)
		fmt.Println("Error shutting down server:", err)
	}
}

func createBucket(c *fiber.Ctx) error {
	logrus.Debug("postCreateBucket!")
	requestData := map[string]interface{}{
		"host":   "example.com",
		"bucket": "my_bucket",
	}
	url := "/beeload/create/bucket"
	res := sendPost(c, url, requestData)
	fmt.Println(res)
	return nil //TODO: Надо что то сделать с отрисовкой ответа и добавить правильные данные
}
func addDataForBeeLoad(c *fiber.Ctx) error {
	logrus.Debug("postAddHomePageForBeeLoad!")
	//var req addconfluence
	body := c.Body()
	fmt.Println(body)
	//json.Unmarshal(c.Body(), &req)
	//pgAddInfo(req.Bucket, req.ID)
	return nil //TODO: реалиация
}

func compareData(c *fiber.Ctx) error {
	logrus.Debug("postCompareData!")
	//var req dataCompareRest
	//c.Body()
	//json.Unmarshal(c.Body(), &req)
	//fmt.Println(req)
	//confTestReport(req.Application, req.Bucket, req.ApplicationC)
	return nil //TODO: реалиация
}

func compareRelease(c *fiber.Ctx) error {
	logrus.Debug("postCompareRelease!")
	requestData := map[string]interface{}{
		"application":  "application1",
		"applicationC": "application2",
		"bucket":       "my_bucket",
	}
	url := "/beeload/compare/release"
	sendPost(c, url, requestData)
	return nil
}

func startPage(c *fiber.Ctx) error {
	logrus.Debug("getStartPage")
	return c.Render("main",
		fiber.Map{"Title": "Мотай вниз", "Description": "Где то тут юзаются шаблоны."})
}

func getMainPage(c *fiber.Ctx) error {
	logrus.Debug("getMainPage")
	return c.Render("index",
		fiber.Map{"Table_reports": get_last_10_reports_table(GetTableDataReports(c, "", 10)),
			"Table_tests":  get_current_tests(GetTableDataTests(c)),
			"Table_status": get_status_table(GetTableDataStatus(c)),
		})
}

func getCurrentTests(c *fiber.Ctx) error {
	logrus.Debug("getCurrentTests")
	url := "/beeload/get/tabledatacurrenttests"
	res := sendGet(c, url)
	dataStr := string(res)

	// Разбиваем на строки
	rows := strings.Split(dataStr, "},{")

	// Подготовка для парсинга
	var data [][]string
	for _, row := range rows {
		// Удаление лишних символов
		row = strings.Trim(row, "[{]}")
		// Разбиваем на элементы
		items := strings.Split(row, ",")
		var itemStrings []string
		for _, item := range items {
			item = strings.Trim(item, "\" ")
			itemStrings = append(itemStrings, item)
		}
		data = append(data, itemStrings)
	}
	//fmt.Println(data)
	return c.Render("current_tests",
		fiber.Map{"CurrentTests": get_test_table(data)})
}

func getSettings(c *fiber.Ctx) error {
	logrus.Debug("getSettings")
	value := c.Locals("user")
	claims, ok := (value).(*jwt.MapClaims)
	username, ok := (*claims)["username"].(string)
	db := getDbConn()
	activeProject, err := GetUserProject(db, username)
	projectsList, err := GetUserProjects(db, username)
	isAdmin, err := hasUserRole(db, username, "admin")
	logrus.Debug("is_admin: ", isAdmin)
	if !ok {
		// Обработка ошибки преобразования
		logrus.Error("getSettings: username conversion failed")
		return fmt.Errorf("username conversion failed")
	}
	if err != nil {
		logrus.Error(err)
	}
	additional := ""
	if isAdmin {
		projectsList, _ = GetAllProjects(db)
		additional = "        <button class=\"l_btn\">Добавить сценарий</button>\n        <button class=\"l_btn\">Администрирование</button>"
		return c.Render("settings",
			fiber.Map{
				"User":          "Текущий пользователь: " + username,
				"Additional":    additional,
				"Versions":      "<option selected>Выберите версию</option>",
				"ActiveProject": make_settings_projects_list(activeProject, projectsList)})
	} else {
		additional = ""
		return c.Render("settings",
			fiber.Map{
				"User":          "Текущий пользователь: " + username,
				"Additional":    additional,
				"Versions":      "<option selected>Выберите версию</option>",
				"ActiveProject": make_settings_projects_list(activeProject, projectsList)})
	}
}

func setActiveUserProject(c *fiber.Ctx) error {
	logrus.Debug("setActiveUserProject")
	value := c.Locals("user")
	claims, _ := (value).(*jwt.MapClaims)
	username, _ := (*claims)["username"].(string)
	db := getDbConn()
	bucket := new(Project)
	if err := c.BodyParser(bucket); err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Debug("Project = ", bucket.Name)
	err := SetActiveProject(db, username, bucket.Name)
	if err != nil {
		logrus.Error("setActiveUserProject: ", err)
	}
	activeProject, _ := GetUserProject(db, username)
	projectsList, _ := GetUserProjects(db, username)
	return c.Render("settings",
		fiber.Map{"ActiveProject": make_settings_projects_list(activeProject, projectsList)})
}

func addVersion(c *fiber.Ctx) error {
	version := new(Version)
	if err := c.BodyParser(version); err != nil {
		logrus.Error(err)
		return err
	}
	url := "/beeload/add/version"
	requestData := map[string]interface{}{
		"version": version.Value,
	}
	res := sendRequest(c, "Post", url, requestData)
	fmt.Print(res)
	return nil
} //TODO: COMPLETE

func addMethodic(c *fiber.Ctx) error {
	methodic := new(MethodicSet)
	if err := c.BodyParser(methodic); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
} //TODO: COMPLETE

func getReportHistory(c *fiber.Ctx) error {
	logrus.Debug("getReportHistory")
	value := c.Locals("user")
	claims, ok := (value).(*jwt.MapClaims)
	username, ok := (*claims)["username"].(string)
	db := getDbConn()
	activeProject, err := GetUserProject(db, username)
	if !ok {
		// Обработка ошибки преобразования
		logrus.Error("getSettings: username conversion failed")
		return fmt.Errorf("username conversion failed")
	}
	if err != nil {
		logrus.Error(err)
	}
	return c.Render("report_history",
		fiber.Map{"Table_reports": get_last_10_reports_table(GetTableDataReports(c, activeProject, 0))})
}

func getTestHistory(c *fiber.Ctx) error {
	logrus.Debug("getTestHistory")
	value := c.Locals("user")
	claims, ok := (value).(*jwt.MapClaims)
	username, ok := (*claims)["username"].(string)
	db := getDbConn()
	activeProject, err := GetUserProject(db, username)
	if !ok {
		// Обработка ошибки преобразования
		logrus.Error("getSettings: username conversion failed")
		return fmt.Errorf("username conversion failed")
	}
	if err != nil {
		logrus.Error(err)
	}
	return c.Render("test_history",
		fiber.Map{"Table_reports": get_last_10_reports_table(GetTableDataReports(c, activeProject, 0))})
}

func getStartTest(c *fiber.Ctx) error {
	logrus.Debug("getStartTest")
	return c.Render("start_test",
		fiber.Map{"Buckets": add_tags(get_bucket_list(c))})
}

func getCompare(c *fiber.Ctx) error {
	logrus.Debug("getCompare")
	return c.Render("compare_tests",
		fiber.Map{"Buckets": add_tags(get_bucket_list(c))})
}

func getCompareRelease(c *fiber.Ctx) error {
	logrus.Debug("getCompareRelease")
	return c.Render("relise_policy",
		fiber.Map{"Buckets": add_tags(get_bucket_list(c)), "Projects": `<option>Выберите бакет</option>`})
}

func getMakeReport(c *fiber.Ctx) error {
	logrus.Debug("getMakeReport")
	return c.Render("make_report",
		fiber.Map{"Buckets": `<option>Выберите проект</option>`, "Projects": add_tags(get_project_list(c))})
}

func getCreateBucket(c *fiber.Ctx) error {
	logrus.Debug("getCreateBucket")
	return c.Render("create_bucket",
		fiber.Map{"Projects": add_tags(get_project_list(c))})
}
func getReportHomePage(c *fiber.Ctx) error {
	logrus.Debug("getReportHomePage")
	return c.Render("set_report_homepage",
		fiber.Map{"Buckets": add_tags(get_bucket_list(c))})
}
func getSetMethodic(c *fiber.Ctx) error {
	logrus.Debug("getSetMethodic")
	return c.Render("set_methodic",
		fiber.Map{"Buckets": add_tags(get_bucket_list(c))})
}

func GetVersionsList(c *fiber.Ctx) error {
	logrus.Debug("GetVersionsList")
	project := new(Project)
	if err := c.BodyParser(project); err != nil {
		return err
	}
	return c.JSON(get_versions_list(c, project.Name))
}

func GetProjectBuckets(c *fiber.Ctx) error {
	logrus.Debug("getProjectBuckets")
	project := new(Project)
	if err := c.BodyParser(project); err != nil {
		return err
	}
	return c.JSON(get_project_buckets(c, project.Name))
}

func GetBucketProjects(c *fiber.Ctx) error {
	logrus.Debug("getBucketProjects")
	bucket := new(Project)
	if err := c.BodyParser(bucket); err != nil {
		return err
	}
	return c.JSON(get_bucket_projects(c, bucket.Name))
}
