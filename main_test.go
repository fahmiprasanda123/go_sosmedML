package main

import (
	"testing"
)

func TestIsEligible(t *testing.T) {
	brand := "Indomilk"
	keyword := "Susu"

	tests := []struct {
		name     string
		post     PostEntity
		expected bool
	}{
		{
			name: "Valid via Caption",
			post: PostEntity{
				Platform: "instagram",
				Caption:  "Aku suka minum Indomilk susu setiap pagi",
			},
			expected: true,
		},
		{
			name: "Valid via TikTok Title",
			post: PostEntity{
				Platform: "tiktok",
				Caption:  "Cuma joget",
				Title:    "Review Indomilk Susu cair",
			},
			expected: true,
		},
		{
			name: "Valid via Single Comment",
			post: PostEntity{
				Platform: "instagram",
				Caption:  "Foto estetik",
				Comments: []string{"Keren bang", "Wah ini Indomilk susu kesukaanku"},
			},
			expected: true,
		},
		{
			name: "Invalid: Brand and Keyword Split Between Fields",
			post: PostEntity{
				Platform: "instagram",
				Caption:  "Ini produk Indomilk",
				Comments: []string{"Rasanya rasa susu"},
			},
			expected: false,
		},
		{
			name: "Valid: Case Insensitive Random Case",
			post: PostEntity{
				Platform: "instagram",
				Caption:  "iNDoMiLk ini rasanya SuSu banget",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEligible(tt.post, brand, keyword); got != tt.expected {
				t.Errorf("IsEligible() = %v, want %v", got, tt.expected)
			}
		})
	}
}
