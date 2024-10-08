package controllers

import (
    "fmt"
	"mongodb/data"
	"mongodb/models"
	"mongodb/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskController interface {
	GetAllTasks(c *gin.Context)
	CreateTasks(c *gin.Context)
	GetTasksById(c *gin.Context)
	UpdateTasksById(c *gin.Context)
	DeleteTasksById(c *gin.Context)
}

type taskController struct {
	taskService data.TaskService
}

func NewTaskController() (*taskController, error) {
	service_reference, err := data.NewTaskService()
	if err != nil {
		return nil, err
	}
	return &taskController{
		taskService: *service_reference,
	}, nil
}

func (tc *taskController) GetAllTasks(c *gin.Context) {
	tasks, err, statusCode := tc.taskService.GetTasks()
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, tasks)
}

func (tc *taskController) CreateTasks(c *gin.Context) {
	var newTak models.Task
	// v := validator.New()
	if err := c.ShouldBindJSON(&newTak); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid data", "error": err.Error()})
		return
	}
	// if err := v.Struct(newTak); err != nil {
	// 	fmt.Printf(err.Error())
	// 	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid or missing data", "error": err.Error()})
	// 	return
	// }
	// err := utils.ValidateStatus(&newTak)
	// if err != nil {
	// 	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid or missing data, allowed options are Pending, In Progress, and Completed", "error": err.Error()})
	// 	return
	// }
	fmt.Println(newTak)
	task, err, statusCode := tc.taskService.CreateTasks(&newTak)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusCreated, task)
	}	
}

func (tc *taskController) GetTasksById(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task, err, statusCode := tc.taskService.GetTasksById(objectID)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(200, task)
	}
}

func (tc *taskController) UpdateTasksById(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedTask models.Task

	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	er := utils.ValidateStatus(&updatedTask)
	if er != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid or missing data, allowed options are Pending, In Progress, and Completed", "error": err.Error()})
		return
	}
	data, er, status := tc.taskService.UpdateTasksById(objectID, updatedTask)
	if er != nil {
		c.IndentedJSON(status, gin.H{"error": er.Error()})
		return
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Task updated", "data": data})
	}
}

func (tc *taskController) DeleteTasksById(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err, statusCode := tc.taskService.DeleteTasksById(objectID)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"success": "Data Deleted"})
	}
}
