package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

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
				continue
			}

			fmt.Println(fmt.Sprintf("%s: %s", cyan(path), line))
		}
	}
}

func findWordInFile(pattern, path string) error {
	inFile, err := os.Open(path)

	if err != nil {
		return err
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)

	findWordInBuffer(pattern, path, scanner)

	return nil
}

func printFile(include, pattern string, excludeDir []string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print("Walking error: ", err)
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
				fmt.Println("File path matching error: ", err)
				return err
			}
			if matched {
				err = findWordInFile(pattern, path)
				if err != nil {
					log.Print("Error finding word in file: ", err)
				}
			}
		}
		return nil
	}
}

func worker(id int, pattern string, dirsChan chan string, wg *sync.WaitGroup) {
	for {
		select {
		case dir := <-dirsChan:
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				log.Fatal("Reading directory error: ", err)
			}

			for _, file := range files {
				if file.IsDir() {
					go func(dirsC chan string, path string) {
						dirsC <- path
					}(dirsChan, fmt.Sprintf("%s/%s", dir, file.Name()))
					wg.Add(1)
				} else {
					findWordInFile(pattern, fmt.Sprintf("%s/%s", dir, file.Name()))
				}
			}
			wg.Done()
		default:
		}
	}
}

func walkParrallel(dir, pattern string) {
	var wg sync.WaitGroup

	numWorkers := 4
	if n := runtime.NumCPU(); n > numWorkers {
		numWorkers = n
	}

	dirsChan := make(chan string, numWorkers)

	for w := 1; w <= numWorkers; w++ {
		go worker(w, pattern, dirsChan, &wg)
	}

	go func(dirsC chan string, path string) {
		dirsC <- path
	}(dirsChan, dir)
	wg.Add(1)

	wg.Wait()
}

func main() {
	info, _ := os.Stdin.Stat()

	if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		excludeDir := flag.String(
			"exclude-dir",
			".bzr,CVS,.git,.hg,.svn",
			"List of coma separated dirs",
		)
		include := flag.String("include", "*", "Include pattern")
		parrallel := flag.Bool("p", false, "Run gogrep in parallel")

		flag.Usage = func() {
			fmt.Fprintf(
				os.Stderr,
				"Usage: %s [flags] [pattern] [file]\n\nFlags:\n",
				os.Args[0],
			)
			flag.PrintDefaults()
		}

		flag.Parse()

		args := flag.Args()
		if len(args) != 2 {
			flag.Usage()
			return
		}

		pattern := flag.Args()[0]
		dir := flag.Args()[1]

		if *parrallel {
			walkParrallel(dir, pattern)
		} else {
			err := filepath.Walk(dir, printFile(*include, pattern, strings.Split(*excludeDir, ",")))
			if err != nil {
				log.Fatal("Walking error: ", err)
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
