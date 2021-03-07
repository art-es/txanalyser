package txanalyser

import (
	"sort"
	"time"
)

type Prioritizer struct {
	transactions []*Transaction
	totalTime    time.Duration
}

func (p Prioritizer) Len() int {
	return len(p.transactions)
}

func (p Prioritizer) Less(i, j int) bool {
	// IDEA: skip the tranasaction with duration more than total time
	//
	// if p.transactions[i].BankLatency > p.totalTime {
	// 	return false
	// }
	return p.transactions[i].IncomePerTime() > p.transactions[j].IncomePerTime()
}

func (p Prioritizer) Swap(i, j int) {
	p.transactions[i], p.transactions[j] = p.transactions[j], p.transactions[i]
}

func Prioritize(transactions []*Transaction, totalTimes ...time.Duration) {
	totalTime := time.Second
	if len(totalTimes) > 0 {
		totalTime = totalTimes[0]
	}

	sort.Sort(Prioritizer{transactions, totalTime})

	// IDEA: shift the transaction with a small duration if it fits into the remaining gap
	//
	// round := totalTime
	// for i := 0; i < len(transactions)-1; i++ {
	// 	if round-transactions[i].BankLatency > 0 {
	// 		round -= transactions[i].BankLatency
	// 		continue
	// 	}
	//
	// 	for j := i + 1; j < len(transactions); j++ {
	// 		if round-transactions[j].BankLatency > 0 {
	// 			// maybe try to use linked list
	// 			next := transactions[j]
	// 			for k := i; k < j+1; i++ {
	// 				transactions[k], next = next, transactions[k]
	// 			}
	// 		}
	// 	}
	// }
}
