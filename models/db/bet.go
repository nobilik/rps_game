package db

import (
	"errors"
	"rps/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetUserPendingBets(userID uint) ([]*models.Bet, error) {
	var bets []*models.Bet
	err := DB.Where("status = ? AND acceptor_id = ?", models.Pending, userID).Find(&bets).Error
	return bets, err
}

func GetBetByID(id uint) (*models.Bet, error) {
	var bet *models.Bet
	err := DB.First(&bet, id).Error
	return bet, err
}

// on new bet
// creates bet and updates the user balance
func Freeze(userID, acceptorID uint, amount float64, item models.Item) error {
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
		bet := &models.Bet{
			UserID:     userID,
			Amount:     amount,
			ItemValue:  item.Value,
			Status:     models.Pending,
			AcceptorID: acceptorID,
		}
		if err := tx.Create(&bet).Error; err != nil {
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

// on bet decline
// updates bet and updates the user balance
func Release(betID uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {

		var bet *models.Bet
		if err := tx.Where("id =?", betID).First(&bet).Error; err != nil {
			return err
		}

		var user *models.User
		// Lock the user record to prevent race conditions
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", bet.UserID).First(&user).Error; err != nil {
			return err
		}

		if err := tx.Where("id =?", betID).Update("status = ?", models.Declined).Error; err != nil {
			return err
		}

		// Update the user balance
		if err := tx.Model(&models.User{}).Where("id = ?", bet.UserID).
			UpdateColumn("balance", gorm.Expr("balance + ?", bet.Amount)).Error; err != nil {
			return err
		}

		return nil
	})
}

// if Bet wins
// 1 Release bet with status = Accepted + Result = Winning, update balance of user with bet amount
// 2 make transaction with BetID, type Lose for an Acceptor, decrease his balance (amount = bet.Amount)
// 3 make transaction with BetID, type Win for an User, increase his balance (amount = bet.Amount)
func OnWin(bet *models.Bet, acceptorItem *models.Item) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var winner *models.User
		// Lock the user record to prevent race conditions
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", bet.UserID).First(&winner).Error; err != nil {
			return err
		}

		var loser *models.User
		// Lock the user record to prevent race conditions
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", bet.AcceptorID).First(&loser).Error; err != nil {
			return err
		}

		if err := tx.Model(&bet).Updates(map[string]interface{}{"status": models.Accepted, "result": models.Winning, "acceptor_item_value": acceptorItem.Value}).Error; err != nil {
			return err
		}

		// winner
		// Create the transaction
		transaction := &models.Transaction{
			UserID: bet.UserID,
			Amount: bet.Amount,
			Type:   models.Win,
			BetID:  &bet.ID,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}
		// Update the winner balance
		if err := tx.Model(&models.User{}).Where("id = ?", bet.UserID).
			UpdateColumn("balance", gorm.Expr("balance + ?", bet.Amount*2)).Error; err != nil {
			return err
		}

		// loser
		// Create the transaction
		transaction = &models.Transaction{
			UserID: bet.AcceptorID,
			Amount: bet.Amount,
			Type:   models.Lose,
			BetID:  &bet.ID,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}
		// Update the loser balance
		if err := tx.Model(&models.User{}).Where("id = ?", bet.AcceptorID).
			UpdateColumn("balance", gorm.Expr("balance - ?", bet.Amount)).Error; err != nil {
			return err
		}
		return nil
	})
}

// if bet equals
// 1 Release bet with status = Accepted + Result = Equal, update balance of user with bet amount
func OnEqual(bet *models.Bet, acceptorItem *models.Item) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var user *models.User
		// Lock the user record to prevent race conditions
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", bet.UserID).First(&user).Error; err != nil {
			return err
		}

		if err := tx.Model(&bet).Updates(map[string]interface{}{"status": models.Accepted, "result": models.Equal, "acceptor_item_value": acceptorItem.Value}).Error; err != nil {
			return err
		}

		// Update the user balance
		if err := tx.Model(&models.User{}).Where("id = ?", bet.UserID).
			UpdateColumn("balance", gorm.Expr("balance + ?", bet.Amount)).Error; err != nil {
			return err
		}

		return nil
	})
}

// if bet lose
// 1 Change bet status = Accepted + Result = Loss
// 2 make transaction with BetID, type Lose for an User (amount = bet.Amount), do not update balance
// 3 make transaction with BetID, type Win for an Acceptor, increase his balance (amount = bet.Amount)
func OnLose(bet *models.Bet, acceptorItem *models.Item) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var winner *models.User
		// Lock the user record to prevent race conditions
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", bet.AcceptorID).First(&winner).Error; err != nil {
			return err
		}

		var loser *models.User
		// Lock the user record to prevent race conditions
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", bet.UserID).First(&loser).Error; err != nil {
			return err
		}

		if err := tx.Model(&bet).Updates(map[string]interface{}{"status": models.Accepted, "result": models.Loss, "acceptor_item_value": acceptorItem.Value}).Error; err != nil {
			return err
		}

		// loser
		// Create the transaction
		transaction := &models.Transaction{
			UserID: bet.UserID,
			Amount: bet.Amount,
			Type:   models.Lose,
			BetID:  &bet.ID,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}
		// we don't update the loser balance

		// winner
		// Create the transaction
		transaction = &models.Transaction{
			UserID: bet.AcceptorID,
			Amount: bet.Amount,
			Type:   models.Win,
			BetID:  &bet.ID,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}
		// Update the loser balance
		if err := tx.Model(&models.User{}).Where("id = ?", bet.AcceptorID).
			UpdateColumn("balance", gorm.Expr("balance + ?", bet.Amount)).Error; err != nil {
			return err
		}
		return nil
	})
}
