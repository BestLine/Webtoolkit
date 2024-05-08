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

func GetCurrentTests(c *fiber.Ctx) string {
	var err error
	var response *http.Response
	var testData TestData
	var result [][]string
	logrus.Debug("GetCurrentTests")
	url := "http://ms-loadrtst026:8001/test/live"
	response, err = http.Get(url)
	if err != nil {
		logrus.Error("Get sending error: ", err)
		c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		return ""
	}

	byte_reader := RespToByteReader(response)
	fmt.Println("GetCurrentTests responce: ", string(byte_reader))
	defer response.Body.Close()
	logrus.Debug("GetCurrentTests targetURL: ", url)
	logrus.Debug("GetCurrentTests response: ", string(byte_reader))

	if response.StatusCode != 200 {
		logrus.Error("GetCurrentTests responce code: ", response.StatusCode)
		c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		return ""
	}

	err = json.Unmarshal(byte_reader, &testData)
	if err != nil {
		logrus.Error("Error decoding JSON:", err)
		return ""
	}

	for _, test := range testData.Tests {
		result = append(result, []string{test.Application, test.Bucket, test.State})
	}

	return get_test_table(result)
}

func get_test_table(data [][]string) string {
	table := "<table>\n"
	table += "<thead>\n"
	table += "<tr>\n"
	table += "<th>Тест</th>\n"
	table += "<th>Проект</th>\n"
	table += "<th>Графана</th>\n"
	table += "<th>Перезапустить</th>\n"
	table += "<th>Остановить</th>\n"
	table += "</tr>\n"
	table += "</thead>\n"
	table += "<tbody>\n"
	for _, row := range data {
		table += "<tr>\n"
		i := 0
		for _, col := range row {
			i += 1
			if i < 3 {
				table += "<td>" + col + "</td>\n"
			}
		}
		table +=
			"<td><button class=\"l_btn\" onclick=\"openNewTab('" + row[1] + "','" + row[0] + "')\">Открыть мониторинг</button></td>\n" +
				"<td><button class=\"l_btn\">Перезапустить</button></td>\n" +
				"<td><button class=\"l_btn\" onclick=\"window.location.href='http://ms-loadrtst038:9999/destroy?state=" + row[2] + "'\">Остановить</button></td>"
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
	err := json.Unmarshal(res, &data)
	fmt.Println(string(res))
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

func checkbox_for_projects(projects []string) string {
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
	generators := "<option value=\"cpu2ram4\">cpu2ram4</option>" +
		"<option value=\"cpu4ram8\">cpu4ram8</option>" +
		"<option value=\"cpu8ram16\">cpu8ram16</option>"
	//последнее изменение списка генов 22.12.2023
	res += generators
	res += "</select>"
	return res
}

func syncBuckets() string {
	logrus.Debug("syncBuckets")
	url := "/bucket"
	//projects := []string{"Project1", "Project3", "Project4", "Project5", "Project5", "Project5", "Project6", "Project7", "Project8"}
	//err := test1(projects)
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

func get_subs_by_username(UserName string) ([]string, string) {
	var subs []string
	var telNumber = ""
	var err error
	if UserName != "" {
		telNumber, err = get_telnumber_by_username(UserName)
		subs, err = get_subs_by_telnumber(telNumber)
	} else {
		logrus.Error("get_subs_by_tel ERROR username: ", UserName)
		return nil, ""
	}
	if err != nil {
		//fmt.Println("Error get_subs_by_username:", err)
		logrus.Error("Error get_subs_by_username: ", err)
		return nil, ""
	}
	return subs, telNumber
}
