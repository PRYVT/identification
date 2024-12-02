package controller

import (
	"net/http"
	"time"

	"github.com/L4B0MB4/PRYVT/identification/pkg/aggregates"
	models "github.com/L4B0MB4/PRYVT/identification/pkg/models/query"
	"github.com/L4B0MB4/PRYVT/identification/pkg/query/store/repository"
	"github.com/L4B0MB4/PRYVT/identification/pkg/query/utils"
	"github.com/PRYVT/utils/pkg/auth"
	"github.com/PRYVT/utils/pkg/hash"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userRepo *repository.UserRepository
}

func NewUserController(userRepo *repository.UserRepository) *UserController {
	return &UserController{userRepo: userRepo}
}

func (ctrl *UserController) GetToken(c *gin.Context) {

	tokenReq := &models.TokenRequest{}
	err := c.BindJSON(tokenReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userUuid := hash.GenerateGUID(tokenReq.UserName)

	user, err := ctrl.userRepo.GetUserById(userUuid)
	if err != nil || user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	hashedPw := hash.HashPassword(tokenReq.Password)
	if user.PasswordHash != hashedPw {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	token, err := auth.CreateToken(userUuid, 12*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.TokenResponse{Token: token})
}
func (ctrl *UserController) RefreshToken(c *gin.Context) {

	tokenStr := auth.GetTokenFromHeader(c)
	userUuid, err := auth.GetUserUuidFromToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	token, err := auth.CreateToken(userUuid, 12*time.Hour)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.TokenResponse{Token: token})

}

func (ctrl *UserController) GetUser(c *gin.Context) {

	userUuid, err := utils.GetUserIdParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ua, err := aggregates.NewUserAggregate(userUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(ua.Events) > 0 {
		m := models.UserInfo{
			DisplayName: ua.DisplayName,
			Name:        ua.Name,
			Email:       ua.Email,
			ChangeDate:  ua.ChangeDate,
		}
		c.JSON(http.StatusOK, m)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}

}

func (ctrl *UserController) GetUsers(c *gin.Context) {

	limit := utils.GetLimit(c)
	offset := utils.GetOffset(c)

	users, err := ctrl.userRepo.GetAllUsers(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)

}
