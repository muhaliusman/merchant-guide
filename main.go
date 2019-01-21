package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// make cLine type struct
type cLine struct {
	currentLine string // To save sentence in current line
	splitFirst  string // To save the first index of results of the split function
	splitLast   string // To save the last index of results of the split function
}

// Global variable to save symbols and value
var symbols = map[string]float64{
	"I": 1,
	"V": 5,
	"X": 10,
	"L": 50,
	"C": 100,
	"D": 500,
	"M": 1000,
}

// Global variable to save aliases
var aliases = map[string]string{}

// Global variable to save credit values
var creditValues = map[string]float64{}

func main() {
	// Open file input
	file, err := os.Open("input.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	// scan line by line
	for scanner.Scan() {
		// Initialize new cLine Type
		currentLine := cLine{}
		// Update currentLine field in currentLine struct (cLine type)
		currentLine.updateCurrentLine(scanner.Text())
		// Execute line
		currentLine.sentenceExecution()
	}
	// if scan process has error then log the error
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (cl *cLine) updateCurrentLine(current string) {
	// update field currentLine in cLine type by reference
	cl.currentLine = current
}

func (cl *cLine) updateSplitResult(first string, last string) {
	// Update splitFirst splitLast field in cLine type by reference
	cl.splitFirst = first
	cl.splitLast = last
}

func (cl cLine) sentenceExecution() {
	// Split the sentence by 'is'
	err := cl.splitSentence()

	// If has error then print to console
	if err != nil {
		fmt.Println(err)
		return
	}

	// Check whether the sentence is a question or statement sentence
	if cl.isQuestion() {
		// If this is a question then answer
		cl.answerQuestion()
	} else {
		// If this is statement then asign value to variable
		cl.assign()
	}
}

func (cl cLine) isQuestion() bool {
	// Trim space before first and and last sentence
	s := strings.TrimSpace(cl.currentLine)

	lastChar := s[len(s)-1:] // Get last char after trimspace

	if lastChar == "?" {
		// If last char is "?" it means this is a question
		return true
	}

	return false
}

func (cl *cLine) splitSentence() error {
	// Split the sentence by "is"
	s := strings.Split(string(cl.currentLine), " is ")

	if len(s) != 2 {
		// If there is no "is" then throw error
		return errors.New("I have no idea what you are talking about")
	}

	// If "is" is exist then update update cLine type by reference
	cl.updateSplitResult(s[0], s[1])
	// then return nil error
	return nil
}

func (cl cLine) assign() {
	// Check whether it is a unit assignment or credit assignment
	if cl.isUnitAssigmnment() {
		assignUnit(cl.splitFirst, cl.splitLast)
	} else if cl.isCreditAssignment() {
		assignCredit(cl.splitFirst, cl.splitLast)
	} else {
		// if this is not unit or credit assignment then throw fatal error
		panic(errors.New("Cannot assign value"))
	}
}

func (cl cLine) isUnitAssigmnment() bool {
	// If symbol exist it means this is unit assignment
	_, exist := symbols[strings.TrimSpace(cl.splitLast)]

	return exist
}

func (cl cLine) isCreditAssignment() bool {
	// If sentence contain "credits" it means this sentence is credit assignment
	if strings.Contains(cl.splitLast, "Credits") {
		return true
	}

	return false
}

func assignUnit(firstSentence string, lastSentence string) {
	// save alias based symbols value and then save to aliases global variable
	aliases[firstSentence] = lastSentence
}

func assignCredit(firstSentence string, lastSentence string) {
	// Split lastsentence by space
	val := strings.Split(lastSentence, " ")

	// Declare variable
	var iCredit int
	var cVal int
	var err error
	var creditVal float64

	for i, v := range val {
		if strings.Contains(strings.ToLower(v), "credit") {
			// Lok for index of "Credit" words
			iCredit = i
		}
	}

	if iCredit > 0 {
		// If credit words is exist it mean the value of credit in index of "credit " words - 1
		cVal, err = strconv.Atoi(val[iCredit-1])

		if err != nil {
			panic(err)
		}
	} else {
		panic(errors.New("Cannot asign the value"))
	}

	sCredit := strings.Split(firstSentence, " ")

	// Find credit unit in sentence
	x, alias := sCredit[len(sCredit)-1], sCredit[:len(sCredit)-1]

	// Convert the unit alias to integer
	totalVal := convertToInteger(alias)

	// To handle divided by zero
	if totalVal > 0 {
		creditVal = float64(cVal) / float64(totalVal)
	}

	// Save credit value to global slice variable
	creditValues[x] = creditVal
}

func convertToInteger(alias []string) float64 {
	// Declare buffer to concat process
	var buffer bytes.Buffer
	for _, v := range alias {
		if _, ok := aliases[v]; ok {
			for i, a := range aliases {
				if i == v {
					// write roman char
					buffer.WriteString(a)
				}
			}
		}
	}
	// Get total based on roman char
	total := getTotal(buffer.String())
	// Reset buffer
	buffer.Reset()

	return total
}

func (cl cLine) answerQuestion() {
	var buffer bytes.Buffer
	// Remove "?" char from sentence
	lastSentence := strings.TrimSuffix(cl.splitLast, "?")
	// Then split
	sString := strings.Split(lastSentence, " ")
	// Get Total
	total := convertToInteger(sString)

	buffer.WriteString(lastSentence)
	buffer.WriteString(" is ")

	if cl.isQuestionCredit() {
		// if is question about credit unit then multiply the total by the value of the unit
		sCredit := strings.Split(lastSentence, " ")
		creditUnit := sCredit[len(sCredit)-1]

		creditValue := creditValues[creditUnit]

		total = total * creditValue
	}

	totalStr := strconv.Itoa(int(total))
	buffer.WriteString(totalStr)

	if cl.isQuestionCredit() {
		buffer.WriteString(" Credits")
	}

	fmt.Println(buffer.String())
	buffer.Reset()
}

func (cl cLine) isQuestionCredit() bool {
	// If the sentence contain "how many words" it means this is question about credit
	if strings.Contains(cl.splitFirst, "how many") {
		return true
	}

	return false
}

func getTotal(str string) float64 {
	total := float64(0)
	for i := 0; i < len(str); i++ {
		if i+1 < len(str) {
			// If the roman symbol has less value than the next roman symbol value, then subtract
			if symbols[str[i+1:(i+2)]] > symbols[str[i:(i+1)]] {
				total += symbols[str[i+1:(i+2)]] - symbols[str[i:(i+1)]]
				// Jump one iteration
				i++
				continue
			}
		}

		total += symbols[str[i:(i+1)]]
	}

	return total
}
