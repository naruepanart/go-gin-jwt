package models

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string		`json:"username"`
	Password string		`json:"password"`
	Profile Profile	
	Roles []UserRole	
}
type APIUser struct {
	Username string				`json:"username"`
}


func CreateUser(db *gorm.DB, newUser *User) (err error) {
	err = db.Create(newUser).Error
	if err != nil {
		return err
	}
	return nil
}
func GetUsers(db *gorm.DB, apusers *[]APIUser) (err error) {
	err = db.Model(&User{}).Find(apusers).Error
	if err != nil {
		return err
	}
	return nil
}
func GetUserByUsername(db *gorm.DB, user *User, username string) (err error) {
	err = db.Where("username = ?", username).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}
func UpdateUser(db *gorm.DB, todo *User) (err error) {
	db.Save(todo)
	return nil 
}
func DeleteUserByUsername(db *gorm.DB, user *User, username string) (err error) {
	GetUserByUsername(db, user, username)

	tx := db.Begin()	
	if tx.Where("username = ?", username).Delete(&User{}); tx.Error != nil {
		tx.Rollback()
    	return tx.Error
	}
	if tx.Where("user_id = ?", user.ID).Delete(&Profile{}); tx.Error != nil {
		tx.Rollback()
    	return tx.Error
	}	

	return tx.Commit().Error	
}