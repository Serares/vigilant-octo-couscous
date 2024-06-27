package exchange_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Bitstarz-eng/event-processing-challenge/exchanger/exchange"
)

func TestSave(t *testing.T) {
	l1 := exchange.AmountsStore{}
	l2 := exchange.AmountsStore{}

	amountToStore := exchange.CurrencyToEuro{From: "USD", AmountEuro: 10, Amount: 12}

	toAmountCached := fmt.Sprintf("%f", amountToStore.Amount)

	l1.Add(amountToStore)
	if l1["USD"][toAmountCached] != amountToStore.AmountEuro {
		t.Errorf("Expected %f, got %f instead.", amountToStore.AmountEuro, l1["USD"][toAmountCached])
	}
	tf, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}
	defer os.Remove(tf.Name())
	if err := l1.Save(tf.Name()); err != nil {
		t.Fatalf("Error saving list to file: %s", err)
	}
	if err := l2.Get(tf.Name()); err != nil {
		t.Fatalf("Error getting list from file: %s", err)
	}

	if l1["USD"][toAmountCached] != l2["USD"][toAmountCached] {
		t.Errorf("Task %f should match %f task.", l1["USD"][toAmountCached], l2["USD"][toAmountCached])
	}
}
