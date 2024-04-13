package main

import (
	"errors"
	"fmt"
	"ghostel/pkg/adapters/json_file_config"
	"ghostel/pkg/adapters/mongo_db_operator"
	"ghostel/pkg/adapters/postgres_db_operator"
	"ghostel/pkg/adapters/pretty_table_logger"
	"ghostel/pkg/app"
	"ghostel/pkg/definitions"
	"ghostel/pkg/utils"
	"ghostel/pkg/values"
	"os"
	"time"
)

var logger = app.GetGlobalLogger()

func createOperator(dbURL string) (definitions.IDBOperator, error) {
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

var start = time.Now()

func exit(err error) {
	if err == nil {
		logger.Info("Done in %.3fs.\n", time.Now().Sub(start).Seconds())
		os.Exit(0)
	} else {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func main() {
	args, err := app.ParseProgramArgs(os.Args[1:])
	if err != nil {
		exit(fmt.Errorf("failed to parse program args: %w", err))
	}
	exit(run(args))
}

func run(args app.ProgramArgs) error {

	cfg := json_file_config.NewJSONFileConfig(values.ConfigFilename)

	if args.Command == app.InitCommand {
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
		logger.Info("Created project \"%s\" with database \"%s\"", projectName, sanitizedDBURL)
		return nil
	}

	if args.Command == app.SelectCommand {
		newSelectedProject, err := args.Options.Get(0, "project name")
		if err != nil {
			return err
		}
		if err := cfg.SelectProject(newSelectedProject); err != nil {
			return err
		}
		logger.Info("Selected project \"%s\"", newSelectedProject)
		return nil
	}

	selectedProject, err := cfg.GetProject(nil)
	if err != nil {
		return err
	}

	if args.Command == app.StatusCommand {
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

	dbOperator, err := createOperator(selectedProject.DBURL)
	if err != nil {
		return err
	}

	var snapshotFn func(string) error
	var snapshotMsg string
	switch args.Command {
	case app.SnapshotCommand:
		snapshotFn = dbOperator.Snapshot
		snapshotMsg = "Snapshot \"%s\" created."
	case app.RestoreCommand:
		snapshotFn = dbOperator.Restore
		snapshotMsg = "Snapshot \"%s\" restored."
	case app.RemoveCommand:
		snapshotFn = dbOperator.Remove
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
		logger.Info(snapshotMsg, snapshotName)
		return nil
	}

	if args.Command == app.ListCommand {
		listItems, err := dbOperator.List()
		if err != nil {
			return err
		}
		listItems.Print(pretty_table_logger.NewPrettyTableLogger())
		return nil
	}

	return fmt.Errorf("unknown command \"%s\"", args.Command)
}
