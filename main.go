package main

import (
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

	reg := c.GhidorahReg{
		Name:             "mail",
		ExternAccessible: false,
		Internal:         true,
		Port:             "10000",
	}

	c.ConnectToGhidorah(reg, os.Getenv("GHIDORAH"))

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
