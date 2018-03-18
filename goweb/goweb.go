package main

import (
	"context"
	"fmt"
	"go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"log"
	// "math"
	"math/rand"
	"net/http"
	"time"
)

// define constants here

// ---- constants

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	//create new exporter
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: "thegrinch-project"}) // modify this for your GCP project
	if err != nil {
		log.Fatal(err)
	}
	// Export to Stackdriver Monitoring.
	view.RegisterExporter(exporter)
	view.SetReportingPeriod(1 * time.Second)

	// Set value to be exported to Stackdriver
	stackdriverMetric, err := stats.Int64("OpenCensus/testMetric", "metric value", r.URL.Path[1:])
	if err != nil {
		log.Fatalf("metric not created: %v", err)
	}

	// send to Stackdriver
	go func() {
		for {
			stats.Record(ctx, stackdriverMetric.M(1))
			<-time.After(time.Millisecond * time.Duration(1+rand.Intn(400)))
		}
	}()

	// Show output on web page
	fmt.Fprintf(w, "%s sent to Stackdriver!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
