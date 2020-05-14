package mysqlservice

import (
	"database/sql"
	"fmt"
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
type Projects struct {
	DB *sql.DB
}

// Insert creates a new project in the db
func (p *Projects) Insert(title, description string) (*Project, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Println(title, err)
		return nil, err
	}
	log.Println(title, id, len(id))

	created := time.Now()
	updated := time.Now()

	query := "INSERT INTO projects (id,title,description,created,updated) VALUES(?,?,?,?,?)"
	infoLogger.Printf("[db] %v (%v)", query, fmt.Sprintf("%v::%v::%v", id, title, description))

	result, err := p.DB.Exec(query, id, title, description, created, updated)
	if err != nil {
		errorLogger.Printf("[db] %v", err)
		return nil, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		errorLogger.Printf("[db] %v", err)
		return nil, err
	}

	newProject := Project{id, title, description, created, updated}
	infoLogger.Printf("[db] successfully inserted project: %v(uuid: %v)", lastID, id)
	return &newProject, nil
}

// Select project(s) from the db where title is matched
func (p *Projects) Select(title string) (*Project, error) {
	idFromTitle, err := uuid.FromBytes([]byte(title))
	if err != nil {
		return nil, err
	}

	query := "SELECT FROM projects (id,title,description,created,updated) WHERE id=?"
	infoLogger.Printf("[db] %v", query)
	row, err := p.DB.Query(query, idFromTitle)
	if err != nil {
		errorLogger.Printf("[db] %v", err)
		return nil, err
	}

	proj := Project{}
	err = row.Scan(&proj.ID, &proj.Title, &proj.Description, &proj.Created, &proj.Updated)
	if err != nil {
		errorLogger.Printf("[db] %v", err)
		return nil, err
	}

	infoLogger.Printf("[db] successfully retrieved project: %v", proj.ID)
	return &proj, nil
}
