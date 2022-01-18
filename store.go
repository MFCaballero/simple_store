package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	fileName   = "store.json"
	statusFile = "status.json"
)

type Incident struct {
	Uuid        uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Status      bool      `json:"status"` //true is open and false is solved
	CreatedAt   time.Time `json:"created_at"`
	SolvedAt    time.Time `json:"solved_at,omitempty"`
}

type Store struct {
	Incidents []*Incident `json:"incidents"`
}

type IncidentStatusResponse struct {
	OpenCases       int     `json:"open_cases"`
	ClosedCases     int     `json:"closed_cases"`
	AverageSolution float64 `json:"average_solution"` //in hours
	MaximumSolution int     `json:"maximun_solution"` //in hours
}

func (store *Store) IncidentStatus(date1, date2 time.Time) (string, error) {

	err := store.decodeFile(fileName, false)
	if err != nil {
		return "", err
	}
	openIncidents, solvedIncidents, averageSolution, maximumSolution := store.findIncidents(date1, date2)

	incidentStatus := &IncidentStatusResponse{
		OpenCases:       openIncidents,
		ClosedCases:     solvedIncidents,
		AverageSolution: averageSolution,
		MaximumSolution: maximumSolution,
	}
	jsonIncidentStatus, err := json.Marshal(incidentStatus)
	if err != nil {
		return "", fmt.Errorf("Error %v marshalling to json", err)
	}
	return string(jsonIncidentStatus), nil
}

func (store *Store) findIncidents(date1, date2 time.Time) (open, solved int, averageSolutions float64, maximumSolution int) {
	openIncidents := 0
	solvedIncidents := 0
	solutionsTime := 0.0
	maximumTime := 0

	for _, incident := range store.Incidents {
		if (incident.CreatedAt.After(date1) && date2.After(incident.CreatedAt)) || (incident.CreatedAt.After(date2) && date1.After(incident.CreatedAt)) {
			if incident.Status {
				openIncidents++
				maximumTime = int(math.Max(float64(maximumTime), float64(time.Since(incident.CreatedAt).Hours())))

			} else {
				solvedIncidents++
				solutionsTime += incident.SolvedAt.Sub(incident.CreatedAt).Hours()
				maximumTime = int(math.Max(float64(maximumTime), float64(incident.SolvedAt.Sub(incident.CreatedAt).Hours())))
			}
		}
	}

	if solvedIncidents > 0 {
		return openIncidents, solvedIncidents, math.Round((float64(solutionsTime)/float64(solvedIncidents))*100) / 100, maximumTime
	}
	return openIncidents, solvedIncidents, 0.0, maximumTime
}

func NewStore() *Store {
	return &Store{}
}

func (store *Store) AddIncident(description string) (string, error) {
	incident := &Incident{Uuid: uuid.NewV4(), Description: description, Status: true, CreatedAt: time.Now(), SolvedAt: time.Time{}}
	err := store.decodeFile(fileName, true)
	store.Incidents = append(store.Incidents, incident)
	err = store.encodeFile(fileName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("added incident with id: %s", incident.Uuid), nil
}

func (store *Store) SolveIncident(id uuid.UUID) (string, error) {
	err := store.decodeFile(fileName, false)
	for _, incident := range store.Incidents {
		if incident.Uuid == id {
			incident.Status = false
			incident.SolvedAt = time.Now()
		}
	}
	err = store.encodeFile(fileName)
	if err != nil {
		return "", err
	}
	return "incident solved", nil
}

func (store *Store) decodeFile(filename string, create bool) error {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		if create {
			file, err := os.Create(fileName)
			if err != nil {
				return fmt.Errorf("Error %v creating file %s", err, fileName)
			}
			defer file.Close()
		} else {
			return fmt.Errorf("No data registered yet. Error: %v", err)
		}
	}
	fileReader, _ := ioutil.ReadFile(fileName)
	err = json.Unmarshal(fileReader, store)
	if err != nil {
		return fmt.Errorf("Error %v decoding file %s", err, fileName)
	}
	return nil
}

func (store *Store) encodeFile(filename string) error {
	fileWriter, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("Error %v opening file %s", err, fileName)
	}
	defer fileWriter.Close()
	err = json.NewEncoder(fileWriter).Encode(store)
	if err != nil {
		return fmt.Errorf("Error %v encoding file %s", err, fileName)
	}
	return nil
}
