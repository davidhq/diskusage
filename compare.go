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
    "strings"
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

  fmt.Println(os.Args)

    argsWithoutProg := strings.Fields(os.Args[1]) //split the string

    var a string
    var b string

    if(len(argsWithoutProg) == 0) {
      a = "snapshot_prev.txt"
      b = "snapshot.txt"
    } else {
      a = argsWithoutProg[0]
      b = argsWithoutProg[1]
    }

    //sizes - has to be more than 100k because snapshots ignore smaller files
    const new_limit uint64 = 1000 * bytefmt.KILOBYTE
    const diff_limit uint64 = 1000 * bytefmt.KILOBYTE

    //read previous snapshot

    prev := make(map[string]uint64)

    file, err := os.Open(a)
    check(err)
    defer file.Close()

    re, err := regexp.Compile(`^\s*(\d+)\s+\.?(.*)`)

    var added uint64

    fmt.Print("Reading previous snapshot... " + a)

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        res := re.FindStringSubmatch(scanner.Text())
        if(len(res) > 0) {
          size_temp, _ := strconv.ParseInt(res[1], 10, 64)
          prev[res[2]] = uint64(size_temp)
        }
    }
    check(scanner.Err())

    fmt.Println()

    //read current snapshot

    current := make(map[string]uint64)

    file, err = os.Open(b)
    check(err)
    defer file.Close()

    fmt.Print("Reading current snapshot... " + b)

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

    fmt.Println()
    fmt.Println()

    //check for new files and size increases

    var lineout string
    diff_file, err := os.Create("diff.txt")
    check(err)
    defer diff_file.Close()

    var lines Lines

    for file_name := range current {
      size := current[file_name]
      if(size > 0) {
        old_size, existed := prev[file_name]
        if(!existed && size > new_limit) {
          lines = append(lines, Line{ Size: size, Info: "New " + bytefmt.ByteSize(size) + ":\t" + file_name } )
          added += size
        }
        var diff uint64
        if(size > old_size) {
          diff = size - old_size
          if(existed && diff > diff_limit) {
            lines = append(lines, Line{ Size: diff, Info: "Inc " + bytefmt.ByteSize(diff) + ":\t" + file_name  } )
            added += diff
          }
        }
      }
    }

    sort.Sort(sort.Reverse(lines))

    for _, line := range lines {
      color.Green(line.Info)
      diff_file.WriteString(line.Info)
      diff_file.WriteString("\n")
    }

    lineout = "= Added: " + bytefmt.ByteSize(uint64(added))
    lineout = fmt.Sprintf("\n%s\n\n", lineout)
    fmt.Print(lineout)
    diff_file.WriteString(lineout)

    //check for removed files and size decreases

    lines = make(Lines, 0)
    var removed uint64

    for file_name := range prev {
      size := prev[file_name]
      if(size > 0) {
        old_size, existed := current[file_name]
        if(!existed && size > new_limit) {
          lines = append(lines, Line{ Size: size, Info: "Rem " + bytefmt.ByteSize(size) + ":\t" + file_name } )
          removed += size
        }
        var diff uint64
        if(size > old_size) {
          diff = size - old_size
          if(existed && diff > diff_limit) {
            lines = append(lines, Line{ Size: diff, Info: "Dec " + bytefmt.ByteSize(diff) + ":\t" + file_name } )
            removed += diff
          }
        }
      }
    }

    sort.Sort(sort.Reverse(lines))

    for _, line := range lines {
      color.Yellow(line.Info)
      diff_file.WriteString(line.Info)
      diff_file.WriteString("\n")
    }

    lineout = "= Removed: " + bytefmt.ByteSize(uint64(removed))
    lineout = fmt.Sprintf("\n%s\n\n", lineout)
    fmt.Print(lineout)
    diff_file.WriteString(lineout)

    // NET

    diff := int64(added) - int64(removed)
    if(diff < 0) {
      lineout = "= Net total: -" + bytefmt.ByteSize(uint64(-diff))
      color.Yellow(lineout)
      diff_file.WriteString(lineout + "\n")
    } else {
      lineout = "= Net total: " + bytefmt.ByteSize(uint64(diff))
      color.Green(lineout)
      diff_file.WriteString(lineout + "\n")
    }
}
