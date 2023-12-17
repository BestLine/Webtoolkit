package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func GetTableDataReports(c *fiber.Ctx, project string, count int) TestsTableData {
	logrus.Debug("getCurrentTests")
	var url string
	var data TestsTableData
	if project != "" {
		url = "/beeload/get/tableDataReports?bucket=" + project + "&count=" + strconv.Itoa(count)
	} else {
		url = "/beeload/get/tableDataReports?&count=10"
	}
	logrus.Debug("url: ", url)
	//res := sendGet(c, url)
	res := sendRequest(c, "Get", url)
	json.Unmarshal(res, &data)
	return data
}

func GetTableDataTests(c *fiber.Ctx) CurrentTestsTableData {
	logrus.Debug("GetTableDataTests")
	url := "/beeload/get/tableDataTests"
	//res := sendGet(c, url)
	var data CurrentTestsTableData
	logrus.Debug("url: ", url)
	res := sendRequest(c, "Get", url)
	err := json.Unmarshal(res, &data)
	if err != nil {
		return CurrentTestsTableData{}
	}
	return data
}

func get_last_10_reports_table(data TestsTableData) string {
	table := "<table>\n"
	table += "<thead>\n"
	table += "<tr>\n"
	table += "<th>Имя теста</th>\n"
	table += "<th>Бакет</th>\n"
	table += "<th>Ссылка на отчёт</th>\n"
	table += "</tr>\n"
	table += "</thead>\n"
	table += "<tbody>\n"
	for _, row := range data.Data {
		table += "<tr>\n"
		table += "<td>" + row.Application + "</td>\n"
		table += "<td>" + row.Bucket + "</td>\n"
		table += "<td><a href=" + row.CfURL + ">Ссылка.</a></td>\n"
		table += "</tr>\n"
	}
	table += "</tbody>\n"
	table += "</table>\n"
	return table
}

func get_current_tests(data CurrentTestsTableData) string {
	table := "<table>\n"
	table += "<thead>\n"
	table += "<tr>\n"
	table += "<th>Проект</th>\n"
	table += "<th>Бакет</th>\n"
	table += "<th>Время старта</th>\n"
	table += "<th>Статус теста</th>\n"
	table += "<th>Тип теста</th>\n"
	table += "</tr>\n"
	table += "</thead>\n"
	table += "<tbody>\n"
	for _, row := range data.Data {
		table += "<tr>\n"
		table += "<td>" + row.Project + "</td>\n"
		table += "<td>" + row.Bucket + "</td>\n"
		table += "<td>" + row.StartTime + "</td>\n"
		table += "<td>" + row.Status + "</td>\n"
		table += "<td>" + row.Type + "</td>\n"
		table += "</tr>\n"
	}
	table += "</tbody>\n"
	table += "</table>\n"
	fmt.Println("awd: ", table)
	return table
}

func get_status_table(data [][]string) string {
	table := "<table>\n"
	table += "<thead>\n"
	table += "<tr>\n"
	table += "<th>Система</th>\n"
	table += "<th>Последняя проверка</th>\n"
	table += "<th>Статус</th>\n"
	table += "</tr>\n"
	table += "</thead>\n"
	table += "<tbody>\n"
	for _, row := range data {
		table += "<tr>\n"
		for _, col := range row {
			table += "<td>" + col + "</td>\n"
		}
		table += "</tr>\n"
	}
	table += "</tbody>\n"
	table += "</table>\n"
	return table
}

func GetTableDataStatus(c *fiber.Ctx) [][]string {
	logrus.Debug("GetTableDataStatus")
	url := "/beeload/get/tableDataStatus"
	//res := sendGet(c, url)
	res := sendRequest(c, "Get", url)
	//fmt.Println(res)
	// расшифровка ответа
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
	return data
}

func get_test_table(data [][]string) string {
	table := "<table>\n"
	table += "<thead>\n"
	table += "<tr>\n"
	table += "<th>Система</th>\n"
	table += "<th>Бакет</th>\n"
	table += "<th>Статус</th>\n"
	table += "<th>Перезапустить</th>\n"
	table += "<th>Остановить</th>\n"
	table += "</tr>\n"
	table += "</thead>\n"
	table += "<tbody>\n"
	for _, row := range data {
		table += "<tr>\n"
		for _, col := range row {
			table += "<td>" + col + "</td>\n"
		}
		table += "<td><button class=\"l_btn\">Перезапустить</button></td>\n<td><button class=\"l_btn\">Остановить</button></td>"
		table += "</tr>\n"
	}
	table += "</tbody>\n"
	table += "</table>\n"
	return table
}

func make_settings_projects_list(active string, projects []string) string {
	var s string
	if active != "" {
		s += "<option selected>" + active + "</option>"
	}
	for i := 0; i < len(projects); i++ {
		d := "<option>" + projects[i] + "</option>"
		s += d
	}
	return s
}

func get_bucket_list(c *fiber.Ctx) []string {
	url := "/beeload/get/bucketList"
	//res := sendGet(c, url)
	res := sendRequest(c, "Get", url)
	var data []map[string]string
	err := json.Unmarshal([]byte(res), &data)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	return dataToListOfStrings(data, "Bucket")
}

func get_project_list(c *fiber.Ctx) []string {
	//value := c.Locals("user")
	//claims, _ := (value).(*jwt.MapClaims)
	//username, _ := (*claims)["username"].(string)
	//db := getDbConn()
	//activeProject, err := GetUserActiveProject(db, username)
	url := "/beeload/get/projectList"
	//res := sendGet(c, url)
	res := sendRequest(c, "Get", url)
	var data []map[string]string
	err := json.Unmarshal([]byte(res), &data)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return dataToListOfStrings(data, "Project")
}

func get_versions_list(c *fiber.Ctx, project string) []string {
	fmt.Println(`Project: `, project)
	url := "/beeload/get/versionList?project=" + project
	//res := sendGet(c, url)
	res := sendRequest(c, "Get", url)
	var data []map[string]string
	err := json.Unmarshal([]byte(res), &data)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	return dataToListOfStrings(data, "Version")
}

func get_host_list(c *fiber.Ctx, project string) []string {
	fmt.Println(`Project: `, project)
	url := "/beeload/get/hostList?project=" + project
	res := sendRequest(c, "Get", url)
	var data map[string][]string
	err := json.Unmarshal([]byte(res), &data)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	hosts := data["host"]
	for _, host := range hosts {
		fmt.Println(host)
	}
	return hosts
}

func get_project_buckets(c *fiber.Ctx, project string) []string {
	logrus.Debug("get_project_buckets Project: ", project)
	url := "/beeload/get/bucketList?project=" + project
	//res := sendGet(c, url)
	res := sendRequest(c, "Get", url)
	logrus.Debug("URL: ", url)
	var data []map[string]string
	err := json.Unmarshal([]byte(res), &data)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	//TODO: возвращает список всех тестов
	//TODO: добавить в форму создания отчёта
	return dataToListOfStrings(data, "Bucket")
}

func get_bucket_projects(c *fiber.Ctx, bucket string) []string {
	logrus.Debug("get_bucket_projects Bucket: ", bucket)
	url := "/beeload/get/bucketList?bucket=" + bucket
	//res := sendGet(c, url)
	res := sendRequest(c, "Get", url)
	logrus.Debug("URL: ", url)
	var data []map[string]string
	err := json.Unmarshal([]byte(res), &data)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	return dataToListOfStrings(data, "Bucket")
}

func dataToListOfStrings(data []map[string]string, name string) []string {
	var res []string
	for _, item := range data {
		if thing, found := item[name]; found {
			res = append(res, thing)
		}
	}
	return res
}

func add_tags(imp_list []string) string {
	var s string
	for i := 0; i < len(imp_list); i++ {
		d := "<option>" + imp_list[i] + "</option>"
		s += d
	}
	return s
}

func RespToByteReader(response *http.Response) []byte {
	var buf []byte
	const chunkSize = 1024 // Размер чанка для чтения данных

	for {
		chunk := make([]byte, chunkSize)
		n, err := response.Body.Read(chunk)
		if err != nil && err != io.EOF {
			logrus.Error(err)
			return nil
		}
		if n == 0 {
			break
		}
		buf = append(buf, chunk[:n]...)
	}
	return buf
}

func select_all_users() string {
	logrus.Debug("select_all_users")
	users, _ := GetAllUsers()
	res := "<select name=\"user\" class=\"user\" required><option value=\"\" disabled selected>Выберите пользователя</option>"
	for _, userName := range users {
		res += "<option value=\"" + userName + "\">" + userName + "</option>"
	}
	res += "</select>"
	return res
}

func select_all_projects() string {
	logrus.Debug("select_all_users")
	projects, _ := GetAllProjects()
	res := "<select class=\"project_options\" name=\"project\" required><option value=\"\" disabled selected>Выберите проект</option>"
	for _, projectName := range projects {
		res += "<option value=\"" + projectName + "\">" + projectName + "</option>"
	}
	res += "</select>"
	return res
}

func checkbox_all_projects() string {
	projects, _ := GetAllProjects()
	res := "<div id=\"projectList\">"
	for _, projectName := range projects {
		res += "<input type=\"checkbox\" id=\"" + projectName + "\" name=\"projects\" value=\"" + projectName + "\">\n"
		res += "<label for=\"" + projectName + "\">" + projectName + "</label><br>\n"
	}
	res += "</div>"
	return res
}

func make_user_project_list() string {
	logrus.Debug("make_user_project_list")
	projects, _ := GetAllProjects()
	users, _ := GetAllUsers()
	res := "<select name=\"user\" required><option value=\"\" disabled selected>Выберите пользователя</option>"
	for _, userName := range users {
		res += "<option value=\"" + userName + "\">" + userName + "</option>"
	}
	res += "</select><select name=\"project\" required>\n\t       <option value=\"\" disabled selected>Выберите проект</option>"
	for _, projectName := range projects {
		res += "<option value=\"" + projectName + "\">" + projectName + "</option>"
	}
	res += "</select>"
	return res
}

func make_generators_list() string {
	logrus.Debug("make_generators_list")
	res := "<select name=\"generator\" class=\"genType\" required><option value=\"\" disabled selected>Выберите Генератор</option>"
	generators := "<option value=\"4ram2cpu\">4ram2cpu</option>" +
		"<option value=\"8ram4cpu\">8ram4cpu</option>" +
		"<option value=\"16ram8cpu\">16ram8cpu</option>"
	res += generators
	res += "</select>"
	return res
}

func syncBuckets() string {
	logrus.Debug("syncBuckets")
	url := "/bucket"
	targetURL := viper.GetString("backend.pure_host") + url
	res, err := http.Get(targetURL)
	if err != nil {
		logrus.Error("syncBuckets Error:", err)
		return ""
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		logrus.Error("SyncRequest ", targetURL, " code: ", res.StatusCode)
	}
	var responseStruct struct {
		Bucket []string `json:"bucket"`
	}
	err = json.Unmarshal(RespToByteReader(res), &responseStruct)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return ""
	}
	stringArray := responseStruct.Bucket
	fmt.Println("Array of Strings:", stringArray)
	err = SyncProjects(stringArray)
	if err != nil {
		fmt.Println("Error SyncProjectsSQL:", err)
		return ""
	}
	fmt.Println("syncBuckets Ответ: ", stringArray)
	logrus.Debug("syncBuckets Ответ: ", stringArray)
	return ""
}
