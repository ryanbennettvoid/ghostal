package definitions

import (
	"ghostel/pkg/utils"
	"time"
)

type Project struct {
	Name      string    `json:"name"`
	DBURL     string    `json:"dbUrl"`
	CreatedAt time.Time `json:"createdAt"`
}

type ProjectsList []Project

func (p ProjectsList) Print(logger ITableLogger, selectedProjectName string) {
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
		rows = append(rows, []string{name, p.DBURL, relativeTime, formattedTime})
	}
	logger.Log(columns, rows)
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
