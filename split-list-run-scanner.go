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

const replaceString = "$LIST"

var sortedList []string

func main() {
	if runtime.GOOS == "windows" {
		reset = ""
		red = ""
		yellow = ""
		green = ""
	}

	pathList, pathOutput, linesLimit, command := parseFlags()

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
		pathOutput += "/"
	}
	if !strings.Contains(command, replaceString) {
		fmt.Printf("%sWarning: -command doesn't contain %s%s", yellow, replaceString, reset)
	}

	sliceList := readLocalFile(pathList)

	fmt.Printf("Splitting %d lines after every %dth line\n", len(sliceList), linesLimit)

	var f *os.File
	var err error
	var paths []string

	for i, line := range sliceList {
		if i%linesLimit == 0 {
			path := fmt.Sprintf("%s%d/", pathOutput, i)
			f, err = os.Create(path)
			if err != nil {
				msg := "Log: " + err.Error() + "\n"
				fmt.Print(msg)
			}
			paths = append(paths, path)
			defer f.Close()
		}
		f.WriteString(line + "\n")
	}

	for _, p := range paths {
		commandNew := strings.Replace(command, replaceString, p, -1)
		cmd := exec.Command(commandNew)
		err := cmd.Run()

		if err != nil {
			fmt.Printf("%sError: %s%s\n", red, err, reset)
		}
	}

	fmt.Printf("%sFinished.%s", green, reset)
}

func readLocalFile(path string) []string {

	w, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("%s%s: %s%s\n", red, path, err.Error(), reset)
		os.Exit(3)
	}

	return strings.Split(string(w), "\n")
}

func parseFlags() (string, string, int, string) {
	var pathList string
	var pathOutput string
	var lines int
	var command string

	flag.StringVar(&pathList, "list", "", "path to the list")
	flag.StringVar(&pathOutput, "output", "", "path to output folder")
	flag.IntVar(&lines, "lines", 1000, "after how many lines should be splitted? Default is 1000")
	flag.StringVar(&pathOutput, "command", "", "command to run. Use $LIST where the path of a splitted list shall be inserted")

	flag.Parse()

	return pathList, pathOutput, lines, command
}
