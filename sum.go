package main

import "github.com/pivotal-golang/bytefmt"
import "github.com/fatih/color"

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strconv"
    "strings"
)

// Reading files requires checking most calls for errors.
// This helper will streamline our error checks below.
func check(e error) {
    if e != nil {
      //log.Fatal(e)
      panic(e)
    }
}

func SizeInBytes(s string) (uint64) {
  magnitude := s[len(s)-1:]
  size := s[:len(s)-1]
  switch(magnitude) {
    case "K":
      n, _ := strconv.ParseFloat(size, 64)
      return uint64(n * bytefmt.KILOBYTE)
    case "M":
      n, _ := strconv.ParseFloat(size, 64)
      return uint64(n * bytefmt.MEGABYTE)
    case "G":
      n, _ := strconv.ParseFloat(size, 64)
      return uint64(n * bytefmt.GIGABYTE)
    case "T":
      n, _ := strconv.ParseFloat(size, 64)
      return uint64(n * bytefmt.TERABYTE)
    default:
      n, _ := strconv.ParseFloat(s, 64)
      return uint64(n)
  }
}

func main() {

    argsWithoutProg := os.Args[1:]

    if(len(argsWithoutProg) == 0) {
      fmt.Println("Usage: go run sum.go filter")
      return
    }
    filter := strings.Join(argsWithoutProg, " ")

    file, err := os.Open("diff.txt")
    check(err)
    defer file.Close()

    var added, removed uint64

    re, err := regexp.Compile(`(.+)\:`)

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if(len(line) > 3 && strings.Contains(line[3:], filter)) {
          if(line[0:3] == "Inc" || line[0:3] == "New") {
            res := re.FindStringSubmatch(line[4:])
            if(len(res) > 0) {
              added += SizeInBytes(res[1])
            }
          } else if (line[0:3] == "Dec" || line[0:3] == "Rem") {
            res := re.FindStringSubmatch(line[4:])
            if(len(res) > 0) {
              removed += SizeInBytes(res[1])
            }
          }

        }
    }

    color.Green("Added: " + bytefmt.ByteSize(added))
    color.Yellow("Removed: " + bytefmt.ByteSize(removed))

    diff := int64(added) - int64(removed)
    if(diff < 0) {
      color.Yellow("= Net total: -" + bytefmt.ByteSize(uint64(-diff)))
    } else {
      color.Green("= Net total: " + bytefmt.ByteSize(uint64(diff)))
    }
}
