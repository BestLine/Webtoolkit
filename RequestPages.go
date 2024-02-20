package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

//Здесь находятся запросы которые реализуют отрисовку страниц в браузере

func startPage(c *fiber.Ctx) error {
	logrus.Debug("getStartPage")
	return c.Render("main",
		fiber.Map{"Title": "Мотай вниз", "Description": "Где то тут юзаются шаблоны."})
}

func getMainPage(c *fiber.Ctx) error {
	logrus.Debug("getMainPage")
	return c.Render("index",
		fiber.Map{"Table_reports": get_last_10_reports_table(GetTableDataReports(c, "", 10)),
			"Table_tests":  GetCurrentTests(c),
			"Table_status": "get_status_table(GetTableDataStatus(c))",
		})
}

func GetProjectBuckets(c *fiber.Ctx) error {
	logrus.Debug("getProjectBuckets")
	project := new(Project)
	if err := c.BodyParser(project); err != nil {
		logrus.Error("getProjectBuckets Parsing error: ", err)
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

func GetHostList(c *fiber.Ctx) error {
	logrus.Debug("GetHostList")
	project := new(Project)
	if err := c.BodyParser(project); err != nil {
		return err
	}
	return c.JSON(get_host_list(c, project.Name))
}

func getCompare(c *fiber.Ctx) error {
	logrus.Debug("getCompare")
	return c.Render("compare_tests",
		fiber.Map{"Buckets": add_tags(get_bucket_list(c))})
}

func getCompareRelease(c *fiber.Ctx) error {
	logrus.Debug("getCompareRelease")
	return c.Render("relise_policy",
		fiber.Map{"Buckets": add_tags(get_bucket_list(c)), "Projects": `<option>Выберите бакет disabled</option>`})
}

func getMakeReport(c *fiber.Ctx) error {
	//TODO: ЧИНИМ НАХУЙ
	//TODO: ЧИНИМ НАХУЙ
	//TODO: ЧИНИМ НАХУЙ
	//TODO: ЧИНИМ НАХУЙ
	logrus.Debug("getMakeReport")
	value := c.Locals("user")
	claims, _ := (value).(*jwt.MapClaims)
	username, _ := (*claims)["username"].(string)
	activeProject, _ := GetUserActiveProject(username)
	projectsList, _ := GetUserProjects(username)
	isAdmin, _ := hasUserRole(username, "admin")
	//fmt.Println("getMakeReport username ", username)
	//fmt.Println("getMakeReport activeProject ", activeProject)
	//fmt.Println("getMakeReport projectsList ", projectsList)
	if isAdmin {
		projectsList, _ = GetAllProjects()
		activeProject = "Выберите проект"
		return c.Render("make_report",
			fiber.Map{"Buckets": `<option style=" display: none;">Выберите проект</option>`,
				"Projects": make_settings_projects_list(activeProject, projectsList)})
	} else {
		return c.Render("make_report",
			fiber.Map{"Buckets": `<option style=" display: none;">Выберите проект</option>`,
				"Projects": make_settings_projects_list(activeProject, projectsList)})
	}
}

func getCreateBucket(c *fiber.Ctx) error {
	logrus.Debug("getCreateBucket")
	return c.Render("create_bucket",
		fiber.Map{"Projects": add_tags(get_project_list(c))})
}

func getAdminPanel(c *fiber.Ctx) error {
	logrus.Debug("getAdminPanel")
	return c.Render("adminPanel",
		fiber.Map{"UserProjectsList": make_user_project_list()})
}

func getStartTest(c *fiber.Ctx) error {
	logrus.Debug("getStartTest")
	return c.Render("start_test",
		fiber.Map{"GeneratorsList": make_generators_list()})
}

func getSettings(c *fiber.Ctx) error {
	logrus.Debug("getSettings")
	value := c.Locals("user")
	claims, ok := (value).(*jwt.MapClaims)
	username, ok := (*claims)["username"].(string)
	activeProject, err := GetUserActiveProject(username)
	projectsList, err := GetUserProjects(username)
	isAdmin, err := hasUserRole(username, "admin")
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
		projectsList, _ = GetAllProjects()
		activeProject = "Выберите проект"
		additional = "<a class=\"l_btn\" href=\"/adminPanel\">Администрирование</a>"
		return c.Render("settings",
			fiber.Map{
				"User":          "Текущий пользователь: " + username,
				"Additional":    additional,
				"Versions":      "<option disabled selected>Выберите версию</option>",
				"HostList":      "<option disabled selected>Выберите хост</option>",
				"ActiveProject": make_settings_projects_list(activeProject, projectsList)})
	} else {
		additional = ""
		return c.Render("settings",
			fiber.Map{
				"User":          "Текущий пользователь: " + username,
				"Additional":    additional,
				"Versions":      "<option disabled selected>Выберите версию</option>",
				"HostList":      "<option disabled selected>Выберите хост</option>",
				"ActiveProject": make_settings_projects_list(activeProject, projectsList)})
	}
}

func getAdminSubscription(c *fiber.Ctx) error {
	logrus.Debug("getAdminSubscription")
	value := c.Locals("user")
	claims, ok := (value).(*jwt.MapClaims)
	username, ok := (*claims)["username"].(string)
	if !ok {
		// Обработка ошибки преобразования
		logrus.Error("getSettings: username conversion failed")
		return fmt.Errorf("username conversion failed")
	}

	subs, telNumber := get_subs_by_username(username)
	fmt.Println(subs)
	if subs != nil {
		return c.Render("adminSubscription",
			fiber.Map{"SelectProjects": checkbox_for_projects(subs), "TelNumber": telNumber})
	} else {
		return c.Render("adminSubscription",
			fiber.Map{"SelectProjects": checkbox_all_projects(), "TelNumber": telNumber})
	}
}

func getCurrentTests(c *fiber.Ctx) error {
	logrus.Debug("getCurrentTests")
	res := GetCurrentTests(c)
	return c.Render("current_tests",
		fiber.Map{"CurrentTests": res})
}

func testView(c *fiber.Ctx) error {
	return c.Render("scenario_generator",
		fiber.Map{"CurrentTests": "res"})
}
