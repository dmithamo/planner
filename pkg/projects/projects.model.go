package projects

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Projects exposes the functionality availed by this pkg
type Projects struct {
	IDB *sql.DB
}

// Model represents the fields present in a project
type Model struct {
	ProjectNumber int64     `json:"project_number"`
	ProjectID     uuid.UUID `json:"id"`
	ProjectSlug   string    `json:"project_slug"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Created       time.Time `json:"created"`
	Updated       time.Time `json:"updated"`
}

// makeSlug creates a slug given a title
func makeSlug(title string) string {
	return strings.Join(strings.Split(title, " "), "-")
}

// Insert creates a new project in the db
func (p *Projects) Insert(title, description string) (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	projectSlug := strings.ToLower(makeSlug(title))

	query := "INSERT INTO projects (projectID,projectSlug,title,description) VALUES(?,?,?,?)"
	result, err := p.IDB.Exec(query, id, projectSlug, title, description)
	if err != nil {
		return "", fmt.Errorf("insert error %v", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("lastInsertID() error %v", err)
	}

	return projectSlug, nil
}

// SelectOne project(s) from the db where title is matched
func (p *Projects) SelectOne(slug string) (*Model, error) {
	query := "SELECT projectNumber,projectSlug,title,description,created,updated FROM projects WHERE projectSlug=?"
	proj := &Model{}
	err := p.IDB.QueryRow(query, slug).Scan(
		&proj.ProjectNumber,
		&proj.ProjectSlug,
		&proj.Title,
		&proj.Description,
		&proj.Created,
		&proj.Updated,
	)
	if err != nil {
		return nil, err
	}

	return proj, nil
}

// SelectAll retrieves all the projects present in the db
func (p *Projects) SelectAll() ([]*Model, error) {
	query := "SELECT projectNumber,projectSlug,title,description,created,updated FROM projects"
	rows, err := p.IDB.Query(query)
	if err != nil {
		return nil, err
	}

	projects := []*Model{}
	for rows.Next() {
		proj := &Model{}
		err = rows.Scan(
			&proj.ProjectNumber,
			&proj.ProjectSlug,
			&proj.Title,
			&proj.Description,
			&proj.Created,
			&proj.Updated,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, proj)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}
