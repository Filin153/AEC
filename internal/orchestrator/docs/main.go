package docs

import (
	_ "AEC/internal/orchestrator/docs/docs" // Это импорт, который содержит аннотации Swagger
	"AEC/internal/orchestrator/transport"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"io/ioutil"
	"net/http"
)

// Создает swagger

// @title           Вычислитель арифметических выражений(AEC)
// @version         1.0
// @description     Сервис для параллельного вычисления арифметических выражений
// @termsOfService  https://t.me/GusGus153
// @contact.name   Goose
// @contact.url    https://t.me/GusGus153
// @host      localhost:8080
// @BasePath  /

// @Summary AllServer
// @Tags Server
// @Description Get all server
// @Router /server/all [get]
func AllServer(c *gin.Context) {
	// Send a request to your API
	apiURL := "http://localhost:9999/server/all" // Update with your actual API URL

	resp, err := http.Get(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond to the client with the API response
	var data transport.Answer
	json.Unmarshal(body, &data)
	c.JSON(http.StatusOK, data)
}

// @Summary DeleteServer
// @Tags Server
// @Description Delete one server
// @Param id path string true "ID" // Add a dummy parameter to make Swagger recognize the route
// @Router /server/del/{id} [delete]
func DeleteServer(c *gin.Context) {
	// Retrieve the 'id' parameter from the path
	serverID := c.Param("id")

	// Construct the DELETE request
	apiURL := "http://localhost:9999/server/del/" + serverID // Update with your actual API URL
	req, err := http.NewRequest("DELETE", apiURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send the DELETE request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond to the client with the API response
	var data transport.Answer
	json.Unmarshal(body, &data)
	c.JSON(http.StatusOK, data)
}

type TaskType struct {
	Task     string `json:"task"`
	UserId   string `json:"user_id"`
	AddTime  string `json:"add_time"`
	SubTime  string `json:"sub_time"`
	MultTime string `json:"mult_time"`
	DevTime  string `json:"dev_time"`
}

// @Summary AddTask
// @Tags Task
// @Accept json
// @Description Add one task
// @Param input body TaskType true "Request body in JSON format"
// @Router / [post]
func AddTask(c *gin.Context) {

	// Construct the DELETE request
	apiURL := "http://localhost:9999/"

	var requestData TaskType
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Marshal the requestData struct to JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the DELETE request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond to the client with the API response
	var data transport.Answer
	json.Unmarshal(body, &data)
	c.JSON(http.StatusOK, data)
}

// @Summary AddWorker
// @Tags Server
// @Description Add some worker
// @Param id path string true "ID" // Add a dummy parameter to make Swagger recognize the route
// @Param add path string true "workers" // Add a dummy parameter to make Swagger recognize the route
// @Router /server/add/{id}/{add} [post]
func AddWorker(c *gin.Context) {
	// Retrieve the 'id' parameter from the path
	serverID := c.Param("id")
	maxWorker := c.Param("add")

	apiURL := "http://localhost:9999/server/add/" + serverID + "/" + maxWorker // Update with your actual API URL
	req, err := http.NewRequest("POST", apiURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send the DELETE request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond to the client with the API response
	var data transport.Answer
	json.Unmarshal(body, &data)
	c.JSON(http.StatusOK, data)
}

// @Summary GetTask
// @Tags Task
// @Description Get one task
// @Param id path string true "ID" // Add a dummy parameter to make Swagger recognize the route
// @Router /task/{id} [get]
func GetTask(c *gin.Context) {
	// Send a request to your API

	taskID := c.Param("id")
	apiURL := "http://localhost:9999/task/" + taskID // Update with your actual API URL

	resp, err := http.Get(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond to the client with the API response
	var data transport.Answer
	json.Unmarshal(body, &data)
	c.JSON(http.StatusOK, data)
}

// @Summary GetUser
// @Tags User
// @Description Get info about user
// @Param id path string true "ID" // Add a dummy parameter to make Swagger recognize the route
// @Router /user/{id} [get]
func GetUser(c *gin.Context) {
	// Send a request to your API

	taskID := c.Param("id")
	apiURL := "http://localhost:9999/user/" + taskID // Update with your actual API URL

	resp, err := http.Get(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond to the client with the API response
	var data transport.Answer
	json.Unmarshal(body, &data)
	c.JSON(http.StatusOK, data)
}

func Swag() {
	r := gin.New()

	// Add handler for the /server/all route
	r.GET("/server/all", AllServer)
	r.DELETE("/server/del/:id", DeleteServer)
	r.POST("/", AddTask)
	r.POST("/server/add/:id/:add", AddWorker)
	r.GET("/task/:id", GetTask)
	r.GET("/user/:id", GetUser)

	// Swagger routes
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
