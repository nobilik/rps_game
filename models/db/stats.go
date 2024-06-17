package db

import "rps/models"

func GetBetStatsByStatus() map[string]interface{} {
	var res map[string]interface{}
	DB.Raw(`
	SELECT 
    COUNT(*) AS total,
    SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) AS pending,
    SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) AS declined,
    SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) AS played
	FROM bets;
	`, models.Pending, models.Declined, models.Accepted).Scan(&res)
	return res
}

func GetBetStatsByResult() map[string]interface{} {
	var res map[string]interface{}
	DB.Raw(`
	SELECT 
    SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) AS losts,
    SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) AS equals,
    SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) AS winnings
	FROM bets;
	`, models.Loss, models.Equal, models.Winning).Scan(&res)
	return res
}

// func GetBetFundStats() map[string]interface{} {
// 	var res map[string]interface{}
// 	DB.Raw(`
// 	SELECT
//     SUM(CASE WHEN status = ? THEN amount ELSE 0 END) AS losts,
//     SUM(CASE WHEN status = ? THEN amount ELSE 0 END) AS winnings
// 	FROM bets;
// 	`, models.Loss, models.Winning).Scan(&res)
// 	return res
// }

// etc...
// also if we need personal stats we can add userID param to queries
