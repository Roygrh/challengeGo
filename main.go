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
	words := []string{"Message-ID", "X-From", "X-To", "X-cc", "X-bcc", "X-Folder", "X-Origin",
		"X-FileName", "Date", "From", "To", "Subject", "Cc", "Mime-Version",
		"Content-Type", "Content-Transfer-Encoding", "Bcc"}
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

func setSimpleValue(line string, index int, words []string, jsonString string) string {
	var result string
	//index := validateKeywords(line, words)

	if len(line) == 0 {
		return result
	}

	if index >= 0 {
		values := strings.Split(line, words[index]+":")
		value := values[1]
		tmp := `"` + words[index] + `":"` + values[1]
		if len(value) > 0 && value[len(value)-1:] == "," {
			result = tmp
		} else {
			result = tmp[:len(tmp)] + `",`
		}
	} else {
		if len(line) > 0 && line[len(line)-1:] == "," {
			result = line
		} else {
			result = line[:len(line)] + `",`
		}
	}
	return result
}

func setValue(line string, nextLine string, words []string, jsonString string) string {
	var result string
	index := validateKeywords(line, words)
	nextIndex := validateKeywords(nextLine, words)
	if index >= 0 {
		values := strings.Split(line, words[index]+":")
		if len(jsonString) > 0 {
			tmp := `"` + words[index] + `":"` + values[1]
			if nextIndex == -1 {
				result = tmp
			} else {
				result = tmp[:len(tmp)] + `",`
			}
		} else {
			tmp := `{"` + words[index] + `":"` + values[1]
			if nextIndex == -1 {
				result = tmp
			} else {
				result = tmp[:len(tmp)] + `",`
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
	index := getIndexOfMessageLine(lines, words)
	for i, line := range lines {
		result := strings.Join(contentString, "")
		trimmedLine := strings.TrimSpace(line)

		if i < index {
			if i < len(lines)-1 {
				tmp := setValue(trimmedLine, lines[i+1], words, result)
				contentString[i] = tmp
			}
		} else {
			if index == i {
				tmp := setValue(trimmedLine, lines[i+1], words, result)
				contentString[i] = `","Message":"` + tmp
			} else if i < len(lines)-1 {
				tmp := setValue(trimmedLine, lines[i+1], words, result)
				contentString[i] = tmp + "\n"
			} else {
				tmp := setValue(trimmedLine, "", words, result)
				contentString[i] = tmp
			}
		}
		/*if index == i {
			tmp := setValue(trimmedLine, lines[i+1], words, result)
			contentString[i] = `","Message":"` + tmp
		} else {
			if i < len(lines)-1 {
				tmp := setValue(trimmedLine, lines[i+1], words, result)
				contentString[i] = tmp
			} else if i == len(lines)-1 {
				tmp := setValue(trimmedLine, "", words, result)
				contentString[i] = tmp
			}
		}*/
	}
	//contentString[contentIndex+1] = `"Message":` + contentString[contentIndex+1]
	contentString[len(lines)] = `"}`
	totalResult := strings.Join(contentString, "")
	/*fmt.Println(totalResult)
	fmt.Println("********************")*/
	returnValue := totalResult[:len(totalResult)]
	return returnValue
}

func setAllSimpleValues(content string) string {
	lines := strings.Split(content, "\n")
	contentString := make([]string, len(lines)+2)
	words := getKeyWords()
	contentString[0] = "{"
	for i, line := range lines {
		index := validateKeywords(line, words)
		result := strings.Join(contentString, "")
		trimmedLine := strings.TrimSpace(line)
		if i < len(lines)-1 {
			tmp := setSimpleValue(trimmedLine, index, words, result)
			contentString[i+1] = tmp
		} else if i == len(lines)-1 {
			tmp := setSimpleValue(trimmedLine, index, words, result)
			contentString[i+1] = tmp
		}
	}
	contentString[len(lines)+1] = "}"
	totalResult := strings.Join(contentString, "")
	returnValue := totalResult[:len(totalResult)]
	return returnValue
}

func getIndexOfMessageLine(lines []string, words []string) int {
	var index int
	for i, line := range lines {
		position := validateKeywords(line, words)
		if position >= 0 && len(line) > 0 {
			index = i
		}
	}
	return index + 1
}

func loadDataFromFile(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}
	return content
}

func transformDataToString(data []byte) string {
	rawString := string(data)
	firstFilter := strings.ReplaceAll(rawString, `\`, `\\`)
	secondFilter := strings.ReplaceAll(firstFilter, "\t", "   ")
	thirdFilter := strings.ReplaceAll(secondFilter, "'", " ")
	fourthFilter := strings.ReplaceAll(thirdFilter, `"`, " ")
	return fourthFilter
}

func readDirectory(directoryPath string) {

	files, err := os.Open(directoryPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer files.Close()

	fileInfo, err := files.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range fileInfo {
		if file.IsDir() {
			// aquí puedes hacer una llamada recursiva para explorar la subcarpeta
		} else {
			fmt.Println(file.Name())
		}
	}
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

	emailContent := transformDataToString(content)

	result := setAllValues(emailContent)
	//result2 := setAllSimpleValues(emailContent)
	//fmt.Println(result2[:len(result2)-2] + "}")
	fmt.Println(result)
	fmt.Println("*************************")
	var jsonData map[string]interface{}
	json.Unmarshal([]byte(result), &jsonData)
	//json.Unmarshal([]byte(result2[:len(result2)-2]+"}"), &jsonData)

	fmt.Println(jsonData)
	fmt.Println("*************************")

	readDirectory("all_documents")

}
