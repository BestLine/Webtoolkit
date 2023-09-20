package main

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

func createSQLiteDB() (*sql.DB, error) {
	// Открываем или создаем файл базы данных SQLite
	db, err := sql.Open("sqlite", "./mydatabase.db")
	if err != nil {
		return nil, err
	}

	// Создаем таблицы (если их нет)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			password TEXT NOT NULL,
			email TEXT
		)
	`)
	if err != nil {
		db.Close()
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		)
	`)
	if err != nil {
		db.Close()
		return nil, err
	}

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
			project_name TEXT NOT NULL
		)
	`)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func registerUser(db *sql.DB, username, password, email string, roles []string) error {
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

func checkUserCredentials(db *sql.DB, username, password string) (bool, error) {
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

func getUserRole(db *sql.DB, username string) string {
	// Проверяем наличие роли у пользователя лениво
	row := db.QueryRow("select name as role from users left join user_roles ur on users.id = ur.user_id left join roles r on r.id = ur.role_id where username=?", username)
	var storedRole string
	err := row.Scan(&storedRole)
	if err != nil {
		return "Ошибка при выполнении запроса"
	}
	return storedRole
}

func hasUserRole(db *sql.DB, username, role string) (bool, error) {
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

func GetUserProject(db *sql.DB, username string) (string, error) {
	query := "SELECT projects.project_name " +
		"FROM user_projects " +
		"LEFT JOIN projects ON user_projects.project_id = projects.id " +
		"LEFT JOIN users ON users.id = user_projects.user_id " +
		"WHERE users.username = '" + username + "' AND user_projects.active = 1;"

	var name string
	err := db.QueryRow(query, username).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func GetUserProjects(db *sql.DB, username string) ([]string, error) {
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

func GetAllProjects(db *sql.DB) ([]string, error) {
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

func SetActiveProject(db *sql.DB, username string, projectToActivate string) error {
	// Сначала помечаем все проекты пользователя как неактивные
	updateQuery := "UPDATE user_projects " +
		"SET active = 0 " +
		"WHERE user_id = (SELECT id FROM users WHERE username = ?)"
	logrus.Debug("QW_1 = ", updateQuery, username)
	_, err := db.Exec(updateQuery, username)
	if err != nil {
		return err
	}

	// Затем помечаем конкретный проект как активный
	updateProjectQuery := "UPDATE user_projects " +
		"SET active = 1 " +
		"WHERE user_id = (SELECT id FROM users WHERE username = ?) " +
		"AND project_id = (SELECT id FROM projects WHERE project_name = ?)"
	logrus.Debug("QW_2 = ", updateProjectQuery, username, projectToActivate)
	_, err = db.Exec(updateProjectQuery, username, projectToActivate)
	if err != nil {
		return err
	}

	return nil
}

// spec
func getDbConn() (db *sql.DB) {
	db, err := createSQLiteDB()
	if err != nil {
		fmt.Println("Ошибка при создании базы данных:", err)
		return
	}
	return db
}
