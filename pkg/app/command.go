package app

import "fmt"

type Command string

const HelpCommand = "help"
const InitCommand = "init"
const SelectCommand = "select"
const StatusCommand = "status"
const SnapshotCommand = "snapshot"
const RestoreCommand = "restore"
const RemoveCommand = "rm"
const ListCommand = "ls"

type CommandInfo struct {
	Command     Command
	Template    string
	Description string
}

func (c CommandInfo) Row() []string {
	return []string{c.Template, c.Description}
}

var AllCommands = func(executable string) []CommandInfo {
	return []CommandInfo{
		{HelpCommand, fmt.Sprintf("%s %s", executable, HelpCommand), "Show the list of commands"},
		{InitCommand, fmt.Sprintf("%s %s <project_name> <database_name>", executable, InitCommand), "Initialize project in current directory"},
		{SelectCommand, fmt.Sprintf("%s %s <project_name>", executable, SelectCommand), "Select a project"},
		{StatusCommand, fmt.Sprintf("%s %s", executable, StatusCommand), "Show all projects in current directory"},
		{SnapshotCommand, fmt.Sprintf("%s %s <snapshot_name>", executable, SnapshotCommand), "Create a snapshot in the selected project"},
		{RestoreCommand, fmt.Sprintf("%s %s <snapshot_name>", executable, RestoreCommand), "Restore a snapshot in the selected project"},
		{RemoveCommand, fmt.Sprintf("%s %s <snapshot_name>", executable, RemoveCommand), "Remove a snapshot in the selected project"},
		{ListCommand, fmt.Sprintf("%s %s", executable, ListCommand), "List all snapshots in the selected project"},
	}
}
