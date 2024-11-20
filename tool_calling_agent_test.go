package main

import "testing"

func Test_toolCallingAgent(t *testing.T) {
	type args struct {
		setup  Setup
		prompt string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toolCallingAgent(tt.args.setup, tt.args.prompt); got != tt.want {
				t.Errorf("toolCallingAgent() = %v, want %v", got, tt.want)
			}
		})
	}
}
