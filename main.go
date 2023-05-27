package main

import (
	"bytes"
	"encoding/json"
  "errors"
	c "github.com/CalebDepatie/go-common"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"net/http"
	"os"
)

func main() {
	defer c.LogInfo("Server shutting down")

	err := godotenv.Load()
	if err != nil {
		c.LogFatal("Error loading .env file")
	}

  connectToGhidorah()
  
	http.HandleFunc("/send", sendMail)

	c.LogFatal(http.ListenAndServe(":10000", nil))

}

func sendMail(w http.ResponseWriter, r *http.Request) {
	var (
		unmarshalErr *json.UnmarshalTypeError
	)

	mail_args := struct {
		To      string "json:to"
		Subject string "json:subject"
		Body    string "json:body"
	}{}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&mail_args)

	if err != nil {
		if errors.As(err, &unmarshalErr) {
			c.LogWarning("JSON Error: ", unmarshalErr.Field)
		} else {
			c.LogWarning("Request Error: ", err.Error())
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	smtpDialer := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL"), os.Getenv("PASS"))

	msg := gomail.NewMessage()
	msg.SetHeader("From", os.Getenv("EMAIL"))
	msg.SetHeader("To", mail_args.To)
	msg.SetHeader("Subject", mail_args.Subject)
	msg.SetBody("text/html", mail_args.Body)

	if err := smtpDialer.DialAndSend(msg); err != nil {
		c.LogError("Could not send message: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

}

// initialize connection to Ghidorah
func connectToGhidorah() {
	reg := struct {
		Name             string `json:name`
		ExternAccessible bool   `json:extern_facing`
		Internal         bool   `json:internal`
		Port             string `json:port`
	}{
		Name:             "mail",
		ExternAccessible: false,
		Internal:         true,
		Port:             "10000",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reg)
	if err != nil {
		c.LogFatal("Could not encode Registration Data", err)
	}

	req, err := http.NewRequest(http.MethodGet, os.Getenv("GHIDORAH")+"/register", &buf)
	if err != nil {
		c.LogFatal("Could not create service registation request", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.LogFatal("Could not connect to Ghidorah", err)
	}

	if resp.StatusCode != 200 {
		c.LogFatal("Could not register service", resp.Status)
	}

	c.LogInfo("Connected to Ghidorah")
}
