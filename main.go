package main

import (
  "fmt"
  "os"
  "os/exec"
  "sync"
  "strconv"
)

var wg sync.WaitGroup // 1

func routine() {
  defer wg.Done() // 3
  fmt.Println("routine finished")
}

func main() {
  wg.Add(1) // 2
  for i:=0;i<1;i++ {
    curl()
  }
  wg.Wait() // 4

  fmt.Println("main finished")
}

func curl() float64 {
  var (
    cmdOut []byte
    err error
  )
  defer wg.Done()

  fmt.Println("Starting curl")
  cmdName := "curl"
  cmdArgs := []string{"localhost:8000/slow-endpoint"}
  cmdArgs = []string{"-w", "@curl-elapsed-time-format.txt", "-o", "/dev/null", "-s","localhost:8000/slow-endpoint"}
  if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
    fmt.Fprintln(os.Stderr, "There was an error running git rev-parse command: ", err)
    os.Exit(1)
  }
  sha := string(cmdOut)
  result, err := strconv.ParseFloat(sha, 64)
  return result
}
