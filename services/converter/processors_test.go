package converter

import (
	"github.com/lmika/day-one-to-hugo/models"
	"github.com/lmika/gopkgs/fp/slices"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSelectOption_EntryMatch(t *testing.T) {
	entries := []models.Entry{
		{
			Date: time.Date(2022, 2, 5, 8, 30, 0, 0, time.UTC),
			Text: "feb 2",
		},
		{
			Date: time.Date(2023, 1, 12, 14, 22, 0, 0, time.UTC),
			Text: "jan 1",
		},
		{
			Date: time.Date(2023, 10, 20, 10, 7, 0, 0, time.UTC),
			Text: "oct 20",
		},
	}

	tests := []struct {
		description string
		options     SelectOption
		wantIndex   []int
	}{
		{
			description: "no guard",
			options:     SelectOption{},
			wantIndex:   []int{0, 1, 2},
		},
		{
			description: "from guard 1",
			options:     SelectOption{FromDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
			wantIndex:   []int{1, 2},
		},
		{
			description: "from guard 2",
			options:     SelectOption{FromDate: time.Date(2023, 1, 12, 14, 22, 0, 0, time.UTC)},
			wantIndex:   []int{1, 2},
		},
		{
			description: "from guard 3",
			options:     SelectOption{FromDate: time.Date(2023, 6, 12, 14, 22, 0, 0, time.UTC)},
			wantIndex:   []int{2},
		},
		{
			description: "to guard 1",
			options:     SelectOption{ToDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
			wantIndex:   []int{0},
		},
		{
			description: "to guard 2",
			options:     SelectOption{ToDate: time.Date(2023, 1, 12, 14, 22, 0, 0, time.UTC)},
			wantIndex:   []int{0},
		},
		{
			description: "to guard 3",
			options:     SelectOption{ToDate: time.Date(2023, 6, 12, 14, 22, 0, 0, time.UTC)},
			wantIndex:   []int{0, 1},
		},
		{
			description: "between guard 1",
			options: SelectOption{
				FromDate: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ToDate:   time.Date(2023, 1, 12, 14, 22, 0, 0, time.UTC),
			},
			wantIndex: []int{0},
		},
		{
			description: "between guard 2",
			options: SelectOption{
				FromDate: time.Date(2023, 1, 12, 14, 22, 0, 0, time.UTC),
				ToDate:   time.Date(2023, 6, 12, 14, 22, 0, 0, time.UTC),
			},
			wantIndex: []int{1},
		},
		{
			description: "between guard 3",
			options: SelectOption{
				FromDate: time.Date(2023, 1, 12, 14, 22, 0, 0, time.UTC),
				ToDate:   time.Date(2024, 6, 12, 14, 22, 0, 0, time.UTC),
			},
			wantIndex: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			wantEntries := slices.Map(tt.wantIndex, func(t int) models.Entry { return entries[t] })
			gotEntries := slices.Filter(entries, tt.options.EntryMatch)

			assert.Equal(t, wantEntries, gotEntries)
		})
	}
}
