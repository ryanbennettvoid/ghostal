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
		_, _ = fmt.Fprintf(os.Stdout, "Done in %.3fs.\n", time.Now().Sub(start).Seconds())
		os.Exit(0)
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "ERR: %s\n", err)
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
		tableLogger := pretty_table_logger.NewPrettyTableLogger()
		tableLogger.Log([]string{"Selected Project"}, [][]string{{newSelectedProject}})
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
	switch args.Command {
	case app.SnapshotCommand:
		snapshotFn = dbOperator.Snapshot
	case app.RestoreCommand:
		snapshotFn = dbOperator.Restore
	case app.RemoveCommand:
		snapshotFn = dbOperator.Remove
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
