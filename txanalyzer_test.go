package txanalyser_test

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/art-es/txanalyser"
)

func Test(t *testing.T) {
	transactions, err := parseTransactions("./testdata/transactions.csv")
	if err != nil {
		t.Fatal(err)
	}
	if err = parseLatencies("./testdata/api_latencies.json", transactions); err != nil {
		t.Fatal(err)
	}

	dur := time.Millisecond * 50

	fmt.Println("income in unsorted list:\n", calcIncomeForTime(transactions, dur))
	txanalyser.Prioritize(transactions, dur)
	fmt.Println("income in sorted list\n", calcIncomeForTime(transactions, dur))
}

func parseTransactions(path string) ([]*txanalyser.Transaction, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	if rec, err := r.Read(); err != nil {
		return nil, err
	} else if !reflect.DeepEqual(rec, []string{"id", "amount", "bank_country_code"}) {
		return nil, err
	}

	transactions := make([]*txanalyser.Transaction, 0)
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(rec) < 3 {
			return nil, fmt.Errorf("invalid record: %v", rec)
		}

		amount, err := strconv.ParseFloat(rec[1], 64)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &txanalyser.Transaction{
			ID:              rec[0],
			Amount:          amount,
			BankCountryCode: rec[2],
		})
	}

	return transactions, nil
}

func parseLatencies(path string, transactions []*txanalyser.Transaction) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	latencies := make(map[string]int)
	if err = json.Unmarshal(data, &latencies); err != nil {
		return err
	}

	for _, tx := range transactions {
		v, ok := latencies[tx.BankCountryCode]
		if !ok {
			return fmt.Errorf("not found latency for the bank with code: %s", tx.BankCountryCode)
		}
		tx.BankLatency = time.Duration(v) * time.Millisecond
	}
	return nil
}

func calcIncomeForTime(transactions []*txanalyser.Transaction, d time.Duration) (income float64) {
	for _, tx := range transactions {
		d -= tx.BankLatency
		if d <= 0 {
			break
		}
		income += tx.Amount
	}
	return income
}
