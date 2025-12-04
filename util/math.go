package util

import (
	"fmt"
	"math"
	"strconv"
)

const (
	CurrencySymbol = "Br"
)

// ParseStringToFloat64 parses a string to float64
func ParseStringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// FormatFare formats a fare amount with currency symbol
func FormatFare(amount float64) string {
	rounded := math.Round(amount)
	return CurrencySymbol + formatNumber(rounded)
}

func formatNumber(num float64) string {
	return fmt.Sprintf("%.0f", num)
}

// ValidateFareAmount validates that a fare amount is within min/max bounds
func ValidateFareAmount(amount, minFare, maxFare float64) error {
	if amount < minFare {
		return fmt.Errorf("fare amount %.0f is below minimum fare %.0f", amount, minFare)
	}
	if amount > maxFare {
		return fmt.Errorf("fare amount %.0f exceeds maximum fare %.0f", amount, maxFare)
	}
	return nil
}

// RoundToNearestBirr rounds an amount to the nearest birr
func RoundToNearestBirr(amount float64) float64 {
	return math.Round(amount)
}

// Max returns the maximum of two integers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
