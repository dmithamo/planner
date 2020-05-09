package mysqlservice

import (
	"log"
	"time"

	"github.com/google/uuid"
)

// Project represents the fields present in a project
type Project struct {
	ID          uuid.UUID
	Title       string
	Description string
	Created     time.Time
	Updated     time.Time
}

// Projects maps to the `projects` table in the db
type Projects struct{}

// Insert creates a new project in the db
func (projects *Projects) Insert(title, description string) (*Project, error) {
	log.Printf("inserting...%v, %v", title, description)
	return nil, nil
}

// Select retrieves (a slice of)  project(s) from the db
func (projects *Projects) Select(title string) ([]*Project, error) {
	return nil, nil
}
