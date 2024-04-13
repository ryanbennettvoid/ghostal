package json_file_config

import (
	"encoding/json"
	"errors"
	"fmt"
	"ghostel/pkg/definitions"
	"ghostel/pkg/utils"
	"io/ioutil"
	"os"
	"time"
)

type JSONFileConfig struct {
	FilePath   string
	ConfigData definitions.ConfigData
}

func NewJSONFileConfig(filePath string) *JSONFileConfig {
	return &JSONFileConfig{
		FilePath: filePath,
	}
}

func (cm *JSONFileConfig) load() error {
	data, err := ioutil.ReadFile(cm.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No file yet, start with empty config
		}
		return err
	}
	return json.Unmarshal(data, &cm.ConfigData)
}

func (cm *JSONFileConfig) save() error {
	data, err := json.MarshalIndent(cm.ConfigData, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cm.FilePath, data, 0644)
}

func (cm *JSONFileConfig) InitProject(name, dbURL string) error {
	if len(name) == 0 {
		return errors.New("name cannot be empty")
	}
	if len(dbURL) == 0 {
		return errors.New("db URL cannot be empty")
	}

	if err := cm.load(); err != nil {
		return err
	}

	for _, p := range cm.ConfigData.Projects {
		if p.Name == name {
			return errors.New("project already exists")
		}
	}

	cm.ConfigData.SelectedProject = utils.ToPointer(name)

	newProject := definitions.Project{
		Name:      name,
		DBURL:     dbURL,
		CreatedAt: time.Now(),
	}

	cm.ConfigData.Projects = append(cm.ConfigData.Projects, newProject)
	return cm.save()
}

func (cm *JSONFileConfig) SelectProject(name string) error {
	if err := cm.load(); err != nil {
		return err
	}

	for _, p := range cm.ConfigData.Projects {
		if p.Name == name {
			cm.ConfigData.SelectedProject = &p.Name
			return cm.save()
		}
	}
	return fmt.Errorf("project \"%s\" not found", name)
}

func (cm *JSONFileConfig) GetProject(name *string) (definitions.Project, error) {
	if err := cm.load(); err != nil {
		return definitions.Project{}, err
	}

	projectToGet := ""
	if name != nil {
		projectToGet = *name
	} else if cm.ConfigData.SelectedProject != nil {
		projectToGet = *cm.ConfigData.SelectedProject
	} else {
		return definitions.Project{}, errors.New("failed to determine project to get")
	}

	for _, p := range cm.ConfigData.Projects {
		if p.Name == projectToGet {
			return p, nil
		}
	}
	return definitions.Project{}, errors.New("project not found")
}

func (cm *JSONFileConfig) SetProject(name *string, project definitions.Project) error {
	if err := cm.load(); err != nil {
		return err
	}

	projectToSet := ""
	if name != nil {
		projectToSet = *name
	} else if cm.ConfigData.SelectedProject != nil {
		projectToSet = *cm.ConfigData.SelectedProject
	} else {
		return errors.New("failed to determine project to set")
	}

	for i, p := range cm.ConfigData.Projects {
		if p.Name == projectToSet {
			cm.ConfigData.Projects[i] = project
			return cm.save()
		}
	}

	return errors.New("project not found")
}

func (cm *JSONFileConfig) GetAllProjects() (definitions.ProjectsList, error) {
	if err := cm.load(); err != nil {
		return nil, err
	}

	if cm.ConfigData.Projects == nil {
		return make([]definitions.Project, 0), nil
	}

	return cm.ConfigData.Projects, nil
}
