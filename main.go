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

func setValue(line string, nextLine string, words []string, jsonString string) string {
	var result string
	index := validateKeywords(line, words)
	nextIndex := validateKeywords(nextLine, words)
	if index >= 0 {
		line = strings.Trim(line, " ")
		values := strings.Split(line, words[index]+":")
		if len(jsonString) > 0 {
			tmp := `"` + words[index] + `":"` + values[1]
			if nextIndex == -1 {
				result = tmp
			} else {
				result = tmp[:len(tmp)-1] + `",`
			}
		} else {
			tmp := `{"` + words[index] + `":"` + values[1]
			if nextIndex == -1 {
				result = tmp
			} else {
				result = tmp[:len(tmp)-1] + `",`
			}
		}
	} else {
		if nextIndex == -1 {
			result = line
		} else {
			result = line + `",`
		}
	}
	return result
}

func setAllValues(content string) string {
	lines := strings.Split(content, "\n")
	contentString := make([]string, len(lines)+1)
	words := getKeyWords()
	for i, line := range lines {
		if i == 123 {
			tmp1 := strings.TrimSpace(line)
			fmt.Println(len(tmp1))
		}
		result := strings.Join(contentString, "")
		if i < len(lines)-1 {
			tmp := setValue(strings.TrimSpace(line), lines[i+1], words, result)
			contentString[i] = tmp
		} else {
			tmp := setValue(strings.TrimSpace(line), "", words, result)
			contentString[i] = tmp
		}
	}
	contentString[len(lines)] = `"}`
	totalResult := strings.Join(contentString, "")
	returnValue := totalResult[:len(totalResult)-1]
	//contentString[len(lines)] = `"}`
	return returnValue
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
	//part := emailContent[:261]

	result := setAllValues(emailContent)
	//part := emailContent[:5210]
	//lines := strings.Split(firsLine, "\n")
	//fmt.Println(part)
	fmt.Println("/******/")
	b := strings.ReplaceAll(result, "	", "")
	//c := strings.ReplaceAll(b, "\n", "")
	//fmt.Println(b)
	var jsonData map[string]interface{}
	json.Unmarshal([]byte(b), &jsonData)
	fmt.Println(jsonData)
	//fmt.Println(result)
}
