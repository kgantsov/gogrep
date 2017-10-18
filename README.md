# gogrep
A simple and probably faster replacement of GNU grep command written in GO


Gogrep searching plain-text data sets for lines that match a specified pattern.


# Speed test

Test were performed on MacBook 2,8 GHz Intel Core i5 with SSD disk

Running gogrep in parallel mode:

    time gogrep --p=true "LineProfiler" ~/workspace  84.83s user 61.87s system 213% cpu 1:08.58 total
    time gogrep --p=true "LineProfiler" ~/workspace  83.33s user 62.78s system 207% cpu 1:10.25 total
    time gogrep --p=true "LineProfiler" ~/workspace  85.08s user 61.91s system 214% cpu 1:08.62 total

Running gogrep in serial mode:

    time gogrep "LineProfiler" ~/workspace  23.79s user 59.87s system 50% cpu 2:46.41 total
    time gogrep "LineProfiler" ~/workspace  24.24s user 65.92s system 50% cpu 2:57.76 total
    time gogrep "LineProfiler" ~/workspace  24.33s user 66.59s system 50% cpu 3:00.90 total

Running GNU grep:

    time grep -R "LineProfiler" ~/workspace   295.22s user 65.64s system 82% cpu 7:17.82 total
    time grep -R "LineProfiler" ~/workspace   291.46s user 68.22s system 81% cpu 7:20.85 total
    time grep -R "LineProfiler" ~/workspace   295.80s user 69.85s system 81% cpu 7:31.33 total


# Installation

## Install binary (MacOS and Linux):

    curl -s -L https://github.com/kgantsov/gogrep/releases/download/v0.1/setup.sh | sh

## Install from the source:

First of all `github.com/fatih/color` library needs to be installed.

    go get github.com/fatih/color

Then it can be easily compiled:

    go build

# Usage

To get some help run program with `--help` flag

    gogrep --help
    Usage: gogrep [flags] [pattern] [file]

    Flags:
    -exclude-dir string
            List of coma separated dirs (default ".bzr,CVS,.git,.hg,.svn")
    -include string
            Include pattern (default "*")
    -p	Run gogrep in parallel


# Some examples

Run gogrep excluding `.bzr, CVS, .git, .hg, .svn, env` directories:

    gogrep --exclude-dir=".bzr,CVS,.git,.hg,.svn,env" "import" ~/workspace/

Search only in `.py` files:

    gogrep --include="*.py" "import" /Users/koss/workspace/iconik


Pipling output from other commands:

    cat main.go | gogrep "func "
    func findWordInBuffer(pattern, path string, scanner *bufio.Scanner) {
    func findWordInFile(pattern, path string) {
    func printFile(include, pattern string, excludeDir []string) filepath.WalkFunc {
    func addToQueue(jobs chan string, path string) {
    func worker(id int, pattern string, jobs chan string, wg *sync.WaitGroup) {
    func walkParrallel(dir, pattern string) {
    func main() {
