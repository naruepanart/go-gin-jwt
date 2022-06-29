package models

import (
	"errors"

	"gorm.io/gorm"
)


type UserRole struct {
	gorm.Model
	Role string	`json:"role"`
	UserID uint	`json:"userid"`
}

func AddUserRole(db *gorm.DB, userRole *UserRole) (err error) {
	var userRoles []UserRole
	err = GetRolesByUserId(db, &userRoles, userRole.UserID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	for _, item := range userRoles {
		if item.Role == userRole.Role {
			return nil
		}
	}
	err = db.Create(userRole).Error
	if err != nil {
		return err
	}
	return nil
}

func GetRolesByUserId(db *gorm.DB, userRoles *[]UserRole, id uint) (err error) {
	err = db.Where("user_id = ?", id).Find(userRoles).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserRole(db *gorm.DB, userRoles *UserRole) (err error) {
	err = db.Unscoped().Where("user_id = ? AND role=? ", 
				userRoles.UserID, userRoles.Role).Delete(userRoles).Error
					
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return err	
}
