package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Response represents a generic API response
type Response struct {
	Message string `json:"message"`
	ResponseCode int `json:"responseCode"`
}

type ReservationFormData struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	PickupDateTime  string `json:"pickupDateTime"`
	PickupLocation  string `json:"pickupLocation"`
	DropoffLocation string `json:"dropoffLocation"`
	TransportType   string `json:"transportType"`
	TripType        string `json:"tripType"`
}

func ReservationHandler(w http.ResponseWriter, r *http.Request) {
	var reservationData ReservationFormData
	err := json.NewDecoder(r.Body).Decode(&reservationData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if reservationData.FirstName == "" {
		http.Error(w, "First name is required", http.StatusBadRequest)
		return
	}
	if reservationData.LastName == "" {
		http.Error(w, "Last name is required", http.StatusBadRequest)
		return
	}
	if reservationData.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	if reservationData.Phone == "" {
		http.Error(w, "Phone number is required", http.StatusBadRequest)
		return
	}
	if reservationData.PickupDateTime == "" {
		http.Error(w, "Pickup date and time is required", http.StatusBadRequest)
		return
	}
	if reservationData.PickupLocation == "" {
		http.Error(w, "Pickup location is required", http.StatusBadRequest)
		return
	}
	if reservationData.DropoffLocation == "" {
		http.Error(w, "Dropoff location is required", http.StatusBadRequest)
		return
	}
	if reservationData.TransportType == "" {
		http.Error(w, "Transport type is required", http.StatusBadRequest)
		return
	}
	if reservationData.TripType == "" {
		http.Error(w, "Trip type is required", http.StatusBadRequest)
		return
	}

	emailBody := fmt.Sprintf(
		"Name: %s %s\n"+
			"Email: %s\n"+
			"Phone: %s\n"+
			"Pickup Date and Time: %s\n"+
			"Pickup Location: %s\n"+
			"Dropoff Location: %s\n"+
			"Transport Type: %s\n"+
			"Trip Type: %s",
		reservationData.FirstName, reservationData.LastName,
		reservationData.Email,
		reservationData.Phone,
		reservationData.PickupDateTime,
		reservationData.PickupLocation,
		reservationData.DropoffLocation,
		reservationData.TransportType,
		reservationData.TripType,
	)
	emailSubject := fmt.Sprintf("Ride Request - %s %s", reservationData.FirstName, reservationData.LastName)

	_, sendEmailErr := SendSimpleEmail(emailSubject, emailBody, nil, "")
	if sendEmailErr != nil {
		log.Printf("Error sending email: %v", sendEmailErr)
		http.Error(w, "Could not receive application. Please contact us directly.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Message:      "Ride request received! We will get back to you soon.",
		ResponseCode: http.StatusOK,
	})
}

func ApplicationHandler(w http.ResponseWriter, r *http.Request) {
	// Check if resume file is present
	file, header, err := r.FormFile("resume")
	if err != nil {
		http.Error(w, "Resume file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get and validate form values
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	experience := r.FormValue("experience")
	availability := r.FormValue("availability")

	// Check if any required field is empty
	if firstName == "" {
		http.Error(w, "First name is required", http.StatusBadRequest)
		return
	}
	if lastName == "" {
		http.Error(w, "Last name is required", http.StatusBadRequest)
		return
	}
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	if phone == "" {
		http.Error(w, "Phone number is required", http.StatusBadRequest)
		return
	}
	if experience == "" {
		http.Error(w, "Experience is required", http.StatusBadRequest)
		return
	}
	if availability == "" {
		http.Error(w, "Availability is required", http.StatusBadRequest)
		return
	}

	// Prepare to send email with application details
	emailBody := fmt.Sprintf(
		"Name: %s %s\n"+
			"Email: %s\n"+
			"Phone: %s\n"+
			"Experience: %s years\n"+
			"Availability: %s",
		firstName, lastName,
		email,
		phone,
		experience,
		availability,
	)
	emailSubject := fmt.Sprintf("New Application Received - %s %s", firstName, lastName)
	
	_, sendEmailErr := SendSimpleEmail(emailSubject, emailBody, file, header.Filename)
	if sendEmailErr != nil {
		log.Printf("Error sending email: %v", sendEmailErr)
		http.Error(w, "Could not receive application. Please contact us directly.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Message:      "Application received! We will get back to you soon.",
		ResponseCode: http.StatusOK,
	})
}
