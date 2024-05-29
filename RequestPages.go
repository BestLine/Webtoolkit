package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"time"
)

//Здесь находятся запросы которые реализуют отрисовку страниц в браузере

func startPage(c *fiber.Ctx) error {
	logrus.Debug("getStartPage")
	return c.Render("main",
		fiber.Map{"Title": "Мотай вниз", "Description": "Где то тут юзаются шаблоны."})
}

func getMainPage(c *fiber.Ctx) error {
	logrus.Debug("getMainPage")
	return c.Render("index", fiber.Map{})
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

func getTests(c *fiber.Ctx) error {
	// реализация таблицы тестов
	logrus.Debug("getTests")
	url := "/beeload/tests?"

	res := sendRequest(c, "Get", url)
	//res := `[{"status":777,"state":"bb80f12f-8724-4c4e-adc9-df23ca64c338","test_data":{"conf_id":0,"application":"","bucket":"","delivery":"","type":"","time_start":0,"time_end":0},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi/scenario/uapi_redis","count":1,"resource":"","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6479"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"890"},{"key":"threads_by_step","value":"89"},{"key":"step_duration","value":"600"},{"key":"step_rump_up","value":"60"},{"key":"rpm_by_thread","value":"600"},{"key":"test_type","value":"MaxPerf_env"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test"},{"key":"version","value":"1.0"}],"testplan":"my-scenario.jmx"}},{"status":777,"state":"44fe9d78-5022-4cf6-8c14-6ba0846cfe60","test_data":{"conf_id":0,"application":"","bucket":"","delivery":"","type":"","time_start":0,"time_end":0},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi/scenario/uapi_redis","count":1,"resource":"","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6479"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"890"},{"key":"threads_by_step","value":"89"},{"key":"step_duration","value":"600"},{"key":"step_rump_up","value":"60"},{"key":"rpm_by_thread","value":"600"},{"key":"test_type","value":"MaxPerf_env"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test"},{"key":"version","value":"1.0"}],"testplan":"my-scenario.jmx"}},{"status":777,"state":"e6d56a1d-7346-495a-8738-984e00ea3f12","test_data":{"conf_id":0,"application":"","bucket":"","delivery":"","type":"","time_start":0,"time_end":0},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi/scenario/uapi_redis","count":1,"resource":"","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6479"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"890"},{"key":"threads_by_step","value":"89"},{"key":"step_duration","value":"600"},{"key":"step_rump_up","value":"60"},{"key":"rpm_by_thread","value":"600"},{"key":"test_type","value":"MaxPerf_env"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test"},{"key":"version","value":"1.0"}],"testplan":"my-scenario.jmx"}},{"status":777,"state":"074464b3-42b1-4db2-8cee-d8fa001ade07","test_data":{"conf_id":0,"application":"","bucket":"","delivery":"","type":"","time_start":0,"time_end":0},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi/scenario/uapi_redis","count":10,"resource":"cpu8ram16","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6479"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"890"},{"key":"threads_by_step","value":"890"},{"key":"step_duration","value":"3600"},{"key":"step_rump_up","value":"600"},{"key":"rpm_by_thread","value":"600"},{"key":"test_type","value":"MaxPerf_env"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test"},{"key":"version","value":"1.0"}],"testplan":"my-scenario.jmx"}},{"status":2,"state":"0cb49070-977f-4114-b609-0b117a8560f1","test_data":{"conf_id":0,"application":"2024-05-21 19-48 Redis_test MaxPerf_env","bucket":"UAPI_Redis","delivery":"1.0","type":"MaxPerf_env","time_start":1716310102,"time_end":1716314480},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi/scenario/uapi_redis","count":10,"resource":"cpu8ram16","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6479"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"890"},{"key":"threads_by_step","value":"890"},{"key":"step_duration","value":"3600"},{"key":"step_rump_up","value":"600"},{"key":"rpm_by_thread","value":"600"},{"key":"test_type","value":"MaxPerf_env"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test"},{"key":"version","value":"1.0"}],"testplan":"my-scenario.jmx"}},{"status":0,"state":"0658432e-0151-4fa4-a50d-8351660c1100","test_data":{"conf_id":0,"application":"","bucket":"","delivery":"","type":"","time_start":0,"time_end":0},"request_data":{"gitlab":"","count":0,"resource":"","data":null,"testplan":""}},{"status":2,"state":"4cd23e46-1d97-47f9-9bb7-c968543979e2","test_data":{"conf_id":0,"application":"2024-05-21 18-09 Redis_test MaxPerf_env","bucket":"UAPI_Redis","delivery":"1.0","type":"MaxPerf_env","time_start":1716304201,"time_end":1716306421},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi_redis/scenario/uapi_redis","count":10,"resource":"cpu8ram16","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6479"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"890"},{"key":"threads_by_step","value":"89"},{"key":"step_duration","value":"600"},{"key":"step_rump_up","value":"60"},{"key":"rpm_by_thread","value":"600"},{"key":"test_type","value":"MaxPerf_env"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test"},{"key":"version","value":"1.0"}],"testplan":"my-scenario.jmx"}},{"status":2,"state":"54307d30-6b74-4254-8e66-ff131ed481e9","test_data":{"conf_id":0,"application":"2024-05-21 17-57 Redis_test MaxPerf_env","bucket":"UAPI_Redis","delivery":"1.0","type":"MaxPerf_env","time_start":1716303453,"time_end":1716306408},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi_redis/scenario/uapi_redis","count":10,"resource":"cpu8ram16","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6479"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"890"},{"key":"threads_by_step","value":"89"},{"key":"step_duration","value":"600"},{"key":"step_rump_up","value":"60"},{"key":"rpm_by_thread","value":"600"},{"key":"test_type","value":"MaxPerf_env"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test"},{"key":"version","value":"1.0"}],"testplan":"my-scenario.jmx"}},{"status":2,"state":"3da9e001-b4d5-4f19-bb87-2a822ff71780","test_data":{"conf_id":0,"application":"2024-05-21 17-40 Redis_test MaxPerf_env","bucket":"UAPI_Redis","delivery":"1.0","type":"MaxPerf_env","time_start":1716302430,"time_end":1716306395},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi_redis/scenario/uapi_redis","count":10,"resource":"cpu8ram16","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6479"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"890"},{"key":"threads_by_step","value":"89"},{"key":"step_duration","value":"600"},{"key":"step_rump_up","value":"60"},{"key":"rpm_by_thread","value":"600"},{"key":"test_type","value":"MaxPerf_env"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test"},{"key":"version","value":"1.0"}],"testplan":"my-scenario.jmx"}},{"status":777,"state":"2751c644-7ee8-4e44-bc36-d58c414d7fe7","test_data":{"conf_id":0,"application":"","bucket":"","delivery":"","type":"","time_start":0,"time_end":0},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi/scenario/uapi_redis","count":10,"resource":"cpu8ram16","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6479"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"890"},{"key":"threads_by_step","value":"89"},{"key":"step_duration","value":"600"},{"key":"step_rump_up","value":"60"},{"key":"rpm_by_thread","value":"300"},{"key":"test_type","value":"MaxPerf"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test"},{"key":"version","value":"1"}],"testplan":"my-scenario.jmx"}},{"status":2,"state":"6ffd51da-63e2-4653-a691-606afca5142c","test_data":{"conf_id":0,"application":"2024-05-21 17-02 Redis_test MaxPerf_env","bucket":"UAPI_Redis","delivery":"1.0","type":"MaxPerf_env","time_start":1716300123,"time_end":1716305423},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi_redis/scenario/uapi_redis","count":1,"resource":"cpu4ram8","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6479"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"890"},{"key":"threads_by_step","value":"89"},{"key":"step_duration","value":"600"},{"key":"step_rump_up","value":"60"},{"key":"rpm_by_thread","value":"600"},{"key":"test_type","value":"MaxPerf_env"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test"},{"key":"version","value":"1.0"}],"testplan":"my-scenario.jmx"}},{"status":2,"state":"3af37648-476f-4e1a-969d-6c7bcae6f7db","test_data":{"conf_id":0,"application":"2024-05-21 17-01 Redis_test MaxPerf_env","bucket":"UAPI_Redis","delivery":"1.0","type":"MaxPerf_env","time_start":1716300108,"time_end":1716305413},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi_redis/scenario/uapi_redis","count":1,"resource":"cpu4ram8","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6479"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"890"},{"key":"threads_by_step","value":"89"},{"key":"step_duration","value":"600"},{"key":"step_rump_up","value":"60"},{"key":"rpm_by_thread","value":"600"},{"key":"test_type","value":"MaxPerf_env"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test"},{"key":"version","value":"1.0"}],"testplan":"my-scenario.jmx"}},{"status":2,"state":"03eb5990-1a49-47ca-8f2f-442f385ea00f","test_data":{"conf_id":0,"application":"2024-05-21 15-26 Redis_test MaxPerf","bucket":"UAPI_Redis","delivery":"1","type":"MaxPerf","time_start":1716294384,"time_end":1716301161},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi/scenario/uapi_redis","count":10,"resource":"cpu8ram16","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6479"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"890"},{"key":"threads_by_step","value":"89"},{"key":"step_duration","value":"600"},{"key":"step_rump_up","value":"60"},{"key":"rpm_by_thread","value":"300"},{"key":"test_type","value":"MaxPerf"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test"},{"key":"version","value":"1"}],"testplan":"my-scenario.jmx"}},{"status":777,"state":"a660ae7e-f906-4f8f-adcf-e1148e832b81","test_data":{"conf_id":0,"application":"","bucket":"","delivery":"","type":"","time_start":0,"time_end":0},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi/scenario/uapi_redis","count":20,"resource":"cpu4ram8","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6479"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"450"},{"key":"threads_by_step","value":"45"},{"key":"step_duration","value":"600"},{"key":"step_rump_up","value":"60"},{"key":"rpm_by_thread","value":"300"},{"key":"test_type","value":"MaxPerf"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test"},{"key":"version","value":"1"}],"testplan":"my-scenario.jmx"}},{"status":0,"state":"646cf27f-bb47-4078-a628-8f43f0399ecf","test_data":{"conf_id":0,"application":"","bucket":"","delivery":"","type":"","time_start":0,"time_end":0},"request_data":{"gitlab":"","count":0,"resource":"","data":null,"testplan":""}},{"status":2,"state":"a6d5cda5-26ba-4c1f-b218-991b47916479","test_data":{"conf_id":0,"application":"2024-05-21 14-37 Redis_test_env MaxPerf_env","bucket":"UAPI_Redis","delivery":"0.5","type":"MaxPerf_env","time_start":1716291427,"time_end":1716291580},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi/scenario/uapi_redis","count":2,"resource":"cpu4ram8","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6379"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"100"},{"key":"threads_by_step","value":"50"},{"key":"step_duration","value":"60"},{"key":"step_rump_up","value":"6"},{"key":"rpm_by_thread","value":"60"},{"key":"test_type","value":"MaxPerf_env"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test_env"},{"key":"version","value":"0.5"}],"testplan":"my-scenario.jmx"}},{"status":2,"state":"e200028c-0171-4420-b0cb-0c33de1dfb58","test_data":{"conf_id":0,"application":"2024-05-21 14-36 Redis_test_env MaxPerf_env","bucket":"UAPI_Redis","delivery":"0.5","type":"MaxPerf_env","time_start":1716291363,"time_end":1716291515},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi/scenario/uapi_redis","count":1,"resource":"cpu4ram8","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6379"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"100"},{"key":"threads_by_step","value":"50"},{"key":"step_duration","value":"60"},{"key":"step_rump_up","value":"6"},{"key":"rpm_by_thread","value":"60"},{"key":"test_type","value":"MaxPerf_env"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test_env"},{"key":"version","value":"0.5"}],"testplan":"my-scenario.jmx"}},{"status":2,"state":"6aec2cb5-cc41-4789-9d08-06e45d07bb4a","test_data":{"conf_id":0,"application":"2024-05-21 14-26 Redis_test MaxPerf_env","bucket":"UAPI_Redis","delivery":"1.0","type":"MaxPerf_env","time_start":1716290802,"time_end":1716292312},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi_redis/scenario/uapi_redis","count":3,"resource":"cpu4ram8","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6479"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"890"},{"key":"threads_by_step","value":"89"},{"key":"step_duration","value":"600"},{"key":"step_rump_up","value":"60"},{"key":"rpm_by_thread","value":"600"},{"key":"test_type","value":"MaxPerf_env"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test"},{"key":"version","value":"1.0"}],"testplan":"my-scenario.jmx"}},{"status":777,"state":"1fe787a8-a04c-491b-9c6c-809f33379b0c","test_data":{"conf_id":0,"application":"","bucket":"","delivery":"","type":"","time_start":0,"time_end":0},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi/scenario/uapi_redis","count":1,"resource":"","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6379"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"100"},{"key":"threads_by_step","value":"50"},{"key":"step_duration","value":"60"},{"key":"step_rump_up","value":"6"},{"key":"rpm_by_thread","value":"60"},{"key":"test_type","value":"MaxPerf_env"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test_env"},{"key":"version","value":"0.5"}],"testplan":"my-scenario.jmx"}},{"status":2,"state":"ed6d9326-8ca7-469d-a2f4-4dc95d73435e","test_data":{"conf_id":0,"application":"2024-05-21 11-06 Redis_test MaxPerf_env","bucket":"UAPI_Redis","delivery":"1.0","type":"MaxPerf_env","time_start":1716289619,"time_end":1716290219},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi_redis/scenario/uapi_redis","count":3,"resource":"cpu4ram8","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6479"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"890"},{"key":"threads_by_step","value":"89"},{"key":"step_duration","value":"600"},{"key":"step_rump_up","value":"60"},{"key":"rpm_by_thread","value":"600"},{"key":"test_type","value":"MaxPerf_env"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test"},{"key":"version","value":"1.0"}],"testplan":"my-scenario.jmx"}},{"status":2,"state":"","test_data":{"conf_id":0,"application":"2024-05-21 10-56 Redis_test MaxPerf_env","bucket":"UAPI_Redis","delivery":"1.0","type":"MaxPerf_env","time_start":1716288988,"time_end":1716292418},"request_data":{"gitlab":"","count":0,"resource":"","data":null,"testplan":""}},{"status":2,"state":"9ee4df17-4d81-4cae-8317-4a882c889341","test_data":{"conf_id":0,"application":"2024-05-21 09-22 Redis test MaxPerf_gui","bucket":"UAPI_Redis","delivery":"0.1","type":"MaxPerf","time_start":1716283347,"time_end":1716290124},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi_redis/scenario/uapi_redis","count":2,"resource":"cpu4ram8","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6479"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"890"},{"key":"threads_by_step","value":"89"},{"key":"step_duration","value":"600"},{"key":"step_rump_up","value":"60"},{"key":"rpm_by_thread","value":"600"},{"key":"test_type","value":"MaxPerf_env"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test"},{"key":"version","value":"1.0"}],"testplan":"my-scenario.jmx"}},{"status":777,"state":"ce081066-523b-41cd-83bd-6759682067e7","test_data":{"conf_id":0,"application":"","bucket":"","delivery":"","type":"","time_start":0,"time_end":0},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi_redis/scenario/uapi_redis","count":2,"resource":"cpu4ram8","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6479"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"890"},{"key":"threads_by_step","value":"89"},{"key":"step_duration","value":"600"},{"key":"step_rump_up","value":"60"},{"key":"rpm_by_thread","value":"600"},{"key":"test_type","value":"MaxPerf_env"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test"},{"key":"version","value":"1.0"}],"testplan":"my-scenario.jmx"}},{"status":2,"state":"208c0bf1-4858-41bc-ad20-87c6bf92423b","test_data":{"conf_id":0,"application":"2024-05-21 08-30 Redis test MaxPerf_gui","bucket":"UAPI_Redis","delivery":"0.1","type":"MaxPerf","time_start":1716280213,"time_end":1716281468},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi_redis/scenario/uapi_redis","count":10,"resource":"cpu4ram8","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6479"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"890"},{"key":"threads_by_step","value":"89"},{"key":"step_duration","value":"600"},{"key":"step_rump_up","value":"60"},{"key":"rpm_by_thread","value":"600"},{"key":"test_type","value":"MaxPerf_env"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test"},{"key":"version","value":"1.0"}],"testplan":"my-scenario.jmx"}},{"status":2,"state":"1120ab16-8c1b-47b2-9f4e-43d7b36c3bcf","test_data":{"conf_id":0,"application":"2024-05-21 08-11 Redis_test MaxPerf_env","bucket":"UAPI_Redis","delivery":"1.0","type":"MaxPerf_env","time_start":1716279067,"time_end":1716281492},"request_data":{"gitlab":"https://git.vimpelcom.ru/products/beeload/uapi_redis/scenario/uapi_redis","count":1,"resource":"cpu4ram8","data":[{"key":"host","value":"10.28.43.6"},{"key":"port","value":"6479"},{"key":"password","value":"GJc1^xRn50+EtB=T"},{"key":"database","value":"14"},{"key":"timeout","value":"2000"},{"key":"total_threads","value":"890"},{"key":"threads_by_step","value":"89"},{"key":"step_duration","value":"600"},{"key":"step_rump_up","value":"60"},{"key":"rpm_by_thread","value":"600"},{"key":"test_type","value":"MaxPerf_env"},{"key":"influx_db","value":"UAPI_Redis"},{"key":"test_name","value":"Redis_test"},{"key":"version","value":"1.0"}],"testplan":"my-scenario.jmx"}}]`
	var testStatuses []TestStatus
	var html string
	err := json.Unmarshal([]byte(res), &testStatuses)
	if err != nil {
		logrus.Error("Error parsing JSON: %v", err)
		return c.Render("show_tests",
			fiber.Map{"Table": html})
	}

	for _, status := range testStatuses {
		if status.Status == 777 {
			html += fmt.Sprintf("<details class=\"test-details\">"+
				"<summary>Тест запустился с ошибкой</summary>"+
				"%s</details>\n", status.RequestData.Gitlab)
		} else if status.Status == 0 {
			html += fmt.Sprintf("<details class=\"test-details\">" +
				"<summary>Тест был запущен из другого источника</summary>" +
				"</details>\n")
		} else {
			t1 := time.Unix(status.TestData.TimeStart, 0)
			t2 := time.Unix(status.TestData.TimeEnd, 0)
			html += fmt.Sprintf("<details class=\"test-details\">"+
				"<summary>"+
				"<table class=\"tests\">"+
				"<tbody>"+
				"<tr>"+
				"<td>%s</td>"+
				"<td>%s</td>"+
				"<td>%s</td>"+
				"<td>%s</td>"+
				"<td>%d</td>"+
				"</tr>"+
				"</tbody>"+
				"</table>"+
				"</summary>", status.TestData.Bucket, t1.Format(time.RFC3339), status.TestData.Type, status.TestData.Application, status.Status)
			i := 0
			html += "<table class=\"tests\">"
			html += "<tbody>"
			html += "<tr>"
			html += fmt.Sprintf("<td>Type = %s</td>", status.TestData.Type)
			html += fmt.Sprintf("<td>TimeStart = %s</td>", t1.Format(time.RFC3339))
			html += fmt.Sprintf("<td>TimeEnd = %s</td>", t2.Format(time.RFC3339))
			html += fmt.Sprintf("<td>Bucket = %s</td>", status.TestData.Bucket)
			html += fmt.Sprintf("<td>Application = %s</td>", status.TestData.Application)
			html += "</tr>"
			html += "<tr>"
			html += fmt.Sprintf("<td>ConfID = %d</td>", status.TestData.ConfID)
			html += fmt.Sprintf("<td>Delivery = %s</td>", status.TestData.Delivery)
			html += "</tr>"
			for _, item := range status.RequestData.Data {
				//println("tag: ", item.Key, "  -  ", item.Value)
				i += 1
				if i%5 == 0 {
					html += "<tr>"
				}
				html += fmt.Sprintf("<td>%s = %s</td>", item.Key, item.Value)
				if i%5 == 0 {
					html += "</tr>"
					i = 0
				}
			}
			html += "</table>"
			html += "</tbody>"
			html += "<div class=\"tableButtonsBlock\">"

			html += "<button class=\"tableButton\">Restart</button>"
			html += "<button class=\"tableButton\">Stop</button>"
			html += "<button class=\"tableButton\" onclick=\"openNewTab('" + status.TestData.Bucket + "','" + status.TestData.Application + "')\">Grafana</button>"
			html += "</div>"
			html += "</details>\n"
		}
	}
	logrus.Debug("getTests html: ", html)
	return c.Render("show_tests",
		fiber.Map{"Table": html})
}

func testView(c *fiber.Ctx) error {
	return c.Render("scenario_generator",
		fiber.Map{"CurrentTests": "res"})
}
