package main
import (
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
)
func TestImageCreateFromURLInvalidStatusCode(t *testing.T) {
	fmt.Println("======================================================================================================================================")
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(404)
						}))
	defer server.Close()
	image := Image{}
	fmt.Println("server >>",server)
	err := image.CreateFromURL(server.URL)
	if err != errImageURLInvalid {
		t.Errorf("Expected errImageURLInvalid but got %s", err)
	}
}
func TestImageCreateFromURLInvalidContentType(t *testing.T) {
	fmt.Println("======================================================================================================================================")
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
														w.WriteHeader(200)
														w.Header().Add("Content-Type", "text/html")
												  }))
	defer server.Close()
	fmt.Println("server >>",server)
	image := Image{}
	err := image.CreateFromURL(server.URL)
	if err != errInvalidImageType {
		t.Errorf("Expected errInvalidImageType but got %s", err)
	}
}
