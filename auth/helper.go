package auth

import (
	"ilmudata/restapisecurity/models"
	"strings"

	"net/http"

	"gorm.io/gorm"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type HelperRepo struct {
	Db *gorm.DB
}

func InitHelper(db *gorm.DB) *HelperRepo {
	db.AutoMigrate(&models.User{})
	return &HelperRepo{Db: db}
}

func (repository *HelperRepo) VerifyToken(c *gin.Context) {   
	authHeader := c.Request.Header.Get("Authorization")
	token := parseToken(authHeader)

	if token != "" {
		var userToken models.UserToken
		userToken.Token = token
		err := models.GetToken(repository.Db, &userToken, userToken.Token)
		if err!=nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})			
			return
		}
		c.Set("CURRENT_JWT_TOKEN", token)
	}else {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
}

func (repository *HelperRepo) CheckRoles(roles []string) gin.HandlerFunc { 
	
	return func(c *gin.Context) {
		payload, _ := c.Get("JWT_PAYLOAD")
		claims := payload.(jwt.MapClaims)

		var user models.User
		err := models.GetUserByUsername(repository.Db, &user, claims["id"].(string))
		if err == nil {
			var userRoles []models.UserRole
			err := models.GetRolesByUserId(repository.Db, &userRoles, user.ID)
			if err == nil {
				for _, userRole := range userRoles {
					for _, role := range roles {
						if strings.ToUpper(userRole.Role) == strings.ToUpper(role) {
							c.Next()
							return
						}
					}
				}
			}
		}
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	
}

func parseToken(authorization string) (token string){ 
	if authorization != "" {
		parts := strings.SplitN(authorization, " ", 2)
		if len(parts)==2 {
			return parts[1]
		}	
	}

	return ""
}