package game

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	gs := NewGame("TestPlayer", "TestPlayer Capital", MediumDifficulty, []string{})
	
	if gs == nil {
		t.Fatal("NewGame returned nil")
	}
	
	if gs.PlayerName != "TestPlayer" {
		t.Errorf("Expected PlayerName 'TestPlayer', got '%s'", gs.PlayerName)
	}
	
	if gs.Portfolio.Cash <= 0 {
		t.Error("Portfolio cash should be greater than 0")
	}
	
	if len(gs.AvailableStartups) == 0 {
		t.Error("Should have available startups")
	}
}

func TestMakeInvestment(t *testing.T) {
	gs := NewGame("TestPlayer", "TestPlayer Capital", MediumDifficulty, []string{})
	
	initialCash := gs.Portfolio.Cash
	investAmount := int64(50000)
	
	err := gs.MakeInvestment(0, investAmount)
	
	if err != nil {
		t.Fatalf("MakeInvestment failed: %v", err)
	}
	
	if gs.Portfolio.Cash != initialCash-investAmount {
		t.Errorf("Expected cash to be %d, got %d", initialCash-investAmount, gs.Portfolio.Cash)
	}
	
	if len(gs.Portfolio.Investments) == 0 {
		t.Error("Should have at least one investment")
	}
}

func TestProcessTurn(t *testing.T) {
	gs := NewGame("TestPlayer", "TestPlayer Capital", MediumDifficulty, []string{})
	
	// Make an investment first
	gs.MakeInvestment(0, 50000)
	
	initialTurn := gs.Portfolio.Turn
	
	messages := gs.ProcessTurn()
	
	if gs.Portfolio.Turn != initialTurn+1 {
		t.Errorf("Expected turn to be %d, got %d", initialTurn+1, gs.Portfolio.Turn)
	}
	
	if messages == nil {
		t.Error("ProcessTurn should return messages")
	}
}

func TestGetFinalScore(t *testing.T) {
	gs := NewGame("TestPlayer", "TestPlayer Capital", MediumDifficulty, []string{})
	
	// Make an investment
	gs.MakeInvestment(0, 50000)
	
	netWorth, roi, successfulExits := gs.GetFinalScore()
	
	if netWorth == 0 {
		t.Error("Final net worth should not be 0")
	}
	
	if roi < 0 {
		t.Error("ROI should not be negative")
	}
	
	if successfulExits < 0 {
		t.Error("Successful exits should not be negative")
	}
}

