package main

import (
	"errors"
	"fmt"
	"ghostel/pkg/adapters/json_file_config"
	"ghostel/pkg/adapters/postgres_db_adapter"
	"ghostel/pkg/app"
	"ghostel/pkg/definitions"
	"ghostel/pkg/utils"
	_ "github.com/olekukonko/tablewriter"
	"os"
	"time"
)

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
	if err := run(args); err != nil {
		exit(err)
	}
}

func run(args app.ProgramArgs) error {

	cfg, err := json_file_config.CreateJSONFileConfig()
	if err != nil {
		return err
	}

	if args.Command == app.InitCommand {
		dbURL := args.Options
		if len(dbURL) == 0 {
			return errors.New("must specify DB url")
		}
		if err := cfg.Set("database_url", dbURL); err != nil {
			return fmt.Errorf("failed to set database_url in config: %w", err)
		}
		return nil
	}

	dbURL, err := cfg.Get("database_url")
	if err != nil {
		return fmt.Errorf("failed to get database URL from config: %w", err)
	}

	var dbOperator definitions.IDBOperator

	scheme, err := utils.GetURLScheme(dbURL)
	if err != nil {
		return fmt.Errorf("failed to get URL scheme: %w", err)
	}
	switch scheme {
	case "postgresql":
		dbOperator, err = postgres_db_adapter.CreatePostgresDBAdapter(dbURL)
		if err != nil {
			return fmt.Errorf("failed to initialize postgres adapter: %w", err)
		}
	default:
		return fmt.Errorf("unsupported database scheme: %s", scheme)
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
		snapshotName := args.Options
		if len(snapshotName) == 0 {
			return errors.New("snapshot name must be specified")
		}
		if err := snapshotFn(snapshotName); err != nil {
			return err
		}
	}

	if args.Command == app.ListCommand {
		listItems, err := dbOperator.List()
		if err != nil {
			return err
		}
		listItems.Print()
	}

	exit(nil)

	return nil
}
