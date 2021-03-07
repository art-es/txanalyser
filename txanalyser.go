package txanalyser

import (
	"time"
)

type Transaction struct {
	ID              string
	Amount          float64
	BankName        string
	BankCountryCode string

	// Additional fields
	BankLatency time.Duration
}

func (t Transaction) IncomePerTime() float64 {
	return t.Amount / float64(t.BankLatency.Microseconds())
}

type Result struct {
	ID         string
	Fraudulent bool
}

func ProcessTransactions(transactions []*Transaction) []*Result {
	results := make([]*Result, 0)
	for i, tx := range transactions {
		results[i] = &Result{
			ID:         tx.ID,
			Fraudulent: processTransaction(tx),
		}
	}
	return results
}

func processTransaction(tx *Transaction) bool {
	time.Sleep(tx.BankLatency)
	return true
}
