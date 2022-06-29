package models

import (
	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	Fullname string	`json:"fullname"`
	Email string  `json:"email"`
	Age int `json:"age"`
	Country string `json:"country"`
	UserID uint	`json:"userid"`
}

func CreateProfile(db *gorm.DB, profile *Profile) (err error) {
	err = db.Create(profile).Error
	if err != nil {
		return err
	}
	return nil
}

func GetProfiles(db *gorm.DB, profiles *[]Profile) (err error) {
	err = db.Find(profiles).Error
	if err != nil {
		return err
	}
	return nil
}

func GetProfile(db *gorm.DB, profile *Profile, id string) (err error) {
	err = db.Where("id = ?", id).First(profile).Error
	if err != nil {
		return err
	}
	return nil
}

func GetProfileByUserId(db *gorm.DB, profile *Profile, id uint) (err error) {
	err = db.Where("user_id = ?", id).First(profile).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateProfile(db *gorm.DB, profile *Profile) (err error) {
	db.Save(profile)
	return nil
}

func DeleteProfile(db *gorm.DB, profile *Profile, id string) (err error) {
	db.Where("id = ?", id).Delete(profile)
	return nil
}



