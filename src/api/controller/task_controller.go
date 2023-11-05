package controller

import (
	"delta-core/bootstrap"
	pubsub "delta-core/services"
	"fmt"

	"net/http"

	"delta-core/domain"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskController struct {
	TaskUsecase domain.TaskUsecase
}

// PingExample godoc
// @Summary Create task
// @Schemes
// @Description Create task
// @Tags Task
// @Security ApiKeyAuth
// @Param task body domain.Task true "Create Task"
// @Accept json
// @Produce json
// @Success 200
// @Router /task [post]
func (tc *TaskController) Create(c *gin.Context) {
	env, ok := c.MustGet("env").(*bootstrap.Env)
	if !ok {
		fmt.Println(ok)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Environment injection failed"})
		return
	}
	var task domain.Task
	err := c.ShouldBind(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	userID := c.GetString("x-user-id")
	task.ID = primitive.NewObjectID()

	task.UserID, err = primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	err = tc.TaskUsecase.Create(c, &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	go startSubscribe(env, task)

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Task created successfully1",
	})
}

func startSubscribe(env *bootstrap.Env, task domain.Task) {
	fmt.Printf("Start a new subscription with subClientID = taskId: %v, topic = task.Title: %s. ", task.ID.Hex(), task.Title)
	// start a new goroutine to receive new messages
	go pubsub.CreateSubClient(env, task.ID.Hex(), task.Title)
}

// PingExample godoc
// @Summary Fetch task
// @Schemes
// @Description Fetch task
// @Tags Task
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200
// @Router /task [get]
func (u *TaskController) Fetch(c *gin.Context) {
	userID := c.GetString("x-user-id")

	tasks, err := u.TaskUsecase.FetchByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
