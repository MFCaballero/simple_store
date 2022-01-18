package store

import (
	"testing"
	"time"
)

func TestStore_IncidentStatus(t *testing.T) {

	tests := []struct {
		name    string
		date1   string
		date2   string
		want    string
		wantErr bool
	}{
		{"success", "2022-01-08 09:00", "2022-01-19 11:00", `{"open_cases":2,"closed_cases":2,"average_solution":96.18,"maximun_solution":241}`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &Store{}
			parsedDate1, err := time.Parse("2006-01-02 15:04", tt.date1)
			parsedDate2, err := time.Parse("2006-01-02 15:04", tt.date2)
			got, err := store.IncidentStatus(parsedDate1, parsedDate2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.IncidentStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Store.IncidentStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
