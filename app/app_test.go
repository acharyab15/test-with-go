package app_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/acharyab15/test-with-go/app"
	"golang.org/x/net/publicsuffix"
)

func TestHome(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	app.Home(w, r)

	resp := w.Result()
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll() err = %s; want nil", err)
	}
	got := string(body)
	want := "<h1>Welcome!</h1>"
	if got != want {
		t.Errorf("GET / = %s; want %s", got, want)
	}
}

func TestApp_V1(t *testing.T) {
	server := httptest.NewServer(&app.Server{})
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("GET / err = %s; want nil", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll() err = %s; want nil", err)
	}
	got := string(body)
	want := "<h1>Welcome!</h1>"
	if got != want {
		t.Errorf("GET / = %s; want %s", got, want)
	}
}

// signedInClient
// create a cookie jar
// create a http client with cookie jar
// login and set the session cookie
func signedInClient(t *testing.T, baseURL string) *http.Client {
	// Our cookiejar will keep and set cookies for us between requests.
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		t.Fatalf("cookiejar.New() err = %s; want nil", err)
	}
	client := &http.Client{
		Jar: jar,
	}

	// Our client has a cookie jar, but it has no session cookie. By logging
	// in we can ensure that it gets set.
	loginURL := baseURL + "/login"
	req, err := http.NewRequest(http.MethodPost, loginURL, nil)
	if err != nil {
		t.Fatalf("NewRequest() err = %s; want nil", err)
	}
	_, err = client.Do(req)
	if err != nil {
		t.Fatalf("POST /login err = %s; want nil", err)
	}
	t.Logf("Cookies: %v", client.Jar.Cookies(req.URL))
	return client
}

type headerClient struct {
	headers map[string]string
}

func signedInRequest(t *testing.T, method, target string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, target, body)
	if err != nil {
		t.Fatalf("http.NewRequest() err = %s; want nil", err)
	}
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: "fake_session_token",
	})
	return req
}

func (hc headerClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for hk, hv := range hc.headers {
		req.Header.Set(hk, hv)
	}
	client := http.Client{}
	return client.Do(req)
}

func TestApp_V2(t *testing.T) {
	server := httptest.NewServer(&app.Server{})
	defer server.Close()

	t.Run("custom build request", func(t *testing.T) {
		req := signedInRequest(t, http.MethodGet, server.URL+"/admin", nil)
		var client http.Client
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("GET /admin err = %s; want nil", err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("ioutil.ReadAll() err = %s; want nil", err)
		}
		got := string(body)
		want := "<h1>Welcome to the admin page!</h1>"
		if got != want {
			t.Errorf("GET / = %s; want %s", got, want)
		}
	})

	t.Run("cookie based auth", func(t *testing.T) {
		client := signedInClient(t, server.URL)
		res, err := client.Get(server.URL + "/admin")
		if err != nil {
			t.Errorf("GET /admin err = %s; want nil", err)
		}
		if res.StatusCode != 200 {
			t.Errorf("GET /admin code = %d; want %d", res.StatusCode, 200)
		}

		res, err = client.Get(server.URL + "/header-admin")
		if err != nil {
			t.Errorf("GET /header-admin err = %s; want nil", err)
		}
		if res.StatusCode != 403 {
			t.Errorf("GET /header-admin code = %d; want %d", res.StatusCode, 403)
		}
	})
	t.Run("header based auth", func(t *testing.T) {
		client := headerClient{
			headers: map[string]string{"api-key": "fake_api_key"},
		}
		res, err := client.Get(server.URL + "/admin")
		if err != nil {
			t.Errorf("GET /admin err = %s; want nil", err)
		}
		// no cookie so 403
		if res.StatusCode != 403 {
			t.Errorf("GET /admin code = %d; want %d", res.StatusCode, 403)
		}
		res, err = client.Get(server.URL + "/header-admin")
		if err != nil {
			t.Errorf("GET /header-admin err = %s; want nil", err)
		}
		// fake_api_key, so 200
		if res.StatusCode != 200 {
			t.Errorf("GET /header-admin code = %d; want %d", res.StatusCode, 200)
		}
	})
}
