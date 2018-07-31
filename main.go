package main
import (
  "fmt"
  "os"
  "os/exec"
  "strconv"
  "time"
  "flag"
)

// Makes Curl requests and returns the response time.
func goCurl(url string, c chan float64) {
    c <- curl(url)
}

// Makes Curl requests and returns the response time.
func curl(url string) float64 {

    // Variables
    var (
      cmdOut []byte
      err error
    )

    // Shout out to Nathan Leclaire for his example
    // URL: https://nathanleclaire.com/blog/2014/12/29/shelled-out-commands-in-golang/
    cmdName := "curl"
    cmdArgs := []string{"-w", "@curl-elapsed-time-format.txt", "-o", "/dev/null", "-s", url}
    if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
      fmt.Fprintln(os.Stderr, "There was an error running git rev-parse command: ", err)
      os.Exit(1)
    }
    sha := string(cmdOut)
    result, err := strconv.ParseFloat(sha, 64)
    return result
}

func printCurlInformation(url string, numCurls int, message string) {
    fmt.Println("Making curl calls against", url)
    fmt.Println("Making", numCurls, message)
}

func printCurlElapsedTime(curlNum int, elapsedTime float64) {
    fmt.Println("Curl", curlNum + 1, "came in with a call time of", elapsedTime)
}

func main() {
    // Declare and Parse Flags
    sumCurlsElapsedTime := 0.0
    numCurlsPtr := flag.Int("numCurls", 1, "Number of curls calls against the url")
    concurrentPtr :=  flag.Bool("concurrent", false, "Run concurrent curl commands")
    flag.Parse()

    // Get URL
    url := flag.Arg(0)

   start := time.Now()

    if(!*concurrentPtr) {
        printCurlInformation(url, *numCurlsPtr, "curl calls sequentially")
        for i:=0;i<*numCurlsPtr;i++ {
            elapsedTime := curl(url)
            printCurlElapsedTime(i, elapsedTime)
            sumCurlsElapsedTime += elapsedTime
        }
    } else {
        printCurlInformation(url, *numCurlsPtr, "curl calls concurrently")
        c := make(chan float64)

        // Start curls
        for i:=0;i<*numCurlsPtr;i++ {
            go goCurl(url, c)
        }

        // Wait for curls to finish
        for i:=0;i<*numCurlsPtr;i++ {
            elapsedTime := <- c
            printCurlElapsedTime(i, elapsedTime)
            sumCurlsElapsedTime += elapsedTime
        }
    }


    // Print average curl time
    fmt.Println("Average Curl Time:", sumCurlsElapsedTime/float64(*numCurlsPtr))

    // Calculate elapsed Tim
    elapsed := time.Since(start)
    fmt.Println("Elapsed time:", elapsed)
}
