package summary

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDaysInRange(t *testing.T) {
	testCases := []struct {
		start, end time.Time
		days       []time.Time
	}{
		{
			start: time.Date(2018, 2, 1, 0, 0, 0, 0, time.UTC),
			end:   time.Date(2018, 2, 3, 0, 0, 0, 0, time.UTC),
			days: []time.Time{
				time.Date(2018, 2, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 2, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 2, 3, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			start: time.Date(2018, 2, 28, 0, 0, 0, 0, time.UTC),
			end:   time.Date(2018, 3, 2, 0, 0, 0, 0, time.UTC),
			days: []time.Time{
				time.Date(2018, 2, 28, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 3, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 3, 2, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			actual := daysInRange(tc.start, tc.end)
			assert.Equal(t, tc.days, actual)
		})
	}
}
