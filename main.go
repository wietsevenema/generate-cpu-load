package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"
)

type loadavg struct {
	LoadAverage1     float64
	LoadAverage5     float64
	LoadAverage10    float64
	RunningProcesses int
	TotalProcesses   int
	LastProcessID    int
}

func getLoadAvg() (*loadavg, error) {
	self := new(loadavg)

	raw, _ := ioutil.ReadFile("/proc/loadavg")

	fmt.Sscanf(string(raw), "%f %f %f %d/%d %d",
		&self.LoadAverage1, &self.LoadAverage5, &self.LoadAverage10,
		&self.RunningProcesses, &self.TotalProcesses,
		&self.LastProcessID)

	return self, nil
}

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/generate-load", load)
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
	name, _ := os.Hostname()
	fmt.Fprintf(w, "<body>")
	fmt.Fprintf(w, "Hostname: %s <br/>", name)
	fmt.Fprintf(w, "This iframe will keep refreshing every second, generating CPU load:<br/>")
	fmt.Fprintf(w, "<iframe src='/generate-load'>")
	fmt.Fprintf(w, "</body>")
}

func load(w http.ResponseWriter, r *http.Request) {
	fact := int64(1100000)
	start := time.Now()
	var f big.Int
	f.MulRange(1, fact)
	end := time.Now()
	hostname, _ := os.Hostname()
	loadavg, _ := getLoadAvg()

	fmt.Fprintf(w, "<body><meta http-equiv='refresh' content='1'>")
	fmt.Fprintf(w, "Factorial %d! calculation finished in %s <br>", fact, end.Sub(start))
	fmt.Fprintf(w, "Hostname: %s<br>", hostname)
	fmt.Fprintf(w, "load averages: %.2f %.2f %.2f", loadavg.LoadAverage1, loadavg.LoadAverage5, loadavg.LoadAverage10)

	fmt.Fprintf(w, "</body>")
}
