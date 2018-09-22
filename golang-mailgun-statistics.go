package main

import (
	"bytes"
	"fmt"
	"gopkg.in/mailgun/mailgun-go.v1"
	"os"
	"strings"
)

var privateKey = "key-xxxxx"
var publicKey = "pubkey-xxxxx"

func main() {
	go process("lorem.com")
	go process("ipsum.com")
	var input string
	fmt.Scanln(&input)
}

func process(domainName string) {
	mg := mailgun.NewMailgun(domainName, privateKey, publicKey)
	ei := mg.NewEventIterator()

	// Open the files
	baseDirectory, _ := os.Getwd()
	exportedDeliversFileName := baseDirectory + "/" + domainName + "-delivers.txt"
	fileDelivers, errFileDelivers := os.OpenFile(exportedDeliversFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	checkError(errFileDelivers)

	exportedOpensFileName := baseDirectory + "/" + domainName + "-opens.txt"
	fileOpens, errFileOpens := os.OpenFile(exportedOpensFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	checkError(errFileOpens)

	exportedClicksFileName := baseDirectory + "/" + domainName + "-clicks.txt"
	fileClicks, errFileClicks := os.OpenFile(exportedClicksFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	checkError(errFileClicks)

	exportedTagsFileName := baseDirectory + "/" + domainName + "-tags.txt"
	fileTags, errFileTags := os.OpenFile(exportedTagsFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	checkError(errFileTags)

	// Defer to close the file when the function ends
	defer fileDelivers.Close()
	defer fileOpens.Close()
	defer fileClicks.Close()
	defer fileTags.Close()

	var proceed bool = true
	var initial bool = true

	for proceed {
		// Fetch data from Mailgun
		if initial {
			err := ei.GetFirstPage(mailgun.GetEventsOptions{})
			checkError(err)
		} else {
			err := ei.GetNext()
			checkError(err)
		}

		// Update the initial status
		initial = false

		// Compile the data
		var bufferDelivers bytes.Buffer
		var bufferOpens bytes.Buffer
		var bufferClicks bytes.Buffer
		var bufferTags bytes.Buffer

		for _, log := range ei.Events() {
			event := log["event"].(string)
			replacer := strings.NewReplacer("[", "", "]", "", "*", "") // Replace "[" and "]" with an empty string
			tags := replacer.Replace(fmt.Sprintf("%v", log["tags"]))   // convert the interface to string

			// Only process those with campaign tags
			if len(tags) > 0 {

				// Get the subject of the email
				subject := "NULL" // So when we imported this CSV to MySQL later, the subject column would have NULL as its value
				message := log["message"].(interface{}).(map[string]interface{})
				for k1, v1 := range message {
					if k1 == "headers" {
						headers := v1.(interface{}).(map[string]interface{})
						if headers["subject"] != nil {
							subject = "\"" + headers["subject"].(string) + "\""
							fmt.Println(subject)
						}
					}
				}

				bufferTags.WriteString("\"" + tags + "\"," + subject)
				bufferTags.WriteString("\n")
				fmt.Println("[", event, "] (", tags, ")", subject)

				if event == "delivered" {
					bufferDelivers.WriteString(log["recipient"].(string))
					bufferDelivers.WriteString("\n")
				}

				if event == "opened" {
					bufferOpens.WriteString(log["recipient"].(string))
					bufferOpens.WriteString("\n")
				}

				if event == "clicked" {
					bufferClicks.WriteString(log["recipient"].(string))
					bufferClicks.WriteString("\n")
				}
			}
		}

		// Write the data to a file
		fileDelivers.WriteString(bufferDelivers.String())
		fileOpens.WriteString(bufferOpens.String())
		fileClicks.WriteString(bufferClicks.String())
		fileTags.WriteString(bufferTags.String())

		// Check to proceed or not
		if len(ei.Events()) == 0 {
			proceed = false
			fmt.Println("Finished " + domainName + "!")
		}
	}

	fmt.Print(initial)
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
