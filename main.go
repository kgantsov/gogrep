package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func findWordInBuffer(pattern, path string, scanner *bufio.Scanner) {
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, pattern) {
			cyan := color.New(color.FgCyan).SprintFunc()
			red := color.New(color.FgRed).SprintFunc()
			line = strings.Replace(line, pattern, red(pattern), -1)

			if len(path) == 0 {
				fmt.Println(fmt.Sprintf("%s", line))
			} else {
				fmt.Println(fmt.Sprintf("%s: %s", cyan(path), line))
			}
		}
	}
}

func findWordInFile(pattern, path string) {
	inFile, _ := os.Open(path)
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)

	findWordInBuffer(pattern, path, scanner)
}

func printFile(include, pattern string, excludeDir []string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}

		if info.IsDir() {
			dir := filepath.Base(path)
			for _, d := range excludeDir {
				if d == dir {
					return filepath.SkipDir
				}
			}
		}
		if !info.IsDir() {
			matched, err := filepath.Match(include, info.Name())
			if err != nil {
				fmt.Println(err)
				return err
			}
			if matched {
				findWordInFile(pattern, path)
			}
		}
		return nil
	}
}

func main() {
	info, _ := os.Stdin.Stat()

	if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		excludeDir := flag.String(
			"exclude-dir",
			".bzr,CVS,.git,.hg,.svn",
			"List of coma separated dirs (Default is: .bzr,CVS,.git,.hg,.svn)",
		)
		include := flag.String("include", "*", "Include pattern (Default is: *)")
		parrallel := flag.Bool("p", false, "Run walk in parallel")

		flag.Parse()

		args := flag.Args()
		if len(args) != 2 {
			log.Print("Not enough arguments")
			return
		}

		pattern := flag.Args()[0]
		dir := flag.Args()[1]

		if *parrallel {
			fmt.Println("Not implemented yet")
		} else {
			err := filepath.Walk(dir, printFile(*include, pattern, strings.Split(*excludeDir, ",")))
			if err != nil {
				log.Fatal(err)
			}
		}
	} else if info.Size() > 0 {
		flag.Parse()

		args := flag.Args()
		if len(args) != 1 {
			log.Print("Not enough arguments")
			return
		}

		pattern := flag.Args()[0]
		scanner := bufio.NewScanner(os.Stdin)

		findWordInBuffer(pattern, "", scanner)
	}
}
