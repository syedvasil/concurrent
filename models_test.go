package main

import (
	"reflect"
	"testing"
)

func TestGetUniqueKey(t *testing.T) {
	validLoc := LocationData{
		Name:    "New York",
		Geo:     "USA",
		Country: "SomeP",
	}
	ValidKey := Key{
		City:    "New York",
		Country: "SomeP",
		Geo:     "USA",
	}
	invalidLoc := LocationData{
		Name:    "New York",
		Country: "USA",
	}
	invalidKey := Key{
		City:    "New York",
		Country: "USA",
	}
	type args struct {
		city LocationData
	}
	tests := []struct {
		name string
		args args
		want Key
	}{
		{
			name: "Unique key with all fields",
			args: args{
				city: validLoc,
			},
			want: ValidKey,
		},
		// Test case with partially empty name
		{
			name: "Unique key with empty name",
			args: args{
				city: invalidLoc,
			},
			want: invalidKey, // Key should be empty for invalid data
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUniqueKey(tt.args.city); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUniqueKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocationData_basicValidate(t *testing.T) {
	type fields struct {
		Latitude     string
		Longitude    string
		Geo          string
		Name         string
		ProvinceIcon string
		Province     string
		CountryIcon  string
		Country      string
	}
	type args struct {
		verifyData LocationData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// Test case with all fields filled
		{
			name: "Basic validation with all fields",
			fields: fields{
				Name:     "New York",
				Province: "New York",
				Country:  "USA",
			},
			args: args{verifyData: LocationData{
				Name:     "New York",
				Province: "New York",
				Country:  "USA",
			}},
			want: true,
		},
		// Test case with empty name
		{
			name: "Basic validation with empty name",
			fields: fields{
				Province: "New York",
				Country:  "USA",
			},
			args: args{verifyData: LocationData{}},
			want: false,
		},
		// Test case with empty country
		{
			name: "Basic validation with empty country",
			fields: fields{
				Name:     "New York",
				Province: "New York",
			},
			args: args{verifyData: LocationData{}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inp := LocationData{
				Latitude:     tt.fields.Latitude,
				Longitude:    tt.fields.Longitude,
				Geo:          tt.fields.Geo,
				Name:         tt.fields.Name,
				ProvinceIcon: tt.fields.ProvinceIcon,
				Province:     tt.fields.Province,
				CountryIcon:  tt.fields.CountryIcon,
				Country:      tt.fields.Country,
			}
			if got := inp.basicValidate(tt.args.verifyData); got != tt.want {
				t.Errorf("basicValidate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocationData_hardValidate(t *testing.T) {
	type fields struct {
		Latitude     string
		Longitude    string
		Geo          string
		Name         string
		ProvinceIcon string
		Province     string
		CountryIcon  string
		Country      string
	}
	type args struct {
		verifyData LocationData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// Test case with all fields valid
		{
			name: "Hard validation with all fields valid",
			fields: fields{
				Name:      "New York",
				Province:  "New York",
				Country:   "USA",
				Latitude:  "-74.0059", // Example valid latitude
				Longitude: "40.7128",  // Example valid longitude
			},
			args: args{verifyData: LocationData{
				Name:      "New York",
				Province:  "New York",
				Country:   "USA",
				Latitude:  "-74.0059", // Example valid latitude
				Longitude: "40.7128",  // Example valid longitude
			}},
			want: true,
		},
		// Test case with invalid latitude format
		{
			name: "Hard validation with invalid latitude format",
			fields: fields{
				Name:      "New York",
				Province:  "New York",
				Country:   "USA",
				Latitude:  "invalid_latitude",
				Longitude: "40.7128",
			},
			args: args{verifyData: LocationData{}},
			want: false,
		},
		// Test case with invalid longitude format
		{
			name: "Hard validation with invalid longitude format",
			fields: fields{
				Name:      "New York",
				Province:  "New York",
				Country:   "USA",
				Latitude:  "-74.0059",
				Longitude: "invalid_longitude",
			},
			args: args{verifyData: LocationData{}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inp := LocationData{
				Latitude:     tt.fields.Latitude,
				Longitude:    tt.fields.Longitude,
				Geo:          tt.fields.Geo,
				Name:         tt.fields.Name,
				ProvinceIcon: tt.fields.ProvinceIcon,
				Province:     tt.fields.Province,
				CountryIcon:  tt.fields.CountryIcon,
				Country:      tt.fields.Country,
			}
			if got := inp.hardValidate(tt.args.verifyData); got != tt.want {
				t.Errorf("hardValidate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHardCheck(t *testing.T) {
	type fields struct {
		Latitude     string
		Longitude    string
		Geo          string
		Name         string
		ProvinceIcon string
		Province     string
		CountryIcon  string
		Country      string
	}
	type args struct {
		verifyData LocationData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// Test case with all fields valid
		{
			name: "Hard validation with all fields valid",
			fields: fields{
				Name:      "New York",
				Province:  "New York",
				Country:   "USA",
				Latitude:  "-74.0059", // Example valid latitude
				Longitude: "40.7128",  // Example valid longitude
			},
			args: args{verifyData: LocationData{
				Name:      "New York",
				Province:  "New York",
				Country:   "USA",
				Latitude:  "-74.0059", // Example valid latitude
				Longitude: "40.7128",  // Example valid longitude
			}},
			want: true,
		},
		// Test case with invalid latitude format
		{
			name: "Hard validation with invalid latitude format",
			fields: fields{
				Name:      "New York",
				Province:  "New York",
				Country:   "USA",
				Latitude:  "invalid_latitude",
				Longitude: "40.7128",
			},
			args: args{verifyData: LocationData{}},
			want: false,
		},
		// Test case with invalid longitude format
		{
			name: "Hard validation with invalid longitude format",
			fields: fields{
				Name:      "New York",
				Province:  "New York",
				Country:   "USA",
				Latitude:  "-74.0059",
				Longitude: "invalid_longitude",
			},
			args: args{verifyData: LocationData{}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inp := LocationData{
				Latitude:     tt.fields.Latitude,
				Longitude:    tt.fields.Longitude,
				Geo:          tt.fields.Geo,
				Name:         tt.fields.Name,
				ProvinceIcon: tt.fields.ProvinceIcon,
				Province:     tt.fields.Province,
				CountryIcon:  tt.fields.CountryIcon,
				Country:      tt.fields.Country,
			}
			if got := hardCheck(tt.args.verifyData, inp); got != tt.want {
				t.Errorf("hardValidate() = %v, want %v", got, tt.want)
			}
		})
	}
}
