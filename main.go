package main

import (
	"fmt"
)

func main() {
	err := loadAuthenticCities("cities.json", loadDataToStruct, GetUniqueKey)
	if err != nil {
		fmt.Println("Error loading authentic cities:", err)
		return
	}

	validated, inValid, unprocessable := ProcessFilesWithoutMutex("tmp",
		HelperUtils{loadDataToStruct, GetUniqueKey, hardCheck, getAllFiles})

	fmt.Println("Successfully Validated Elements:", validated)
	fmt.Println("Unsuccessfully Validated Elements:", inValid)
	fmt.Println("Unprocessable Files:", unprocessable)
	fmt.Println("Successfully Validated Elements:", len(validated))
	fmt.Println("Unsuccessfully Validated Elements:", len(inValid))
	fmt.Println("Unprocessable Files:", len(unprocessable))
}
