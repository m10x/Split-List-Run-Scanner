package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	reset  = "\033[0m"
	red    = "\033[31m"
	yellow = "\033[33m"
	green  = "\033[32m"
)

/*
go run .\split-list-run-scanner.go -list "C:\Users\max\Documents\git\Web-Cache-Vulnerability-Scanner\test\adobe\domains" -lines 500 -outputpath "C:\Users\max\Documents\git\Web-Cache-Vulnerability-Scanner\test\adobe\new"
*/

const replaceString = "$LIST"

var sortedList []string

func main() {
	if runtime.GOOS == "windows" {
		reset = ""
		red = ""
		yellow = ""
		green = ""
	}

	pathList, pathOutput, nameOutput, linesLimit, command := parseFlags()

	if pathList == "" {
		fmt.Printf("%sError: -list wasn't specified%s\n", red, reset)
		os.Exit(1)
	}
	if pathOutput == "" {
		fmt.Printf("%sError: -output wasn't specified%s\n", red, reset)
		os.Exit(2)
	} else {
		pathOutput = strings.TrimSuffix(pathOutput, "/")
		pathOutput = strings.TrimSuffix(pathOutput, "\\")
		pathOutput += "\\"
	}
	if !strings.Contains(command, replaceString) {
		fmt.Printf("%sWarning: -command doesn't contain %s%s\n", yellow, replaceString, reset)
	}

	sliceList := readLocalFile(pathList)

	fmt.Printf("Splitting %d lines after every %dth line\n", len(sliceList), linesLimit)

	var f *os.File
	var err error
	var paths []string

	for i, line := range sliceList {
		if i%linesLimit == 0 {
			path := fmt.Sprintf("%s%d\\", pathOutput, i)
			err = os.MkdirAll(path, 0755)
			if err != nil {
				fmt.Printf("%sError MkDir:%s%s\n", red, err.Error(), reset)
			}
			f, err = os.Create(path + nameOutput)
			if err != nil {
				fmt.Printf("%sError CreateFile:%s%s\n", red, err.Error(), reset)
			}
			paths = append(paths, path)
			defer f.Close()
		}
		f.WriteString(line)
		if i%linesLimit != linesLimit-1 && i+1 < len(sliceList) {
			f.WriteString("\n")
		}
	}

	for _, p := range paths {
		commandNew := strings.Replace(command, replaceString, p+nameOutput, -1)
		commandNew = strings.Replace(commandNew, "$PATH", p, -1)
		fmt.Println(commandNew)

		cmd := exec.Command("powershell", "start-process", "powershell.exe", "-argument", "'"+commandNew+"'")
		cmd.Start()
	}

	fmt.Printf("%sFinished.%s", green, reset)
}

func readLocalFile(path string) []string {

	w, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("%sError while Reading list %s: %s%s\n", red, path, err.Error(), reset)
		os.Exit(3)
	}

	return strings.Split(string(w), "\n")
}

func parseFlags() (string, string, string, int, string) {
	var pathList string
	var pathOutput string
	var lines int
	var command string
	var nameOutput string

	flag.StringVar(&pathList, "list", "", "path to the list")
	flag.StringVar(&pathOutput, "outputpath", "", "path to output folder")
	flag.StringVar(&nameOutput, "outputname", "list", "name for output file. Default is 'list'")
	flag.IntVar(&lines, "lines", 1000, "after how many lines should be splitted? Default is 1000")
	flag.StringVar(&command, "command", "", "command to run. Use $LIST where the path of a splitted list shall be inserted")

	flag.Parse()

	return pathList, pathOutput, nameOutput, lines, command
}
