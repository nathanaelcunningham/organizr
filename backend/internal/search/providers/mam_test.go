package providers

import (
	"testing"

	"github.com/nathanael/organizr/internal/models"
)

func Test_parseSeriesInfo(t *testing.T) {
	tests := []struct {
		name       string
		seriesInfo string
		want       []models.SeriesInfo
	}{
		{
			name:       "Single series with book number",
			seriesInfo: `{"30281":["Awaken Online","10",10.000000]}`,
			want: []models.SeriesInfo{
				{ID: "30281", Name: "Awaken Online", Number: "10"},
			},
		},
		{
			name:       "Empty series info",
			seriesInfo: "",
			want:       []models.SeriesInfo{},
		},
		{
			name:       "Multiple series",
			seriesInfo: `{"123":["Series A","1",1.0],"456":["Series B","2",2.0]}`,
			want: []models.SeriesInfo{
				{ID: "123", Name: "Series A", Number: "1"},
				{ID: "456", Name: "Series B", Number: "2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseSeriesInfo(tt.seriesInfo)
			if len(got) != len(tt.want) {
				t.Errorf("parseSeriesInfo() returned %d items, want %d", len(got), len(tt.want))
				return
			}
			// Note: map iteration order is non-deterministic, so we check by ID
			gotMap := make(map[string]models.SeriesInfo)
			for _, s := range got {
				gotMap[s.ID] = s
			}
			for _, want := range tt.want {
				got, ok := gotMap[want.ID]
				if !ok {
					t.Errorf("parseSeriesInfo() missing series ID %s", want.ID)
					continue
				}
				if got.Name != want.Name || got.Number != want.Number {
					t.Errorf("parseSeriesInfo() ID %s = {Name: %s, Number: %s}, want {Name: %s, Number: %s}",
						want.ID, got.Name, got.Number, want.Name, want.Number)
				}
			}
		})
	}
}
