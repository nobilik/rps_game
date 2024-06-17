package db

import "rps/models"

// all but self
func GetUsers(id interface{}) ([]*models.User, error) {
	var users []*models.User
	err := DB.Where("id!=?", id).Find(&users).Error
	return users, err
}

// the expected error is login not unique.
// we can add errors checker later
func CreateUser(user *models.User) error {
	return DB.Create(&user).Error
}

func GetUserByLogin(login string) (*models.User, error) {
	var user *models.User
	err := DB.Where("login =?", login).First(&user).Error
	return user, err
}

func GetUserByID(id uint) (*models.User, error) {
	var user *models.User
	err := DB.First(&user, id).Error
	return user, err
}

// here we have two possibilities.
// 1. we request DB on every http request
// 2. we restrict user to has ONLY session
// because of balance issues
func SessionUser(token string) (*models.User, error) {
	id, err := RedisClient.Get("SESSION:" + token).Result()
	if err != nil {
		return nil, err
	}
	var user *models.User
	err = DB.First(&user, id).Error
	return user, err
}
