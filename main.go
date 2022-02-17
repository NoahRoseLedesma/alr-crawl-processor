package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"

	"github.com/korovkin/limiter"
	"github.com/schollz/progressbar/v3"
	"github.com/slyrz/warc"
)

// Regular expression to match aria landmarks in HTML
var alr_regex, _ = regexp.Compile("role=\"(.+)\"")

// Landmark counters
// Each is stored as a global variable to allow for atomic access
var n_banner uint64 = 0
var n_navigation uint64 = 0
var n_main uint64 = 0
var n_complementary uint64 = 0
var n_contentinfo uint64 = 0
var n_search uint64 = 0
var n_form uint64 = 0
var n_application uint64 = 0
var n_other uint64 = 0

// Count the total number of documents with landmarks
var num_documents_with_landmarks uint64 = 0

// Determine the WARC output file path for a WARC input file
func getOutputFilePath(in_path string) string {
	// Remove .gz from the end of the path
	filename := strings.TrimSuffix(in_path, ".gz")
	// Get the filename
	filename = filepath.Base(filename)
	// Return the output file path
	return "output/" + filename
}

// Write records corresponding to HTML files with landmarks to a WARC file
func recordWriter(records chan *warc.Record, path string) {
	// Open the WARC file and create a writer interface
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	writer := warc.NewWriter(f)

	// Write the records to the WARC file
	for record := range records {
		writer.WriteRecord(record)
	}
}

func processDocument(record *warc.Record, out_records chan *warc.Record) bool {
	buf := new(strings.Builder)
	io.Copy(buf, record.Content)
	content := buf.String()
	match := alr_regex.FindAllStringSubmatch(content, -1)

	has_landmark := false
	for _, m := range match {
		switch m[1] {
		case "banner":
			atomic.AddUint64(&n_banner, 1)
			has_landmark = true
		case "navigation":
			atomic.AddUint64(&n_navigation, 1)
			has_landmark = true
		case "main":
			atomic.AddUint64(&n_main, 1)
			has_landmark = true
		case "complementary":
			atomic.AddUint64(&n_complementary, 1)
			has_landmark = true
		case "contentinfo":
			atomic.AddUint64(&n_contentinfo, 1)
			has_landmark = true
		case "search":
			atomic.AddUint64(&n_search, 1)
			has_landmark = true
		case "form":
			atomic.AddUint64(&n_form, 1)
			has_landmark = true
		case "application":
			atomic.AddUint64(&n_application, 1)
			has_landmark = true
		default:
			atomic.AddUint64(&n_other, 1)
		}
	}

	if has_landmark {
		atomic.AddUint64(&num_documents_with_landmarks, 1)
		num_documents_with_landmarks++

		// Save the record to a WARC file
		out_records <- record
	}

	return has_landmark
}

// Process a single WARC file and return the number of HTML documents.
func processWarc(path string) int {
	// Open the warc file and create a reader interface
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	reader, err := warc.NewReader(f)
	if err != nil {
		panic(err)
	}

	// Create a channel to receive records
	out_records := make(chan *warc.Record)
	// Create a worker to handle writing records to a WARC file
	go recordWriter(out_records, getOutputFilePath(path))

	// Count the total number of HTML documents in the archive
	num_documents := 0

	// Create a waitgroup
	limit := limiter.NewConcurrencyLimiter(150)
	bar := progressbar.Default(-1)

	for {
		record, err := reader.ReadRecord()
		if err != nil {
			break
		}

		// Check the content type
		if record.Header["content-type"] == "application/http; msgtype=response" && record.Header["warc-identified-payload-type"] == "text/html" {
			num_documents++
			limit.Execute(func() {
				processDocument(record, out_records)
				bar.Add(1)
			})
		}
	}

	limit.Wait()
	bar.Finish()

	return num_documents
}

// Main function
func main() {
	// Count the total number of HTML documents
	num_documents := 0

	// Create the output directory
	os.Mkdir("output", os.ModePerm)

	// Discover all of the gunzip files in the data directory
	filepath.WalkDir("data", func(path string, dir os.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}

		if filepath.Ext(path) == ".gz" {
			fmt.Println("Processing", path)
			num_documents += processWarc(path)
		}
		return nil
	})

	fmt.Println("Total number of documents:", num_documents)
	fmt.Println("Total number of documents with landmarks:", num_documents_with_landmarks)
	// Print the landmark counts
	fmt.Println("Landmark counts:")
	fmt.Println()
	fmt.Println("Banner:", n_banner)
	fmt.Println("Navigation:", n_navigation)
	fmt.Println("Main:", n_main)
	fmt.Println("Complementary:", n_complementary)
	fmt.Println("ContentInfo:", n_contentinfo)
	fmt.Println("Search:", n_search)
	fmt.Println("Form:", n_form)
	fmt.Println("Application:", n_application)
	fmt.Println("Other:", n_other)
}
