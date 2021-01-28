package stripe_test

import (
	"encoding/json"
	"testing"

	"github.com/acharyab15/test-with-go/stripe"
)

var errorJSON = []byte(`{
  "error": {
    "code": "resource_missing",
    "doc_url": "https://stripe.com/docs/error-codes/resource-missing",
    "message": "No such customer: 'cus_'",
    "param": "customer",
    "type": "invalid_request_error"
  }
}`)

func TestError_Unmarshal(t *testing.T) {
	var se stripe.Error
	err := json.Unmarshal(errorJSON, &se)
	if err != nil {
		t.Fatalf("Unmarshal() err = %s; want nil", err)
	}
	wantCode := "resource_missing"
	if se.Code != wantCode {
		t.Errorf("Code = %s; want %s", se.Code, wantCode)
	}
	wantType := "invalid_request_error"
	if se.Type != wantType {
		t.Errorf("Type = %s; want %s", se.Type, wantType)
	}

}

func TestError_Marshal(t *testing.T) {
	se := stripe.Error{
		Code:    "test-code",
		DocURL:  "test-doc-url",
		Message: "test-message",
		Param:   "test-param",
	}
	data, err := json.Marshal(se)
	if err != nil {
		t.Fatalf("Marshal() err = %v; want nil", err)
	}
	var got stripe.Error
	err = json.Unmarshal(data, &got)
	if err != nil {
		t.Fatalf("Unmarshal() err = %v; want nil", err)
	}
	if got != se {
		t.Errorf("got = %v; want %v", got, se)
		t.Log("Is Unmarshal working? It is required for this test to pass.")
	}
}
