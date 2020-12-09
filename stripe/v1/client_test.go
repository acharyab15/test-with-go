package stripe_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/acharyab15/test-with-go/stripe/v1"
)

func TestApp(t *testing.T) {
	client, mux, teardown := stripe.TestClient(t)
	defer teardown()

	mux.HandleFunc("/v1/charges", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id":"ch_1HwbX02eZvKYlo2Cck108DY4","amount":2000,
		"description":"Charge for demo purposes.","status":"succeeded"}`)
	})

	charge, err := client.Charge(123, "doesnt_matter", "something else")
	if err != nil {
		t.Errorf("Charge() err = %s; want nil", err)
	}
	if charge.ID != "ch_1HwbX02eZvKYlo2Cck108DY4" {
		t.Errorf("Charge() id = %s; want %s", charge.ID, "ch_1HwbX02eZvKYlo2Cck108DY4")
	}
	if charge.Status != "succeeded" {
		t.Errorf("Charge() status = %s; want %s", charge.Status, "succeeded")
	}
}
