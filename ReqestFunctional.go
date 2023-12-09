package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

//Здесь находятся запросы которые реализуют какой либо функционал кроме отрисовки страниц в браузере

func getReportHistory(c *fiber.Ctx) error {
	logrus.Debug("getReportHistory")
	value := c.Locals("user")
	claims, ok := (value).(*jwt.MapClaims)
	username, ok := (*claims)["username"].(string)
	activeProject, err := GetUserActiveProject(username)
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
	activeProject, err := GetUserActiveProject(username)
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

func addMethodic(c *fiber.Ctx) error {
	logrus.Debug("addMethodic")
	value := c.Locals("user")
	claims, _ := (value).(*jwt.MapClaims)
	username, _ := (*claims)["username"].(string)
	activeProject, err := GetUserActiveProject(username)
	methodic := new(MethodicSet)
	if err = c.BodyParser(methodic); err != nil {
		logrus.Error("addMethodic parse ERROR: ", err)
		return err
	}
	page, err := strconv.Atoi(methodic.Page)
	AddProjectMethodic(page, methodic.Version, activeProject)
	if err != nil {
		logrus.Error("addMethodic ERROR: ", err)
		return err
	}
	return c.JSON("OK")
}

func addVersion(c *fiber.Ctx) error {
	logrus.Debug("addVersion")
	value := c.Locals("user")
	claims, _ := (value).(*jwt.MapClaims)
	username, _ := (*claims)["username"].(string)
	activeProject, err := GetUserActiveProject(username)
	version := new(Version)
	if err = c.BodyParser(version); err != nil {
		logrus.Error(err)
		return err
	}
	err = AddProjectVersion(version.Value, activeProject)
	if err != nil {
		logrus.Error("addVersion ERROR: ", err)
		return err
	}
	return c.SendString("OK")
}

func getVersion(c *fiber.Ctx) error {
	logrus.Debug("getVersion")
	value := c.Locals("user")
	claims, _ := (value).(*jwt.MapClaims)
	username, _ := (*claims)["username"].(string)
	activeProject, err := GetUserActiveProject(username)
	res, err := GetProjectVersion(activeProject)
	if err != nil {
		logrus.Error("getVersion ERROR: ", err)
	}
	return c.JSON(res)
}

func createBucket(c *fiber.Ctx) error {
	logrus.Debug("postCreateBucket!")
	//TODO: ГОВНА ПИРОГА!!!!
	value := c.Locals("user")
	claims, _ := (value).(*jwt.MapClaims)
	username, _ := (*claims)["username"].(string)
	logrus.Debug("username = ", username)

	newBucket := new(NewBucket)
	if err := c.BodyParser(newBucket); err != nil {
		logrus.Error(err)
		return err
	}
	activeProject, err := GetUserActiveProject(username)
	err = AddProjectBucket(newBucket.Bucket, newBucket.Host, activeProject)
	if err != nil {
		logrus.Error("createBucket ERROR: ", err)
		return err
	}
	return c.SendString("OK")
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
	//sendPost(c, url, requestData)
	sendRequest(c, "Post", url, requestData)
	return nil
}

func getCurrentTests(c *fiber.Ctx) error {
	logrus.Debug("getCurrentTests")
	url := "/beeload/get/tabledatacurrenttests"
	//res := sendGet(c, url)
	res := sendRequest(c, "Get", url)
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

func addProject(c *fiber.Ctx) error {
	logrus.Debug("addProject")
	value := c.Locals("user")
	claims, _ := (value).(*jwt.MapClaims)
	username, _ := (*claims)["username"].(string)
	logrus.Debug("username = ", username)
	logrus.Debug("raw = ", string(c.Body()))

	project := new(Project)
	if err := c.BodyParser(project); err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Debug("Project = ", project.Name)
	err := AddNewProject(project.Name)
	if err != nil {
		return err
	}
	return c.SendString("OK")
}

func addConflPage(c *fiber.Ctx) error {
	logrus.Debug("addConflPage")
	value := c.Locals("user")
	claims, _ := (value).(*jwt.MapClaims)
	username, _ := (*claims)["username"].(string)
	logrus.Debug("username = ", username)

	page := new(NewRootPage)
	if err := c.BodyParser(page); err != nil {
		logrus.Error("addConflPage parse ERROR: ", err)
		return err
	}
	activeProject, err := GetUserActiveProject(username)
	int_page, err := strconv.Atoi(page.Page)
	AddProjectRootPage(int_page, activeProject)
	if err != nil {
		logrus.Error("addConflPage ERROR: ", err)
		return err
	}
	return c.JSON("OK")
}

func setActiveUserProject(c *fiber.Ctx) error {
	logrus.Debug("setActiveUserProject")
	value := c.Locals("user")
	claims, _ := (value).(*jwt.MapClaims)
	username, _ := (*claims)["username"].(string)
	bucket := new(Project)
	if err := c.BodyParser(bucket); err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Debug("Project = ", bucket.Name)
	isAdmin, err := hasUserRole(username, "admin")
	//err = SetActiveProject(db, username, bucket.Name, isAdmin)
	err = SetActiveProject(username, bucket.Name, isAdmin)
	if err != nil {
		logrus.Error("setActiveUserProject ERROR: ", err)
		return err
	}
	activeProject, err := GetUserActiveProject(username)
	if err != nil {
		logrus.Error("GetUserActiveProject ERROR: ", err)
		return err
	}
	logrus.Debug("setActiveUserProject ACTIVE: ", activeProject)
	//projectsList, _ := GetUserProjects(db, username)
	return c.SendString("OK")
}

func assignProjects(c *fiber.Ctx) error {
	logrus.Debug("assignProjects")
	username := c.FormValue("user")
	projectNames := strings.Split(c.FormValue("projects"), ",")
	//projectNames := c.FormValue("projects")
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Projects: %v\n", projectNames)
	return c.SendStatus(fiber.StatusOK)
}