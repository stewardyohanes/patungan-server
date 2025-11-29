package utils

import (
	"errors"
	"math"
	"patungan-server/internal/models"

	"github.com/google/uuid"
)

type SplitCalculator struct{}

func NewSplitCalculator() *SplitCalculator {
	return &SplitCalculator{}
}

// CalculateEqualSplit divides total amount equally among participants
func (sc *SplitCalculator) CalculateEqualSplit(totalAmount float64, participantCount int) (float64, error) {
	if participantCount <= 0 {
		return 0, errors.New("participant count must be greater than 0")
	}

	// Round to 2 decimal places
	amount := math.Round((totalAmount/float64(participantCount))*100) / 100
	return amount, nil
}

// CalculateItemBasedSplit calculates split based on items consumed
func (sc *SplitCalculator) CalculateItemBasedSplit(items []models.BillItem) (map[uuid.UUID]float64, error) {
	// Map to store user ID -> total amount owed
	userAmounts := make(map[uuid.UUID]float64)

	for _, item := range items {
		totalItemCost := item.Price * float64(item.Quantity)
		participantCount := len(item.Participants)

		if participantCount == 0 {
			continue // Skip items with no participants
		}

		// Calculate total split ratio
		var totalRatio float64
		for _, ip := range item.Participants {
			totalRatio += ip.SplitRatio
		}

		if totalRatio == 0 {
			totalRatio = float64(participantCount) // Default: equal split
		}

		// Calculate amount for each participant based on their ratio
		for _, ip := range item.Participants {
			ratio := ip.SplitRatio
			if ratio == 0 {
				ratio = 1.0 // Default ratio
			}

			amount := (totalItemCost * ratio) / totalRatio
			amount = math.Round(amount*100) / 100

			userAmounts[ip.UserID] += amount
		}
	}

	return userAmounts, nil
}

// CalculateCustomSplit validates custom amounts equal total
func (sc *SplitCalculator) CalculateCustomSplit(totalAmount float64, customAmounts map[uuid.UUID]float64) error {
	var sum float64

	for _, amount := range customAmounts {
		sum += amount
	}

	sum = math.Round(sum*100) / 100
	totalAmount = math.Round(totalAmount*100) / 100

	if sum != totalAmount {
		return errors.New("custom amounts do not equal total amount")
	}

	return nil
}

// DistributeRemainder handles rounding errors by distributing remainder
func (sc *SplitCalculator) DistributeRemainder(totalAmount float64, amounts map[uuid.UUID]float64) map[uuid.UUID]float64 {
	var sum float64
	for _, amount := range amounts {
		sum += amount
	}

	remainder := math.Round((totalAmount-sum)*100) / 100

	if remainder != 0 && len(amounts) > 0 {
		// Add remainder to first participant
		for userID := range amounts {
			amounts[userID] += remainder
			break
		}
	}

	return amounts
}

// CalculateSplitPercentages calculates what percentage each participant owes
func (sc *SplitCalculator) CalculateSplitPercentages(amounts map[uuid.UUID]float64, totalAmount float64) map[uuid.UUID]float64 {
	percentages := make(map[uuid.UUID]float64)

	for userID, amount := range amounts {
		percentage := (amount/totalAmount) * 100
		percentage = math.Round(percentage*100) / 100
		percentages[userID] = percentage
	}

	return percentages
}