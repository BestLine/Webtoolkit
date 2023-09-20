package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func loginHandler(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	db := getDbConn()
	check, err := checkUserCredentials(db, username, password)
	if !check {
		logrus.Error("Invalid credentials! Username: ", username, " PWD: ", password)
		return c.Render("login",
			fiber.Map{"error": "Не верное имя пользователя или пароль!"})
	}

	role := getUserRole(db, username)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // токен действителен в течение 24 часов

	tokenString, err := token.SignedString([]byte("secret"))
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	if err != nil {
		logrus.Error("Internal server error: ", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	return c.Redirect("/")
}

func jwtMiddleware() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		cookie := c.Cookies("jwt")
		token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil // используем секретный ключ, который был сгенерирован при создании токена
		})

		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			logrus.Error("StatusUnauthorized: ", err)
			return c.Redirect("/login")
		}

		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok {
			c.Status(fiber.StatusUnauthorized)
			logrus.Error("StatusUnauthorized: ", err)
			return c.Redirect("/login")
		}

		// Проверяем, есть ли поле "role" в утверждениях токена и проверяем, что роль равна "admin"
		role, ok := (*claims)["role"].(string)
		logrus.Debug("role = ", role)
		logrus.Debug("role check!")
		c.Locals("user", claims)
		return c.Next()
	}
}

// ///////////// регистрация и ролевая модель ///////////////
func getRegister(c *fiber.Ctx) error {
	logrus.Debug("getRegister")
	return c.Render("register",
		fiber.Map{"Title": "Мотай вниз", "Description": "Где то тут юзаются шаблоны."})
}

func getLogin(c *fiber.Ctx) error {
	logrus.Debug("getLogin")
	return c.Render("login",
		fiber.Map{"msg": "Для продолжения требуется авторизация."})
}

func logout(c *fiber.Ctx) error {
	logrus.Debug("getLogout")
	c.ClearCookie("jwt") // удаляем куки пользователя
	return c.Render("login", fiber.Map{"msg": "Для продолжения требуется авторизация."})
}

func register(c *fiber.Ctx) error {
	logrus.Debug("postRegister")
	username := c.FormValue("username")
	password := c.FormValue("password")
	email := c.FormValue("email")
	logrus.Debug("username: ", username, " password: ", password, " email: ", email)
	db := getDbConn()
	err := registerUser(db, username, password, email, []string{"user"})
	if err != nil {
		logrus.Error("Ошибка при регистрации пользователя: ", err)
		return c.Render("login",
			fiber.Map{"msg": "Ошибка при регистрации пользователя"})
	}
	return c.Render("login",
		fiber.Map{"msg": "Теперь вы можете авторизоваться используя новую учётную запись."})
}

func checkUserPermission(c *fiber.Ctx) error {
	data := c.Locals("user")
	logrus.Debug("UserData: ", data)
	logrus.Debug("UserData: ", data.(jwt.MapClaims)["role"].(string))
	return nil
	//TODO: сделать проверку роли пользователя на соответствие доступу
}
