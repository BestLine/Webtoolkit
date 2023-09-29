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

	for _, arg := range args {
		switch v := arg.(type) {
		case map[string]interface{}:
			body = v
		case string:
			strArgs = append(strArgs, v)
		}
		fmt.Printf("Type: %T, Value: %v\n", arg, arg)
	}
	targetURL := viper.GetString("backend.host") + strArgs[1]
	if strArgs[0] == "Get" {
		response, err = http.Get(targetURL)
	} else if strArgs[0] == "Post" {
		logrus.Debug("body: ", body)
		jsonBody, err := json.Marshal(body)
		if err != nil {
			logrus.Error(err)
			c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
			return nil
		}
		response, err = http.Post(targetURL, "application/json", bytes.NewBuffer(jsonBody))
	}
	if err != nil {
		logrus.Error(err)
		c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		return nil
	}
	defer response.Body.Close()

	if err != nil {
		logrus.Error(err)
		c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		return nil
	}
	logrus.Debug("targetURL: ", targetURL)
	logrus.Debug("response: ", response)
	//fmt.Println("resp code: ", response.StatusCode)
	if response.StatusCode != 200 {
		logrus.Error("sendRequest responce code: ", response.StatusCode)
		c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		return nil
	}
	return RespToByteReader(response)
}
