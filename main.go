package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Payment struct {
	Id       string    `json:"Id"`
	Title    string    `json:"Title"`
	Price    int       `json:"Price"`
	Date     time.Time `json:"Date"`
	Type     string    `json:"Type"`
	Comment  string    `json:"Comment"`
	Category string    `json:"Category"`
}

var Payments []Payment

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Homepage!")
	fmt.Println("Homepage Endpoint Hit")
}

func createNewPayment(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	var payment Payment
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data")
	}
	json.Unmarshal(reqBody, &payment)
	Payments = append(Payments, payment)

	json.NewEncoder(w).Encode(payment)
}

func returnSinglePayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	for _, payment := range Payments {
		if payment.Id == key {
			json.NewEncoder(w).Encode(payment)
		}
	}
}

func returnAllPayments(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: All Payments Endpoint")
	json.NewEncoder(w).Encode(Payments)
}

func updatePayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var updatedEvent Payment
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data to update")
	}
	json.Unmarshal(reqBody, &updatedEvent)
	for i, payment := range Payments {
		if payment.Id == id {

			payment.Title = updatedEvent.Title
			payment.Price = updatedEvent.Price
			payment.Date = updatedEvent.Date
			payment.Type = updatedEvent.Type
			payment.Comment = updatedEvent.Comment
			payment.Category = updatedEvent.Category
			Payments = append(Payments[:i], payment)
			json.NewEncoder(w).Encode(payment)
		}
	}
}

func deletePayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for index, payment := range Payments {
		if payment.Id == id {
			Payments = append(Payments[:index], Payments[index+1:]...)
			fmt.Fprintf(w, "The payment with ID %v has been deleted successfully", id)
		}
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/payment", createNewPayment).Methods("POST")
	myRouter.HandleFunc("/payments", returnAllPayments).Methods("GET")
	myRouter.HandleFunc("/payments/{id}", returnSinglePayment).Methods("GET")
	myRouter.HandleFunc("/payments/{id}", updatePayment).Methods("PATCH")
	myRouter.HandleFunc("/payments/{id}", deletePayment).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")

	Payments = []Payment{
		Payment{Id: "1", Title: "Salary", Price: 120,
			Date: time.Date(2006, time.February, 1, 15, 04, 05, 0, time.UTC),
			Type: "Income", Comment: "comment for id1", Category: "work"},
		Payment{Id: "2", Title: "Rent", Price: 50,
			Date: time.Date(2022, time.February, 18, 11, 15, 0, 0, time.UTC),
			Type: "Expense", Comment: "comment for id2", Category: "house"},
		Payment{Id: "3", Title: "Magnum", Price: 20,
			Date: time.Date(2022, time.February, 15, 13, 50, 0, 0, time.UTC),
			Type: "Expense", Comment: "comment for id2", Category: "food"},
	}
	fmt.Println("Total", Total(Payments))

	// var incomes []Payment

	// for _, payment := range Payments {
	// 	if isIncome(payment) {
	// 		incomes = append(incomes, payment)

	// 	}
	// }

	handleRequests()

}
func Total(Payments []Payment) (result int) {
	result = 0
	for _, payment := range Payments {
		if isIncome(payment) {
			result += payment.Price
		} else {
			result -= payment.Price
		}
	}
	return result
}
func isIncome(payment Payment) bool {
	return payment.Type == "Income"
}
