package founder

import (
	"testing"
)

func TestNewFounderGame(t *testing.T) {
	template := StartupTemplate{
		ID:               "test",
		Name:             "Test Startup",
		Type:             "SaaS",
		Description:      "A test startup",
		InitialCash:      100000,
		MonthlyBurn:      10000,
		InitialCustomers: 5,
		InitialMRR:       5000,
		AvgDealSize:      1000,
		BaseChurnRate:    0.05,
		BaseCAC:          1000,
		TargetMarketSize: 10000,
		CompetitionLevel: "medium",
		InitialTeam: map[string]int{
			"engineers":        2,
			"sales":            1,
			"customer_success": 1,
			"marketing":        0,
		},
	}

	fs := NewFounderGame("TestFounder", template, []string{})

	if fs == nil {
		t.Fatal("NewFounderGame returned nil")
	}

	if fs.FounderName != "TestFounder" {
		t.Errorf("Expected FounderName 'TestFounder', got '%s'", fs.FounderName)
	}

	if fs.Cash <= 0 {
		t.Error("Cash should be greater than 0")
	}

	if fs.MRR <= 0 {
		t.Error("MRR should be greater than 0")
	}
}

func TestAddCustomer(t *testing.T) {
	template := StartupTemplate{
		ID:               "test",
		Name:             "Test Startup",
		Type:             "SaaS",
		InitialCash:      100000,
		InitialCustomers: 5,
		InitialMRR:       5000,
		AvgDealSize:      1000,
		BaseChurnRate:    0.05,
		BaseCAC:          1000,
		TargetMarketSize: 10000,
		CompetitionLevel: "medium",
		InitialTeam: map[string]int{
			"engineers":        2,
			"sales":            1,
			"customer_success": 1,
			"marketing":        0,
		},
	}

	fs := NewFounderGame("TestFounder", template, []string{})

	initialCustomerListLen := len(fs.CustomerList)

	customer := fs.addCustomer(1500, "direct")

	if customer.DealSize != 1500 {
		t.Errorf("Expected deal size 1500, got %d", customer.DealSize)
	}

	if customer.Source != "direct" {
		t.Errorf("Expected source 'direct', got '%s'", customer.Source)
	}

	if len(fs.CustomerList) != initialCustomerListLen+1 {
		t.Errorf("Expected %d customers in list, got %d", initialCustomerListLen+1, len(fs.CustomerList))
	}
}

func TestProcessMonth(t *testing.T) {
	template := StartupTemplate{
		ID:               "test",
		Name:             "Test Startup",
		Type:             "SaaS",
		InitialCash:      100000,
		InitialCustomers: 5,
		InitialMRR:       5000,
		AvgDealSize:      1000,
		BaseChurnRate:    0.05,
		BaseCAC:          1000,
		TargetMarketSize: 10000,
		CompetitionLevel: "medium",
		InitialTeam: map[string]int{
			"engineers":        2,
			"sales":            1,
			"customer_success": 1,
			"marketing":        0,
		},
	}

	fs := NewFounderGame("TestFounder", template, []string{})

	initialTurn := fs.Turn

	messages := fs.ProcessMonth()

	if fs.Turn != initialTurn+1 {
		t.Errorf("Expected turn to be %d, got %d", initialTurn+1, fs.Turn)
	}

	if messages == nil {
		t.Error("ProcessMonth should return messages")
	}
}

func TestGetFinalScore(t *testing.T) {
	template := StartupTemplate{
		ID:               "test",
		Name:             "Test Startup",
		Type:             "SaaS",
		InitialCash:      100000,
		InitialCustomers: 5,
		InitialMRR:       5000,
		AvgDealSize:      1000,
		BaseChurnRate:    0.05,
		BaseCAC:          1000,
		TargetMarketSize: 10000,
		CompetitionLevel: "medium",
		InitialTeam: map[string]int{
			"engineers":        2,
			"sales":            1,
			"customer_success": 1,
			"marketing":        0,
		},
	}

	fs := NewFounderGame("TestFounder", template, []string{})

	// Set up game state to have MRR for valuation calculation
	fs.MRR = 50000
	fs.MonthlyGrowthRate = 0.15

	outcome, valuation, founderEquity := fs.GetFinalScore()

	// Outcome can be empty for an ongoing game - that's OK
	_ = outcome

	if valuation < 0 {
		t.Error("Valuation should not be negative")
	}

	if founderEquity < 0 || founderEquity > 100 {
		t.Errorf("Founder equity should be between 0 and 100, got %.2f", founderEquity)
	}
}
