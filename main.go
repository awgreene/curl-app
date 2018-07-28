package main
import (
  "fmt"
  "os"
  "os/exec"
  "strconv"
)

// Makes Curl requests and returns the response time.
func goCurl(c chan float64) {
    c <- curl()
}

// Makes Curl requests and returns the response time.
func curl() float64 {

    // Variables
    var (
      cmdOut []byte
      err error
    )

    // Shout out to Nathan Leclaire for his example
    // URL: https://nathanleclaire.com/blog/2014/12/29/shelled-out-commands-in-golang/
    cmdName := "curl"
    cmdArgs := []string{"-w", "@curl-elapsed-time-format.txt", "-o", "/dev/null", "-s","localhost:8000/slow-endpoint"}
    if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
      fmt.Fprintln(os.Stderr, "There was an error running git rev-parse command: ", err)
      os.Exit(1)
    }
    sha := string(cmdOut)
    result, err := strconv.ParseFloat(sha, 64)
    return result
}

func main() {
    // Declare Variables
    sum := 0.0
    numCurls,_ := strconv.Atoi(os.Args[1])
    serial,_ := strconv.ParseBool(os.Args[2])

    if(serial) {
        fmt.Println("Making", numCurls, "curl calls sequentially")
        for i:=0;i<numCurls;i++ {
            temp := curl()
            fmt.Println("Curl", i + 1, "came in with a call time of", temp)
            sum += temp
        }
    } else {
        // go routines
        fmt.Println("Making", numCurls, "curl calls in parallel")
        c := make(chan float64)

        // Start curls
        for i:=0;i<numCurls;i++ {
            go goCurl(c)
        }

        // Wait for curls to finish
        for i:=0;i<numCurls;i++ {
            fmt.Printf("Waiting")
            temp := <- c
            fmt.Println("Curl", i + 1, "came in with a call time of", temp)
            sum += temp
        }
    }

    // Print average curl time
    fmt.Println("Average Curl Time:", sum/float64(numCurls))
}
