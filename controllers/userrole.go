package controllers

import (
	"ilmudata/restapisecurity/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRoleRepo struct {
	Db *gorm.DB
}

func NewUserRoleController(db *gorm.DB) *UserRoleRepo {
	db.AutoMigrate(&models.User{}, &models.UserRole{})
	return &UserRoleRepo{Db: db}
}


func (repository *UserRoleRepo) AddRoleToUser(c *gin.Context) {
	var userRole models.UserRole
	var userRoleApi models.UserRoleApi

	if c.BindJSON(&userRoleApi) == nil {
		// check existing username
		var existingUser models.User
		models.GetUserByUsername(repository.Db, &existingUser, userRoleApi.Username)		
		if (existingUser.ID == 0) {
			c.JSON(http.StatusNotFound, gin.H{"message": "username not found"})
			return
		}
		// add to userRole
		userRole.UserID = existingUser.ID
		userRole.Role = strings.ToUpper(userRoleApi.Role)
		
		err := models.AddUserRole(repository.Db, &userRole)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user role was created successfully"})
	}else{		
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
	}	
}

func (repository *UserRoleRepo) DeleteUserRole(c *gin.Context) {
	var userRole models.UserRole
	var userRoleApi models.UserRoleApi

	if c.BindJSON(&userRoleApi) == nil {
		// check existing username
		var existingUser models.User
		models.GetUserByUsername(repository.Db, &existingUser, userRoleApi.Username)		
		if (existingUser.ID == 0) {
			c.JSON(http.StatusNotFound, gin.H{"message": "username not found"})
			return
		}
		// delete userrole
		userRole.UserID = existingUser.ID
		userRole.Role = strings.ToUpper(userRoleApi.Role)
		
		err := models.DeleteUserRole(repository.Db, &userRole)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user role was deleted successfully"})
	}else{		
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
	}	
}
