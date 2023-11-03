package converter_test

import (
	"github.com/lmika/day-one-to-hugo/models"
	"github.com/lmika/day-one-to-hugo/services/converter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_WriteToHugo(t *testing.T) {
	scenarios := []struct {
		description     string
		source          string
		expectedTitle   string
		expectedContent string
	}{
		{
			description:     "entry without title",
			source:          "This is an entry\n\nThis is another entry\n",
			expectedContent: "This is an entry\n\nThis is another entry\n",
		},
		{
			description:     "entry with title",
			source:          "# Entry Title\n\nThis is an entry\n\nThis is another entry\n",
			expectedTitle:   "Entry Title",
			expectedContent: "This is an entry\n\nThis is another entry\n",
		},
		{
			description:     "entry with title and heading",
			source:          "# Entry Title\n\nThis is an entry\n\n# Sub Entry\n\nThis is another entry\n",
			expectedTitle:   "Entry Title",
			expectedContent: "This is an entry\n\n# Sub Entry\n\nThis is another entry\n",
		},
		{
			description:     "entry without unnecessary backslashes",
			source:          "RA\\-V missions\\.",
			expectedContent: "RA-V missions.\n",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			entry := models.Entry{
				Text: scenario.source,
			}

			svc := converter.New(nil)
			post, err := svc.ConvertToPost(entry)

			assert.Nil(t, err)
			assert.Equal(t, scenario.expectedTitle, post.Title)
			assert.Equal(t, scenario.expectedContent, post.Content)
		})
	}
}
