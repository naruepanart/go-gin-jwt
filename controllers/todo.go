package controllers

import (
	"errors"
	"fmt"
	"ilmudata/restapisecurity/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TodoRepo struct {
	Db *gorm.DB
}

func NewTodoController(db *gorm.DB) *TodoRepo {
	db.AutoMigrate(&models.Todo{})
	return &TodoRepo{Db: db}
}


func (repository *TodoRepo) CreateTodo(c *gin.Context) {
	var todo models.Todo
	if c.BindJSON(&todo) == nil {
		err := models.CreateTodo(repository.Db, &todo)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, todo)
	}else{
		fmt.Println(todo)
		c.JSON(http.StatusBadRequest, todo)
	}	
}


func (repository *TodoRepo) GetTodos(c *gin.Context) {
	var todo []models.Todo
	err := models.GetTodos(repository.Db, &todo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (repository *TodoRepo) GetTodo(c *gin.Context) {
	id, _ := c.Params.Get("id")
	idn, _ := strconv.Atoi(id)
	var todo models.Todo
	
	err := models.GetTodoById(repository.Db, &todo, idn)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (repository *TodoRepo) UpdateTodo(c *gin.Context) {
	var todo models.Todo
	var updatedTodo models.Todo

	id, _ := c.Params.Get("id")
	idn, _ := strconv.Atoi(id)	
	err := models.GetTodoById(repository.Db, &updatedTodo, idn)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if c.BindJSON(&todo) == nil {
		updatedTodo.Task = todo.Task
		updatedTodo.Completed = todo.Completed
		updatedTodo.StartDate = todo.StartDate
		updatedTodo.EndDate = todo.EndDate

		err = models.UpdateTodo(repository.Db, &updatedTodo)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, updatedTodo)
	}else{
		c.JSON(http.StatusBadRequest, todo)
	}	
}

func (repository *TodoRepo) DeleteTodo(c *gin.Context) {
	var todo models.Todo
	id, _ := c.Params.Get("id")
	idn, _ := strconv.Atoi(id)

	err := models.DeleteTodoById(repository.Db, &todo, idn)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo was deleted successfully"})
}
