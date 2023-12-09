package main

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

func DBInit() {
	var err error

	// Инициализация пула подключений к SQLite3
	db, err = sql.Open("sqlite", "./mydatabase.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	//defer db.Close()

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
		    FOREIGN KEY (user_id) REFERENCES users(user_id),
		    FOREIGN KEY (project_id) REFERENCES projects(project_id),
		    PRIMARY KEY (user_id, project_id)
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
		return err
	}

	if count > 0 {
		return fmt.Errorf("Пользователь с именем пользователя %s уже существует", username)
	}

	// Сохраняем данные нового пользователя в таблице users
	result, err := db.Exec("INSERT INTO users (username, password, email) VALUES (?, ?, ?)", username, password, email)
	if err != nil {
		return err
	}

	// Получаем ID нового пользователя
	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Получаем ID ролей из таблицы roles по их именам
	var roleIDs []int
	for _, roleName := range roles {
		var roleID int
		err := db.QueryRow("SELECT id FROM roles WHERE name = ?", roleName).Scan(&roleID)
		if err != nil {
			return err
		}
		roleIDs = append(roleIDs, roleID)
	}

	// Сохраняем связи пользователей и их ролей в таблице user_roles
	for _, roleID := range roleIDs {
		_, err := db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", userID, roleID)
		if err != nil {
			return err
		}
	}

	return nil
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

func AddProjectRootPage(PageId int, ProjectName string) error {
	Query := "INSERT OR IGNORE INTO project_root_page (project_id, page_id) SELECT p.id, ? FROM projects p WHERE p.project_name = ?;"
	logrus.Debug("QW = ", Query, " ProjectName=", ProjectName, " PageId=", PageId)
	_, err := db.Exec(Query, PageId, ProjectName)
	if err != nil {
		return err
	}
	return nil
	//TODO: new feature
}

func AddUserSubscriptions(UserName, Projects) error {
	for _, project := range Projects {
		Query := "INSERT OR IGNORE INTO user_subscriptions"
	}
	return nil
}
