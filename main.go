package main

import (
	"fmt"
	"ilmudata/restapisecurity/auth"
	"ilmudata/restapisecurity/controllers"
	"ilmudata/restapisecurity/database"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


func main() {
	err := godotenv.Load(".env")
	if err!=nil {
		log.Fatalf("Error loading .env file")
	}
	r := gin.Default()

	// cors
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8088","http://localhost:9090"}
	r.Use(cors.New(config))

	// init db and controllers
	db := database.InitDb()
	todoController := controllers.NewTodoController(db)
	userController := controllers.NewUserController(db)	
	userRoleController := controllers.NewUserRoleController(db)	
	

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "welcome todo app",
		})
	})
	r1 := r.Group("/api")
	{
		// todo
		r1.POST("/todo", todoController.CreateTodo)
		r1.GET("/todo", todoController.GetTodos)
		r1.GET("/todo/:id", todoController.GetTodo)
		r1.PUT("/todo/:id", todoController.UpdateTodo)
		r1.DELETE("/todo/:id", todoController.DeleteTodo)

		// register
		r1.POST("/register", userController.Register)
		r1.GET("/users", userController.GetUsers)
		r1.GET("/user/:username", userController.GetUser)

		// change pass and profile
		r1.PUT("/changepassword", userController.ChangePassword)
		r1.PUT("/changeprofile", userController.ChangeProfile)

		
		// delete user
		r1.DELETE("/user/:username", userController.DeleteUser)
	}

	// basic auth
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"admin": "pass123",
		"user1": "pass123",
	}))
	authorized.GET("/todo", todoController.GetTodos)
	authorized.POST("/todo", todoController.CreateTodo)

	// basic auth with db
	basicAuth := auth.InitBasicAuth(db)
	authorized2 := r.Group("/admin2", basicAuth.BasicAuth())
	{
		authorized2.GET("/todo", todoController.GetTodos)
		authorized2.POST("/todo", todoController.CreateTodo)
	}	

	// jwt
	jwtMiddleware, _ := auth.InitJwt(db)
	authHelper := auth.InitHelper(db)	

	// login
	r.POST("/login", jwtMiddleware.LoginHandler)
	// logout for persistence token
	r.GET("/logout", authHelper.VerifyToken, jwtMiddleware.LogoutHandler)
	
	// Refresh time can be longer than token timeout
	//r.GET("/refresh_token", jwtMiddleware.RefreshHandler)
	// Refresh for persistence token
	r.GET("/refresh_token", authHelper.VerifyToken, jwtMiddleware.RefreshHandler)

	//jwtRoute := r.Group("/member", jwtMiddleware.MiddlewareFunc())
	jwtRoute := r.Group("/member", authHelper.VerifyToken, jwtMiddleware.MiddlewareFunc())
	{
		jwtRoute.GET("/todo", todoController.GetTodos)
		
		
		// roles based
		jwtRoute.GET("/todo2", authHelper.CheckRoles([]string{"USER", "ADMIN", "MANAGER"}), todoController.GetTodos)
		jwtRoute.POST("/todo", authHelper.CheckRoles([]string{"USER", "ADMIN", "MANAGER"}), todoController.CreateTodo)
		jwtRoute.GET("/admin", authHelper.CheckRoles([]string{"ADMIN"}), func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "this api for admin member",
			})
		})
		jwtRoute.GET("/manager", authHelper.CheckRoles([]string{"MANAGER"}), func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "this api for manager member",
			})
		})
		jwtRoute.GET("/adminmanager", authHelper.CheckRoles([]string{"ADMIN","MANAGER"}), func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "this api for admin and manager member",
			})
		})


		// roles 
		jwtRoute.POST("/addrole", userRoleController.AddRoleToUser)
		jwtRoute.POST("/removerole", userRoleController.DeleteUserRole)
	}

	
		
	r.Run("localhost:8080")
	fmt.Println("Server is running")
}