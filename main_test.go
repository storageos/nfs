package main

import (
	"os"
	"testing"
)

func Test_getEnv(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		defaultVal string
		setVal     string
		want       string
	}{
		{
			name:   "set",
			key:    "Test_getEnv",
			setVal: "stringVal",
			want:   "stringVal",
		},
		{
			name:       "not set with default",
			key:        "Test_getEnv",
			defaultVal: "defaultVal",
			want:       "defaultVal",
		},
		{
			name: "not set no default",
			key:  "Test_getEnv",
			want: "",
		},
		{
			name:       "set with default",
			key:        "Test_getEnv",
			defaultVal: "defaultVal",
			setVal:     "stringVal",
			want:       "stringVal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			os.Clearenv()

			if tt.setVal != "" {
				os.Setenv(tt.key, tt.setVal)
				defer os.Setenv(tt.key, "")
			}

			if got := getEnv(tt.key, tt.defaultVal); got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBoolEnv(t *testing.T) {
	type args struct {
		key        string
		defaultVal bool
	}
	tests := []struct {
		name       string
		key        string
		defaultVal bool
		setVal     string
		want       bool
		wantErr    bool
	}{
		{
			name:   "set true",
			key:    "Test_getBoolEnv",
			setVal: "true",
			want:   true,
		},
		{
			name:   "set false",
			key:    "Test_getBoolEnv",
			setVal: "false",
			want:   false,
		},
		{
			name:       "not set with default true",
			key:        "Test_getBoolEnv",
			defaultVal: true,
			want:       true,
		},
		{
			name:       "not set with default false",
			key:        "Test_getBoolEnv",
			defaultVal: false,
			want:       false,
		},
		{
			name: "not set no default",
			key:  "Test_getBoolEnv",
			want: false,
		},
		{
			name:       "set with default",
			key:        "Test_getBoolEnv",
			defaultVal: false,
			setVal:     "true",
			want:       true,
		},
		{
			name:   "set 1",
			key:    "Test_getBoolEnv",
			setVal: "1",
			want:   true,
		},
		{
			name:   "set 0",
			key:    "Test_getBoolEnv",
			setVal: "0",
			want:   false,
		},
		{
			name:    "set non-bool",
			key:     "Test_getBoolEnv",
			setVal:  "non-bool",
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			os.Clearenv()

			if tt.setVal != "" {
				os.Setenv(tt.key, tt.setVal)
				defer os.Setenv(tt.key, "")
			}

			got, err := getBoolEnv(tt.key, tt.defaultVal)
			if err == nil && tt.wantErr {
				t.Error("getBoolEnv(): got no error even though we wanted one")
			} else if err != nil && !tt.wantErr {
				t.Errorf("getBoolEnv(): got an error even though we wanted none, got: %v", err)
			}

			if got != tt.want {
				t.Errorf("getBoolEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
