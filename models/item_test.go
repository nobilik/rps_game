package models

import "testing"

// TestBeats function
// just simple test for game core logic
func TestBeats(t *testing.T) {
	// Определение элементов
	rock := &Item{Name: "rock", Value: 0}
	paper := &Item{Name: "paper", Value: 1}
	scissors := &Item{Name: "scissors", Value: 2}

	// test cases
	tests := []struct {
		first   *Item
		second  *Item
		want    ResultType
		message string
	}{
		{rock, rock, Equal, "rock vs rock should be equal"},
		{rock, paper, Loss, "rock vs paper should be loss"},
		{rock, scissors, Winning, "rock vs scissors should be winning"},
		{paper, rock, Winning, "paper vs rock should be winning"},
		{paper, paper, Equal, "paper vs paper should be equal"},
		{paper, scissors, Loss, "paper vs scissors should be loss"},
		{scissors, rock, Loss, "scissors vs rock should be loss"},
		{scissors, paper, Winning, "scissors vs paper should be winning"},
		{scissors, scissors, Equal, "scissors vs scissors should be equal"},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			got := tt.first.Beats(tt.second)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
