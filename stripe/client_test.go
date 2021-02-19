package stripe_test

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/acharyab15/test-with-go/stripe"
)

var (
	apiKey string
	update bool
)

const (
	tokenAmex        = "tok_amex"
	tokenInvalid     = "tok_alsdfkhas"
	tokenExpiredCard = "tok_chargeDeclinedExpiredCard"
)

func init() {
	flag.StringVar(&apiKey, "key", "", "Your TEST secret key for the Stripe API. If present, integration tests will be run using this key.")
	flag.BoolVar(&update, "update", false, "Set this flag to update the response used in local tests. This requires that the key flag is set so that we can interact with the Stripe API.")
}

func TestClient_Local(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, `{
			  "id": "cus_Irsj9PF2LvU4gR",
			  "object": "customer",
			  "address": null,
			  "balance": 0,
			  "created": 1612211129,
			  "currency": "usd",
			  "default_source": null,
			  "delinquent": false,
			  "description": null,
			  "discount": null,
			  "email": null,
			  "invoice_prefix": "A8BE20B",
			  "invoice_settings": {
				"custom_fields": null,
				"default_payment_method": null,
				"footer": null
			  },
			  "livemode": false,
			  "metadata": {},
			  "name": null,
			  "next_invoice_sequence": 1,
			  "phone": null,
			  "preferred_locales": [],
			  "shipping": null,
			  "tax_exempt": "none"
		}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	c := stripe.Client{
		Key:     "gibberish-key",
		BaseURL: server.URL,
	}
	_, err := c.Customer("randomtoken", "randomemail")
	if err != nil {
		t.Fatalf("err = %v; want nil", err)
	}
}

func stripeClient(t *testing.T) (*stripe.Client, func()) {
	teardown := make([]func(), 0)
	c := stripe.Client{
		Key: apiKey,
	}
	if update {
		rc := &recorderClient{}
		c.HttpClient = rc
		teardown = append(teardown, func() {
			t.Logf("Length of responses = %d", len(rc.responses))
			for _, res := range rc.responses {
				t.Logf("Pretending to save res: %v", res)
			}
		})
	}
	return &c, func() {
		for _, fn := range teardown {
			fn()
		}
	}
}

func TestClient_Customer(t *testing.T) {
	if apiKey == "" {
		t.Skip("No API key provided")
	}

	type checkFn func(t *testing.T, cus *stripe.Customer, err error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoErr := func() checkFn {
		return func(t *testing.T, cus *stripe.Customer, err error) {
			if err != nil {
				t.Fatalf("err = %v; want nil", err)
			}
		}
	}

	hasErrType := func(typee string) checkFn {
		return func(t *testing.T, cus *stripe.Customer, err error) {
			se, ok := err.(stripe.Error)
			if !ok {
				t.Fatalf("err isn't a stripe.Error")
			}
			if se.Type != typee {
				t.Errorf("err.Type = %s; want %s", se.Type, typee)
			}
		}
	}

	hasIDPrefix := func() checkFn {
		return func(t *testing.T, cus *stripe.Customer, err error) {
			if !strings.HasPrefix(cus.ID, "cus_") {
				t.Errorf("ID = %s; want prefix %q", cus.ID, "cus_")
			}
		}
	}
	hasCardDefaultSource := func() checkFn {
		return func(t *testing.T, cus *stripe.Customer, err error) {
			if !strings.HasPrefix(cus.DefaultSource, "card_") {
				t.Errorf("DefaultSource = %s; want prefix %q", cus.DefaultSource, "card_")
			}
		}
	}
	hasEmail := func(email string) checkFn {
		return func(t *testing.T, cus *stripe.Customer, err error) {
			if cus.Email != email {
				t.Errorf("Customer() Email = %s; want %s", cus.Email, email)
			}
		}
	}
	tests := map[string]struct {
		token  string
		email  string
		checks []checkFn
	}{
		"valid customer with amex": {
			token:  tokenAmex,
			email:  "test@testwithgo.com",
			checks: check(hasNoErr(), hasIDPrefix(), hasCardDefaultSource(), hasEmail("test@testwithgo.com")),
		},
		"invalid token": {
			token:  tokenInvalid,
			email:  "test@testwithgo.com",
			checks: check(hasErrType(stripe.ErrTypeInvalidRequest)),
		},
		"expired token": {
			token:  tokenExpiredCard,
			email:  "test@testwithgo.com",
			checks: check(hasErrType(stripe.ErrTypeCardError)),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c, teardown := stripeClient(t)
			defer teardown()
			cus, err := c.Customer(tc.token, tc.email)
			for _, check := range tc.checks {
				check(t, cus, err)
			}
		})
	}
}

func TestClient_Charge(t *testing.T) {
	if apiKey == "" {
		t.Skip("No API key provided")
	}

	type checkFn func(t *testing.T, charge *stripe.Charge, err error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasAmount := func(amount int) checkFn {
		return func(t *testing.T, charge *stripe.Charge, err error) {
			if charge.Amount != amount {
				t.Errorf("Amount = %d; want %d", charge.Amount, amount)
			}
		}
	}

	hasNoErr := func() checkFn {
		return func(t *testing.T, charge *stripe.Charge, err error) {
			if err != nil {
				t.Fatalf("err = %v; want nil", err)
			}
		}
	}

	hasErrType := func(typee string) checkFn {
		return func(t *testing.T, charge *stripe.Charge, err error) {
			se, ok := err.(stripe.Error)
			if !ok {
				t.Fatalf("err isn't a stripe.Error")
			}
			if se.Type != typee {
				t.Errorf("err.Type = %s; want %s", se.Type, typee)
			}
		}
	}

	c, teardown := stripeClient(t)
	defer teardown()
	// Create a customer for the test
	tok := "tok_amex"
	email := "test@test.com"
	cus, err := c.Customer(tok, email)
	if err != nil {
		t.Fatalf("Customer() err = %v' want nil", err)
	}

	tests := map[string]struct {
		customerID string
		amount     int
		checks     []checkFn
	}{
		"valid charge": {
			customerID: cus.ID,
			amount:     1234,
			checks:     check(hasNoErr(), hasAmount(1234)),
		},
		"invalid customer id": {
			customerID: "cus_missing",
			amount:     1234,
			checks:     check(hasErrType(stripe.ErrTypeInvalidRequest)),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c, teardown := stripeClient(t)
			defer teardown()
			charge, err := c.Charge(tc.customerID, tc.amount)
			for _, check := range tc.checks {
				check(t, charge, err)
			}
		})
	}
}
