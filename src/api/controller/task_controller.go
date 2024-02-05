package controller

import (
	"delta-core/bootstrap"
	"delta-core/internal"
	"fmt"
	"strings"

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
// @Param tradeItem query string false "Trade item"
// @Param tradeItemCategory query string false "Trade item category"
// @Param tradeSignal query string false "Trade signal"
// // @Param task body domain.Task true "Create Task"
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
	// var task domain.Task
	// err := c.ShouldBind(&task)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
	// 	return
	// }
	tradeItem := c.Query("tradeItem")
	tradeItemCategory := c.Query("tradeItemCategory")
	tradeSignal := c.Query("tradeSignal")
	userID := c.GetString("x-user-id")

	var items []domain.TradeItem
	var signals []domain.TradeSignal
	if strings.EqualFold(tradeItemCategory, "ALL") || strings.EqualFold(tradeItem, "ALL") {
		items = internal.GetAllTradeItems()
	} else if tradeItemCategory != "" {
		items = internal.GetTradeItemByCategory(tradeItemCategory)
	} else if tradeItem != "" {
		item, err := internal.GetTradeItemByName(tradeItem)
		if err != nil {
			c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
			return
		}
		items = append(items, item)
	} else {
		// do nothing and items is empty
	}
	if strings.EqualFold(tradeSignal, "ALL") {
		signals = internal.GetAllTradeSignals()
	} else if tradeSignal != "" {
		sig, err := internal.GetSignalByName(tradeSignal)
		if err != nil {
			c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
			return
		}
		signals = append(signals, sig)
	} else {
		// do nothing and signal is empty
	}

	userProfile, err := tc.ProfileUsecase.GetProfileByID(c, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	// loop through items with signals and create tasks
	for _, item := range items {
		for _, signal := range signals {
			var task domain.Task
			task.UserID, err = primitive.ObjectIDFromHex(userID)
			if err != nil {
				c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
				return
			}
			task.ID = primitive.NewObjectID()
			task.Title = internal.BuildTaskTitle(item, signal)
			err = tc.TaskUsecase.Create(c, &task)
			if err != nil {
				c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
				return
			}
			go tc.SignalSubUsecase.Subscribe(env, task, userProfile)
		}
	}
	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: fmt.Sprintf("%d signals subscribed", len(items)*len(signals)),
	})
}

// PingExample godoc
// @Summary Fetch task
// @Schemes
// @Description Fetch task
// @Tags Task
// @Security ApiKeyAuth
// @Param tradeItem query string false "Trade item"
// @Param tradeItemCategory query string false "Trade item category"
// @Param tradeSignal query string false "Trade signal"
// @Accept json
// @Produce json
// @Success 200
// @Router /task [get]
func (u *TaskController) Fetch(c *gin.Context) {
	userID := c.GetString("x-user-id")
	tradeItem := c.Query("tradeItem")
	tradeItemCategory := c.Query("tradeItemCategory")
	tradeSignal := c.Query("tradeSignal")

	tasks, err := u.TaskUsecase.FetchByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	tasksOut := internal.FilterTasks(tradeItem, tradeItemCategory, tradeSignal, tasks)
	c.JSON(http.StatusOK, tasksOut)
}

// PingExample godoc
// @Summary Delete task - unsubsribe
// @Schemes
// @Description Unsubsribe
// @Tags Task
// @Security ApiKeyAuth
// @Param tradeItem query string false "Trade item"
// @Param tradeItemCategory query string false "Trade item category"
// @Param tradeSignal query string false "Trade signal"
// @Accept json
// @Produce json
// @Success 200
// @Router /task [delete]
func (u *TaskController) Cancel(c *gin.Context) {
	userID := c.GetString("x-user-id")
	// taskId := c.Query("taskId")
	// task, err := u.TaskUsecase.FetchByTaskID(c, taskId)
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
	// 	return
	// }
	tradeItem := c.Query("tradeItem")
	tradeItemCategory := c.Query("tradeItemCategory")
	tradeSignal := c.Query("tradeSignal")

	tasks, err := u.TaskUsecase.FetchByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	tasksOut := internal.FilterTasks(tradeItem, tradeItemCategory, tradeSignal, tasks)

	for _, task := range tasksOut {
		fmt.Printf(">> going to unsubsribe task ID: %s Title: %s", task.ID.Hex(), task.Title)
		u.SignalSubUsecase.Unsubscribe(&task)
		err = u.TaskUsecase.Delete(c, &task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Unsubscribed successfully",
	})
}
