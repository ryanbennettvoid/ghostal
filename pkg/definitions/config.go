package definitions

import (
	"ghostal/pkg/utils"
	"ghostal/pkg/values"
	"net/url"
	"strings"
	"time"
)

type Project struct {
	Name        string    `json:"name"`
	DBURL       string    `json:"dbUrl"`
	FastRestore *bool     `json:"fastRestore"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (p Project) DBName() string {
	parsedURL, err := url.Parse(p.DBURL)
	if err != nil {
		return values.Unknown
	}
	return strings.TrimPrefix(parsedURL.Path, "/")
}

func (p Project) DBType(dbOperatorBuilders []IDBOperatorBuilder) string {
	for _, builder := range dbOperatorBuilders {
		if _, err := builder.BuildOperator(p.DBURL); err == nil {
			return builder.ID()
		}
	}
	return values.Unknown
}

type ProjectsList []Project

func (p ProjectsList) TableInfo(selectedProjectName string, dbOperatorBuilders []IDBOperatorBuilder) ([]string, [][]string) {
	columns := []string{"Project", "Database", "Type", "Age"}
	rows := make([][]string, 0)
	for _, p := range p {
		projectName := p.Name
		if p.Name == selectedProjectName {
			projectName = "* " + projectName
		} else {
			projectName = "  " + projectName
		}
		dbName := p.DBName()
		dbType := p.DBType(dbOperatorBuilders)
		relativeTime := utils.ToRelativeTime(p.CreatedAt, time.Now())
		rows = append(rows, []string{projectName, dbName, dbType, relativeTime})
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
