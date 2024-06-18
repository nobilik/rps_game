package core

import (
	"log"
	"math/rand"
	"rps/models"
	"rps/models/db"
	"time"
)

var (
	AI           *models.User
	BetInterval  = 60
	MaxBetAmount = 100
)

// here is AI player logic

func init() {
	var err error
	AI, err = db.GetUserByID(10000)
	if err != nil {
		initAI()
	}

}
func initAI() {
	AI = &models.User{
		Login:   "Alberto",
		Balance: 1000,
		IsAI:    true,
	}
	AI.ID = 10000 // for now we use one AI and restrict it by DB id
	db.CreateUser(AI)
}

func AcceptBetByAI(acceptorID uint, initialItem *models.Item, bet *models.Bet) error {
	if !isAI(acceptorID) {
		return nil
	}
	time.Sleep(getSec(10))
	item := getAIMove()
	result := initialItem.Beats(item)

	// onLose && onWin from bet initiator side
	var err error
	switch result {
	case models.Loss: // means that acceptor wins so we send notification (respond to acceptor)
		err = db.OnLose(bet, item)
		AI.Balance += bet.Amount
	case models.Equal:
		err = db.OnEqual(bet, item)
	case models.Winning:
		err = db.OnWin(bet, item)
		AI.Balance -= bet.Amount
		checkAIBalance()
	}
	return err
}

// we call it from goroutine
// there are no restrictions for multi bet
func PlaceBetByAI() {
	for {
		time.Sleep(getSec(BetInterval))
		user := getRandomUser()
		if user == nil {
			return
		}

		item := getAIMove()
		err := db.Freeze(AI.ID, user.ID, getBetAmount(), *item)
		if err != nil {
			log.Printf("error on bet %v", err)
		}
		checkAIBalance()
	}
}

func getAIMove() *models.Item {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return models.Items[r.Intn(len(models.Items))]
}

func getSec(topLimit int) time.Duration {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return time.Duration(r.Intn(topLimit)) * time.Second
}

func isAI(id uint) bool {
	return id == AI.ID
}

func checkAIBalance() error {
	if AI.Balance > 1000 {
		return nil
	}
	err := db.TopUp(AI.ID, 1000)
	if err == nil {
		AI.Balance += 1000
	}
	return err
}

func getRandomUser() *models.User {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	users, err := db.GetUsers(AI.ID)
	if err != nil || users == nil {
		return nil
	}
	return users[r.Intn(len(users))]
}

func getBetAmount() float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	intAmount := r.Intn(MaxBetAmount-1) + 1
	return float64(intAmount)
}
