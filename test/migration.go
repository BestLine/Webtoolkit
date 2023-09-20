package ''
import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func createBucket(c *fiber.Ctx) error {
	fmt.Println("CreateBucket!")
	//var req createBucketApi
	//c.Body()
	//
	//json.Unmarshal(c.Body(), &req)
	//createBucketPost(req.Bucket, req.Host)
	return nil
}
func addDataForBeeLoad(c *fiber.Ctx) error {
	//var req addconfluence
	//c.Body()
	//json.Unmarshal(c.Body(), &req)
	//pgAddInfo(req.Bucket, req.ID)
	return nil
}

func compareData(c *fiber.Ctx) error {
	//var req dataCompareRest
	//c.Body()
	//json.Unmarshal(c.Body(), &req)
	//fmt.Println(req)
	//confTestReport(req.Application, req.Bucket, req.ApplicationC)
	return nil
}

func compareRelease(c *fiber.Ctx) error {
	//var req dataCompareRest
	//c.Body()
	//json.Unmarshal(c.Body(), &req)
	//fmt.Println(req)
	//confTestReport(req.Application, req.Bucket, req.ApplicationC)
	return nil
}

func GetTableDataCurrentTests() [][]string {
	res := [][]string{
		{"IDP", "Jmeter_IDP", "В процессе", "<button class=\"l_btn\">Перезапустить</button>", "<button class=\"l_btn\">Остановить</button>"},
		{"BackCRM", "Jmeter_BackCRM", "Присутствуют ошибки", "<button class=\"l_btn\">Перезапустить</button>", "<button class=\"l_btn\">Остановить</button>"},
	}
	// тут хедеры из которых понять что за таблица
	//	table += "<th>Система</th>\n"
	//	table += "<th>Бакет</th>\n"
	//	table += "<th>Статус</th>\n"
	//	table += "<th>Перезапустить</th>\n"
	//	table += "<th>Остановить</th>\n"
	return res
}

func get_bucket_list() []string {
	// общий запрос на список бакетов
	// понадобится ещё один такой с возможностью фильтра по юзеру
	buckets := []string{ // Тут крч список бакетов которые нужны для разного
		"Bucket_1",
		"Bucket_2",
		"Bucket_3",
	}
	return buckets
}

func get_project_list() []string {
	// общий запрос на список проектов
	// понадобится ещё один такой с возможностью фильтра по юзеру
	buckets := []string{ // Тут крч список бакетов которые нужны для разного
		"Project_1",
		"Project_2",
		"Project_3",
		"Project_4",
		"Project_5",
		"Project_5",
		"Project_5",
	}
	return buckets
}

func get_project_buckets(project string) []string {
	// запрос на список бакетов с фильтром по проекту
	fmt.Println(`Project: `, project)
	buckets := []string{ // Тут крч список бакетов в зависимости от проекта
		"Bucket_1",
		"Bucket_2",
		"Bucket_3",
	}
	return buckets
}

func get_bucket_projects(bucket string) []string {
	// запрос на список проектов с фильтром по бакету
	fmt.Println(`Bucket: `, bucket)
	buckets := []string{ // Тут крч список бакетов в зависимости от проекта
		"Project_1",
		"Project_2",
		"Project_3",
		"Project_4",
		"Project_5",
		"Project_5",
		"Project_5",
	}
	return buckets
}

func GetTableDataReports() [][]string {
	// запрос списка последних тестов по юзеру
	res := [][]string{
		{"IDP", "Jmeter_IDP", "<a href=\"report1\">Ссылка.</a>", "тестовый тест"},
		{"BackCRM", "Jmeter_BackCRM", "<a href=\"report1\">Ссылка.</a>", "тестовый тест"},
		{"BackCRM", "Jmeter_BackCRM", "<a href=\"report1\">Ссылка.</a>", "тестовый тест"},
		{"BackCRM", "Jmeter_BackCRM", "<a href=\"report1\">Ссылка.</a>", "тестовый тест"},
		{"IDP", "Jmeter_IDP", "<a href=\"report1\">Ссылка.</a>", "тестовый тест"},
		{"IDP", "Jmeter_IDP", "<a href=\"report1\">Ссылка.</a>", "тестовый тест"},
		{"Mapic", "Jmeter_Mapic", "<a href=\"report1\">Ссылка.</a>", "тестовый тест"},
		{"Mapic", "Jmeter_Mapic", "<a href=\"report1\">Ссылка.</a>", "тестовый тест"},
		{"Mapic", "Jmeter_Mapic", "<a href=\"report1\">Ссылка.</a>", "тестовый тест"},
	}
	return res
}

func GetTableDataTests() [][]string {
	// история недавно проведённых тестов
	res := [][]string{
		{"IDP", "Jmeter_IDP", "5:56 28.07.2023", "В процессе проведения", "MaxPerf"},
		{"BackCRM", "Jmeter_BackCRM", "3:20 28.07.2023", "Остановлен", "Stability"},
		{"BackCRM", "Jmeter_BackCRM", "0:12 28.07.2023", "Завершён с ошибками", "MaxPerf"},
	}
	return res
}

func GetTableDataStatus() [][]string {
	// запрос по статусам на главную
	res := [][]string{
		{"Генераторы", "5 минут назад", "Штатный режим"},
		{"Автоотчёт", "3 минуты назад", "Штатный режим"},
		{"Jenkins", "243 минуты назад", "Не отвечает"},
	}
	return res
}
