package app

import (
	"errors"
	"fmt"
	"ghostel/pkg/adapters/json_file_config"
	"ghostel/pkg/definitions"
	"ghostel/pkg/utils"
	"ghostel/pkg/values"
)

type App struct {
	version            string
	dbOperatorBuilders []definitions.IDBOperatorBuilder
	logger             definitions.ILogger
	tableBuilder       definitions.ITableBuilder
}

func NewApp(
	version string,
	dbOperatorBuilders []definitions.IDBOperatorBuilder,
	logger definitions.ILogger,
	tableBuilder definitions.ITableBuilder,
) *App {
	return &App{
		version:            version,
		dbOperatorBuilders: dbOperatorBuilders,
		logger:             logger,
		tableBuilder:       tableBuilder,
	}
}

func (a *App) createOperator(dbURL string) (definitions.IDBOperator, error) {
	for _, builder := range a.dbOperatorBuilders {
		dbOperator, err := builder.BuildOperator(dbURL)
		if err != nil {
			if err != values.UnsupportedURLSchemeError {
				a.logger.Warning("unexpected error while attempting to build %s operator: %s", builder.ID(), err)
			}
			continue
		}
		return dbOperator, nil
	}
	return nil, errors.New("no supported database operator found for the given database URL")
}

func (a *App) parseProgramArgs(args []string) (ProgramArgs, error) {
	if len(args) == 0 || len(args[0]) == 0 {
		return ProgramArgs{}, values.NoProgramArgsProvidedError
	}
	command := Command(args[0])
	options := make([]string, 0)
	if len(args) > 1 {
		options = args[1:]
	}
	return ProgramArgs{
		Command: command,
		Options: options,
	}, nil
}

func (a *App) printVersion(executable string) error {
	a.logger.Passthrough("%s version %s\n", executable, a.version)
	return nil
}

func (a *App) printHelp(executable string) error {
	appDescription := "\nGhostel (gho) is a database snapshot/restore tool for MongoDB and Postgres."
	a.logger.Passthrough(appDescription)
	a.logger.Passthrough("")
	columns := []string{"Command", "Description"}
	rows := make([][]string, 0)
	for _, command := range AllCommands(executable) {
		rows = append(rows, command.Row())
	}
	a.logger.Passthrough(a.tableBuilder.BuildTable(columns, rows))
	a.logger.Passthrough("")
	return nil
}

func (a *App) initProject(cfg definitions.IConfig, args ProgramArgs) error {
	projectName, err := args.Options.Get(0, "project name")
	if err != nil {
		return err
	}
	dbURL, err := args.Options.Get(1, "database URL")
	if err != nil {
		return err
	}
	// attempt to create operator just to see if DB url is valid
	if _, err = a.createOperator(dbURL); err != nil {
		return err
	}
	if err := cfg.InitProject(projectName, dbURL); err != nil {
		return err
	}
	sanitizedDBURL, err := utils.SanitizeDBURL(dbURL)
	if err != nil {
		return fmt.Errorf("failed to sanitize database url: %w", err)
	}
	a.logger.Passthrough("Created project \"%s\" with database \"%s\"\n", projectName, sanitizedDBURL)
	return nil
}

func (a *App) selectProject(cfg definitions.IConfig, args ProgramArgs) error {
	projectName, err := args.Options.Get(0, "project name")
	if err != nil {
		return err
	}
	if err := cfg.SelectProject(projectName); err != nil {
		return err
	}
	a.logger.Passthrough("Selected project \"%s\"\n", projectName)
	return nil
}

func (a *App) setProjectConfigKeyValue(cfg definitions.IConfig, args ProgramArgs) error {
	selectedProject, err := cfg.GetProject(nil)
	if err != nil {
		return err
	}
	key, err := args.Options.Get(0, "project config key")
	if err != nil {
		return err
	}
	value, err := args.Options.Get(1, "project config value")
	if err != nil {
		return err
	}
	switch key {
	case "fastRestore":
		{
			asBool, err := utils.StringAsBool(value)
			if err != nil {
				return err
			}
			selectedProject.FastRestore = utils.ToPointer(asBool)
		}
	default:
		return fmt.Errorf("invalid key: \"%s\"", key)
	}
	if err := cfg.SetProject(utils.ToPointer(selectedProject.Name), selectedProject); err != nil {
		return err
	}
	return nil
}

func (a *App) printStatus(cfg definitions.IConfig) error {
	allProjects, err := cfg.GetAllProjects()
	if err != nil {
		return err
	}
	for idx := range allProjects {
		sanitizedDBURL, err := utils.SanitizeDBURL(allProjects[idx].DBURL)
		if err != nil {
			return fmt.Errorf("failed to sanitize database url: %w", err)
		}
		allProjects[idx].DBURL = sanitizedDBURL
	}
	selectedProject, err := cfg.GetProject(nil)
	if err != nil {
		return err
	}
	columns, rows := allProjects.TableInfo(selectedProject.Name)
	a.logger.Passthrough(a.tableBuilder.BuildTable(columns, rows))
	return nil
}

func (a *App) getDBOperator(cfg definitions.IConfig) (definitions.IDBOperator, error) {
	selectedProject, err := cfg.GetProject(nil)
	if err != nil {
		return nil, err
	}
	dbOperator, err := a.createOperator(selectedProject.DBURL)
	if err != nil {
		return nil, err
	}
	return dbOperator, nil
}

func (a *App) snapshotCommand(cfg definitions.IConfig, args ProgramArgs, operation string) error {
	selectedProject, err := cfg.GetProject(nil)
	if err != nil {
		return nil
	}
	snapshotName, err := args.Options.Get(0, "project name")
	if err != nil {
		return err
	}
	dbOperator, err := a.getDBOperator(cfg)
	if err != nil {
		return err
	}
	switch operation {
	case "create":
		if err := dbOperator.Snapshot(snapshotName); err != nil {
			return err
		}
	case "restore":
		fastRestore := false
		if selectedProject.FastRestore != nil && *selectedProject.FastRestore == true {
			fastRestore = true
		}
		if err := dbOperator.Restore(snapshotName, fastRestore); err != nil {
			return err
		}
	case "delete":
		if err := dbOperator.Delete(snapshotName); err != nil {
			return err
		}
	default:
		return errors.New("invalid operation")
	}
	a.logger.Passthrough("Snapshot \"%s\" %sd.\n", snapshotName, operation)
	return nil
}

func (a *App) listSnapshots(cfg definitions.IConfig) error {
	dbOperator, err := a.getDBOperator(cfg)
	if err != nil {
		return err
	}
	listItems, err := dbOperator.List()
	if err != nil {
		return err
	}
	columns, rows := listItems.TableInfo()
	a.logger.Passthrough(a.tableBuilder.BuildTable(columns, rows))
	return nil
}

func (a *App) Run(dataStore definitions.IDataStore, executable string, programArgs []string) error {
	args, err := a.parseProgramArgs(programArgs)
	if err != nil {
		if err == values.NoProgramArgsProvidedError {
			return a.printHelp(executable)
		}
		return err
	}

	switch args.Command {
	case VersionCommand:
		return a.printVersion(executable)
	case HelpCommand:
		return a.printHelp(executable)
	}

	cfg := json_file_config.NewJSONFileConfig(dataStore)

	switch args.Command {
	case InitCommand:
		return a.initProject(cfg, args)
	case SelectCommand:
		return a.selectProject(cfg, args)
	case SetCommand:
		return a.setProjectConfigKeyValue(cfg, args)
	case StatusCommand:
		return a.printStatus(cfg)
	}

	switch args.Command {
	case SnapshotCommand:
		return a.snapshotCommand(cfg, args, "create")
	case RestoreCommand:
		return a.snapshotCommand(cfg, args, "restore")
	case DeleteCommand:
		return a.snapshotCommand(cfg, args, "delete")
	case ListCommand:
		return a.listSnapshots(cfg)
	}

	fullHelpCommand := fmt.Sprintf("%s help", executable)
	return fmt.Errorf("unknown command \"%s\" - run \"%s\" for help", args.Command, fullHelpCommand)
}
