package router

import (
	"mongodb/controllers"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	controller_reference, err := controllers.NewTaskController()
	if err != nil {
		return nil
	}
	taskController := *controller_reference

	r.GET("/tasks", taskController.GetAllTasks)
	r.POST("/tasks", taskController.CreateTasks)
	r.GET("/tasks/:id", taskController.GetTasksById)
	r.PUT("/tasks/:id", taskController.UpdateTasksById)
	r.DELETE("/tasks/:id", taskController.DeleteTasksById)

	return r
}