package domain

import (
	"testing"
	"time"
)

func defaultValidStarModel() StarModel {
	return StarModel{
		TokenId:     "12345",
		Coordinates: "012345.67+012345.67",
		Name:        "alpha centauri",
		Price:       "100.00",
		IsForSale:   true,
		Date:        time.Now(),
		Wallet:      &WalletModel{Address: "0x1234567890123456789012345678901234567890"},
		Action:      Create,
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		star     StarModel
		expected bool
	}{
		{"Valid Star - Create", defaultValidStarModel(), true},
		{"Invalid Star - Create with Invalid TokenId", func() StarModel {
			s := defaultValidStarModel()
			s.TokenId = ""
			return s
		}(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.star.Validate()
			if (err == nil) != tt.expected {
				t.Errorf("TestValidate(%v): expected %v, got %v", tt.star, tt.expected, err == nil)
			}
		})
	}
}

func TestValidateTokenId(t *testing.T) {
	tests := []struct {
		name     string
		tokenId  string
		expected bool
	}{
		{"Valid TokenId", "12345", true},
		{"Invalid TokenId - Empty", "", false},
		{"Invalid TokenId - Non-Numeric", "abcde", false},
		{"Invalid TokenId - Leading Zero", "01234", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := defaultValidStarModel()
			s.TokenId = tt.tokenId
			err := s.ValidateTokenId()
			if (err == nil) != tt.expected {
				t.Errorf("TestValidateTokenId(%s): expected %v, got %v", tt.tokenId, tt.expected, err == nil)
			}
		})
	}
}

func TestValidateCoordinates(t *testing.T) {
	tests := []struct {
		name        string
		coordinates string
		expected    bool
	}{
		{"Valid Coordinates", "012345.67+012345.67", true},
		{"Invalid Coordinates - Format", "123.456", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := defaultValidStarModel()
			s.Coordinates = tt.coordinates
			err := s.ValidateCoordinates()
			if (err == nil) != tt.expected {
				t.Errorf("TestValidateCoordinates(%s): expected %v, got %v", tt.coordinates, tt.expected, err == nil)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name     string
		starName string
		expected bool
	}{
		{"Valid Name", "alpha centauri", true},
		{"Invalid Name - Too Short", "a", false},
		{"Invalid Name - Too Long", "lorem ipsum dolor sit amet consectetur adipiscing elit", false},
		{"Invalid Name - Invalid Characters", "alpha123", false},
		{"Invalid Name - Empty", "", false},
		{"Invalid Name - Leading Space", " alpha centauri", false},
		{"Invalid Name - Trailing Space", "alpha centauri ", false},
		{"Invalid Name - Multiple Spaces", "alpha  centauri", false},
		{"Invalid Name - Upper Case", "Alpha centauri", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := defaultValidStarModel()
			s.Name = tt.starName
			err := s.ValidateName()
			if (err == nil) != tt.expected {
				t.Errorf("TestValidateName(%s): expected %v, got %v", tt.starName, tt.expected, err == nil)
			}
		})
	}
}

func TestValidatePrice(t *testing.T) {
	tests := []struct {
		name      string
		starPrice string
		expected  bool
	}{
		{"Valid Price - Integer", "100", true},
		{"Valid Price - Decimal", "99.99", true},
		{"Invalid Price - Non-Numeric", "abc", false},
		{"Invalid Price - Too Long Integer Part", "1234567890123", false},
		{"Invalid Price - Too Long Fraction Part", "100.1234567890123456789", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := defaultValidStarModel()
			s.Price = tt.starPrice
			err := s.ValidatePrice()
			if (err == nil) != tt.expected {
				t.Errorf("TestValidatePrice(%s): expected %v, got %v", tt.starPrice, tt.expected, err == nil)
			}
		})
	}
}
