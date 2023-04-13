package main

import (
	"context"
	"encoding/json"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var (
	STATE_STORE_NAME = "statestore"
	daprClient       dapr.Client
)

type MyValues struct {
	Values []string
}

func writeHandler(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	daprClient, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}

	value := r.URL.Query().Get("message")

	result, _ := read(ctx, "values")
	myValues := MyValues{}
	if result.Value != nil {
		json.Unmarshal(result.Value, &myValues)
	}

	if myValues.Values == nil || len(myValues.Values) == 0 {
		myValues.Values = []string{value}
	} else {
		myValues.Values = append(myValues.Values, value)
	}

	jsonData, err := json.Marshal(myValues)

	err = save(ctx, "values", jsonData)
	if err != nil {
		panic(err)
	}

	respondWithJSON(w, http.StatusOK, myValues)
}

func save(ctx context, key string, data []byte) error {
	return daprClient.SaveState(ctx, STATE_STORE_NAME, key, data, nil)
}

func read(ctx context, key string) ([]byte, error) {
	return daprClient.GetState(ctx, STATE_STORE_NAME, key, nil)
}

func readHandler(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	daprClient, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}

	result, err := read(ctx, "values")
	myValues := MyValues{}
	json.Unmarshal(result.Value, &myValues)

	respondWithJSON(w, http.StatusOK, myValues)

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	r := mux.NewRouter()

	// Dapr subscription routes orders topic to this route
	r.HandleFunc("/write", writeHandler).Methods("POST")
	r.HandleFunc("/read", readHandler).Methods("GET")

	// Add handlers for readiness and liveness endpoints
	r.HandleFunc("/health/{endpoint:readiness|liveness}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(os.Getenv("KO_DATA_PATH"))))
	http.Handle("/", r)

	log.Printf("Rejekts Frontend App Started in port 8080!")
	// Start the server; this is a blocking call
	log.Fatal(http.ListenAndServe(":"+appPort, nil))

}
