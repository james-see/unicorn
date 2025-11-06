package game

import (
	"fmt"
	"math/rand"
)


func (gs *GameState) ProcessManagementFees() []string {
	messages := []string{}

	// Check for fee_waiver upgrade - skip fees for first 12 months
	hasFeeWaiver := false
	for _, upgradeID := range gs.PlayerUpgrades {
		if upgradeID == "fee_waiver" {
			hasFeeWaiver = true
			break
		}
	}

	// Skip fees if Fee Waiver is active and we're in first 12 months
	if hasFeeWaiver && gs.Portfolio.Turn <= 12 {
		return messages // No fees charged
	}

	// Charge management fee monthly (annual rate / 12)
	monthlyFeeRate := gs.Portfolio.AnnualManagementFee / 12.0
	fee := int64(float64(gs.Portfolio.InitialFundSize) * monthlyFeeRate)

	if fee > 0 && gs.Portfolio.Cash >= fee {
		gs.Portfolio.Cash -= fee
		gs.Portfolio.ManagementFeesCharged += fee

		// Also charge AI players
		for i := range gs.AIPlayers {
			aiFee := int64(float64(gs.AIPlayers[i].Portfolio.InitialFundSize) * monthlyFeeRate)
			if gs.AIPlayers[i].Portfolio.Cash >= aiFee {
				gs.AIPlayers[i].Portfolio.Cash -= aiFee
				gs.AIPlayers[i].Portfolio.ManagementFeesCharged += aiFee
			}
		}

		// Only show message every 12 months (annually)
		if gs.Portfolio.Turn%12 == 0 {
			annualFee := fee * 12
			messages = append(messages, fmt.Sprintf(
				"?? Annual management fee charged: $%d (2%% of fund size)",
				annualFee,
			))
		}
	}

	return messages
}

func (gs *GameState) ProcessDramaticEvents() []string {
	messages := []string{}

	for _, event := range gs.DramaticEventQueue {
		if event.ScheduledTurn == gs.Portfolio.Turn {
			// Find the company
			for i := range gs.AvailableStartups {
				if gs.AvailableStartups[i].Name == event.CompanyName {
					startup := &gs.AvailableStartups[i]

					oldValuation := startup.Valuation
					startup.Valuation = int64(float64(startup.Valuation) * event.ImpactPercent)

					// Generate message based on event type
					var eventMsg string
					var emoji string

					switch event.EventType {
					case "cofounder_split":
						emoji = "💔"
						if event.Severity == "severe" {
							eventMsg = "Co-founders had MAJOR falling out! CEO resigned. Board in chaos."
						} else if event.Severity == "moderate" {
							eventMsg = "Co-founder conflict! One founder left with equity dispute."
						} else {
							eventMsg = "Minor co-founder disagreement resolved, but caused delays."
						}
					case "scandal":
						emoji = "🔥"
						if event.Severity == "severe" {
							eventMsg = "MAJOR SCANDAL! CEO involved in workplace harassment allegations."
						} else if event.Severity == "moderate" {
							eventMsg = "PR scandal! Questionable business practices exposed."
						} else {
							eventMsg = "Minor controversy in the press, manageable."
						}
					case "lawsuit":
						emoji = "⚖️"
						if event.Severity == "severe" {
							eventMsg = "Class-action lawsuit filed! Facing $50M+ in liabilities."
						} else if event.Severity == "moderate" {
							eventMsg = "Patent infringement lawsuit. Legal costs mounting."
						} else {
							eventMsg = "Small legal dispute, expected to settle."
						}
					case "fraud":
						emoji = "🚨"
						if event.Severity == "severe" {
							eventMsg = "FRAUD DISCOVERED! CFO cooking books. SEC investigation."
						} else {
							eventMsg = "Financial irregularities found. Auditors called in."
						}
					case "data_breach":
						emoji = "🔓"
						if event.Severity == "severe" {
							eventMsg = "MASSIVE DATA BREACH! Customer data leaked. GDPR fines incoming."
						} else if event.Severity == "moderate" {
							eventMsg = "Security breach! Customer trust damaged."
						} else {
							eventMsg = "Minor security incident, quickly patched."
						}
					case "key_hire_quit":
						emoji = "👋"
						if event.Severity == "severe" {
							eventMsg = "CTO quit and joined competitor! Taking team with them."
						} else if event.Severity == "moderate" {
							eventMsg = "VP Engineering resigned. Product roadmap delayed."
						} else {
							eventMsg = "Senior engineer left. Minor setback."
						}
					case "regulatory_issue":
						emoji = "📋"
						if event.Severity == "severe" {
							eventMsg = "Regulatory crackdown! Business model under threat."
						} else {
							eventMsg = "New compliance requirements. Extra costs."
						}
					case "pivot_fail":
						emoji = "🔄"
						if event.Severity == "severe" {
							eventMsg = "Pivot FAILED! Lost key customers and burning cash fast."
						} else {
							eventMsg = "Pivot struggling. Market not responding well."
						}
					case "competitor_attack":
						emoji = "⚔️"
						if event.Severity == "severe" {
							eventMsg = "Competitor launched predatory pricing! Market share plummeting."
						} else {
							eventMsg = "New competitor with better product. Losing customers."
						}
					case "product_failure":
						emoji = "💥"
						if event.Severity == "severe" {
							eventMsg = "Major product launch FLOPPED! Customers demanding refunds."
						} else {
							eventMsg = "Product update buggy. Customer complaints rising."
						}
					default:
						emoji = "⚠️"
						eventMsg = "Unexpected crisis hit the company."
					}

					// Check if player invested
					for j := range gs.Portfolio.Investments {
						if gs.Portfolio.Investments[j].CompanyName == event.CompanyName {
							inv := &gs.Portfolio.Investments[j]
							inv.CurrentValuation = startup.Valuation

							valuationDrop := oldValuation - startup.Valuation
							dropPercent := float64(valuationDrop) / float64(oldValuation) * 100

							messages = append(messages, fmt.Sprintf(
								"%s %s: %s (Valuation: $%s → $%s, -%.0f%%)",
								emoji,
								event.CompanyName,
								eventMsg,
								formatCurrency(oldValuation),
								formatCurrency(startup.Valuation),
								dropPercent,
							))
							break
						}
					}

					// Update AI investments
					for k := range gs.AIPlayers {
						for j := range gs.AIPlayers[k].Portfolio.Investments {
							if gs.AIPlayers[k].Portfolio.Investments[j].CompanyName == event.CompanyName {
								gs.AIPlayers[k].Portfolio.Investments[j].CurrentValuation = startup.Valuation
							}
						}
					}
				}
			}
		}
	}

	return messages
}

func (gs *GameState) ProcessFundingRounds() []string {
	messages := []string{}

	for _, event := range gs.FundingRoundQueue {
		if event.ScheduledTurn == gs.Portfolio.Turn {
			// Find the company
			for i := range gs.AvailableStartups {
				if gs.AvailableStartups[i].Name == event.CompanyName {
					startup := &gs.AvailableStartups[i]

					var preMoneyVal int64
					var postMoneyVal int64
					var dilutionFactor float64

					if event.IsDownRound {
						// Down round: pre-money is 60-90% of current valuation
						downFactor := 0.6 + rand.Float64()*0.3 // 60%-90%
						preMoneyVal = int64(float64(startup.Valuation) * downFactor)
						postMoneyVal = preMoneyVal + event.RaiseAmount
						dilutionFactor = float64(preMoneyVal) / float64(postMoneyVal)

						// Check if any investor has board seat - down rounds require board approval
						if gs.HasAnyBoardSeat(event.CompanyName) {
							// Only create vote if player has board seat (player votes, AI votes are simulated)
							if gs.HasBoardSeat(event.CompanyName) {
								// Create board vote for down round
								vote := BoardVote{
									CompanyName:  event.CompanyName,
									VoteType:     "down_round",
									Title:        fmt.Sprintf("Down Round: $%s at $%s pre-money", formatCurrency(event.RaiseAmount), formatCurrency(preMoneyVal)),
									Description:  fmt.Sprintf("%s proposes raising $%s in a DOWN ROUND at $%s pre-money (down from $%s). This will significantly dilute your equity.", event.CompanyName, formatCurrency(event.RaiseAmount), formatCurrency(preMoneyVal), formatCurrency(startup.Valuation)),
									OptionA:      "Approve",
									OptionB:      "Reject",
									ConsequenceA: fmt.Sprintf("Down round approved. Company raises $%s at reduced valuation.", formatCurrency(event.RaiseAmount)),
									ConsequenceB: "Down round rejected. Company must find alternative funding or accept worse terms.",
									RequiresVote: true,
									Turn:         gs.Portfolio.Turn,
									Metadata: map[string]interface{}{
										"raiseAmount":      event.RaiseAmount,
										"preMoneyVal":      preMoneyVal,
										"postMoneyVal":     postMoneyVal,
										"currentValuation": startup.Valuation,
									},
								}
								gs.PendingBoardVotes = append(gs.PendingBoardVotes, vote)
								messages = append(messages, fmt.Sprintf(
									"🏛️  BOARD VOTE REQUIRED: %s proposes a DOWN ROUND. Vote will be required.",
									event.CompanyName,
								))
								continue // Skip processing this round until vote is complete
							}
						}
					} else {
						// Normal round
						preMoneyVal = startup.Valuation
						postMoneyVal = preMoneyVal + event.RaiseAmount
						dilutionFactor = float64(preMoneyVal) / float64(postMoneyVal)
					}

					// Check for Portfolio Insurance upgrade
					hasPortfolioInsurance := false
					for _, upgradeID := range gs.PlayerUpgrades {
						if upgradeID == "portfolio_insurance" {
							hasPortfolioInsurance = true
							break
						}
					}

					// Update player's investment if they invested in this company
					for j := range gs.Portfolio.Investments {
						if gs.Portfolio.Investments[j].CompanyName == event.CompanyName {
							inv := &gs.Portfolio.Investments[j]

							// Check if Portfolio Insurance protects this investment from down rounds
							shouldProtect := false
							if event.IsDownRound && hasPortfolioInsurance && !gs.InsuranceUsed {
								// First down round hit by Portfolio Insurance is protected
								shouldProtect = true
								gs.InsuranceUsed = true
								gs.ProtectedCompany = event.CompanyName
							}

							// If follow-on investment was made this turn, equity was already recalculated
							// in MakeFollowOnInvestment based on post-money valuation, so we don't dilute again
							if !inv.FollowOnThisTurn {
								// Normal case: dilute existing equity (unless protected by Portfolio Insurance)
								oldEquity := inv.EquityPercent
								if !shouldProtect {
									inv.EquityPercent *= dilutionFactor
								} else {
									// Portfolio Insurance protects this investment - no dilution
									messages = append(messages, fmt.Sprintf(
										"🛡️  PORTFOLIO INSURANCE: %s protected from down round dilution! Equity remains at %.2f%%",
										event.CompanyName,
										oldEquity,
									))
								}

								// Record the round
								inv.Rounds = append(inv.Rounds, FundingRound{
									RoundName:        event.RoundName,
									PreMoneyVal:      preMoneyVal,
									InvestmentAmount: event.RaiseAmount,
									PostMoneyVal:     postMoneyVal,
									Month:            gs.Portfolio.Turn,
								})

								// Only show dilution messages if not protected by Portfolio Insurance
								if !shouldProtect {
									if event.IsDownRound {
										messages = append(messages, fmt.Sprintf(
											"⚠️  %s raised $%s in DOWN ROUND (%s)! Valuation dropped. Equity: %.2f%% → %.2f%%",
											event.CompanyName,
											formatCurrency(event.RaiseAmount),
											event.RoundName,
											oldEquity,
											inv.EquityPercent,
										))
									} else {
										messages = append(messages, fmt.Sprintf(
											"🚀 %s raised $%s in %s round! Your equity diluted from %.2f%% to %.2f%%",
											event.CompanyName,
											formatCurrency(event.RaiseAmount),
											event.RoundName,
											oldEquity,
											inv.EquityPercent,
										))
									}
								}
							} else {
								// Follow-on investment case: equity already calculated correctly, just record the round
								oldEquity := inv.EquityPercent
								inv.Rounds = append(inv.Rounds, FundingRound{
									RoundName:        event.RoundName,
									PreMoneyVal:      preMoneyVal,
									InvestmentAmount: event.RaiseAmount,
									PostMoneyVal:     postMoneyVal,
									Month:            gs.Portfolio.Turn,
								})

								// Reset flag for next turn
								inv.FollowOnThisTurn = false

								messages = append(messages, fmt.Sprintf(
									"🚀 %s raised $%s in %s round! Your equity: %.2f%% (includes your follow-on investment)",
									event.CompanyName,
									formatCurrency(event.RaiseAmount),
									event.RoundName,
									oldEquity,
								))
							}
						}
					}

					// Update AI players' investments
					for k := range gs.AIPlayers {
						for j := range gs.AIPlayers[k].Portfolio.Investments {
							if gs.AIPlayers[k].Portfolio.Investments[j].CompanyName == event.CompanyName {
								inv := &gs.AIPlayers[k].Portfolio.Investments[j]
								inv.EquityPercent *= dilutionFactor
								inv.Rounds = append(inv.Rounds, FundingRound{
									RoundName:        event.RoundName,
									PreMoneyVal:      preMoneyVal,
									InvestmentAmount: event.RaiseAmount,
									PostMoneyVal:     postMoneyVal,
									Month:            gs.Portfolio.Turn,
								})
							}
						}
					}

					// Update company valuation
					startup.Valuation = postMoneyVal

					// Also update current valuation for all investors
					for j := range gs.Portfolio.Investments {
						if gs.Portfolio.Investments[j].CompanyName == event.CompanyName {
							gs.Portfolio.Investments[j].CurrentValuation = postMoneyVal
						}
					}
					for k := range gs.AIPlayers {
						for j := range gs.AIPlayers[k].Portfolio.Investments {
							if gs.AIPlayers[k].Portfolio.Investments[j].CompanyName == event.CompanyName {
								gs.AIPlayers[k].Portfolio.Investments[j].CurrentValuation = postMoneyVal
							}
						}
					}
				}
			}
		}
	}

	return messages
}

func (gs *GameState) ProcessAcquisitions() []string {
	messages := []string{}

	for _, event := range gs.AcquisitionQueue {
		if event.ScheduledTurn == gs.Portfolio.Turn {
			// Find the company
			for i := range gs.AvailableStartups {
				if gs.AvailableStartups[i].Name == event.CompanyName {
					startup := &gs.AvailableStartups[i]

					// Calculate EBITDA (approximation: annual net income)
					annualEBITDA := startup.NetIncome * 12
					if annualEBITDA < 0 {
						// For unprofitable companies, use revenue multiple instead
						annualEBITDA = startup.MonthlyRevenue * 12
						event.OfferMultiple *= 0.3 // Lower multiple for revenue-based
					}

					// Calculate acquisition offer
					offerValue := int64(float64(annualEBITDA) * event.OfferMultiple)

					// Ensure minimum offer value
					if offerValue < startup.Valuation/2 {
						offerValue = startup.Valuation / 2
					}

					// Check if player invested in this company
					for j := range gs.Portfolio.Investments {
						if gs.Portfolio.Investments[j].CompanyName == event.CompanyName {
							inv := &gs.Portfolio.Investments[j]

							// Calculate payout
							payout := int64((inv.EquityPercent / 100.0) * float64(offerValue))
							returnMultiple := float64(payout) / float64(inv.AmountInvested)

							// If player has board seat, require board vote for acquisitions (unless bad due diligence)
							if inv.Terms.HasBoardSeat && event.DueDiligence != "bad" {
								// Create board vote
								vote := BoardVote{
									CompanyName:  event.CompanyName,
									VoteType:     "acquisition",
									Title:        fmt.Sprintf("Acquisition Offer: $%s", formatCurrency(offerValue)),
									Description:  fmt.Sprintf("Acquirer offers $%s (%.1fx EBITDA) for %s. Your payout would be $%s (%.1fx return).", formatCurrency(offerValue), event.OfferMultiple, event.CompanyName, formatCurrency(payout), returnMultiple),
									OptionA:      "Accept",
									OptionB:      "Reject",
									ConsequenceA: fmt.Sprintf("Acquisition proceeds. You receive $%s.", formatCurrency(payout)),
									ConsequenceB: "Acquisition rejected. Company continues operating independently.",
									RequiresVote: true,
									Turn:         gs.Portfolio.Turn,
									Metadata: map[string]interface{}{
										"offerValue":       offerValue,
										"currentValuation": startup.Valuation,
										"dueDiligence":     event.DueDiligence,
										"offerMultiple":    event.OfferMultiple,
									},
								}
								gs.PendingBoardVotes = append(gs.PendingBoardVotes, vote)
								messages = append(messages, fmt.Sprintf(
									"🏛️  BOARD VOTE REQUIRED: %s received acquisition offer of $%s. Vote will be required.",
									event.CompanyName,
									formatCurrency(offerValue),
								))
								break // Don't execute acquisition yet - wait for vote
							}

							// Add acquisition message based on due diligence
							switch event.DueDiligence {
							case "bad":
								messages = append(messages, fmt.Sprintf(
									"⚠️  %s acquisition FELL THROUGH! Due diligence issues. Offer was $%s (%.1fx EBITDA)",
									event.CompanyName,
									formatCurrency(offerValue),
									event.OfferMultiple,
								))
							case "good":
								messages = append(messages, fmt.Sprintf(
									"🎉 %s ACQUIRED for $%s (%.1fx EBITDA)! Your %.2f%% = $%s (%.1fx return)",
									event.CompanyName,
									formatCurrency(offerValue),
									event.OfferMultiple,
									inv.EquityPercent,
									formatCurrency(payout),
									returnMultiple,
								))
								// Execute acquisition
								gs.Portfolio.Cash += payout
								// Remove investment from portfolio
								gs.Portfolio.Investments = append(gs.Portfolio.Investments[:j], gs.Portfolio.Investments[j+1:]...)
							default: // normal
								messages = append(messages, fmt.Sprintf(
									"💰 %s ACQUIRED for $%s (%.1fx EBITDA)! Your %.2f%% = $%s (%.1fx return)",
									event.CompanyName,
									formatCurrency(offerValue),
									event.OfferMultiple,
									inv.EquityPercent,
									formatCurrency(payout),
									returnMultiple,
								))
								// Execute acquisition
								gs.Portfolio.Cash += payout
								// Remove investment from portfolio
								gs.Portfolio.Investments = append(gs.Portfolio.Investments[:j], gs.Portfolio.Investments[j+1:]...)
							}
							break
						}
					}

					// Handle AI player acquisitions
					if event.DueDiligence != "bad" {
						for k := range gs.AIPlayers {
							for j := len(gs.AIPlayers[k].Portfolio.Investments) - 1; j >= 0; j-- {
								if gs.AIPlayers[k].Portfolio.Investments[j].CompanyName == event.CompanyName {
									inv := &gs.AIPlayers[k].Portfolio.Investments[j]
									payout := int64((inv.EquityPercent / 100.0) * float64(offerValue))
									gs.AIPlayers[k].Portfolio.Cash += payout
									// Remove from AI portfolio
									gs.AIPlayers[k].Portfolio.Investments = append(
										gs.AIPlayers[k].Portfolio.Investments[:j],
										gs.AIPlayers[k].Portfolio.Investments[j+1:]...,
									)
									break
								}
							}
						}
					}
				}
			}
		}
	}

	return messages
}