package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"strings"
)

//Здесь находятся запросы которые реализуют какой либо функционал кроме отрисовки страниц в браузере

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

func addMethodic(c *fiber.Ctx) error {
	methodic := new(MethodicSet)
	if err := c.BodyParser(methodic); err != nil {
		logrus.Error(err)
		return err
	}
	url := "/beeload/add/methodic"
	requestData := map[string]interface{}{
		"bucket":  methodic.Bucket,
		"version": methodic.Version,
		"page":    methodic.Page,
	}
	res := sendRequest(c, "Post", url, requestData)
	return c.SendString(string(res))
} //TODO: COMPLETE

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
	//fmt.Print(string(res))
	return c.SendString(string(res))
}

func createBucket(c *fiber.Ctx) error {
	logrus.Debug("postCreateBucket!")
	requestData := map[string]interface{}{
		"host":   "example.com",
		"bucket": "my_bucket",
	}
	url := "/beeload/create/bucket"
	//res := sendPost(c, url, requestData)
	res := sendRequest(c, "Post", url, requestData)
	c.SendString(string(res))
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
		logrus.Error("setActiveUserProject ERROR: ", err)
		return err
	}
	activeProject, err := GetUserProject(db, username)
	if err != nil {
		logrus.Error("GetUserProject ERROR: ", err)
		return err
	}
	logrus.Debug("setActiveUserProject ACTIVE: ", activeProject)
	//projectsList, _ := GetUserProjects(db, username)
	return c.SendString("OK")
}
