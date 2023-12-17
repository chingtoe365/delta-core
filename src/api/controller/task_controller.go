package controller

import (
	"delta-core/bootstrap"
	"fmt"

	"net/http"

	"delta-core/domain"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskController struct {
	TaskUsecase      domain.TaskUsecase
	SignalSubUsecase domain.SignalSubUsecase
	ProfileUsecase   domain.ProfileUsecase
}

// PingExample godoc
// @Summary Subsribe
// @Schemes
// @Description Subsribe to signal
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

	userProfile, err := tc.ProfileUsecase.GetProfileByID(c, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

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
	go tc.SignalSubUsecase.Subscribe(env, task, userProfile)

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Task created successfully",
	})
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

// PingExample godoc
// @Summary Delete task - unsubsribe
// @Schemes
// @Description Unsubsribe
// @Tags Task
// @Param taskId query string true "Unsubscribe signal"
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200
// @Router /task [delete]
func (u *TaskController) Cancel(c *gin.Context) {
	// userID := c.GetString("x-user-id")
	taskId := c.Query("taskId")
	task, err := u.TaskUsecase.FetchByTaskID(c, taskId)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
		return
	}
	fmt.Printf(">> going to unsubsribe task ID: %s Title: %s", task.ID.Hex(), task.Title)
	u.SignalSubUsecase.Unsubscribe(&task)
	err = u.TaskUsecase.Delete(c, &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Unsubscribed successfully",
	})
}
