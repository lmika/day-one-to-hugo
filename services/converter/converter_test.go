package converter_test

import (
	"github.com/lmika/day-one-to-hugo/models"
	"github.com/lmika/day-one-to-hugo/services/converter"
	"github.com/lmika/gopkgs/fp/slices"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_WriteToHugo(t *testing.T) {
	scenarios := []struct {
		description string
		source      string
		wantTitle   string
		wantContent string
		wantImages  []int
	}{
		{
			description: "entry without title",
			source:      "This is an entry\n\nThis is another entry\n",
			wantContent: "This is an entry\n\nThis is another entry\n",
			wantImages:  []int{},
		},
		{
			description: "entry with title",
			source:      "# Entry Title\n\nThis is an entry\n\nThis is another entry\n",
			wantTitle:   "Entry Title",
			wantContent: "This is an entry\n\nThis is another entry\n",
			wantImages:  []int{},
		},
		{
			description: "entry with title and heading",
			source:      "# Entry Title\n\nThis is an entry\n\n# Sub Entry\n\nThis is another entry\n",
			wantTitle:   "Entry Title",
			wantContent: "This is an entry\n\n# Sub Entry\n\nThis is another entry\n",
			wantImages:  []int{},
		},
		{
			description: "entry without unnecessary backslashes",
			source:      "RA\\-V missions\\.",
			wantContent: "RA-V missions.\n",
			wantImages:  []int{},
		},
		{
			description: "convert image url",
			source:      "![My image](dayone-moment://91E303B8B3FB4345AE028CE8E0752935)",
			wantContent: "![My image](/images/bla.jpeg)\n",
			wantImages:  []int{0},
		},
		{
			description: "convert image url 2",
			source:      "![My image](dayone-moment://91E303B8B3FB4345AE028CE8E0752935)\n\n![My image](dayone-moment://abc123)",
			wantContent: "![My image](/images/bla.jpeg)\n\n![My image](/images/fla.png)\n",
			wantImages:  []int{0, 1},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			entry := models.Entry{
				Text: scenario.source,
				Photos: []models.Moment{
					{
						ID:   "91E303B8B3FB4345AE028CE8E0752935",
						MD5:  "bla",
						Type: "jpeg",
					},
					{
						ID:   "abc123",
						MD5:  "fla",
						Type: "png",
					},
				},
			}

			svc := converter.New(nil, false)
			post, err := svc.ConvertToPost(entry)

			assert.Nil(t, err)
			assert.Equal(t, scenario.wantTitle, post.Title)
			assert.Equal(t, scenario.wantContent, post.Content)
			assert.Equal(t, post.Moments, slices.Map(scenario.wantImages, func(i int) models.Moment { return entry.Photos[i] }))
		})
	}
}
