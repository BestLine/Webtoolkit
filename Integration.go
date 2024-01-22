package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

func sendRequest(c *fiber.Ctx, args ...interface{}) []byte {
	// 1st Get\Post
	// 2st url
	// 3rd body
	logrus.Debug("sendRequest")
	var strArgs []string
	var response *http.Response
	var err error
	var body map[string]interface{}
	var byte_body []byte
	var jsonBody []byte

	for _, arg := range args {
		switch v := arg.(type) {
		case map[string]interface{}:
			body = v
		case []byte:
			byte_body = v
		case string:
			strArgs = append(strArgs, v)
		}
		//fmt.Printf("Type: %T, Value: %v\n", arg, arg)
	}
	targetURL := viper.GetString("backend.pure_host") + strArgs[1]
	if strArgs[0] == "Get" {
		response, err = http.Get(targetURL)
		if err != nil {
			logrus.Error("Get sending error: ", err)
			c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
			return nil
		}
	} else if strArgs[0] == "Post" {
		logrus.Debug("sendRequest body: ", body)
		jsonBody, err = json.Marshal(body)
		if err != nil {
			logrus.Error("Post json.Marshal error: ", err)
			c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
			return nil
		}
		response, err = http.Post(targetURL, "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			logrus.Error("Post sending error: ", err)
			c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
			return nil
		}
	} else if strArgs[0] == "Post2" {
		logrus.Debug("sendRequest body: ", string(byte_body))
		response, err = http.Post(targetURL, "application/json", bytes.NewBuffer(byte_body))
		if err != nil {
			logrus.Error("Post sending error: ", err)
			c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
			return nil
		}
	} else if strArgs[0] == "Post3" {
		logrus.Debug("sendRequest body: ", string(byte_body))
		targetURL = strArgs[1]
		response, err = http.Post(targetURL, "application/json", bytes.NewBuffer(byte_body))
		if err != nil {
			logrus.Error("Post sending error: ", err)
			c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
			return nil
		}
	}
	if err != nil {
		logrus.Error("sendRequest unknown error: ", err)
		c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		return nil
	}
	fmt.Println("responce: ", response)
	fmt.Println("err: ", err)
	defer response.Body.Close()

	logrus.Debug("sendRequest targetURL: ", targetURL)
	logrus.Debug("sendRequest response: ", response)
	if response.StatusCode != 200 {
		logrus.Error("sendRequest responce code: ", response.StatusCode)
		c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		return nil
	}
	return RespToByteReader(response)
}
