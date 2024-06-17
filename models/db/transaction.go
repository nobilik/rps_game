package db

import (
	"errors"
	"rps/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetTransactionsByUserID(userID uint) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	err := DB.Where("user_id =?", userID).Find(&transactions).Error
	return transactions, err
}

func TopUp(userID uint, amount float64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var user *models.User
		// Lock the user record to prevent race conditions
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", userID).First(&user).Error; err != nil {
			return err
		}

		// Create the transaction
		transaction := &models.Transaction{
			UserID: userID,
			Amount: amount,
			Type:   models.TopUp,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		// Update the user balance
		if err := tx.Model(&models.User{}).Where("id = ?", userID).
			UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		return nil
	})
}

func Withdraw(userID uint, amount float64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var user *models.User
		// Lock the user record to prevent race conditions
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", userID).First(&user).Error; err != nil {
			return err
		}

		// Check if the user has enough balance
		if user.Balance < amount { // add all active bets amounts
			return errors.New("insufficient balance")
		}

		// Create the transaction
		transaction := &models.Transaction{
			UserID: userID,
			Amount: amount,
			Type:   models.Withdrawal,
		}
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// Update the user balance
		if err := tx.Model(&models.User{}).Where("id = ?", userID).
			UpdateColumn("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}

		return nil
	})
}
