package auth

import (
	"log"
	"net/http"
	"os"
	"time"

	"ilmudata/restapisecurity/models"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)


var identityKey = "id"
var identityName = "name"

func InitJwt(db *gorm.DB) (*jwt.GinJWTMiddleware, error) {
	JWT_KEY := os.Getenv("JWT_KEY")
	JWT_REALM := os.Getenv("JWT_REALM")
	
	db.AutoMigrate(&models.User{}, &models.Profile{},&models.UserToken{})

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       JWT_REALM,
		Key:         []byte(JWT_KEY),
		Timeout:     7*24*time.Hour, // token expiration, 7 days
		MaxRefresh:  10*24*time.Hour, // refresh token expiration, 10 days
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.UserPayload); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
					identityName: v.Fullname,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.UserPayload{
				UserName: claims[identityKey].(string),
				Fullname: claims[identityName].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals models.Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			var user models.User
			var profile models.Profile

			err := models.GetUserByUsername(db, &user,username)
			if err == nil {
				
				err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
				if err == nil {

					err := models.GetProfileByUserId(db, &profile,user.ID)
					if err == nil {
						c.Set("CURRENT_USERNAME", username)
						return &models.UserPayload{
							UserName:  username,
							Fullname:  profile.Fullname,
						}, nil
					}
					
				}
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*models.UserPayload); ok {				
				return true
			}

			return false
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {			
			// c.JSON(code, gin.H{
			// 	"access_token":  token,
			// 	"expires_in": expire.Format(time.RFC3339),
			// })


			// persistence token
			username,_ := c.Get("CURRENT_USERNAME")

			var userToken models.UserToken
			userToken.Token = token
			userToken.Username = username.(string)

			models.DeleteTokenByUsername(db, &userToken)
			err := models.SetToken(db, &userToken)
			if err == nil {
				c.JSON(code, gin.H{
					"access_token":  token,
					"expires_in": expire.Format(time.RFC3339),
				})
			}else{
				c.JSON(http.StatusInternalServerError, gin.H{
					"message":  "there was a problem in token processing",					
				})
			}

			
		},
		LogoutResponse: func(c *gin.Context, code int) {
			// persistence token
			var userToken models.UserToken			

			oldToken,_ := c.Get("CURRENT_JWT_TOKEN")
			userToken.Token = oldToken.(string)
			models.DeleteToken(db, &userToken)

			c.JSON(http.StatusOK, gin.H{
				"message":  "logout process was succeed",					
			})
			
		},
		RefreshResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			// c.JSON(code, gin.H{
			// 	"access_token":  token,
			// 	"expires_in": expire.Format(time.RFC3339),
			// })

			// persistence token
			var userToken models.UserToken			

			oldToken,_ := c.Get("CURRENT_JWT_TOKEN")
			userToken.Token = oldToken.(string)
			models.DeleteToken(db, &userToken)

			userToken.Token = token
			err := models.SetToken(db, &userToken)
			if err == nil {
				c.JSON(code, gin.H{
					"access_token":  token,
					"expires_in": expire.Format(time.RFC3339),
				})
			}else{
				c.JSON(http.StatusInternalServerError, gin.H{
					"message":  "there was a problem in token processing",					
				})
			}
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": "unauthorized",
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",
		// no cookie
		SendCookie: false,

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}
	return authMiddleware, err
}