package founder

import (
	"fmt"
	"math/rand"
)

// InitializeSecurity initializes the security system
func InitializeSecurity(fs *FounderState) {
	if fs.SecurityPosture == nil {
		fs.SecurityPosture = &SecurityPosture{
			SecurityScore:    50, // Start at 50/100
			ComplianceCerts: []string{},
			SecurityTeamSize: 0,
			SecurityBudget:   0,
			LastAudit:        0,
			Vulnerabilities:  0,
			BugBountyActive:  false,
			BugBountyBudget:  0,
		}
	}
	if fs.SecurityIncidents == nil {
		fs.SecurityIncidents = []SecurityIncident{}
	}
}

// CanHaveSecurityIncidents checks if company is large enough to be a target
func (fs *FounderState) CanHaveSecurityIncidents() bool {
	// Unlock: $500k+ MRR OR 100+ customers
	arr := fs.MRR * 12
	return arr >= 500000 || fs.Customers >= 100
}

// SpawnSecurityIncident generates a security incident
func (fs *FounderState) SpawnSecurityIncident() *SecurityIncident {
	if !fs.CanHaveSecurityIncidents() {
		return nil
	}

	// Probability: 2% per month, increases with low security score
	baseProbability := 0.02
	if fs.SecurityPosture.SecurityScore < 50 {
		baseProbability = 0.05 // 5% if security score < 50
	}
	if fs.SecurityPosture.SecurityScore < 30 {
		baseProbability = 0.10 // 10% if security score < 30
	}

	if rand.Float64() > baseProbability {
		return nil
	}

	// Incident types
	incidentTypes := []string{"data_breach", "ransomware", "ddos", "insider_threat", "vulnerability"}
	incidentType := incidentTypes[rand.Intn(len(incidentTypes))]

	// Severity based on security score
	severityRoll := rand.Float64()
	severity := "low"
	if fs.SecurityPosture.SecurityScore < 40 {
		if severityRoll < 0.3 {
			severity = "critical"
		} else if severityRoll < 0.6 {
			severity = "high"
		} else {
			severity = "medium"
		}
	} else if fs.SecurityPosture.SecurityScore < 60 {
		if severityRoll < 0.2 {
			severity = "high"
		} else if severityRoll < 0.5 {
			severity = "medium"
		}
	}

	// Customers affected
	customersAffected := 0
	switch severity {
	case "critical":
		customersAffected = int(float64(fs.Customers) * (0.15 + rand.Float64()*0.15)) // 15-30%
	case "high":
		customersAffected = int(float64(fs.Customers) * (0.05 + rand.Float64()*0.10)) // 5-15%
	case "medium":
		customersAffected = int(float64(fs.Customers) * (0.01 + rand.Float64()*0.05)) // 1-5%
	case "low":
		customersAffected = int(float64(fs.Customers) * (0.001 + rand.Float64()*0.01)) // 0.1-1%
	}

	// Data exposed
	dataTypes := []string{"PII", "financial", "health", "none"}
	dataExposed := dataTypes[rand.Intn(len(dataTypes))]

	// Response costs
	responseCost := int64(50000 + rand.Int63n(150000)) // $50-200k
	legalCosts := int64(0)
	if severity == "critical" || severity == "high" {
		legalCosts = int64(200000 + rand.Int63n(300000)) // $200-500k
	}

	// Reputation damage
	reputationDamage := 0.1
	if severity == "critical" {
		reputationDamage = 0.3
	} else if severity == "high" {
		reputationDamage = 0.2
	} else if severity == "medium" {
		reputationDamage = 0.1
	}

	incident := SecurityIncident{
		Type:             incidentType,
		Severity:         severity,
		Month:            fs.Turn,
		CustomersAffected: customersAffected,
		DataExposed:      dataExposed,
		ResponseCost:     responseCost,
		LegalCosts:       legalCosts,
		ReputationDamage: reputationDamage,
		Resolved:         false,
		ResolutionMonth:  0,
		ResponseActions:  []string{},
	}

	fs.SecurityIncidents = append(fs.SecurityIncidents, incident)
	fs.ActiveSecurityIncident = &incident

	return &incident
}

// RespondToSecurityIncident handles security incident response
func (fs *FounderState) RespondToSecurityIncident(action string) error {
	if fs.ActiveSecurityIncident == nil {
		return fmt.Errorf("no active security incident")
	}

	incident := fs.ActiveSecurityIncident

	// Response actions
	switch action {
	case "contain":
		cost := int64(50000 + rand.Int63n(100000)) // $50-150k
		if cost > fs.Cash {
			return fmt.Errorf("insufficient cash (need $%s)", formatCurrency(cost))
		}
		fs.Cash -= cost
		incident.ResponseActions = append(incident.ResponseActions, "containment")
		incident.ResponseCost += cost
		incident.ReputationDamage *= 0.8 // Reduce damage by 20%

	case "investigate":
		cost := int64(100000 + rand.Int63n(200000)) // $100-300k
		if cost > fs.Cash {
			return fmt.Errorf("insufficient cash (need $%s)", formatCurrency(cost))
		}
		fs.Cash -= cost
		incident.ResponseActions = append(incident.ResponseActions, "forensic_investigation")
		incident.ResponseCost += cost
		incident.ReputationDamage *= 0.7 // Reduce damage by 30%

	case "notify":
		cost := int64(20000 + rand.Int63n(80000)) // $20-100k
		if cost > fs.Cash {
			return fmt.Errorf("insufficient cash (need $%s)", formatCurrency(cost))
		}
		fs.Cash -= cost
		incident.ResponseActions = append(incident.ResponseActions, "customer_notification")
		incident.ResponseCost += cost
		// Transparency reduces churn impact

	case "defend":
		cost := int64(200000 + rand.Int63n(300000)) // $200-500k
		if cost > fs.Cash {
			return fmt.Errorf("insufficient cash (need $%s)", formatCurrency(cost))
		}
		fs.Cash -= cost
		incident.ResponseActions = append(incident.ResponseActions, "legal_defense")
		incident.LegalCosts += cost

	case "resolve":
		// Resolve incident (requires all actions taken)
		if len(incident.ResponseActions) < 2 {
			return fmt.Errorf("must take containment and investigation actions first")
		}
		incident.Resolved = true
		incident.ResolutionMonth = fs.Turn
		fs.ActiveSecurityIncident = nil

		// Improve security score after resolution (lessons learned)
		fs.SecurityPosture.SecurityScore += 5
		if fs.SecurityPosture.SecurityScore > 100 {
			fs.SecurityPosture.SecurityScore = 100
		}

	default:
		return fmt.Errorf("invalid action: %s", action)
	}

	return nil
}

// InvestInSecurity invests in security improvements
func (fs *FounderState) InvestInSecurity(monthlyBudget int64) error {
	if monthlyBudget > fs.Cash {
		return fmt.Errorf("insufficient cash")
	}

	fs.Cash -= monthlyBudget
	fs.SecurityPosture.SecurityBudget += monthlyBudget

	// Improve security score
	scoreIncrease := int(monthlyBudget / 10000) // $10k = +1 score
	fs.SecurityPosture.SecurityScore += scoreIncrease
	if fs.SecurityPosture.SecurityScore > 100 {
		fs.SecurityPosture.SecurityScore = 100
	}

	// Reduce vulnerabilities
	if fs.SecurityPosture.Vulnerabilities > 0 {
		vulnsFixed := int(monthlyBudget / 50000) // $50k fixes 1 vulnerability
		fs.SecurityPosture.Vulnerabilities -= vulnsFixed
		if fs.SecurityPosture.Vulnerabilities < 0 {
			fs.SecurityPosture.Vulnerabilities = 0
		}
	}

	return nil
}

// HireSecurityTeam hires security team members
func (fs *FounderState) HireSecurityTeam(count int) error {
	cost := int64(count) * 8333 // $100k/year = $8,333/month
	if cost > fs.Cash {
		return fmt.Errorf("insufficient cash")
	}

	fs.Cash -= cost
	fs.SecurityPosture.SecurityTeamSize += count

	// Improve security score
	fs.SecurityPosture.SecurityScore += count * 5
	if fs.SecurityPosture.SecurityScore > 100 {
		fs.SecurityPosture.SecurityScore = 100
	}

	return nil
}

// GetComplianceCertification achieves a compliance certification
func (fs *FounderState) GetComplianceCertification(cert string) error {
	validCerts := map[string]bool{
		"SOC2":     true,
		"ISO27001": true,
		"HIPAA":    true,
		"GDPR":     true,
	}
	if !validCerts[cert] {
		return fmt.Errorf("invalid certification: %s", cert)
	}

	// Check if already have it
	for _, c := range fs.SecurityPosture.ComplianceCerts {
		if c == cert {
			return fmt.Errorf("already have %s certification", cert)
		}
	}

	// Cost: $50-200k depending on cert
	cost := int64(50000)
	switch cert {
	case "SOC2":
		cost = 50000
	case "ISO27001":
		cost = 100000
	case "HIPAA":
		cost = 150000
	case "GDPR":
		cost = 200000
	}

	if cost > fs.Cash {
		return fmt.Errorf("insufficient cash (need $%s)", formatCurrency(cost))
	}

	fs.Cash -= cost
	fs.SecurityPosture.ComplianceCerts = append(fs.SecurityPosture.ComplianceCerts, cert)

	// Improve security score
	fs.SecurityPosture.SecurityScore += 10
	if fs.SecurityPosture.SecurityScore > 100 {
		fs.SecurityPosture.SecurityScore = 100
	}

	return nil
}

// ProcessSecurityIncidents processes active security incidents
func (fs *FounderState) ProcessSecurityIncidents() []string {
	var messages []string

	if fs.ActiveSecurityIncident != nil {
		incident := fs.ActiveSecurityIncident

		// Unresolved incidents cause ongoing damage
		if !incident.Resolved {
			// Churn from affected customers
			churnRate := 0.05 // 5% churn per month
			if incident.Severity == "critical" {
				churnRate = 0.20 // 20% churn
			} else if incident.Severity == "high" {
				churnRate = 0.10 // 10% churn
			}

			customersLost := int(float64(incident.CustomersAffected) * churnRate)
			if customersLost > 0 {
				mrrLost := int64(customersLost) * fs.AvgDealSize
				fs.Customers -= customersLost
				fs.DirectCustomers -= customersLost
				fs.DirectMRR -= mrrLost
				fs.syncMRR()

				messages = append(messages, fmt.Sprintf("🔒 Security incident: Lost %d customers (churn: %.0f%%)", customersLost, churnRate*100))
			}

			// CAC increase due to reputation damage
			cacMultiplier := 1.0 + incident.ReputationDamage
			fs.BaseCAC = int64(float64(fs.BaseCAC) * cacMultiplier)
		}
	}

	return messages
}

