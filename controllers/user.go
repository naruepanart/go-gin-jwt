package controllers

import (
	"errors"
	"fmt"
	"ilmudata/restapisecurity/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepo struct {
	Db *gorm.DB
}

func NewUserController(db *gorm.DB) *UserRepo {
	db.AutoMigrate(&models.User{}, &models.Profile{})
	return &UserRepo{Db: db}
}


func (repository *UserRepo) Register(c *gin.Context) {
	var user models.User
	var profile models.Profile
	var reg models.UserRegister

	if c.BindJSON(&reg) == nil {
		// check existing username
		var existingUser models.User
		models.GetUserByUsername(repository.Db, &existingUser, reg.Username)		
		if (existingUser.ID > 0) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already taken"})
			return
		}

		user.Username = reg.Username
		profile.Fullname = reg.Fullname
		profile.Email = reg.Email
		user.Profile = profile

		// hashing user password
		passwordBytes, _ := bcrypt.GenerateFromPassword([]byte(reg.Password), 10)
		user.Password = string(passwordBytes)

		err2 := models.CreateUser(repository.Db, &user)
		if err2 != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err2})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user was created successfully"})
	}else{
		fmt.Println(user)
		c.JSON(http.StatusBadRequest, user)
	}	
}


func (repository *UserRepo) GetUsers(c *gin.Context) {
	var profiles []models.Profile
	err := models.GetProfiles(repository.Db, &profiles)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, profiles)
}

func (repository *UserRepo) GetUser(c *gin.Context) {
	username, _ := c.Params.Get("username")	
	var user models.User
	
	err := models.GetUserByUsername(repository.Db, &user, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var profile models.Profile
	err2 := models.GetProfileByUserId(repository.Db, &profile, uint(user.ID))
	if err2 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err2})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (repository *UserRepo) ChangePassword(c *gin.Context) {
	var err error
	var user models.User
	var updatePass models.NewPassword

	if c.BindJSON(&updatePass) == nil {
		// check username
		err = models.GetUserByUsername(repository.Db, &user, updatePass.Username)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		// verify current password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(updatePass.Password))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		// hashing a new password
		passwordBytes, _ := bcrypt.GenerateFromPassword([]byte(updatePass.NewPassword), 10)
		user.Password = string(passwordBytes)

		err = models.UpdateUser(repository.Db, &user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "password has be changed successfully"})
	}else{
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	}	
}

func (repository *UserRepo) ChangeProfile(c *gin.Context) {
	var err error
	var user models.User
	var profile models.Profile
	var updateProfile models.UserProfile

	if c.BindJSON(&updateProfile) == nil {
		// check username
		err = models.GetUserByUsername(repository.Db, &user, updateProfile.Username)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		// get profile
		err = models.GetProfileByUserId(repository.Db, &profile, uint(user.ID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		// update profile
		profile.Fullname = updateProfile.Fullname
		profile.Email = updateProfile.Email
		profile.Age = updateProfile.Age
		profile.Country = updateProfile.Country
		
		err = models.UpdateProfile(repository.Db, &profile)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user profile has be changed successfully"})
	}else{
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	}	
}

func (repository *UserRepo) DeleteUser(c *gin.Context) {
	var user models.User
	username, _ := c.Params.Get("username")	
	err := models.DeleteUserByUsername(repository.Db, &user, username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User was deleted successfully"})
}
