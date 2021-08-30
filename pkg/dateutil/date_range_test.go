package dateutil

import (
	"testing"
	"time"
)

func TestNilDateRange(t *testing.T) {
	dateRange := NewDateRange(nil, nil)
	if dateRange.Start != nil || dateRange.End != nil {
		t.Errorf("unexpected value")
	}
}

func TestDateRange(t *testing.T) {
	start := time.Date(2021, 06, 01, 00, 00, 00, 00, time.Local)
	end := time.Date(2021, 06, 02, 00, 00, 00, 00, time.Local)
	expectedEnd := time.Date(2021, 06, 02, 23, 59, 59, 59, time.Local)
	dateRange := NewDateRange(&start, &end)
	if dateRange.Start.Unix() != start.Unix() {
		t.Errorf("wrong start date")
	}
	if dateRange.End.Unix() != expectedEnd.Unix() {
		t.Errorf("wrong end date")
	}
}
