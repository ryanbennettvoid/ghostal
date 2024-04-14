package app

import (
	"errors"
	"fmt"
	"ghostel/pkg/adapters/json_file_config"
	"ghostel/pkg/adapters/mongo_db_operator"
	"ghostel/pkg/adapters/postgres_db_operator"
	"ghostel/pkg/adapters/pretty_table_logger"
	"ghostel/pkg/definitions"
	"ghostel/pkg/utils"
	"ghostel/pkg/values"
)

var NoArgsProvidedError = errors.New("no arguments provided")

type App struct {
	version     string
	logger      definitions.ILogger
	tableLogger definitions.ITableLogger
}

func NewApp(version string, logger definitions.ILogger, tableLogger definitions.ITableLogger) *App {
	return &App{
		version:     version,
		logger:      logger,
		tableLogger: tableLogger,
	}
}

func (a *App) createOperator(dbURL string) (definitions.IDBOperator, error) {
	scheme, err := utils.GetURLScheme(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get URL scheme: %w", err)
	}
	switch scheme {
	case "postgresql":
		return postgres_db_operator.CreatePostgresDBOperator(dbURL)
	case "mongodb":
		return mongo_db_operator.CreateMongoDBOperator(dbURL)
	default:
		return nil, fmt.Errorf("unsupported database scheme: %s", scheme)
	}
}

func (a *App) parseProgramArgs(args []string) (ProgramArgs, error) {
	if len(args) == 0 || len(args[0]) == 0 {
		return ProgramArgs{}, NoArgsProvidedError
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
	fmt.Printf("%s version %s\n", executable, a.version)
	return nil
}

func (a *App) printHelp(executable string) error {
	appDescription := "\nGhostel (gho) is a database snapshot/restore tool for MongoDB and Postgres."
	fmt.Println(appDescription)
	fmt.Println()
	columns := []string{"Command", "Description"}
	rows := make([][]string, 0)
	for _, command := range AllCommands(executable) {
		rows = append(rows, command.Row())
	}
	a.tableLogger.Log(columns, rows)
	fmt.Println()
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
	if err := cfg.InitProject(projectName, dbURL); err != nil {
		return err
	}
	sanitizedDBURL, err := utils.SanitizeDBURL(dbURL)
	if err != nil {
		return fmt.Errorf("failed to sanitize database url: %w", err)
	}
	a.logger.Info("Created project \"%s\" with database \"%s\"", projectName, sanitizedDBURL)
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
	a.logger.Info("Selected project \"%s\"", projectName)
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
	allProjects.Print(
		pretty_table_logger.NewPrettyTableLogger(),
		selectedProject.Name,
	)
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
		if err := dbOperator.Restore(snapshotName); err != nil {
			return err
		}
	case "delete":
		if err := dbOperator.Delete(snapshotName); err != nil {
			return err
		}
	default:
		return errors.New("invalid operation")
	}
	a.logger.Info("Snapshot \"%s\" created.")
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
	listItems.Print(pretty_table_logger.NewPrettyTableLogger())
	return nil
}

func (a *App) Run(executable string, programArgs []string) error {
	args, err := a.parseProgramArgs(programArgs)
	if err != nil {
		if err == NoArgsProvidedError {
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

	cfg := json_file_config.NewJSONFileConfig(values.ConfigFilename)

	switch args.Command {
	case InitCommand:
		return a.initProject(cfg, args)
	case SelectCommand:
		return a.selectProject(cfg, args)
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
