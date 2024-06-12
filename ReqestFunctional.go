package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

func GetListOfTests(c *fiber.Ctx) error {
	logrus.Debug("GetListOfTests")
	url := "/get_project_buckets"
	body := c.Body()
	logrus.Debug("GetListOfTests body: ", string(body))
	res := sendRequest(c, "Post2", url, body)
	logrus.Debug("GetListOfTests res: ", string(res))
	return c.SendString(string(res))
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

func addUserToProject(c *fiber.Ctx) error {
	logrus.Debug("addUserToProject")
	logrus.Debug("raw = ", string(c.Body()))
	userProject := new(UserProject)
	if err := c.BodyParser(userProject); err != nil {
		logrus.Error(err)
		return err
	}
	err := AddUserToProject(userProject.User, userProject.Project)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
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
	value := c.Locals("user")
	claims, _ := (value).(*jwt.MapClaims)
	username, _ := (*claims)["username"].(string)
	phone := c.FormValue("phone")
	projectNames := c.FormValue("projects")
	err := set_telnumber_to_username(username, phone)
	if err != nil {
		logrus.Error("assignProjects set_telnumber_to_username ERROR: ", err)
		return err
	}
	logrus.Debug("projectNames: ", projectNames, "\n ", strings.Split(projectNames, ","))
	err = set_subs_by_telnumber(phone, strings.Split(projectNames, ","))
	if err != nil {
		logrus.Error("assignProjects set_subs_by_telnumber ERROR: ", err)
		return err
	}
	//TODO: добавить отправку запроса с подписками на бек
	//TODO: необходимо добавить синхронизацию подписок с беком
	return c.SendStatus(fiber.StatusOK)
}

func testCreate(c *fiber.Ctx) error {
	logrus.Debug("testCreate")
	url := viper.GetString("backend.test_starter") + "/create"
	res := sendRequest(c, "Post3", url, c.Body())
	if res != nil {
		return c.SendStatus(fiber.StatusOK)
	} else {
		logrus.Error("testCreate SendStatus error!")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
}

func getSyncBuckets(c *fiber.Ctx) error {
	logrus.Debug("getSyncBuckets")
	res := syncBuckets()
	fmt.Println(res)
	logrus.Debug(res)
	return c.SendStatus(fiber.StatusOK)
}

func makeReport(c *fiber.Ctx) error {
	logrus.Debug("makeReport")
	logrus.Debug("makeReport OriginalURL: ", c.OriginalURL())
	logrus.Debug("makeReport Body: ", string(c.Body()))
	res := sendRequest(c, "Post2", c.OriginalURL(), c.Body())
	if res != nil {
		dataStr := string(res)
		logrus.Debug("makeReport dataStr: ", dataStr)
		return c.SendString(dataStr)
	} else {
		logrus.Debug("makeReport sendRequest error! ")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	//TODO: Проверить работу отправки отчёта
}

func startTestParseEnv(c *fiber.Ctx) error {
	logrus.Debug("startTestParseEnv")
	url := c.OriginalURL()
	logrus.Debug("startTestParseEnv OriginalURL: ", url)
	logrus.Debug("startTestParseEnv Body: ", string(c.Body()))
	if strings.HasPrefix(url, "/parse/env/custom") {
		parts := strings.SplitN(url, "/custom", 2)
		trimmedPath := parts[0]
		logrus.Debug("trimmedPath: ", trimmedPath)
		url = trimmedPath
		return c.Render("scenario_generator",
			fiber.Map{"CurrentTests": "res"})
	}
	res := sendRequest(c, "Post3", viper.GetString("backend.test_starter")+url, c.Body())
	logrus.Debug("startTestParseEnv res: ", string(res))
	envs := new(GitEnvData)
	if err := json.Unmarshal(res, envs); err != nil {
		logrus.Error("startTestParseEnv Unmarshal error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	gitUrl := new(GitLabUrl)
	if err := json.Unmarshal(c.Body(), gitUrl); err != nil {
		logrus.Error("startTestParseEnv Unmarshal gitUrl error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	fmt.Println("ENVS: ", envs.Data)
	fmt.Println("gitUrl: ", gitUrl.Gitlab)
	resp := "<form class=\"formTestStart\" id=\"TestStartForm\">\n"
	resp += "<input class=\"back_button\" type=\"submit\" value=\"Назад\" onclick=\"goBack()\"/>\n"
	resp += "<div class=\"Envs_top\" id=\"env\">\n"
	resp +=
		"<div class=\"input_field\">\n " +
			"<p class=\"area_label\">GIT URL</p>\n " +
			"<span>\n " +
			"<input type=\"text\" class=\"GitUrl version\" id=\"url\" name=\"url\" value=\"" + gitUrl.Gitlab + "\">\n " +
			"</span>\n " +
			"</div>\n" +
			"<div class=\"input_field left\">\n " +
			"<p class=\"area_label\">Количество генераторов</p>\n " +
			"<span>\n " +
			"<i class=\"fa fa-clock\"></i>\n " +
			"<input type=\"number\" class=\"genCount version\" id=\"quantity\" name=\"quantity\" min=\"1\" max=\"10\" value=\"1\" required>\n " +
			"</span>\n " +
			"</div>\n " +
			"<div class=\"input_field\">\n " +
			"<p class=\"area_label\">Аппаратная конфигурация</p>\n " +
			"<span>\n " + make_generators_list() + "\n " + "</span>\n " +
			"</div>\n" +
			"<div class=\"input_field\">\n " +
			"<p class=\"area_label\">Выбор файла для запуска</p>\n " +
			"<span>\n " +
			"<select class=\"fileName\" name=\"filename\">"
	for _, item := range envs.TestPlan {
		resp += "<option>" + item + "</option>"
	}
	resp += "</select>" +
		"</span>\n " +
		"</div>\n"
	resp += "</div>\n"
	resp += "<div class=\"Envs\" id=\"env\">\n"
	for _, env := range envs.Data {
		resp += "<div class=\"form-row\">\n"
		resp += "<p class=\"area_label\">" + env.Key + "</p>\n"
		resp += "<input type=\"text\" name=\"" + env.Key + "\" value=\"" + env.Value + "\"/>\n"
		resp += "</div>\n"
		fmt.Println("ENV key: ", env.Key, "   ENV value: ", env.Value)
	}
	resp += "</div>\n"
	resp += "<input class=\"submit_button\" type=\"submit\" value=\"Скопировать комманду CURL\" onclick=\"createCurlUrl()\"/>\n"
	resp += "<input class=\"submit_button\" type=\"submit\" value=\"Запустить\" onclick=\"startTest()\"/>\n"
	resp += "</form>\n"
	if res != nil {
		logrus.Debug("startTestParseEnv resp: ", resp)
		return c.SendString(resp)
	} else {
		logrus.Debug("makeReport sendRequest error! ")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
}

func hiddenAddRole(c *fiber.Ctx) error {
	role := c.Query("role")
	if role == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Role parameter is required")
	}
	err := addRole(role)
	if err != nil {
		logrus.Error(err)
	}
	return c.SendString("Role added: " + role)
}

func addMonitoringUrl(c *fiber.Ctx) error {
	//TODO: переделать запрос с проброской на бэк
	logrus.Debug("addMonitoringUrl")
	newUrl := new(MonitoringUrl)
	if err := json.Unmarshal(c.Body(), newUrl); err != nil {
		logrus.Error("addMonitoringUrl Unmarshal newUrl error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if err := addMonUrl(newUrl.Project, newUrl.NewUrl); err != nil {
		logrus.Error("addMonitoringUrl addMonUrl error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendString("OK")
}

func getMonitoringUrl(c *fiber.Ctx) error {
	//TODO: переделать запрос с проброской на бэк
	logrus.Debug("getMonitoringUrl")
	project := new(GetMonitoringUrl)
	if err := json.Unmarshal(c.Body(), project); err != nil {
		logrus.Error("getMonitoringUrl Unmarshal project error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if urls, err := getMonUrls(project.Project); err != nil {
		logrus.Error("addMonitoringUrl addMonUrl error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	} else {
		return c.JSON(urls)
	}
}

func getDestroy(c *fiber.Ctx) error {
	logrus.Debug("getDestroy")
	url := "http://ms-loadrtst038:9999"
	return proxyReq(c, url)
}
