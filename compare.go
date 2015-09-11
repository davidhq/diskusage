// Reading and writing files are basic tasks needed for
// many Go programs. First we'll look at some examples of
// reading files.

package main

import "github.com/pivotal-golang/bytefmt"
import "github.com/fatih/color"

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strconv"
    "log"
    //"strings"
    "sort"
)

type Line struct {
    Size uint64
    Info string
}

type Lines []Line

func (slice Lines) Len() int {
    return len(slice)
}

func (slice Lines) Less(i, j int) bool {
    return slice[i].Size < slice[j].Size;
}

func (slice Lines) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}


// Reading files requires checking most calls for errors.
// This helper will streamline our error checks below.
func check(e error) {
    if e != nil {
      //log.Fatal(e)
      panic(e)
    }
}

func IsDirectory(path string) (bool) {
    if info, err := os.Stat(path); err == nil && info.IsDir() {
       return true
    } else {
      return false
    }
}

func main() {

    a := "snapshot_prev.txt"
    b := "snapshot.txt"
    //TESTING
    a = "test_a.txt"
    b = "test_b.txt"
    const new_limit uint64 = 1000 //kilo
    const diff_limit uint64 = 1000 //kilo

    //read previous snapshot

    prev := make(map[string]uint64)

    file, err := os.Open(a)
    check(err)
    defer file.Close()

    re, err := regexp.Compile(`(\d+)\s+\.?(.*)`)

    var total uint64

    fmt.Print("Reading previous snapshot... ")

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        res := re.FindStringSubmatch(scanner.Text())
        if(len(res) > 0) {
          size_temp, _ := strconv.ParseInt(res[1], 10, 64)
          prev[res[2]] = uint64(size_temp)
        }
    }
    check(scanner.Err())

    fmt.Println("done")

    //read current snapshot

    current := make(map[string]uint64)

    file, err = os.Open(b)
    check(err)
    defer file.Close()

    fmt.Print("Reading current snapshot... ")

    scanner = bufio.NewScanner(file)
    for scanner.Scan() {
        res := re.FindStringSubmatch(scanner.Text())
        if(len(res) > 0) {
          size_temp, _ := strconv.ParseInt(res[1], 10, 64)
          current[res[2]] = uint64(size_temp)
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    fmt.Println("done")

    fmt.Println()

    //check for new files and size increases

    var lines Lines

    for file_name := range current {
      size := current[file_name]
      if(size > 0) {
        old_size, existed := prev[file_name]
        if(!existed && size > new_limit) {
          lines = append(lines, Line{ Size: size, Info: "New : " + file_name + " | Size: " + bytefmt.ByteSize(size * bytefmt.KILOBYTE) } )
          total += size
        }
        var diff uint64
        if(size > old_size) {
          diff = size - old_size
          if(existed && diff > diff_limit) {
            lines = append(lines, Line{ Size: diff, Info: " Inc: " + file_name + " | Diff: " + bytefmt.ByteSize(diff * bytefmt.KILOBYTE) } )
            total += diff
          }
        }
      }
    }

    sort.Sort(sort.Reverse(lines))

    for _, line := range lines {
      color.Green(line.Info)
    }

    fmt.Println()
    fmt.Println("Total: " + bytefmt.ByteSize(uint64(total * bytefmt.KILOBYTE)))

    fmt.Println()

    //check for removed files and size decreases

    lines = make(Lines, 0)
    total = 0

    for file_name := range prev {
      size := prev[file_name]
      if(size > 0) {
        old_size, existed := current[file_name]
        if(!existed && size > new_limit) {
          lines = append(lines, Line{ Size: size, Info: "Del : " + file_name + " | Size: " + bytefmt.ByteSize(size * bytefmt.KILOBYTE) } )
          total += size
        }
        var diff uint64
        if(size > old_size) {
          diff = size - old_size
          if(existed && diff > diff_limit) {
            lines = append(lines, Line{ Size: diff, Info: " Dec: " + file_name + " | Diff: " + bytefmt.ByteSize(diff * bytefmt.KILOBYTE) } )
            total += diff
          }
        }
      }
    }

    sort.Sort(sort.Reverse(lines))

    for _, line := range lines {
      color.Red(line.Info)
    }

    fmt.Println()
    fmt.Println("Total: " + bytefmt.ByteSize(uint64(total * bytefmt.KILOBYTE)))
}
