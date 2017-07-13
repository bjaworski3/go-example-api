package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNameHandlerPrint(t *testing.T) {
	req, err := http.NewRequest("GET", "/hello/test name", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(nameHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Hello, test name!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Check that the map updates properly when getting a single name
	expected_map := map[string]int{"test name": 1}
	eq := reflect.DeepEqual(nameMap, expected_map)
	if !eq {
		t.Errorf("handler returned unexpected body: got %v want %v",
			nameMap, expected_map)
	}

	// Reset nameMap for subsequent tests
	nameMap = make(map[string]int)
}

func TestNameHandlerMap(t *testing.T) {
	req, err := http.NewRequest("GET", "/hello/test name", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(nameHandler)

	handler.ServeHTTP(rr, req)
	handler.ServeHTTP(rr, req)
	handler.ServeHTTP(rr, req)

	req, err = http.NewRequest("GET", "/hello/test name2", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(nameHandler)

	handler.ServeHTTP(rr, req)

	// Check that the map updates properly when getting a single name
	expected_map := map[string]int{"test name": 3, "test name2": 1}
	eq := reflect.DeepEqual(nameMap, expected_map)
	if !eq {
		t.Errorf("handler returned unexpected body: got %v want %v",
			nameMap, expected_map)
	}

	// Reset nameMap for subsequent tests
	nameMap = make(map[string]int)
}

func TestNameHandlerBadReq(t *testing.T) {
	requestTypes := []string{"POST", "PUT", "DELETE"}
	for _, requestType := range requestTypes {
		req, err := http.NewRequest(requestType, "/hello/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(nameHandler)

		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusMethodNotAllowed)
		}

		// Check the output is what we expect
		expected := "Invalid request method.\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}
}

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// TODO Update this expected when health check is ready
	body := []byte(rr.Body.String())
	var raw map[string]interface{}
	err = json.Unmarshal(body, &raw)
	if err != nil {
		t.Errorf("handler returned unexpected body, not JSON:%v",
			rr.Body.String())
	}
}

func TestHealthHandlerBadReq(t *testing.T) {
	requestTypes := []string{"POST", "PUT", "DELETE"}
	for _, requestType := range requestTypes {
		req, err := http.NewRequest(requestType, "/health-check", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(healthHandler)

		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusMethodNotAllowed)
		}

		// Check the output is what we expect
		expected := "Invalid request method.\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}
}

func TestCountHandlerGet(t *testing.T) {
	// Add some data to the Name count slice map
	nameMap = map[string]int{"John Smith": 6, "Jane Doe": 25}

	req, err := http.NewRequest("GET", "/counts", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(countHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body := []byte(rr.Body.String())
	var data []map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		t.Errorf("handler returned unexpected body, not JSON:\n%v.\n Error:%v",
			rr.Body.String(), err)
	}

	for _, element := range data {
		if element["name"] == "John Smith" {
			if element["count"] != float64(6) {
				t.Errorf("Count for John Smith should be 6:\n%v. \n Count: %v",
					rr.Body.String(), element["count"])
			}
		} else if element["name"] == "Jane Doe" {
			if element["count"] != float64(25) {
				t.Errorf("Count for Jane Doe should be 25:\n%v. \n Count: %v",
					rr.Body.String(), element["count"])
			}
		} else {
			t.Errorf("Name should be either John Smith or Jane Doe:\n%v.",
				rr.Body.String())
		}
	}

	// Reset nameMap for subsequent tests
	nameMap = make(map[string]int)
}

func TestCountHandlerDelete(t *testing.T) {
	// Add some data to the Name count slice map
	nameMap = map[string]int{"John Smith": 6, "Jane Doe": 25}

	req, err := http.NewRequest("DELETE", "/counts", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(countHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Expects output indicating data has been removed
	expected := "Count data has been removed.\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Reset nameMap for subsequent tests
	nameMap = make(map[string]int)
}

func TestCountHandlerBadReq(t *testing.T) {
	requestTypes := []string{"POST", "PUT"}
	for _, requestType := range requestTypes {
		req, err := http.NewRequest(requestType, "/counts", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(countHandler)

		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusMethodNotAllowed)
		}

		// Check the output is what we expect
		expected := "Invalid request method.\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}
}
