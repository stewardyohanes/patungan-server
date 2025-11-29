package utils

import (
	"errors"
	"math"
)

type CurrencyConverter struct {
	rates map[string]map[string]float64
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]map[string]float64),
	}
}

// SetRate sets conveersion rate from one currency to another
func (cc *CurrencyConverter) SetRate(from, to string, rate float64) {
	if cc.rates[from] == nil {
		cc.rates[from] = make(map[string]float64)
	}
	cc.rates[from][to] = rate
}

// Convert converts amount from one currency to another
func (cc *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	if from == to {
		return amount, nil
	}

	if cc.rates[from] == nil || cc.rates[from][to] == 0 {
		return 0, errors.New("conversion rate not available")
	}

	converted := amount * cc.rates[from][to]
	return math.Round(converted*100) / 100, nil
}

// GetRate returns conversion rate between two currencies
func (cc *CurrencyConverter) GetRate(from, to string) (float64, error) {
	if from == to {
		return 1.0, nil
	}

	if cc.rates[from] == nil || cc.rates[from][to] == 0 {
		return 0, errors.New("conversion rate not available")
	}

	return cc.rates[from][to], nil
}

// LoadDefaultRates loads some common conversion rates (for demo purposes)
func (cc *CurrencyConverter) LoadDefaultRates() {
	// IDR rates
	cc.SetRate("IDR", "USD", 0.000063)
	cc.SetRate("IDR", "EUR", 0.000059)
	cc.SetRate("IDR", "SGD", 0.000085)
	cc.SetRate("IDR", "MYR", 0.00029)

	// USD rates
	cc.SetRate("USD", "IDR", 15900)
	cc.SetRate("USD", "EUR", 0.93)
	cc.SetRate("USD", "SGD", 1.35)
	cc.SetRate("USD", "MYR", 4.62)

	// EUR rates
	cc.SetRate("EUR", "IDR", 17100)
	cc.SetRate("EUR", "USD", 1.08)
	cc.SetRate("EUR", "SGD", 1.45)
	cc.SetRate("EUR", "MYR", 4.97)

	// Add reverse rates for convenience
	// Note: In production, you'd fetch these from an API
}