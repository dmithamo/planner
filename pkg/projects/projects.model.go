package projects

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// Projects exposes the functionality availed by this pkg
type Projects struct {
	IDB         *sql.DB
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}

// Model represents the fields present in a project
type Model struct {
	ID          uuid.UUID
	Title       string
	Description string
	Created     time.Time
	Updated     time.Time
}

// Insert creates a new project in the db
func (p *Projects) Insert(title, description string) (*Model, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		p.ErrorLogger.Println(err)
		return nil, err
	}
	created := time.Now()
	updated := time.Now()

	query := "INSERT INTO projects (id,title,description,created,updated) VALUES(?,?,?,?,?)"
	p.InfoLogger.Printf("[db::projects] %v (%v)", query, fmt.Sprintf("%v::%v::%v", id, title, description))

	result, err := p.IDB.Exec(query, id, title, description, created, updated)
	if err != nil {
		p.ErrorLogger.Printf("[db::projects] %v", err)
		return nil, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		p.ErrorLogger.Printf("[db::projects] %v", err)
		return nil, err
	}

	newProject := Model{id, title, description, created, updated}
	p.InfoLogger.Printf("[db::projects] successfully inserted project: %v(uuid: %v)", lastID, id)
	return &newProject, nil
}

// Select project(s) from the db where title is matched
func (p *Projects) Select(title string) (*Model, error) {
	idFromTitle, err := uuid.FromBytes([]byte(title))
	if err != nil {
		return nil, err
	}

	query := "SELECT FROM projects (id,title,description,created,updated) WHERE id=?"
	p.InfoLogger.Printf("[db::projects] %v", query)
	row, err := p.IDB.Query(query, idFromTitle)
	if err != nil {
		p.ErrorLogger.Printf("[db::projects] %v", err)
		return nil, err
	}

	proj := Model{}
	err = row.Scan(&proj.ID, &proj.Title, &proj.Description, &proj.Created, &proj.Updated)
	if err != nil {
		p.ErrorLogger.Printf("[db::projects] %v", err)
		return nil, err
	}

	p.InfoLogger.Printf("[db::projects] successfully retrieved project: %v", proj.ID)
	return &proj, nil
}
