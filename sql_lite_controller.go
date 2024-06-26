package main

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

func DBInit() {
	var err error
	// Инициализация пула
	db, err = sql.Open("sqlite", "./mydatabase.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}

	// Установка максимального числа открытых соединений в пуле
	db.SetMaxOpenConns(20)

	// Установка максимального числа соединений в пуле, которые могут быть использованы одновременно
	db.SetMaxIdleConns(10)
}

func createSQLiteDB() error {
	var err error
	// Создаем таблицы (если их нет)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			password TEXT NOT NULL,
			email TEXT
		)
	`)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		)
	`)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user_roles (
			user_id INTEGER,
			role_id INTEGER,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(role_id) REFERENCES roles(id)
		)
	`)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user_projects (
			user_id INTEGER,
			project_id INTEGER,
			active INTEGER,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(project_id) REFERENCES projects(id)
		)
	`)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_name TEXT NOT NULL UNIQUE
		)
	`)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS project_versions (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id INTEGER NOT NULL,
			version TEXT NOT NULL UNIQUE
		)
	`)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS project_methodics (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id INTEGER NOT NULL,
			version TEXT NOT NULL,
			methodic_conf_id TEXT NOT NULL UNIQUE
		)
	`)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS project_root_page (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id INTEGER NOT NULL,
			page_id INTEGER NOT NULL UNIQUE
		)
	`)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS project_buckets (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id INTEGER NOT NULL,
			bucket_url TEXT NOT NULL,
			bucket_name TEXT NOT NULL UNIQUE 
		)
	`)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user_subscriptions (
		    user_id INTEGER NOT NULL,
		    project_id INTEGER NOT NULL,
		    FOREIGN KEY (user_id) REFERENCES users(id),
		    FOREIGN KEY (project_id) REFERENCES projects(id),
		    PRIMARY KEY (user_id, project_id)
		)
	`)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user_phones (
			user_id INTEGER NOT NULL,
			phone_number VARCHAR(15) NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id),
			PRIMARY KEY (user_id, phone_number)
		)
	`)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS subscription_lists (
			phone_number VARCHAR(15) NOT NULL,
			project_id INTEGER NOT NULL,
			FOREIGN KEY (phone_number) REFERENCES user_phones(phone_number),
			FOREIGN KEY (project_id) REFERENCES projects(id),
			PRIMARY KEY (phone_number, project_id)
		)
	`)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS monitoring_urls (
		    project_name VARCHAR(100) NOT NULL,
		    url VARCHAR(300) NOT NULL
		)
	`)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'admin'").Scan(&count)
	if count == 0 {
		_, err = db.Exec(`
    		INSERT OR IGNORE INTO users (username, password) VALUES ('admin', '123')
		`)

		_, err = db.Exec(`
    		INSERT OR IGNORE INTO roles (name) VALUES ('admin')
		`)

		_, err = db.Exec(`
    		INSERT OR IGNORE INTO user_roles (user_id, role_id)
    		VALUES (
        		(SELECT id FROM users WHERE username = 'admin'),
        		(SELECT id FROM roles WHERE name = 'admin')
    		)
		`)
	}

	if err != nil {
		return err
	}
	return nil
}

func registerUser(username, password, email string, roles []string) error {
	// Проверяем, что пользователь с таким именем пользователя еще не зарегистрирован
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		logrus.Error("registerUser: Проверка на существование")
		return err
	}

	if count > 0 {
		return fmt.Errorf("Пользователь с именем пользователя %s уже существует", username)
	}

	// Сохраняем данные нового пользователя в таблице users
	result, err := db.Exec("INSERT INTO users (username, password, email) VALUES (?, ?, ?)", username, password, email)
	if err != nil {
		logrus.Error("registerUser: Сохранение пользователя")
		return err
	}

	// Получаем ID нового пользователя
	userID, err := result.LastInsertId()
	if err != nil {
		logrus.Error("registerUser: Получение ID")
		return err
	}

	// Получаем ID ролей из таблицы roles по их именам
	var roleIDs []int
	for _, roleName := range roles {
		var roleID int
		err := db.QueryRow("SELECT id FROM roles WHERE name = ?", roleName).Scan(&roleID)
		if err != nil {
			logrus.Error("registerUser: Получение ID роли")
			return err
		}
		roleIDs = append(roleIDs, roleID)
	}

	// Сохраняем связи пользователей и их ролей в таблице user_roles
	for _, roleID := range roleIDs {
		_, err := db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", userID, roleID)
		if err != nil {
			logrus.Error("registerUser: Сохранение связей")
			return err
		}
	}

	return nil
}

func addRole(role string) error {
	_, err := db.Exec("INSERT INTO roles (name) VALUES (?)", role)
	if err != nil {
		logrus.Error("addRole: Добавление роли")
		return err
	}
	return nil
}

func addMonUrl(project string, url string) error {
	_, err := db.Exec("INSERT INTO monitoring_urls (project_name, url) VALUES (?, ?)", project, url)
	if err != nil {
		logrus.Error("addRole: Добавление роли")
		return err
	}
	return nil
}

func getMonUrls(project string) ([]string, error) {
	query := "SELECT url FROM monitoring_urls WHERE project_name = ?"
	fmt.Println("Project:", project)
	rows, err := db.Query(query, project)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var url string
		err := rows.Scan(&url)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	fmt.Println("URLs:", urls)
	return urls, nil
}

func checkUserCredentials(username, password string) (bool, error) {
	row := db.QueryRow("SELECT username, password FROM users WHERE username = ?", username)

	var storedUsername, storedPassword string
	err := row.Scan(&storedUsername, &storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// Пользователь с указанным логином не найден, возвращаем false
			return false, nil
		}
		// Возникла ошибка при выполнении запроса
		return false, err
	}

	// Проверяем совпадение паролей
	if storedPassword == password {
		// Пароль верен, возвращаем true
		return true, nil
	}

	// Неверный пароль, возвращаем false
	return false, nil
}

func getUserRole(username string) string {
	// Проверяем наличие роли у пользователя лениво
	row := db.QueryRow("select name as role from users left join user_roles ur on users.id = ur.user_id left join roles r on r.id = ur.role_id where username=?", username)
	var storedRole string
	err := row.Scan(&storedRole)
	if err != nil {
		return "Ошибка при выполнении запроса"
	}
	return storedRole
}

func hasUserRole(username, role string) (bool, error) {
	// Проверяем наличие роли у пользователя серьёзно
	//query := "SELECT COUNT(*) FROM users WHERE username = ? AND role = ?"
	query := "SELECT COUNT(*) FROM users left join user_roles on users.id = user_roles.user_id left join roles on user_roles.role_id = roles.id WHERE users.username = ? AND roles.name = ?"
	var count int
	err := db.QueryRow(query, username, role).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetAllUsers() ([]string, error) {
	query := "SELECT username FROM users;"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userNames []string
	for rows.Next() {
		var userName string
		err := rows.Scan(&userName)
		if err != nil {
			return nil, err
		}
		userNames = append(userNames, userName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return userNames, nil
}

func GetUserActiveProject(username string) (string, error) {
	query := "SELECT projects.project_name " +
		"FROM user_projects " +
		"LEFT JOIN projects ON user_projects.project_id = projects.id " +
		"LEFT JOIN users ON users.id = user_projects.user_id " +
		"WHERE users.username = '" + username + "' AND user_projects.active = 1;"

	var name string
	err := db.QueryRow(query).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func GetUserProjects(username string) ([]string, error) {
	query := "SELECT projects.project_name " +
		"FROM user_projects " +
		"LEFT JOIN projects ON user_projects.project_id = projects.id " +
		"LEFT JOIN users ON users.id = user_projects.user_id " +
		"WHERE users.username = ? AND user_projects.active = 0;"

	rows, err := db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projectNames []string
	for rows.Next() {
		var projectName string
		err := rows.Scan(&projectName)
		if err != nil {
			return nil, err
		}
		projectNames = append(projectNames, projectName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projectNames, nil
}

func GetAllProjects() ([]string, error) {
	query := "SELECT project_name FROM projects;"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projectNames []string
	for rows.Next() {
		var projectName string
		err := rows.Scan(&projectName)
		if err != nil {
			return nil, err
		}
		projectNames = append(projectNames, projectName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projectNames, nil
}

func SetActiveProject(username string, projectToActivate string, isAdmin bool) error {
	// Если пользователь является администратором, устанавливаем указанный проект как активный
	if isAdmin {
		// Получаем ID пользователя по его имени
		var userID int
		err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
		if err != nil {
			return err
		}

		// Получаем ID проекта по его имени
		var projectID int
		err = db.QueryRow("SELECT id FROM projects WHERE project_name = ?", projectToActivate).Scan(&projectID)
		if err != nil {
			return err
		}

		// Проверяем, существует ли уже запись о активном проекте для данного пользователя
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM user_projects WHERE user_id = ? AND active = 1", userID).Scan(&count)
		if err != nil {
			return err
		}

		// Если запись уже существует, обновляем ее
		if count > 0 {
			_, err = db.Exec("UPDATE user_projects SET project_id = ? WHERE user_id = ? AND active = 1", projectID, userID)
			if err != nil {
				return err
			}
		} else {
			// Если записи нет, добавляем новую запись
			_, err = db.Exec("INSERT INTO user_projects (user_id, project_id, active) VALUES (?, ?, 1)", userID, projectID)
			if err != nil {
				return err
			}
		}
	} else {
		// Сначала помечаем все проекты пользователя как неактивные
		updateQuery := "UPDATE user_projects " +
			"SET active = 0 " +
			"WHERE user_id = (SELECT id FROM users WHERE username = ?)"
		logrus.Debug("QW_1 = ", updateQuery, " username=", username)
		_, err := db.Exec(updateQuery, username)
		if err != nil {
			return err
			logrus.Debug("ERR1")
		}
		// Для обычного пользователя установите выбранный проект как активный, предварительно деактивировав остальные
		updateAllProjectsQuery := "UPDATE user_projects " +
			"SET active = CASE WHEN project_id = (SELECT id FROM projects WHERE project_name = ?) THEN 1 ELSE 0 END " +
			"WHERE user_id = (SELECT id FROM users WHERE username = ?)"
		logrus.Debug("QW_2 = ", updateAllProjectsQuery, " username=", username, " projectToActivate=", projectToActivate)
		_, err = db.Exec(updateAllProjectsQuery, projectToActivate, username)
	}
	return nil
}

func AddNewProject(ProjectName string) error {
	Query := "INSERT OR IGNORE INTO projects (project_name) VALUES (?);"
	logrus.Debug("QW = ", Query, " ProjectName=", ProjectName)
	_, err := db.Exec(Query, ProjectName)
	if err != nil {
		return err
	}

	return nil
	//TODO: new feature
}

func AddProjectVersion(Version string, ProjectName string) error {
	Query := "INSERT OR IGNORE INTO project_versions (project_id, version) SELECT p.id, ? FROM projects p WHERE p.project_name = ?;"
	logrus.Debug("QW = ", Query, " ProjectName=", ProjectName, " Version=", Version)
	_, err := db.Exec(Query, Version, ProjectName)
	if err != nil {
		return err
	}

	return nil
	//TODO: new feature
}

func GetProjectVersion(projectName string) ([]string, error) {
	rows, err := db.Query("SELECT version FROM project_versions WHERE project_id=(SELECT p.id FROM projects p WHERE p.project_name = ?)", projectName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []string
	for rows.Next() {
		var version string
		err = rows.Scan(&version)
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return versions, nil
}

func AddProjectMethodic(PageId int, Version string, ProjectName string) error {
	Query := "INSERT OR IGNORE INTO project_methodics (project_id, version, methodic_conf_id) SELECT p.id,?, ? FROM projects p WHERE p.project_name = ?;"
	_, err := db.Exec(Query, PageId, Version, ProjectName)
	logrus.Debug("QW = ", Query, " ProjectName=", ProjectName, " PageId=", PageId, " Version=", Version)
	if err != nil {
		return err
	}
	return nil
	//TODO: new feature
}

func AddProjectBucket(BucketName string, BucketUrl string, ProjectName string) error {
	Query := "INSERT OR IGNORE INTO project_buckets (project_id, bucket_url, bucket_name) SELECT p.id, ?, ? FROM projects p WHERE p.project_name = ?;"
	logrus.Debug("QW = ", Query, " ProjectName=", ProjectName, " BucketUrl=", BucketUrl)
	_, err := db.Exec(Query, BucketUrl, BucketName, ProjectName)
	if err != nil {
		return err
	}
	return nil
	//TODO: new feature
}

func AddUserSubscriptions(UserName string, Projects []string) error {
	Query := `DELETE FROM user_subscriptions
		WHERE user_id = (SELECT user_id FROM users WHERE username = ?)
	`
	_, err := db.Exec(Query, UserName)
	for _, project := range Projects {
		Query = `INSERT OR IGNORE INTO user_subscriptions (user_id, project_id)
    		SELECT users.id, projects.id
    		FROM users, projects
    		WHERE users.username = ? AND projects.project_name = ?
		`
		_, err = db.Exec(Query, UserName, project)
		if err != nil {
			return err
		}
	}
	return nil
}

func AddUserToProject(UserName string, Project string) error {
	logrus.Debug("AddUserToProject")
	logrus.Debug("UserName: ", UserName)
	logrus.Debug("Project: ", Project)
	Query := `INSERT INTO user_projects (user_id, project_id, active)
		SELECT users.id, projects.id, 0
		FROM users, projects
		WHERE users.username = ? AND projects.project_name = ?
	`
	_, err := db.Exec(Query, UserName, Project)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func SyncProjects(Projects []string) error {
	logrus.Debug("func SyncProjects started!")
	//TODO: добавить проверку на дубликаты
	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		logrus.Debug("Transaction error: ", err)
		return err
	}
	defer func() {
		if err != nil {
			logrus.Debug("ROLLBACK execution reason: ", err)
			tx.Rollback()
		} else {
			logrus.Debug("COMMIT execution")
			tx.Commit()
		}
	}()

	deleteQuery := `DELETE FROM projects WHERE project_name NOT IN (?)`
	deleteArgs := make([]interface{}, len(Projects))
	for i, project := range Projects {
		deleteArgs[i] = project
	}
	_, err = tx.Exec(deleteQuery, deleteArgs...)
	if err != nil {
		logrus.Debug("SQL Execution error: ", err)
		return err
	}

	insertQuery := `INSERT OR IGNORE INTO projects (project_name) VALUES (?)`
	for _, project := range Projects {
		_, err = tx.Exec(insertQuery, project)
		if err != nil {
			logrus.Debug("SQL Insert error: ", err)
			return err
		}
	}
	return nil
}

func set_telnumber_to_username(UserName string, telNumber string) error {
	logrus.Debug("set_telnumber_to_username")
	fmt.Printf("set_telnumber_to_username telNumber: %v, тип: %T", telNumber, telNumber)
	fmt.Printf("set_telnumber_to_username UserName: %v, тип: %T", UserName, UserName)
	Query := `INSERT OR IGNORE INTO user_phones (user_id, phone_number)
		SELECT users.id, ?
		FROM users
		WHERE users.username = ?
	`
	_, err := db.Exec(Query, telNumber, UserName)
	if err != nil {
		logrus.Error("set_telnumber_to_username ERROR: ", err)
		return err
	}
	return nil
}

func get_telnumber_by_username(UserName string) (string, error) {
	logrus.Debug("get_telnumber_by_username")
	var phoneNumber string

	query := `
        SELECT user_phones.phone_number
        FROM user_phones
        JOIN users ON user_phones.user_id = users.id
        WHERE users.username = ?
    `

	err := db.QueryRow(query, UserName).Scan(&phoneNumber)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	return phoneNumber, nil
}

func set_subs_by_telnumber(phoneNumber string, projectNames []string) error {
	//TODO: НЕ РАБОТАИТ!!!
	//TODO: НЕ РАБОТАИТ!!!
	//TODO: НЕ РАБОТАИТ!!!
	//TODO: НЕ РАБОТАИТ!!!
	//TODO: НЕ РАБОТАИТ!!!
	//TODO: НЕ РАБОТАИТ!!!
	//TODO: НЕ РАБОТАИТ!!!
	//TODO: НЕ РАБОТАИТ!!!
	logrus.Debug("set_subs_by_telnumber")
	fmt.Printf("set_subs_by_telnumber phoneNumber: %v, тип: %T\n", phoneNumber, phoneNumber)
	fmt.Printf("set_subs_by_telnumber projectNames: %v, тип: %T\n", projectNames, projectNames)
	tx, err := db.Begin()
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			logrus.Error(err)
		} else {
			err = tx.Commit()
			if err != nil {
				logrus.Error(err)
			}
		}
	}()

	// Удаляем существующие записи для данного номера телефона
	_, err = tx.Exec("DELETE FROM subscription_lists WHERE phone_number = ?", phoneNumber)
	if err != nil {
		return err
	}

	// Получаем id проектов по их именам
	projectIDs, err := getProjectIDsByName(projectNames)
	if err != nil {
		return err
	}

	// Вставляем новые записи
	stmt, err := tx.Prepare("INSERT OR IGNORE INTO subscription_lists (phone_number, project_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, projectID := range projectIDs {
		_, err = stmt.Exec(phoneNumber, projectID)
		if err != nil {
			return err
		}
	}

	return nil
}

func get_subs_by_telnumber(phoneNumber string) ([]string, error) {
	logrus.Debug("get_subs_by_telnumber")
	var projectNames []string

	// Запрос для получения имен проектов по номеру телефона
	query := `
        SELECT projects.project_name
        FROM subscription_lists
        JOIN projects ON subscription_lists.project_id = projects.id
        WHERE subscription_lists.phone_number = ?
    `

	rows, err := db.Query(query, phoneNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var projectName string
		if err := rows.Scan(&projectName); err != nil {
			return nil, err
		}
		projectNames = append(projectNames, projectName)
	}

	return projectNames, nil
}

func getProjectIDsByName(projectNames []string) ([]int, error) {
	logrus.Debug("getProjectIDsByName")
	var projectIDs []int

	// Создаем строку с плейсхолдерами для IN оператора
	placeholders := make([]interface{}, len(projectNames))
	for i := range projectNames {
		placeholders[i] = "?"
	}

	// Запрос для получения id проектов по их именам
	query := fmt.Sprintf("SELECT id FROM projects WHERE project_name IN (?)")

	// Получаем id проектов
	rows, err := db.Query(query, placeholders...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var projectID int
		if err := rows.Scan(&projectID); err != nil {
			return nil, err
		}
		projectIDs = append(projectIDs, projectID)
	}
	logrus.Debug("getProjectIDsByName finished")
	return projectIDs, nil
}

func test1(Projects []string) error {
	logrus.Debug("test1")
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		logrus.Debug("test1: ", err)
		return err
	}
	defer func() {
		// Rollback the transaction if an error occurs, otherwise commit
		if err != nil {
			logrus.Debug("test2: ", err)
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Delete unnecessary projects
	deleteQuery := `DELETE FROM projects WHERE project_name NOT IN (?)`
	deleteArgs := make([]interface{}, len(Projects))
	for i, project := range Projects {
		deleteArgs[i] = project
	}
	_, err = tx.Exec(deleteQuery, deleteArgs...)
	//_, err = tx.Exec(deleteQuery, Projects)
	if err != nil {
		logrus.Debug("test3: ", err)
		return err
	}

	// Insert new projects
	insertQuery := `INSERT OR IGNORE INTO projects (project_name) VALUES (?)`
	for _, project := range Projects {
		_, err = tx.Exec(insertQuery, project)
		if err != nil {
			logrus.Debug("test4: ", err)
			return err
		}
	}

	return nil
}
