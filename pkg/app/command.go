package app

import "fmt"

type Command string

const VersionCommand = "version"
const HelpCommand = "help"
const InitCommand = "init"
const SelectCommand = "select"
const SetCommand = "set"
const StatusCommand = "status"
const SnapshotCommand = "snapshot"
const RestoreCommand = "restore"
const DeleteCommand = "rm"
const ListCommand = "ls"

type CommandInfo struct {
	Template    string
	Description string
}

func (c CommandInfo) Row() []string {
	return []string{c.Template, c.Description}
}

var AllCommands = func(executable string) []CommandInfo {
	return []CommandInfo{
		{fmt.Sprintf("%s %s", executable, VersionCommand), "Show the version of the program"},
		{fmt.Sprintf("%s %s", executable, HelpCommand), "Show the list of commands"},
		{fmt.Sprintf("%s %s <project_name> <database_name>", executable, InitCommand), "Initialize project in current directory"},
		{fmt.Sprintf("%s %s <project_name>", executable, SelectCommand), "Select a project"},
		{fmt.Sprintf("%s %s", executable, StatusCommand), "Show all projects in current directory"},
		{fmt.Sprintf("%s %s <snapshot_name>", executable, SnapshotCommand), "Create a snapshot in the selected project"},
		{fmt.Sprintf("%s %s <snapshot_name>", executable, RestoreCommand), "Restore a snapshot in the selected project"},
		{fmt.Sprintf("%s %s <snapshot_name>", executable, DeleteCommand), "Delete a snapshot in the selected project"},
		{fmt.Sprintf("%s %s", executable, ListCommand), "List all snapshots in the selected project"},
		{fmt.Sprintf("%s %s", executable, SetCommand), "Sets a configuration value on the selected project"},
	}
}
