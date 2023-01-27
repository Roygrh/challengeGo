package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	words := []string{"Message-ID", "Date", "From", "To", "Subject",
		"Cc", "Mime_Version", "Content_Type", "Content_Transfer_Encoding",
		"Bcc", "X_From", "X_To", "X_bcc", "X_Folder", "X_Origin", "X_FileName"}
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
		result := strings.Contains(line, word+":")
		if result {
			return i
		}
	}
	return index
}

func setValue(line string, words []string, jsonString string) string {
	result := ``
	index := validateKeywords(line, words)

	if index >= 0 {
		line = strings.Trim(line, " ")
		values := strings.Split(line, words[index]+":")
		if len(jsonString) > 0 {
			tmp := `","` + words[index] + `":"` + values[1]
			result = result + tmp
		} else {
			tmp := `{"` + words[index] + `":"` + values[1]
			result = result + tmp
		}
	} else {
		result = line
	}
	return result
}

func setAllValues(content string) string {
	lines := strings.Split(content, "\n")
	contentString := make([]string, len(lines)+2)
	//contentString[0] = "{"
	words := getKeyWords()
	for i, line := range lines {
		result := strings.Join(contentString, "")
		tmp := setValue(line, words, result)
		contentString[i] = tmp
		if i == 0 {
			contentString[i] = tmp + "\n"
		} else {
			contentString[i] = tmp + "\n"
		}
	}
	contentString[len(lines)] = `"}`
	return strings.Join(contentString, "")
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
	part := emailContent[:100]

	result := setAllValues(part)
	//part := emailContent[:5210]
	//lines := strings.Split(firsLine, "\n")

	/*re := regexp.MustCompile("Message-ID: (.*)\nDate: (.*)\nFrom: (.*)\nTo: (.*)")
	match := re.FindStringSubmatch(firsLine)
	if match != nil {
		//fmt.Println(firsLine)
	}

	line := "Message-ID: <1293396.1075840371988.JavaMail.evans@thyme>"
	parts := strings.Split(line, "Message-ID:")*/

	fmt.Println("/******/")
	fmt.Println(result)
	fmt.Println(len(result))
	/*a := `{"Message-ID":" <1293396.1075840371988.JavaMail.evans@thyme>
	","Date":" Wed, 6 Feb 2002 16:09:37 -0800 (PST)
	"}`*/
	b := strings.ReplaceAll(result, " ", "")
	//c := strings.ReplaceAll(b, "\n", "")
	fmt.Println(b)
	var jsonData map[string]interface{}
	json.Unmarshal([]byte(b), &jsonData)
	fmt.Println(jsonData)
	//fmt.Println(result)
}
