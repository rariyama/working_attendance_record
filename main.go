package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func outputShutdownRecords() []string {
	output, err := exec.Command("last", "shutdown").Output()
	if err != nil {
		fmt.Println(err)
	}
	headerRemoved := strings.Replace(string(output), "shutdown  ~                         ", string(' '), -1)
	paddingZero := strings.Replace(headerRemoved, "  ", " 0", -1)
	suffixRemoved := strings.Replace(paddingZero, "wtmp begins ", string(' '), -1)
	shutdownStr := strings.Split(strings.TrimSpace(suffixRemoved), "\n")
	shutdownRecs := shutdownStr[0 : len(shutdownStr)-2]
	return shutdownRecs
}

func exportAlphabetMonth(month string) string {
	switch month {
	case "01":
		return "Jan"
	case "02":
		return "Feb"
	case "03":
		return "Mar"
	case "04":
		return "Apr"
	case "05":
		return "May"
	case "06":
		return "Jun"
	case "07":
		return "Jul"
	case "08":
		return "Aug"
	case "09":
		return "Sep"
	case "10":
		return "Oct"
	case "11":
		return "Nov"
	case "12":
		return "Dec"
	}
	return "target_month is not returned"
}

func main() {

	flag.Parse()
	args := flag.Args()
	//use working year and month as command line arguments
	targetYear := args[0]
	targetMonth := args[1]
	alphabetMonth := exportAlphabetMonth(targetMonth)

	shutdownRecords := outputShutdownRecords()
	var targetTerms = [][]string{}
	for _, record := range shutdownRecords {
		//the example of recordList : ["Thu", "Mar", "26", "19:36"]
		recordList := strings.Split(strings.TrimSpace(record), string(' '))
		if recordList[1] == alphabetMonth {
			targetDate := targetYear + "-" + targetMonth + "-" + recordList[2]
			targetRecord := []string{targetDate, "09:30", recordList[3]}
			targetTerms = append(targetTerms, targetRecord)
		}
	}
	//last shutdown can only show results in descendant order, so need to arrange slice in ascendant order.
	var targetTermsAsc = [][]string{}
	for i := 1; i <= len(targetTerms); i++ {
		targetTermsAsc = append(targetTermsAsc, targetTerms[len(targetTerms)-i])
	}

	// write the contents of slice on csv file
	file, err := os.Create("./record.csv")
	if err != nil {
		fmt.Println(err)
	}

	writer := csv.NewWriter(file)
	for _, term := range targetTermsAsc {
		if err := writer.Write(term); err != nil {
			fmt.Println(err)
		}
	}
	writer.Flush()
}
