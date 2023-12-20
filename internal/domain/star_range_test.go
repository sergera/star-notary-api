package domain

import (
	"testing"
)

func TestValidateStarRange(t *testing.T) {
	tests := []struct {
		name        string
		starRange   StarRangeModel
		expectedErr bool
	}{
		{"Valid Range", StarRangeModel{Start: "100", End: "200", OldestFirst: true}, false},
		{"Invalid Range - Start", StarRangeModel{Start: "0abc", End: "200", OldestFirst: true}, true},
		{"Invalid Range - End", StarRangeModel{Start: "100", End: "0abc", OldestFirst: true}, true},
		{"Invalid Range - Both Start and End", StarRangeModel{Start: "0abc", End: "0xyz", OldestFirst: true}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.starRange.Validate()
			if (err != nil) != tt.expectedErr {
				t.Errorf("TestValidateStarRange(%s): expected error: %v, got error: %v", tt.name, tt.expectedErr, err != nil)
			}
		})
	}
}
