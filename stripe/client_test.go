package stripe_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/acharyab15/test-with-go/stripe"
)

var (
	apiKey string
	update bool
)

const (
	tokenAmex               = "tok_amex"
	tokenVisaDebit          = "tok_visa_debit"
	tokenMastercardPrepaid  = "tok_mastercard_prepaid"
	tokenInvalid            = "tok_alsdfkhas"
	tokenExpiredCard        = "tok_chargeDeclinedExpiredCard"
	tokenIncorrectCVC       = "tok_chargeDeclinedIncorrectCvc"
	tokenInsufficientFunds  = "tok_chargeDeclinedInsufficient_funds"
	tokenChargeCustomerFail = "tok_chargeCustomerFail"
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
	if apiKey == "" {
		count := 0
		handler := func(w http.ResponseWriter, r *http.Request) {
			resp := readResponse(t, count)
			w.WriteHeader(resp.StatusCode)
			w.Write(resp.Body)
			count++
		}
		server := httptest.NewServer(http.HandlerFunc(handler))
		c.BaseURL = server.URL
		teardown = append(teardown, server.Close)
	}
	if update {
		rc := &recorderClient{}
		c.HttpClient = rc
		teardown = append(teardown, func() {
			t.Logf("Length of responses = %d", len(rc.responses))
			for idx, res := range rc.responses {
				recordResponse(t, res, idx)
			}
		})
	}
	return &c, func() {
		for _, fn := range teardown {
			fn()
		}
	}
}

func responsePath(t *testing.T, count int) string {
	return filepath.Join("testdata", filepath.FromSlash(fmt.Sprintf("%s.%d.json", t.Name(), count)))
}

func readResponse(t *testing.T, count int) response {
	var resp response
	path := responsePath(t, count)
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open the response file: %s. err = %v", path, err)
	}
	defer f.Close()
	jsonBytes, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read the response file: %s. err = %v", path, err)
	}
	err = json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		t.Fatalf("failed to json unmarshal the response file: %s. err = %v", path, err)
	}
	return resp
}

func recordResponse(t *testing.T, resp response, count int) {
	path := responsePath(t, count)
	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		t.Fatalf("failed to create the response dir: %s. err = %v", filepath.Dir(path), err)
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create the response file: %s. err = %v", path, err)
	}
	defer f.Close()
	jsonBytes, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal JSON for response file: %s. err = %v", path, err)
	}
	_, err = f.Write(jsonBytes)
	if err != nil {
		t.Fatalf("Failed to write json bytes for response file: %s. err = %v", path, err)
	}
}

func TestClient_Customer(t *testing.T) {
	if apiKey == "" {
		t.Log("No API key provided. Running unit tests using recorded responses. Be sure to run against the real API before committing.")
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
		"expired card": {
			token:  tokenExpiredCard,
			email:  "test@testwithgo.com",
			checks: check(hasErrType(stripe.ErrTypeCardError)),
		},
		"incorrect cvc": {
			token:  tokenIncorrectCVC,
			email:  "test@testwithgo.com",
			checks: check(hasErrType(stripe.ErrTypeCardError)),
		},
		"insufficient funds": {
			token:  tokenInsufficientFunds,
			email:  "test@testwithgo.com",
			checks: check(hasErrType(stripe.ErrTypeInvalidRequest)),
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
		t.Log("No API key provided. Running unit tests using recorded responses. Be sure to run against the real API before committing.")
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

	customerViaToken := func(token string) func(*testing.T, *stripe.Client) string {
		return func(t *testing.T, c *stripe.Client) string {
			email := "test@test.com"
			cus, err := c.Customer(token, email)
			if err != nil {
				t.Fatalf("error creating customer with token %s. err = %v; want nil", token, err)
			}
			return cus.ID
		}
	}

	tests := map[string]struct {
		customerID func(*testing.T, *stripe.Client) string
		amount     int
		checks     []checkFn
	}{
		"valid charge with amex": {
			customerID: customerViaToken(tokenAmex),
			amount:     1234,
			checks:     check(hasNoErr(), hasAmount(1234)),
		},
		"valid charge with debit": {
			customerID: customerViaToken(tokenAmex),
			amount:     1111,
			checks:     check(hasNoErr(), hasAmount(1111)),
		},
		"valid charge with mastercard prepaid": {
			customerID: customerViaToken(tokenAmex),
			amount:     8878,
			checks:     check(hasNoErr(), hasAmount(8878)),
		},
		"invalid customer id": {
			customerID: func(*testing.T, *stripe.Client) string {
				return "cus_missing"
			},
			amount: 1234,
			checks: check(hasErrType(stripe.ErrTypeInvalidRequest)),
		},
		"charge failure": {
			customerID: customerViaToken(tokenChargeCustomerFail),
			amount:     1234,
			checks:     check(hasErrType(stripe.ErrTypeCardError)),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c, teardown := stripeClient(t)
			defer teardown()
			cusID := tc.customerID(t, c)
			charge, err := c.Charge(cusID, tc.amount)
			for _, check := range tc.checks {
				check(t, charge, err)
			}
		})
	}
}
