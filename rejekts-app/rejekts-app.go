package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/gorilla/mux"
)

var (
	STATE_STORE_NAME = "statestore"
	daprClient       dapr.Client
)

type MyValues struct {
	Values []string
}

func writeHandler(w http.ResponseWriter, r *http.Request) {

	value := r.URL.Query().Get("message")

	values, _ := read("values")

	if values.Values == nil || len(values.Values) == 0 {
		values.Values = []string{value}
	} else {
		values.Values = append(values.Values, value)
	}

	jsonData, err := json.Marshal(values)

	err = save("values", jsonData)
	if err != nil {
		panic(err)
	}

	respondWithJSON(w, http.StatusOK, values)
}

func save(key string, data []byte) error {
	ctx := context.Background()
	daprClient, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	return daprClient.SaveState(ctx, STATE_STORE_NAME, key, data, nil)
}

func read(key string) (MyValues, error) {
	ctx := context.Background()
	daprClient, err := dapr.NewClient()
	if err != nil {
		return MyValues{}, err
	}
	result, err := daprClient.GetState(ctx, STATE_STORE_NAME, key, nil)
	if err != nil {
		return MyValues{}, err
	}
	myValues := MyValues{}
	if result.Value != nil {
		json.Unmarshal(result.Value, &myValues)
	}
	return myValues, nil
}

func readHandler(w http.ResponseWriter, r *http.Request) {

	values, _ := read("values")

	respondWithJSON(w, http.StatusOK, values)

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
