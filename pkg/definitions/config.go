package definitions

import (
	"time"
)

type Project struct {
	Name      string    `json:"name"`
	DBURL     string    `json:"dbUrl"`
	CreatedAt time.Time `json:"createdAt"`
}

type ConfigData struct {
	SelectedProject *string   `json:"selectedProject"`
	Projects        []Project `json:"projects"`
}

type IConfig interface {
	InitProject(name, DBURL string) error
	SelectProject(name string) error
	GetProject(name *string) (Project, error)
	SetProject(name *string, value Project) error
	GetAllProjects() ([]Project, error)
}
