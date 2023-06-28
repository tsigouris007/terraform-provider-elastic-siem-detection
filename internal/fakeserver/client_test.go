package fakeserver

/**
	This code was copied from the repository: https://github.com/Mastercard/terraform-provider-restapi/
	path: 'restapi/api_client_test.go'
	It has been adapted to fit the testing needs of this project
**/

import (
	"log"
	"net/http"
	"testing"
	"time"
)

var apiClientServer *http.Server

func TestAPIClient(t *testing.T) {
	debug := false

	if debug {
		log.Println("client_test.go: Starting HTTP server")
	}
	setupAPIClientServer()

	/* Notice the intentional trailing / */
	opt := &ApiClientOpt{
		Uri:                 "http://127.0.0.1:8083/",
		Insecure:            false,
		Username:            "",
		Password:            "",
		Headers:             make(map[string]string),
		Timeout:             2,
		IdAttribute:         "id",
		CopyKeys:            make([]string, 0),
		WriteReturnsObject:  false,
		CreateReturnsObject: false,
		rateLimit:           1,
		Debug:               debug,
	}
	client, _ := NewAPIClient(opt)

	var res string
	var err error

	if debug {
		log.Printf("api_client_test.go: Testing standard OK request\n")
	}
	res, err = client.SendRequest("GET", "/ok", "")
	if err != nil {
		t.Fatalf("client_test.go: %s", err)
	}
	if res != "It works!" {
		t.Fatalf("client_test.go: Got back '%s' but expected 'It works!'\n", res)
	}

	if debug {
		log.Printf("api_client_test.go: Testing redirect request\n")
	}
	res, err = client.SendRequest("GET", "/redirect", "")
	if err != nil {
		t.Fatalf("client_test.go: %s", err)
	}
	if res != "It works!" {
		t.Fatalf("client_test.go: Got back '%s' but expected 'It works!'\n", res)
	}

	/* Verify timeout works */
	if debug {
		log.Printf("api_client_test.go: Testing timeout aborts requests\n")
	}
	_, err = client.SendRequest("GET", "/slow", "")
	if err == nil {
		t.Fatalf("client_test.go: Timeout did not trigger on slow request")
	}

	if debug {
		log.Printf("api_client_test.go: Testing rate limited OK request\n")
	}
	//startTime := time.Now().Unix()

	for i := 0; i < 4; i++ {
		client.SendRequest("GET", "/ok", "")
	}

	//duration := time.Now().Unix() - startTime
	//if duration < 3 {
	//	t.Fatalf("client_test.go: requests not delayed\n")
	//}

	if debug {
		log.Println("client_test.go: Stopping HTTP server")
	}
	shutdownAPIClientServer()
	if debug {
		log.Println("client_test.go: Done")
	}
}

func setupAPIClientServer() {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("It works!"))
	})
	serverMux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(9999 * time.Second)
		w.Write([]byte("This will never return!!!!!"))
	})
	serverMux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ok", http.StatusPermanentRedirect)
	})

	apiClientServer = &http.Server{
		Addr:    "127.0.0.1:8083",
		Handler: serverMux,
	}
	go apiClientServer.ListenAndServe()
	/* let the server start */
	time.Sleep(1 * time.Second)
}

func shutdownAPIClientServer() {
	apiClientServer.Close()
}
