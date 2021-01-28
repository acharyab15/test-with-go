package stripe_test

import (
	"strings"
	"testing"

	"github.com/acharyab15/test-with-go/stripe"
)

func TestClient_Customer(t *testing.T) {
	c := stripe.Client{}
	tok := "tok_amex"
	cus, err := c.Customer(tok)
	if err != nil {
		t.Fatalf("Customer() err = %s; want i%v", err, nil)
	}
	if cus == nil {
		t.Fatalf("Customer() = nil; want non-nil value")
	}
	if !strings.HasPrefix(cus.ID, "cus_") {
		t.Errorf("Customer() ID = %s; want prefix %q", cus.ID, "cus_")
	}
}
