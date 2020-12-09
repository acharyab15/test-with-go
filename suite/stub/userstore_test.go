package stub_test

import (
	"testing"

	"github.com/acharyab15/test-with-go/suite"
	"github.com/acharyab15/test-with-go/suite/stub"
	"github.com/acharyab15/test-with-go/suite/suitetest"
)

var _ suite.UserStore = &stub.UserStore{}

func TestUserStore(t *testing.T) {
	us := &stub.UserStore{}
	suitetest.UserStore(t, us, nil, nil)
}

func TestUserStore_withStruct(t *testing.T) {
	us := &stub.UserStore{}
	tests := suitetest.UserStoreSuite{
		UserStore: us,
	}
	tests.All(t)
}
