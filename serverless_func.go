package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// authenticCities Cache to check and validate the input data
var authenticCities map[Key]LocationData

// readData reads the file and returns the contents.
// A successful call returns err == nil, not err == EOF.
func readData(filepath string) ([]byte, error) {
	return os.ReadFile(filepath)
}

func loadDataToStruct(filepath string) ([]LocationData, error) {
	data, err := readData(filepath)
	if err != nil {
		return nil, err
	}

	var cities []LocationData
	err = json.Unmarshal(data, &cities)
	if err != nil {
		return nil, err
	}

	return cities, nil
}

func loadAuthenticCities(
	filepath string,
	loadDataFunc func(string) ([]LocationData, error),
	getUniqueKeyFunc func(data LocationData) Key,
) error {
	cities, err := loadDataFunc(filepath)
	if err != nil {
		return err
	}

	authenticCities = make(map[Key]LocationData, len(cities))

	for _, city := range cities {
		key := getUniqueKeyFunc(city)

		if _, ok := authenticCities[key]; ok {
			fmt.Printf("Got a duplicate with details %+v \n", city)
		} else {
			authenticCities[key] = city
		}
	}

	return nil
}

func getAllFiles(tmpFolder string) ([]string, error) {
	var allFiles []string
	files, err := os.ReadDir(tmpFolder)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			allFiles = append(allFiles, filepath.Join(tmpFolder, file.Name()))
		}
	}
	return allFiles, nil
}

type HelperUtils struct {
	loadDataFunc     func(string) ([]LocationData, error)
	getUniqueKeyFunc func(data LocationData) Key
	hardValidateFunc func(LocationData, LocationData) bool
	getAllFiles      func(string) ([]string, error)
}

func ProcessFiles(
	tmpFolder string,
	helpers HelperUtils,
) ([]LocationData, []LocationData, []string) {
	var successfullyValidated, unsuccessfullyValidated []LocationData
	var unprocessableFiles []string

	allFiles, err := helpers.getAllFiles(tmpFolder)
	if err != nil {
		fmt.Println("Error reading tmp folder:", err)
		return nil, nil, nil
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, fileP := range allFiles {
		wg.Add(1)
		go processFile(fileP, &wg, &mu, &unprocessableFiles, &successfullyValidated, &unsuccessfullyValidated, helpers)
	}

	wg.Wait()

	return successfullyValidated, unsuccessfullyValidated, unprocessableFiles
}

func processFile(
	tmpPath string,
	wg *sync.WaitGroup,
	mu *sync.Mutex,
	unprocessableFiles *[]string,
	successfullyValidated *[]LocationData,
	unsuccessfullyValidated *[]LocationData,
	helper HelperUtils,
) {
	defer wg.Done()

	cities, err := helper.loadDataFunc(tmpPath)
	if err != nil {
		mu.Lock()
		*unprocessableFiles = append(*unprocessableFiles, tmpPath)
		mu.Unlock()
		return
	}

	for _, element := range cities {
		verifyData, ok := authenticCities[helper.getUniqueKeyFunc(element)]
		if ok && helper.hardValidateFunc(verifyData, element) {
			mu.Lock()
			*successfullyValidated = append(*successfullyValidated, element)
			mu.Unlock()
		} else {
			mu.Lock()
			*unsuccessfullyValidated = append(*unsuccessfullyValidated, element)
			mu.Unlock()
		}
	}
}

func processFileUsingChannels(
	tmpPath string,
	wg *sync.WaitGroup,
	unprocessableFiles chan<- string,
	successfullyValidated chan<- LocationData,
	unsuccessfullyValidated chan<- LocationData,
	utils HelperUtils,
) {
	defer wg.Done()

	cities, err := utils.loadDataFunc(tmpPath)
	if err != nil {
		unprocessableFiles <- tmpPath
		return
	}

	for _, element := range cities {
		verifyData, ok := authenticCities[utils.getUniqueKeyFunc(element)]
		if ok && utils.hardValidateFunc(element, verifyData) {
			successfullyValidated <- element
		} else {
			unsuccessfullyValidated <- element
		}
	}
}

func ProcessFilesWithoutMutex(tmpFolder string, utils HelperUtils) ([]LocationData, []LocationData, []string) {
	var successfullyValidated, unsuccessfullyValidated []LocationData
	var unprocessableFiles []string

	allFiles, err := utils.getAllFiles(tmpFolder)
	if err != nil {
		fmt.Println("Error reading tmp folder:", err)
		return nil, nil, nil
	}

	var wg sync.WaitGroup

	unprocessableChan := make(chan string, len(allFiles))
	authentic := make(chan LocationData, len(allFiles)*10) // Assuming each file has up to 10 cities
	inauthentic := make(chan LocationData, len(allFiles)*10)

	for _, fileP := range allFiles {
		wg.Add(1)
		go processFileUsingChannels(fileP, &wg, unprocessableChan, authentic, inauthentic, utils)
	}

	// Close channels when all goroutines are done
	go func() {
		wg.Wait()
		close(unprocessableChan)
		close(authentic)
		close(inauthentic)
	}()

	for errMsg := range unprocessableChan {
		unprocessableFiles = append(unprocessableFiles, errMsg)
	}

	for data := range authentic {
		successfullyValidated = append(successfullyValidated, data)
	}

	for data := range inauthentic {
		unsuccessfullyValidated = append(unsuccessfullyValidated, data)
	}

	return successfullyValidated, unsuccessfullyValidated, unprocessableFiles
}
