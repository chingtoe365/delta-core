package controller

import (
	"net/http"

	"delta-core/domain"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	ProfileUsecase domain.ProfileUsecase
}

// PingExample godoc
// @Summary Get user profile
// @Schemes
// @Description Get user profile
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200
// @Router /profile [get]
func (pc *ProfileController) Fetch(c *gin.Context) {
	userID := c.GetString("x-user-id")

	profile, err := pc.ProfileUsecase.GetProfileByID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}
