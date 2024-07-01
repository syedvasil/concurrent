package main

import "testing"

func Test_processFiles(t *testing.T) {
	authenticCities = map[Key]LocationData{key1: city1}
	type args struct {
		tmpFolder string
		utils     HelperUtils
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "successful",
			args: args{
				tmpFolder: "valid/path",
				utils:     HelperUtils{loadDataFunc: mockLoadDataToStruct, getUniqueKeyFunc: mockGetUniqueKey, hardValidateFunc: mockHardValidate, getAllFiles: mockGetAllFiles},
			},
		},
		{
			name: "no files present",
			args: args{
				tmpFolder: "invalid/path",
				utils:     HelperUtils{loadDataFunc: mockLoadDataToStruct, getUniqueKeyFunc: mockGetUniqueKey, hardValidateFunc: mockHardValidate, getAllFiles: getAllFiles},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			success, unsuccessful, invalid := ProcessFiles(tt.args.tmpFolder, tt.args.utils)

			valid := tt.args.tmpFolder == "valid/path"
			if valid && len(success) == 0 {
				t.Errorf("ProcessFiles() return empty for sucess")
			}
			if valid && len(unsuccessful) == 0 {
				t.Errorf("ProcessFiles() return empty for unsuccessful")
			}
			if valid && len(invalid) == 0 {
				t.Errorf("ProcessFiles() return empty for invalid")
			}
		})
	}
}

func Test_processFilesC(t *testing.T) {
	authenticCities = map[Key]LocationData{key1: city1}
	type args struct {
		tmpFolder string
		utils     HelperUtils
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "successful",
			args: args{
				tmpFolder: "valid/path",
				utils:     HelperUtils{loadDataFunc: mockLoadDataToStruct, getUniqueKeyFunc: mockGetUniqueKey, hardValidateFunc: mockHardValidate, getAllFiles: mockGetAllFiles},
			},
		},
		{
			name: "no files present",
			args: args{
				tmpFolder: "invalid/path",
				utils:     HelperUtils{loadDataFunc: mockLoadDataToStruct, getUniqueKeyFunc: mockGetUniqueKey, hardValidateFunc: mockHardValidate, getAllFiles: getAllFiles},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			success, unsuccessful, invalid := ProcessFilesWithoutMutex(tt.args.tmpFolder, tt.args.utils)
			valid := tt.args.tmpFolder == "valid/path"
			if valid && len(success) == 0 {
				t.Errorf("ProcessFiles() return empty for sucess")
			}
			if valid && len(unsuccessful) == 0 {
				t.Errorf("ProcessFiles() return empty for unsuccessful")
			}
			if valid && len(invalid) == 0 {
				t.Errorf("ProcessFiles() return empty for invalid")
			}

		})
	}
}
