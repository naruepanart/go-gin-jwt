package models

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Task string				`json:"task"`
	Completed bool			`json:"completed"`
	StartDate  time.Time 	`json:"startdate"`
	EndDate  time.Time 		`json:"enddate"`
}


func CreateTodo(db *gorm.DB, newTodo *Todo) (err error) {
	err = db.Create(newTodo).Error
	if err != nil {
		return err
	}
	return nil
}
func GetTodos(db *gorm.DB, todos *[]Todo) (err error) {
	err = db.Find(todos).Error
	if err != nil {
		return err
	}
	return nil
}
func GetTodoById(db *gorm.DB, todo *Todo, id int) (err error) {
	err = db.Where("id = ?", id).First(todo).Error
	if err != nil {
		return err
	}
	return nil
}
func UpdateTodo(db *gorm.DB, todo *Todo) (err error) {
	db.Save(todo)
	return nil
}
func DeleteTodoById(db *gorm.DB, todo *Todo, id int) (err error) {
	db.Where("id = ?", id).Delete(todo)
	return nil
}