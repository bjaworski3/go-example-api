package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

// Mutex used for nameMap modification
var mutex = &sync.Mutex{}

// Map used to store name counts from /hello/:name calls
var nameMap = make(map[string]int)

// nameHandler takes an http GET request on /hello/:name and will return a
// message based on the name given. nameHandler will also increment a count of
// the number of times each name has been called. If a bad HTTP request is given
// a 405 will be returned.
func nameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { // Verify the method is valid GET
		name := r.URL.Path[len("/hello/"):]
		fmt.Fprintf(w, "Hello, %s!", name)
		if _, ok := nameMap[name]; ok {
			mutex.Lock()
			nameMap[name]++
			mutex.Unlock()
		} else {
			mutex.Lock()
			nameMap[name] = 1
			mutex.Unlock()
		}
	} else { // Not a supported method
		http.Error(w, "Invalid request method.", 405)
	}
}

// healthHandler takes an HTTP GET request on /health and will return a JSON
// formatted string that contains system information including virtual Memory
// swap Memory, CPU usage, and load. If a bad HTTP request is given a 405 will
// be returned.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { // Verify the method is valid GET
		result := make(map[string]interface{}, 0)

		// Get Virtual Memory Stats
		vMem, _ := mem.VirtualMemory()
		// Get Swap Memory Stats
		sMem, _ := mem.SwapMemory()
		// Get CPU Stats
		cpu, _ := cpu.Times(true)
		// Get load information
		load, _ := load.Avg()

		//Put all stats into the map
		result["virtual_memory_info"] = vMem
		result["swap_memory_info"] = sMem
		result["cpu_info"] = cpu
		result["load"] = load

		// Format the map and then print to page
		b, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Fprintln(w, string(b))

	} else { // Not a supported method
		http.Error(w, "Invalid request method.", 405)
	}
}

// countHandler takes an HTTP GET or DELETE request on /counts. The GET request
// will return a JSON formatted string with the counts of each name added with
// the /hello/:name url. The DELETE request will remove all data. If a bad HTTP
// request is given a 405 will be returned.
func countHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { // Verify the method is valid GET
		// Create a map to use as the JSON output
		output := make([]map[string]interface{}, 0)
		// Add all the name and count information in correct format
		for name, count := range nameMap {
			output = append(output, map[string]interface{}{"name": name, "count": count})
		}
		// Convert to json
		jsonString, err := json.MarshalIndent(output, "", "  ")
		// TODO Add a unit test to cover this error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonString)
		// Verify the method is valid DELETE
	} else if r.Method == "DELETE" {
		mutex.Lock()
		nameMap = make(map[string]int)
		mutex.Unlock()
		fmt.Fprintln(w, "Count data has been removed.")
	} else { // Not a supported method
		http.Error(w, "Invalid request method.", 405)
	}
}

// handleRequests sets up all of the possible urls that can be accessed in the
// web-api application. It then serves them on the port 8080. Any other urls
// will return a 404.
func handleRequests() {
	http.HandleFunc("/hello/", nameHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/counts", countHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}
