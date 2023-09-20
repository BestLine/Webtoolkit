package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

func sendPost(c *fiber.Ctx, url string, body map[string]interface{}) []byte {
	logrus.Debug("sendPost")
	jsonBody, err := json.Marshal(body)
	if err != nil {
		c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		return nil
	}

	// Отправляем POST-запрос на целевой URL
	targetURL := viper.GetString("backend.host") + url // Замените на свой URL
	response, err := http.Post(targetURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		return nil
	}
	defer response.Body.Close()

	// Читаем ответ
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		return nil
	}
	logrus.Debug("targetURL: ", targetURL)
	logrus.Debug("body: ", body)
	logrus.Debug("response: ", response)
	fmt.Println("Response =", string(responseBody))
	return responseBody
}

func sendGet(c *fiber.Ctx, url string) []byte {
	// Отправляем GET-запрос на целевой URL
	//targetURL := "http://127.0.0.1:7778" + url // Замените на свой URL
	logrus.Debug("sendGet")
	targetURL := viper.GetString("backend.host") + url // Замените на свой URL
	response, err := http.Get(targetURL)
	if err != nil {
		c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		return nil
	}
	defer response.Body.Close()

	// Читаем ответ
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		return nil
	}
	logrus.Debug("targetURL: ", targetURL)
	logrus.Debug("response: ", response)
	fmt.Println("Response =", string(responseBody))
	return responseBody
}
