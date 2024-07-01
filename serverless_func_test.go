package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
)

// Mock implementations of the functions

func mockLoadDataToStruct(path string) ([]LocationData, error) {

	// Handle valid paths with different data sets
	switch path {
	case "valid/Unsuccessful":
		return []LocationData{{Name: "ValidCity1"}, {Name: "invalid"}}, nil
	case "valid/path":
		return []LocationData{{Name: "ValidCity1"}, {Name: "ValidCity2"}}, nil
	default:
		return nil, fmt.Errorf("failed to load data")
	}
}

func mockGetUniqueKey(data LocationData) Key {
	return Key{
		City:    data.Name,
		Country: data.Country,
		Geo:     data.Geo,
	}
}

func mockHardValidate(city1, city2 LocationData) bool {
	return city1.Name == city2.Name
}

func mockLoadDataFuncSuccess(path string) ([]LocationData, error) {
	return []LocationData{city1}, nil
}

func mockLoadDataFuncFailure(path string) ([]LocationData, error) {
	return mockLoadDataToStruct("")
}

func mockLoadDataPartialFailure(path string) ([]LocationData, error) {
	return []LocationData{{Name: "invalid"}}, nil
}

func mockGetAllFiles(path string) ([]string, error) {
	return []string{"valid/Unsuccessful", "valid/path", "invalid/path"}, nil
}

var city1 = LocationData{Name: "ValidCity1"}
var city2 = LocationData{Name: "ValidCity2"}
var key1 = GetUniqueKey(city1)
var key2 = GetUniqueKey(city2)

var mockAuthenticCities = map[Key]LocationData{
	key1: city1,
	key2: city2,
}

func TestReadData(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up the temp file after the test

	// Write some data to the temporary file
	content := []byte("Hello, Gopher!")
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Use the readData function to read the data back
	data, err := readData(tmpfile.Name())
	if err != nil {
		t.Fatalf("readData returned an error: %v", err)
	}

	// Verify the content is what we expect
	if string(data) != string(content) {
		t.Errorf("readData returned unexpected content: got %v, want %v", string(data), string(content))
	}
}

func TestProcessFile(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	var mu sync.Mutex
	authenticCities = mockAuthenticCities

	// Test cases
	tests := []struct {
		name                 string
		tmpPath              string
		expectedUnprocessed  int
		expectedSuccessful   int
		expectedUnsuccessful int
	}{
		{"Valid file", "valid/path", 0, 2, 0},
		{"Unsuccessful", "valid/Unsuccessful", 0, 1, 1},
		{"Invalid file", "invalid/path", 1, 0, 0},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			var unprocessableFiles []string
			var successfullyValidated []LocationData
			var unsuccessfullyValidated []LocationData
			authenticCities = mockAuthenticCities
			mockUtils := HelperUtils{mockLoadDataToStruct, mockGetUniqueKey, mockHardValidate, nil}
			wg.Add(1)
			go processFile(
				tt.tmpPath, &wg, &mu,
				&unprocessableFiles,
				&successfullyValidated,
				&unsuccessfullyValidated,
				mockUtils,
			)
			wg.Wait()

			if len(unprocessableFiles) != tt.expectedUnprocessed {
				t.Errorf("Expected %d unprocessable files, got %d", tt.expectedUnprocessed, len(unprocessableFiles))
			}
			if len(successfullyValidated) != tt.expectedSuccessful {
				t.Errorf("Expected %d successfully validated cities, got %d", tt.expectedSuccessful, len(successfullyValidated))
			}
			if len(unsuccessfullyValidated) != tt.expectedUnsuccessful {
				t.Errorf("Expected %d unsuccessfully validated cities, got %d", tt.expectedUnsuccessful, len(unsuccessfullyValidated))
			}
		})
	}
}

func TestLoadDataToStruct_InvalidJSON(t *testing.T) {
	// Prepare invalid JSON data
	invalidJSON := []byte(`{ "invalid_data": true }`)

	// Create a temporary file with the invalid JSON data
	f, err := os.CreateTemp("", "invalid_data.json")
	if err != nil {
		t.Errorf("Failed to create temporary file: %v", err)
	}
	defer func() {
		_ = os.Remove(f.Name())
	}()
	_, err = f.Write(invalidJSON)
	if err != nil {
		t.Errorf("Failed to write data to temporary file: %v", err)
	}
	defer f.Close()

	// Call the function with the temporary file path
	_, err = loadDataToStruct(f.Name())

	// Assert error on unmarshalling invalid JSON
	if err == nil {
		t.Errorf("Expected error when unmarshalling invalid JSON")
	}
}

func TestLoadDataToStruct_ReadFileError(t *testing.T) {
	// Mock a non-existent file path
	filePath := "non-existent-file.json"

	// Call the function with the non-existent file path
	_, err := loadDataToStruct(filePath)

	// Assert error on reading a non-existent file
	if err == nil {
		t.Errorf("Expected error when reading non-existent file")
	}
}

func TestLoadAuthenticCities(t *testing.T) {
	t.Parallel()

	// Reset the global authenticCities map before each test
	authenticCities = nil

	// Test cases
	tests := []struct {
		name        string
		filepath    string
		expectError bool
		expectedMap map[Key]LocationData
	}{
		{
			name:        "Valid file with duplicates",
			filepath:    "valid/path",
			expectError: false,
			expectedMap: mockAuthenticCities,
		},
		{
			name:        "Invalid file path",
			filepath:    "invalid/path",
			expectError: true,
			expectedMap: nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			err := loadAuthenticCities(tt.filepath, mockLoadDataToStruct, mockGetUniqueKey)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}

			if !tt.expectError {
				if len(authenticCities) != len(tt.expectedMap) {
					t.Errorf("Expected map length: %d, got: %d", len(tt.expectedMap), len(authenticCities))
				}

				for key, expectedCity := range tt.expectedMap {
					if city, ok := authenticCities[key]; !ok || city != expectedCity {
						t.Errorf("Expected city: %+v for key: %s, got: %+v", expectedCity, key, city)
					}
				}
			}
		})
	}
}

// Mock the readData function
var mockReadData func(filepath string) ([]byte, error)

func TestLoadDataToStruct_ValidJSON(t *testing.T) {
	// Prepare valid JSON data
	jsonData := []byte(`[
       {
           "City": "New York",
           "Country": "USA",
           "Latitude": "-74.0059",
           "Longitude": "40.7128"
       },
       {
           "City": "London",
           "Country": "UK",
           "Latitude": "51.505",
           "Longitude": "-0.09"
       }
   ]`)

	// Create a temporary file with the valid JSON data
	f, err := os.CreateTemp("", "location_data.json")
	if err != nil {
		t.Errorf("Failed to create temporary file: %v", err)
	}
	defer func() {
		_ = os.Remove(f.Name())
	}()
	_, err = f.Write(jsonData)
	if err != nil {
		t.Errorf("Failed to write data to temporary file: %v", err)
	}
	defer f.Close()

	// Call the function with the temporary file path
	cities, err := loadDataToStruct(f.Name())

	// Assert successful execution
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Assert expected number of cities loaded
	if len(cities) != 2 {
		t.Errorf("Expected 2 cities, got %d", len(cities))
	}

	// Assert first city data
	if cities[0].Name != "New York" || cities[0].Country != "USA" {
		t.Errorf("Unexpected city data for New York")
	}

	// Assert second city data
	if cities[1].Name != "London" || cities[1].Country != "UK" {
		t.Errorf("Unexpected city data for London")
	}
}

func TestLoadDataToStruct(t *testing.T) {
	t.Parallel()

	// Define test cases
	tests := []struct {
		name      string
		filepath  string
		mockData  []byte
		mockError error
		want      []LocationData
		wantErr   bool
	}{
		{
			name:      "file read error",
			filepath:  "invalid.json",
			mockData:  nil,
			mockError: errors.New("file not found"),
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "invalid JSON",
			filepath:  "invalid.json",
			mockData:  []byte(`invalid json`),
			mockError: nil,
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			// Set up the mock function
			mockReadData = func(filepath string) ([]byte, error) {
				if filepath == tt.filepath {
					return tt.mockData, tt.mockError
				}
				return nil, errors.New("unexpected filepath")
			}

			// Call the function under test
			got, err := loadDataToStruct(tt.filepath)

			if (err != nil) != tt.wantErr {
				t.Errorf("loadDataToStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadDataToStruct() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessFileUsingChannels_Success(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	unprocessableFiles := make(chan string, 1)
	successfullyValidated := make(chan LocationData, 1)
	unsuccessfullyValidated := make(chan LocationData, 1)

	authenticCities = map[Key]LocationData{key1: city1}
	mockUtils := HelperUtils{mockLoadDataFuncSuccess, mockGetUniqueKey, mockHardValidate, nil}

	go processFileUsingChannels(
		"mock/path",
		&wg,
		unprocessableFiles,
		successfullyValidated,
		unsuccessfullyValidated,
		mockUtils,
	)

	wg.Wait()
	close(unprocessableFiles)
	close(successfullyValidated)
	close(unsuccessfullyValidated)

	select {
	case data := <-successfullyValidated:
		if data.Name != city1.Name {
			t.Errorf("Expected %s successfully validated , got %s", city1.Name, data.Name)
		}
	default:
		t.Error("Expected successfully validated data but got none")
	}
}

func TestProcessFileUsingChannels_Failure(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	unprocessableFiles := make(chan string, 1)
	successfullyValidated := make(chan LocationData, 1)
	unsuccessfullyValidated := make(chan LocationData, 1)
	mockUtils := HelperUtils{mockLoadDataFuncFailure, mockGetUniqueKey, mockHardValidate, nil}

	go processFileUsingChannels(
		"mock/path",
		&wg,
		unprocessableFiles,
		successfullyValidated,
		unsuccessfullyValidated,
		mockUtils,
	)

	wg.Wait()
	close(unprocessableFiles)
	close(successfullyValidated)
	close(unsuccessfullyValidated)

	select {
	case file := <-unprocessableFiles:
		if file != "mock/path" {
			t.Errorf("Expected unprocessable file path to be 'mock/path' but got %v", file)
		}
	default:
		t.Error("Expected unprocessable file but got none")
	}
}

func TestProcessFileUsingChannels_PartialFailure(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	unprocessableFiles := make(chan string, 1)
	successfullyValidated := make(chan LocationData, 1)
	unsuccessfullyValidated := make(chan LocationData, 1)
	mockUtils := HelperUtils{mockLoadDataPartialFailure, mockGetUniqueKey, mockHardValidate, nil}

	go processFileUsingChannels(
		"valid/Unsuccessful",
		&wg,
		unprocessableFiles,
		successfullyValidated,
		unsuccessfullyValidated,
		mockUtils,
	)

	wg.Wait()
	close(unprocessableFiles)
	close(successfullyValidated)
	close(unsuccessfullyValidated)

	select {
	case _ = <-unsuccessfullyValidated:
	default:
		t.Error("Expected unprocessable file but got none")
	}
}

func TestGetAllFiles_EmptyDir(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "getAllFiles_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // Clean up temporary directory

	// Call the function with the empty directory
	files, err := getAllFiles(tmpDir)

	// Check for expected results (no files, no error)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("Expected no files, got %d", len(files))
	}
}

func TestGetAllFiles_WithJsonFiles(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "getAllFiles_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // Clean up temporary directory

	// Create some sample JSON files
	for i := 0; i < 3; i++ {
		fileName := filepath.Join(tmpDir, "file"+string(rune(i+65))+".json") // A-C.json
		_, err := os.Create(fileName)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Call the function with the directory containing JSON files
	files, err := getAllFiles(tmpDir)

	// Check for expected results (list of JSON files, no error)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(files))
	}
	for _, file := range files {
		if filepath.Ext(file) != ".json" {
			t.Errorf("Expected only JSON files, found %s", file)
		}
	}
}

func TestGetAllFiles_ErrorReadingDir(t *testing.T) {
	// Simulate an error by providing an invalid directory path
	invalidPath := "/invalid/path"

	// Call the function with the invalid path
	files, err := getAllFiles(invalidPath)

	// Check for expected results (nil slice, error)
	if err == nil {
		t.Error("Expected error when reading directory")
	}
	if files != nil {
		t.Error("Expected nil slice on error")
	}
}
