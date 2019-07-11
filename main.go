package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	helpers "github.com/ww-tech/aws-sts-proxy/helpers"
)

var (
	assumeRole        = os.Getenv("EKS_ASSUME_ROLE")
	externalID        = os.Getenv("EXTERNAL_ID")
	stringRequirement = helpers.GetEnv("STRING_REQUIREMENT", "")
	port              = helpers.GetEnv("PORT", "8080")
	hc                = helpers.GetEnv("HEALTHCHECK", "/hc")
)

func main() {

	if assumeRole == "" || externalID == "" {
		fmt.Println("Must export EKS_ASSUME_ROLE")
		fmt.Println("Must export EXTERNAL_ID")
		os.Exit(1)
	}

	helper := helpers.Helper{
		AssumeRole:        assumeRole,
		ExternalID:        externalID,
		StringRequirement: stringRequirement,
	}

	http.HandleFunc(hc, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/sts/token", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		role := r.URL.Query().Get("RoleArn")
		durationString := r.URL.Query().Get("Duration")
		externalID := r.URL.Query().Get("ExternalID")
		if durationString == "" {
			durationString = "60"
		}
		duration, err := strconv.ParseInt(durationString, 10, 64)
		if err != nil {
			fmt.Println("ERROR => " + err.Error())
			errorHandler(w, r, http.StatusUnauthorized)
			return
		}

		stsSession, err := helper.GetSTSToken(token, role, duration, externalID)
		if err != nil {
			fmt.Println("ERROR => " + err.Error())
			errorHandler(w, r, http.StatusUnauthorized)
			return
		}

		stsJSON, err := json.Marshal(stsSession)
		if err != nil {
			fmt.Println("ERROR => " + err.Error())
			errorHandler(w, r, http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(stsJSON))
	})
	log.Printf("listening on http://%s/", "0.0.0.0:"+port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	var b map[string]string
	var err error
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		err = json.Unmarshal([]byte(`{"error": "Custom 404"}`), &b)
	} else if status == http.StatusTeapot {
		err = json.Unmarshal([]byte(`{"error": "Something went wrong"}`), &b)
	} else if status == http.StatusUnauthorized {
		err = json.Unmarshal([]byte(`{"error": "Not Authorized"}`), &b)
	} else {
		err = json.Unmarshal([]byte(`{"error": "Bad Request"}`), &b)
	}
	if err != nil {
		fmt.Println("ERROR => ", err)
	}
	json.NewEncoder(w).Encode(b)
}
