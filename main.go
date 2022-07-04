package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//some minor deviation from first crud practice

type Disease struct {
	ID         string      `json:"ID"`
	Name       string      `json:"name"`
	Symptom    *Symptom    `json:"symptom"`
	RiskFactor *RiskFactor `json:"riskfactor"`
}

//how to do a dynamic struct?
type Symptom struct {
	PrimarySx   string `json:"primarysx"`
	SecondarySx string `json:"secondarysx"`
}

type RiskFactor struct {
	PrimaryRx   string `json:"primaryrx"`
	SecondaryRx string `json:"secondaryrx"`
}

var diseases []Disease

// w refers to a response writer interface, r refers to a request struct
func getDiseases(w http.ResponseWriter, r *http.Request) {
	// setting content type as json
	w.Header().Set("Content-Type", "application/json")

	// encode the response to json, and now just passing in the slice
	json.NewEncoder(w).Encode(diseases)
}

//deleting the disease by ID
func deleteDisease(w http.ResponseWriter, r *http.Request) {
	// still set the header as json
	w.Header().Set("Content-Type", "application/json")

	//id passed in
	params := mux.Vars(r)
	//for range need to pass index + item
	for index, item := range diseases {
		if item.ID == params["id"] {
			// whatever ID that matches, replace itself with all the other data that comes after it
			diseases = append(diseases[:index], diseases[index+1:]...)
			break

		}
	}
	//return remaining diseases after deletion
	json.NewEncoder(w).Encode(diseases)
}

//get disease byID
func getDiseasebyID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	// blank identifier becos we wont be using index for this loop
	//looping through diseases slice and return the matching ID
	for _, item := range diseases {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func getDiseasebyName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range diseases {
		if item.Name == params["name"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode("No match")
}

// adds a new disease
func createDisease(w http.ResponseWriter, r *http.Request) {

	//again setting header as json
	w.Header().Set("Content-Type", "application/json")

	var newDisease Disease
	_ = json.NewDecoder(r.Body).Decode(&newDisease)
	newDisease.ID = strconv.Itoa(rand.Intn(1000000))
	diseases = append(diseases, newDisease)
	json.NewEncoder(w).Encode(newDisease)

}

// deletes ID that we specify and creates a new disease
// how this currently works is delete the movie first and then add a new movie
func updateDisease(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application.json")
	params := mux.Vars(r)

	for index, item := range diseases {
		if item.ID == params["id"] {
			//deletes the movie just as in the delete Disease
			diseases = append(diseases[:index], diseases[index+1:]...)

			var newDisease Disease
			_ = json.NewDecoder(r.Body).Decode(&newDisease)
			newDisease.ID = params["id"]
			diseases = append(diseases, newDisease)
			json.NewEncoder(w).Encode(newDisease)
			return
		}
	}

}

func main() {
	r := mux.NewRouter()

	//hardcode some diseases to start
	diseases = append(diseases, Disease{ID: "1", Name: "Angina Pectoris", Symptom: &Symptom{PrimarySx: "Chest pain", SecondarySx: "Pain radiating to jaw and left arm"}, RiskFactor: &RiskFactor{PrimaryRx: "Hypertension", SecondaryRx: "Diabetes Mellitus"}})
	diseases = append(diseases, Disease{ID: "2", Name: "Left Ventricular Failure", Symptom: &Symptom{PrimarySx: "Fatigue", SecondarySx: "Exertional Dyspnoea"}, RiskFactor: &RiskFactor{PrimaryRx: "Ischaemic Heart Disease", SecondaryRx: "Idiopathic Dilated Cardiomyopathy"}})
	diseases = append(diseases, Disease{ID: "3", Name: "Aortic Stenosis", Symptom: &Symptom{PrimarySx: "Angina", SecondarySx: "Dyspnoea"}, RiskFactor: &RiskFactor{PrimaryRx: "Senile Calcification", SecondaryRx: "Rheumatoid Fever"}})

	r.HandleFunc("/disease", getDiseases).Methods("GET")
	r.HandleFunc("/disease/id/{id}", getDiseasebyID).Methods("GET")
	r.HandleFunc("/disease/name/{name}", getDiseasebyName).Methods("GET")
	r.HandleFunc("/disease", createDisease).Methods("POST")
	r.HandleFunc("/disease/{id}", updateDisease).Methods("PUT")
	r.HandleFunc("/disease/{id}", deleteDisease).Methods("DELETE")

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}
