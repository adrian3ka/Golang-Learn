package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"fmt"
)

func TestRequestNewImageUnauthenticated(t *testing.T) {
	fmt.Println("======================================================================================================================================")
	request, _ := http.NewRequest("GET", "/images/new", nil)
	fmt.Println("Request >> ",request)
	recorder := httptest.NewRecorder()

	app := NewApp()
	app.ServeHTTP(recorder, request)

	fmt.Println("recorder.Code >>",recorder.Code ,",",http.StatusFound)
	
	if recorder.Code != http.StatusFound {
		t.Error("Expected a redirect code, but got", recorder.Code)
	}

	loc := recorder.HeaderMap.Get("Location")
	fmt.Println("loc >> ", loc)
	if loc != "/login?next=%252Fimages%252Fnew" {
		t.Error("Expected Location to redirect to sign in, but got", loc)
	}
}

type MockSessionStore struct {
	Session *Session
}

func (store MockSessionStore) Find(string) (*Session, error) {
	return store.Session, nil
}

func (store MockSessionStore) Save(*Session) error {
	return nil
}

func (store MockSessionStore) Delete(*Session) error {
	return nil
}

func TestRequestNewImageAuthenticated(t *testing.T) {
	fmt.Println("======================================================================================================================================")
	// Replace the user store temporarily
	oldUserStore := globalUserStore
	defer func() {
		globalUserStore = oldUserStore
	}()
	globalUserStore = &MockUserStore{
		findUser: &User{},
	}

	expiry := time.Now().Add(time.Hour)

	// Replace the session store temporarily
	oldSessionStore := globalSessionStore
	defer func() {
		globalSessionStore = oldSessionStore
	}()
	globalSessionStore = &MockSessionStore{
		Session: &Session{
			ID:     "session_123",
			UserID: "user_123",
			Expiry: expiry,
		},
	}

	// Create a cookie for the
	authCookie := &http.Cookie{
		Name:    sessionCookieName,
		Value:   "session_123",
		Expires: expiry,
	}
	
	request, _ := http.NewRequest("GET", "/images/new", nil)
	request.AddCookie(authCookie)
	fmt.Println("request >> ",request)
	recorder := httptest.NewRecorder()

	app := NewApp()
	app.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Error("Expected a redirect code, but got", recorder.Code)
	}
}
func BenchmarkRequestNewImageUnauthenticated(b *testing.B) {
	fmt.Println("================Benchmark1===================")
	request, _ := http.NewRequest("GET", "/images/new", nil)
	recorder := httptest.NewRecorder()
	app := NewApp()
	for i := 0; i < b.N; i++ {
		app.ServeHTTP(recorder, request)
	}
}
