package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type Email struct {
	Message_ID                string
	Date                      string
	From                      string
	To                        string
	Subject                   string
	Cc                        string
	Mime_Version              string
	Content_Type              string
	Content_Transfer_Encoding string
	Bcc                       string
	X_From                    string
	X_To                      string
	X_bcc                     string
	X_Folder                  string
	X_Origin                  string
	X_FileName                string
	Message                   string
}

func getKeyWords() []string {
	words := []string{"Message-ID:", "Date:", "From:", "To:", "Subject:",
		"Cc:", "Mime_Version:", "Content_Type:", "Content_Transfer_Encoding:",
		"Bcc:", "X_From:", "X_To:", "X_bcc:", "X_Folder:", "X_Origin:", "X_FileName:"}
	return words
}

func getDataFromFile(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	emailContent := string(content)

	lines := strings.Split(emailContent, "\n")

	return lines
}

func validateKeywords(line string, words []string) int {
	index := -1
	for i, word := range words {
		result := strings.Contains(line, word)
		if result {
			return i
		}
	}
	return index
}

func setValue(line string, words []string, jsonString string) string {
	index := validateKeywords(line, words)
	if index >= 0 {
		line = strings.Trim(line, " ")
		jsonString += words[index]
	}
	return jsonString
}

func main() {
	file, err := os.Open("resources/email.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	emailContent := string(content)
	firsLine := emailContent[:5200]
	//lines := strings.Split(firsLine, "\n")

	re := regexp.MustCompile("Message-ID: (.*)\nDate: (.*)\nFrom: (.*)\nTo: (.*)")
	match := re.FindStringSubmatch(firsLine)
	if match != nil {
		//fmt.Println(firsLine)
	}

	line := "Message-ID: <1293396.1075840371988.JavaMail.evans@thyme>"
	parts := strings.Split(line, "Message-ID:")

	fmt.Println(strings.Trim(parts[1], " "))
}
