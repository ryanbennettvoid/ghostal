package definitions

import (
	"ghostal/pkg/utils"
	"time"
)

type Project struct {
	Name        string    `json:"name"`
	DBURL       string    `json:"dbUrl"`
	FastRestore *bool     `json:"fastRestore"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ProjectsList []Project

func (p ProjectsList) TableInfo(selectedProjectName string) ([]string, [][]string) {
	columns := []string{"Projects", "Database URL", "Created", ""}
	rows := make([][]string, 0)
	for _, p := range p {
		name := p.Name
		if p.Name == selectedProjectName {
			name = "* " + name
		} else {
			name = "  " + name
		}
		relativeTime := utils.ToRelativeTime(p.CreatedAt, time.Now())
		formattedTime := p.CreatedAt.Format("2006-01-02 15:04:05")
		sanitizedDBURL, err := utils.SanitizeDBURL(p.DBURL)
		if err != nil {
			sanitizedDBURL = "<PARSE ERROR>"
		}
		rows = append(rows, []string{name, sanitizedDBURL, relativeTime, formattedTime})
	}
	return columns, rows
}

type ConfigData struct {
	SelectedProject string    `json:"selectedProject"`
	Projects        []Project `json:"projects"`
}

type IConfig interface {
	InitProject(name, DBURL string) error
	SelectProject(name string) error
	GetProject(name *string) (Project, error)
	SetProject(name *string, value Project) error
	GetAllProjects() (ProjectsList, error)
}
