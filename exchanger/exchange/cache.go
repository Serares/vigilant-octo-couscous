package exchange

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var (
	// error not found
	ErrNotFound = errors.New("value not found")
	// error not found
	ErrInternal = errors.New("internal error")
)

type CurrencyToEuro struct {
	From       string  `json:"from"`
	Amount     float64 `json:"amount"`
	AmountEuro float64 `json:"amounteuro"`
}

// should be {"USD": {"35": 30 } }
type AmountsStore map[string]map[string]float64

// cache results from the api
func (l *AmountsStore) Add(item CurrencyToEuro) {
	// using the dereference notation
	// to point to the actual value that the pointer holds
	amountToString := fmt.Sprintf("%f", item.Amount)
	// if the key does not exist, create it
	if _, ok := (*l)[item.From]; !ok {
		(*l)[item.From] = map[string]float64{amountToString: item.AmountEuro}
		return
	}
	// if the key exists, add the value
	(*l)[item.From][amountToString] = item.AmountEuro
}

func (l *AmountsStore) GetOne(from string, amount float64) (float64, error) {
	amountToString := fmt.Sprintf("%f", amount)
	if val, ok := (*l)[from][amountToString]; ok {
		return val, nil
	}

	return 0, ErrNotFound
}

func (a *AmountsStore) Get(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	if len(file) == 0 {
		return nil
	}
	return json.Unmarshal(file, a)
}

func (a *AmountsStore) Save(filename string) error {
	js, err := json.Marshal(*a)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, js, 0644)
}
