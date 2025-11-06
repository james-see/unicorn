package game

import (
	"fmt"
	"math/rand"
	"strings"
)


func (gs *GameState) HasBoardSeat(companyName string) bool {
	for _, inv := range gs.Portfolio.Investments {
		if inv.CompanyName == companyName && inv.Terms.HasBoardSeat {
			return true
		}
	}
	return false
}

func (gs *GameState) HasAnyBoardSeat(companyName string) bool {
	// Check player
	if gs.HasBoardSeat(companyName) {
		return true
	}

	// Check AI players
	for _, ai := range gs.AIPlayers {
		for _, inv := range ai.Portfolio.Investments {
			if inv.CompanyName == companyName && inv.Terms.HasBoardSeat {
				return true
			}
		}
	}

	return false
}

func (gs *GameState) GetPendingBoardVotes() []BoardVote {
	pending := []BoardVote{}
	for _, vote := range gs.PendingBoardVotes {
		if gs.HasBoardSeat(vote.CompanyName) {
			pending = append(pending, vote)
		}
	}
	return pending
}

func (gs *GameState) GetNextBoardVotePreview() string {
	nextTurn := gs.Portfolio.Turn + 1

	// Check for upcoming acquisitions
	for _, event := range gs.AcquisitionQueue {
		if event.ScheduledTurn == nextTurn {
			// Check if player has board seat
			for _, inv := range gs.Portfolio.Investments {
				if inv.CompanyName == event.CompanyName && inv.Terms.HasBoardSeat {
					// Find startup to get current valuation
					for _, startup := range gs.AvailableStartups {
						if startup.Name == event.CompanyName {
							annualEBITDA := startup.NetIncome * 12
							offerMultiple := event.OfferMultiple
							if annualEBITDA < 0 {
								annualEBITDA = startup.MonthlyRevenue * 12
								offerMultiple *= 0.3
							}
							offerValue := int64(float64(annualEBITDA) * offerMultiple)
							if offerValue < startup.Valuation/2 {
								offerValue = startup.Valuation / 2
							}
							return fmt.Sprintf("⚠️  Next turn: %s will receive an acquisition offer (~$%s). Board vote required!",
								event.CompanyName, formatCurrency(offerValue))
						}
					}
				}
			}
		}
	}

	// Check for upcoming down rounds
	for _, event := range gs.FundingRoundQueue {
		if event.ScheduledTurn == nextTurn && event.IsDownRound {
			for _, inv := range gs.Portfolio.Investments {
				if inv.CompanyName == event.CompanyName && inv.Terms.HasBoardSeat {
					return fmt.Sprintf("⚠️  Next turn: %s proposes a DOWN ROUND. Board vote required!",
						event.CompanyName)
				}
			}
		}
	}

	return ""
}

func (gs *GameState) ProcessBoardVote(voteIndex int, playerVote string) (string, bool, error) {
	if voteIndex < 0 || voteIndex >= len(gs.PendingBoardVotes) {
		return "", false, fmt.Errorf("invalid vote index")
	}

	vote := &gs.PendingBoardVotes[voteIndex]
	if !gs.HasBoardSeat(vote.CompanyName) {
		return "", false, fmt.Errorf("you do not have a board seat for %s", vote.CompanyName)
	}

	// Normalize vote
	playerVote = strings.ToLower(strings.TrimSpace(playerVote))
	var votedForA bool
	if playerVote == "a" || playerVote == "accept" || playerVote == "approve" || playerVote == "yes" || playerVote == "1" {
		votedForA = true
	} else if playerVote == "b" || playerVote == "reject" || playerVote == "disapprove" || playerVote == "no" || playerVote == "2" {
		votedForA = false
	} else {
		return "", false, fmt.Errorf("invalid vote choice")
	}

	// Simulate board vote: player vote + AI board members vote
	// Player vote counts based on board seat multiplier
	playerVoteWeight := 1
	for _, inv := range gs.Portfolio.Investments {
		if inv.CompanyName == vote.CompanyName && inv.Terms.HasBoardSeat {
			playerVoteWeight = inv.Terms.BoardSeatMultiplier
			if playerVoteWeight == 0 {
				playerVoteWeight = 1 // Default to 1 if not set
			}
			break
		}
	}

	aiVotesA := 0
	aiVotesB := 0

	// Count AI board members (simulate other investors with board seats)
	numAIBoardMembers := 2 + rand.Intn(2) // 2-3 AI board members

	for i := 0; i < numAIBoardMembers; i++ {
		// AI votes based on their strategy
		voteChance := 0.5
		if vote.VoteType == "acquisition" {
			// AI more likely to accept acquisitions if good terms
			if offerValue, ok := vote.Metadata["offerValue"].(int64); ok {
				if currentVal, ok := vote.Metadata["currentValuation"].(int64); ok {
					if offerValue >= currentVal {
						voteChance = 0.7 // 70% chance to accept good offers
					}
				}
			}
		} else if vote.VoteType == "down_round" {
			// AI less likely to accept down rounds
			voteChance = 0.3
		}

		if rand.Float64() < voteChance {
			aiVotesA++
		} else {
			aiVotesB++
		}
	}

	// Count votes
	totalVotesA := aiVotesA
	totalVotesB := aiVotesB
	if votedForA {
		totalVotesA += playerVoteWeight
	} else {
		totalVotesB += playerVoteWeight
	}

	// Determine outcome
	votePassed := totalVotesA > totalVotesB

	// Store vote result in metadata for execution
	voteCopy := *vote
	voteCopy.Metadata["votePassed"] = votePassed
	voteCopy.Metadata["playerVotedForA"] = votedForA

	// Remove vote from pending list
	gs.PendingBoardVotes = append(gs.PendingBoardVotes[:voteIndex], gs.PendingBoardVotes[voteIndex+1:]...)

	// Generate result message
	voteOutcome := fmt.Sprintf("Board Vote: %d/%d voted for %s, %d/%d voted for %s. ",
		totalVotesA,
		totalVotesA+totalVotesB,
		vote.OptionA,
		totalVotesB,
		totalVotesA+totalVotesB,
		vote.OptionB)

	if votePassed {
		voteOutcome += vote.ConsequenceA
	} else {
		voteOutcome += vote.ConsequenceB
	}

	return voteOutcome, votePassed, nil
}

func (gs *GameState) ExecuteBoardVoteOutcome(vote BoardVote, passed bool) []string {
	messages := []string{}

	switch vote.VoteType {
	case "acquisition":
		if passed {
			// Acquisition approved - execute it
			if offerValue, ok := vote.Metadata["offerValue"].(int64); ok {
				companyName := vote.CompanyName
				for j := range gs.Portfolio.Investments {
					if gs.Portfolio.Investments[j].CompanyName == companyName {
						inv := &gs.Portfolio.Investments[j]
						payout := int64((inv.EquityPercent / 100.0) * float64(offerValue))
						returnMultiple := float64(payout) / float64(inv.AmountInvested)

						messages = append(messages, fmt.Sprintf(
							"🎉 %s ACQUIRED (Board Approved)! Your %.2f%% = $%s (%.1fx return)",
							companyName,
							inv.EquityPercent,
							formatCurrency(payout),
							returnMultiple,
						))

						gs.Portfolio.Cash += payout
						gs.Portfolio.Investments = append(gs.Portfolio.Investments[:j], gs.Portfolio.Investments[j+1:]...)
						break
					}
				}
			}
		} else {
			messages = append(messages, fmt.Sprintf(
				"❌ %s acquisition REJECTED by board. Company continues operating.",
				vote.CompanyName,
			))
		}
	case "down_round":
		if passed {
			// Down round approved - it proceeds (already handled in ProcessFundingRounds)
			messages = append(messages, fmt.Sprintf(
				"✅ Board approved down round for %s. Round proceeds.",
				vote.CompanyName,
			))
		} else {
			// Down round rejected - company must find alternative or accept worse terms
			messages = append(messages, fmt.Sprintf(
				"❌ Board rejected down round for %s. Company must find alternative funding.",
				vote.CompanyName,
			))
		}
	}

	return messages
}