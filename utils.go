package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func GetTableDataReports(c *fiber.Ctx, project string, count int) JSONData {
	logrus.Debug("getCurrentTests")
	var url string
	var data JSONData
	if project != "" {
		url = "/beeload/get/tableDataReports?bucket=" + project + "&count=" + strconv.Itoa(count)
	} else {
		url = "/beeload/get/tableDataReports?&count=10"
	}
	//url := "/beeload/get/tableDataReports"
	logrus.Debug("url: ", url)
	res := sendGet(c, url)
	json.Unmarshal([]byte(res), &data)
	//
	//dataStr := string(res)
	//
	//// Разбиваем на строки
	//rows := strings.Split(dataStr, "},{")
	//
	//// Подготовка для парсинга
	//var data [][]string
	//for _, row := range rows {
	//	// Удаление лишних символов
	//	row = strings.Trim(row, "[{]}")
	//	// Разбиваем на элементы
	//	items := strings.Split(row, ",")
	//	var itemStrings []string
	//	for _, item := range items {
	//		item = strings.Trim(item, "\" ")
	//		itemStrings = append(itemStrings, item)
	//	}
	//	data = append(data, itemStrings)
	//}
	//fmt.Println("DATA:   ", data.Data[0].Application)
	return data
}

func GetTableDataTests(c *fiber.Ctx) [][]string {
	logrus.Debug("GetTableDataTests")
	url := "/beeload/get/tableDataTests"
	res := sendGet(c, url)
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

func get_last_10_reports_table(data JSONData) string {
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

func get_current_tests(data [][]string) string {
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
	res := sendGet(c, url)
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
	res := sendGet(c, url)
	var data []map[string]string
	err := json.Unmarshal([]byte(res), &data)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	// Извлечение значений "Bucket" и создание []string
	var buckets []string
	for _, item := range data {
		if bucket, found := item["Bucket"]; found {
			buckets = append(buckets, bucket)
		}
	}

	// Вывод результатов
	fmt.Println(buckets)
	return buckets
}

func get_project_list(c *fiber.Ctx) []string {
	url := "/beeload/get/projectList"
	res := sendGet(c, url)
	var data []map[string]string
	err := json.Unmarshal([]byte(res), &data)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	// Извлечение значений "Bucket" и создание []string
	var projects []string
	for _, item := range data {
		if project, found := item["Project"]; found {
			projects = append(projects, project)
		}
	}

	// Вывод результатов
	fmt.Println(projects)
	return projects
}

func get_versions_list(c *fiber.Ctx, project string) []string {
	fmt.Println(`Project: `, project)
	url := "/beeload/get/versionList?project=" + project
	res := sendGet(c, url)
	var data []map[string]string
	err := json.Unmarshal([]byte(res), &data)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	var versions []string
	for _, item := range data {
		if version, found := item["Version"]; found {
			versions = append(versions, version)
		}
	}
	fmt.Println(versions)
	return versions
}

func get_project_buckets(c *fiber.Ctx, project string) []string {
	fmt.Println(`Project: `, project)
	url := "/beeload/get/bucketList?project=" + project
	res := sendGet(c, url)
	var data []map[string]string
	err := json.Unmarshal([]byte(res), &data)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	// Извлечение значений "Bucket" и создание []string
	var buckets []string
	for _, item := range data {
		if bucket, found := item["Bucket"]; found {
			buckets = append(buckets, bucket)
		}
	}

	// Вывод результатов
	fmt.Println(buckets)
	return buckets
}

func get_bucket_projects(c *fiber.Ctx, bucket string) []string {
	fmt.Println(`Bucket: `, bucket)
	url := "/beeload/get/bucketList?bucket=" + bucket
	res := sendGet(c, url)
	var data []map[string]string
	err := json.Unmarshal([]byte(res), &data)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	// Извлечение значений "Bucket" и создание []string
	var buckets []string
	for _, item := range data {
		if bucket, found := item["Bucket"]; found {
			buckets = append(buckets, bucket)
		}
	}

	// Вывод результатов
	fmt.Println(buckets)
	return buckets
}

func add_tags(imp_list []string) string {
	var s string
	for i := 0; i < len(imp_list); i++ {
		d := "<option>" + imp_list[i] + "</option>"
		s += d
	}
	return s
}
