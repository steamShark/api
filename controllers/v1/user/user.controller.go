package controllers

/*
import (
	"net/http"

	"steamshark-api/services"
	"steamshark-api/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{Service: service}
}

func (ctrl *UserController) GetUser(c *gin.Context) {
	steamID := c.Param("steamid64")
	user, err := ctrl.Service.GetUserBySteamID(steamID)

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "User not found")
		return
	}

	if user == nil {
		utils.Error(c, http.StatusNotFound, "User not found")
		return
	}

	utils.Success(c, "User found", user)
}
*/
