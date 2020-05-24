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
	ProjectNumber int64     `json:"project_number"`
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Created       time.Time `json:"created"`
	Updated       time.Time `json:"updated"`
}

// Insert creates a new project in the db
func (p *Projects) Insert(title, description string) (*Model, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		p.ErrorLogger.Println(err)
		return nil, err
	}

	query := "INSERT INTO projects (id,title,description) VALUES(?,?,?,?,?)"
	p.InfoLogger.Printf("[db::projects] %v (%v)", query, fmt.Sprintf("%v::%v::%v", id, title, description))

	result, err := p.IDB.Exec(query, id, title, description)
	if err != nil {
		p.ErrorLogger.Printf("[db::projects] %v", err)
		return nil, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		p.ErrorLogger.Printf("[db::projects] %v", err)
		return nil, err
	}

	newProject := Model{
		ProjectNumber: lastInsertID,
		ID:            id,
		Title:         title,
		Description:   description,
	}

	p.InfoLogger.Printf("[db::projects] insert project %v?:: OK", lastInsertID)
	return &newProject, nil
}

// SelectOne project(s) from the db where title is matched
func (p *Projects) SelectOne(id string) (*Model, error) {
	query := "SELECT nu,title,description,created,updated FROM projects WHERE id=?"
	p.InfoLogger.Printf("[db::projects] %v %v", query, id)
	proj := &Model{}
	err := p.IDB.QueryRow(query, id).Scan(
		&proj.ProjectNumber,
		&proj.Title,
		&proj.Description,
		&proj.Created,
		&proj.Updated,
	)
	if err != nil {
		return nil, err
	}

	p.InfoLogger.Printf("[db::projects] select project %v?:: OK", proj.ID)
	return proj, nil
}

// SelectAll retrieves all the projects present in the db
func (p *Projects) SelectAll() ([]*Model, error) {
	query := "SELECT nu,title,description,created,updated FROM projects"
	p.InfoLogger.Printf("[db::projects] %v", query)
	row, err := p.IDB.Query(query, idFromTitle)
	if err != nil {
		p.ErrorLogger.Printf("[db::projects] %v", err)
		return nil, err
	}

	projects := []*Model{}
	for rows.Next() {
		proj := &Model{}
		err = rows.Scan(
			&proj.ProjectNumber,
			&proj.Title,
			&proj.Description,
			&proj.Created,
			&proj.Updated,
		)
		if err != nil {
			p.ErrorLogger.Printf("[db::projects] %v", err)
			return nil, err
		}
		projects = append(projects, proj)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	p.InfoLogger.Printf("[db::projects] retrieved project::%v::?: OK", proj.ID)
	return &proj, nil
}
