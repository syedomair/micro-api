package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	testdata "github.com/syedomair/api_micro/testdata"
)

type testCaseType struct {
	method         string
	url            string
	path           string
	pathParam      string
	requestBody    string
	responseResult string
	responseData   string
}

var testCases []testCaseType

func main() {
	arg := ""
	if len(os.Args) > 1 {
		arg = os.Args[1]
		fmt.Println("Command line argument: ", arg)
	}

	kongAdminURL, _ := minikubeServiceURL("kong-admin")
	kongProxyURL, _ := minikubeServiceURL("kong-proxy")
	publicURL, _ := minikubeServiceURL("public-srvc")
	userURL, _ := minikubeServiceURL("users-srvc")
	roleURL, _ := minikubeServiceURL("roles-srvc")
	batchTaskURL, _ := minikubeServiceURL("batch-tasks-srvc")

	fmt.Println("kong-admin:", kongAdminURL)
	fmt.Println("kong-proxy:", kongProxyURL)
	fmt.Println("public-srvc:", publicURL)
	fmt.Println("users-srvc:", userURL)
	fmt.Println("roles-srvc:", roleURL)
	fmt.Println("batch-tasks-srvc:", batchTaskURL)

	if arg == "kong" {
		publicURL = kongProxyURL
		userURL = kongProxyURL
		roleURL = kongProxyURL
		batchTaskURL = kongProxyURL
	}
	currentTime := time.Now()
	uniqueEmail := "email_" + strconv.FormatInt(currentTime.UnixNano(), 10) + "@gmail.com"

	testCases = []testCaseType{
		{"POST", publicURL, "/v1/register", "", `{"first_name":"` + testdata.ValidFirstName + `", "last_name":"` + testdata.ValidLastName + `", "email":"` + uniqueEmail + `", "password":"` + testdata.ValidPassword + `"}`, `"success"`, ``},
		{"POST", publicURL, "/v1/authenticate", "", `{"email":"` + uniqueEmail + `", "password":"` + testdata.ValidPassword + `"}`, `"success"`, ``},
		/*
			{"POST", roleURL, "/v1/roles", "", `{"title":"` + testdata.RoleTitle1 + `","role_type":"` + testdata.RoleType + `"}`, `"success"`, ``},
			{"GET", roleURL, "/v1/roles/", "role_id", ``, `"success"`, ``},
			{"GET", roleURL, "/v1/roles", "", ``, `"success"`, ``},
			{"PATCH", roleURL, "/v1/roles/", "role_id", `{"title":"` + testdata.RoleTitle2 + `","role_type":"` + testdata.RoleType + `"}`, `"success"`, ``},
			{"GET", roleURL, "/v1/roles/", "role_id", ``, `"success"`, ``},
			{"DELETE", roleURL, "/v1/roles/", "role_id", ``, `"success"`, ``},
			{"GET", userURL, "/v1/users/", "user_id", ``, `"success"`, ``},
			{"GET", userURL, "/v1/users", "", ``, `"success"`, ``},
			{"PATCH", userURL, "/v1/users/", "user_id", `{"first_name":"` + testdata.ValidFirstName + `"}`, `"success"`, ``},
			{"GET", userURL, "/v1/users/", "user_id", ``, `"success"`, ``},
		*/
		{"GET", batchTaskURL, "/v1/batch/users", "", ``, `"success"`, ``},
		{"GET", batchTaskURL, "/v1/batch/users/status/", "batch_task_id", ``, `"success"`, ``},
		{"GET", batchTaskURL, "/v1/batch/users/output/", "batch_task_id", ``, `"success"`, ``},
		{"DELETE", userURL, "/v1/users/", "user_id", ``, `"success"`, ``},
	}
	i := 0
	userId := ""
	roleId := ""
	batchTaskId := ""
	token := ""
	for _, testCase := range testCases {
		url := ""
		if testCase.pathParam == "user_id" {
			url = testCase.url + testCase.path + userId
		} else if testCase.pathParam == "role_id" {
			url = testCase.url + testCase.path + roleId
		} else if testCase.pathParam == "batch_task_id" {
			url = testCase.url + testCase.path + batchTaskId
		} else {
			url = testCase.url + testCase.path
		}
		req, err := http.NewRequest(testCase.method, url, strings.NewReader(testCase.requestBody))

		if err != nil {
			print(err)
		}
		if i > 1 {
			req.Header.Set("apikey", token)
			req.Header.Set("authorization", token)
		} else {
			req.Header.Set("authorization", testdata.TestValidPublicToken)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			print(err)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		var bodyInterface map[string]interface{}

		json.Unmarshal(body, &bodyInterface)
		jsonData, _ := json.Marshal(bodyInterface["data"])
		jsonResult, _ := json.Marshal(bodyInterface["result"])
		fmt.Println("---------------------------------------------------------------------------")
		fmt.Println(strconv.Itoa(i) + " " + testCase.method + " " + string(testCase.url) + " " + testCase.path + " " + testCase.requestBody)
		fmt.Println(string(body))
		fmt.Println(string(jsonData))
		//Green Color output
		fmt.Println("\033[32m" + string(jsonResult) + "\033[39m")

		if testCase.method == "POST" && testCase.path == "/v1/register" {
			var userIdInterface map[string]interface{}
			json.Unmarshal(jsonData, &userIdInterface)
			jsonUserId, _ := json.Marshal(userIdInterface["user_id"])
			userId = string(bytes.Trim(jsonUserId, `"`))
			fmt.Println("userId:", userId)
		}
		if testCase.method == "POST" && testCase.path == "/v1/authenticate" {
			var tokenInterface map[string]interface{}
			json.Unmarshal(jsonData, &tokenInterface)
			jsonToken, _ := json.Marshal(tokenInterface["token"])
			token = string(bytes.Trim(jsonToken, `"`))
			fmt.Println("token:", token)
		}
		if testCase.method == "POST" && testCase.path == "/v1/roles" {
			var roleIdInterface map[string]interface{}
			json.Unmarshal(jsonData, &roleIdInterface)
			jsonUserId, _ := json.Marshal(roleIdInterface["role_id"])
			roleId = string(bytes.Trim(jsonUserId, `"`))
			fmt.Println("roleId:", roleId)
		}
		if testCase.method == "GET" && testCase.path == "/v1/batch/users" {
			var batchTaskIdInterface map[string]interface{}
			json.Unmarshal(jsonData, &batchTaskIdInterface)
			jsonUserId, _ := json.Marshal(batchTaskIdInterface["id"])
			batchTaskId = string(bytes.Trim(jsonUserId, `"`))
			fmt.Println("batchTaskId:", batchTaskId)
		}
		time.Sleep(2 * time.Second)
		fmt.Println("---------------------------------------------------------------------------")

		i++
	}
}

func minikubeServiceURL(serviceName string) (string, error) {
	minikube := "minikube"
	service := "service"
	url := "--url"

	serviceURL, err := exec.Command(minikube, service, serviceName, url).Output()
	if err != nil {
		return "", err
	}
	return strings.Trim(string(serviceURL), "\n"), nil
}
