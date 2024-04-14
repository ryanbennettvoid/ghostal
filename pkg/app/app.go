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
	logger      definitions.ILogger
	tableLogger definitions.ITableLogger
}

func NewApp(logger definitions.ILogger, tableLogger definitions.ITableLogger) *App {
	return &App{
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

func (a *App) help(executable string) {
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
}

func (a *App) Run(executable string, programArgs []string) error {

	args, err := a.parseProgramArgs(programArgs)
	if err != nil {
		if err == NoArgsProvidedError {
			a.help(executable)
			return nil
		}
		return err
	}

	if args.Command == HelpCommand {
		a.help(executable)
		return nil
	}

	cfg := json_file_config.NewJSONFileConfig(values.ConfigFilename)

	if args.Command == InitCommand {
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

	if args.Command == SelectCommand {
		newSelectedProject, err := args.Options.Get(0, "project name")
		if err != nil {
			return err
		}
		if err := cfg.SelectProject(newSelectedProject); err != nil {
			return err
		}
		a.logger.Info("Selected project \"%s\"", newSelectedProject)
		return nil
	}

	selectedProject, err := cfg.GetProject(nil)
	if err != nil {
		return err
	}

	if args.Command == StatusCommand {
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
		allProjects.Print(
			pretty_table_logger.NewPrettyTableLogger(),
			selectedProject.Name,
		)
		return nil
	}

	dbOperator, err := a.createOperator(selectedProject.DBURL)
	if err != nil {
		return err
	}

	var snapshotFn func(string) error
	var snapshotMsg string
	switch args.Command {
	case SnapshotCommand:
		snapshotFn = dbOperator.Snapshot
		snapshotMsg = "Snapshot \"%s\" created."
	case RestoreCommand:
		snapshotFn = dbOperator.Restore
		snapshotMsg = "Snapshot \"%s\" restored."
	case RemoveCommand:
		snapshotFn = dbOperator.Delete
		snapshotMsg = "Snapshot \"%s\" removed."
	}

	if snapshotFn != nil {
		snapshotName, err := args.Options.Get(0, "snapshot name")
		if err != nil {
			return err
		}
		if len(snapshotName) == 0 {
			return errors.New("snapshot name must be specified")
		}
		if err := snapshotFn(snapshotName); err != nil {
			return err
		}
		a.logger.Info(snapshotMsg, snapshotName)
		return nil
	}

	if args.Command == ListCommand {
		listItems, err := dbOperator.List()
		if err != nil {
			return err
		}
		listItems.Print(pretty_table_logger.NewPrettyTableLogger())
		return nil
	}

	fullHelpCommand := fmt.Sprintf("%s help", executable)
	return fmt.Errorf("unknown command \"%s\" - run \"%s\" for help", args.Command, fullHelpCommand)
}
