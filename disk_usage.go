// Reading and writing files are basic tasks needed for
// many Go programs. First we'll look at some examples of
// reading files.

package main

import "github.com/pivotal-golang/bytefmt"

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

    //existing := make([]string, 0)
    //var existing []string

    existing := make(map[string]uint64)

    file, err := os.Open("snapshot_prev.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    re, err := regexp.Compile(`(\d+)\s+\.?(.*)`)

    var total uint64
    var size uint64

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        res := re.FindStringSubmatch(scanner.Text())
        if(len(res) > 0) {
          size_temp, _ := strconv.ParseInt(res[1], 10, 64)
          existing[res[2]] = uint64(size_temp)
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    fmt.Println("Finished reading the old snapshot...")

    // read new snapshot

    file_new, err := os.Open("snapshot.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file_new.Close()

    var lines Lines

    scanner = bufio.NewScanner(file_new)
    for scanner.Scan() {
        res := re.FindStringSubmatch(scanner.Text())
        if(len(res) > 0) {
          var size_temp, _ = strconv.ParseInt(res[1], 10, 64)
          size = uint64(size_temp)
          var file_name = res[2]
          old_size, existed := existing[file_name]
          // !strings.Contains(file_name, "Spotlight-V100")
          if(!existed && size > 100) {
            lines = append(lines, Line{ Size: size, Info: "New file: " + file_name + " | Size: " + bytefmt.ByteSize(size * bytefmt.KILOBYTE) } )
            total += size
          }
          var diff uint64
          if(size > old_size) {
            diff = size - old_size
            if(existed && diff > 10) {
              lines = append(lines, Line{ Size: diff, Info: "Increased in size: " + file_name + " | Size diff: " + bytefmt.ByteSize(diff * bytefmt.KILOBYTE) } )
              total += diff
            }
          }
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    sort.Sort(sort.Reverse(lines))

    for _, line := range lines {
      fmt.Println(line.Info)
    }

    fmt.Println()
    fmt.Println("Total: " + bytefmt.ByteSize(uint64(total * bytefmt.KILOBYTE)))
    //fmt.Println("%#v", lines)
}
