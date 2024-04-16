# Ghostel

Ghostel is a database snapshot/restore tool for Postgres and MongoDB. It is architected to make additional database support simple to implement. It is intended for the local development environment, not production.

Inspired by [Stellar](https://github.com/fastmonkeys/stellar)– a PG snapshot tool written in Python. I decided to implement it in Go and make it database-agnostic.

## Features

- ✅ Supports Postgres and MongoDB
- ✅ Saves database config in directory
- ✅ Can switch between multiple saved databases
- ✅ Restores snapshots without data loss**

** A temporary copy of the original database is created, and only deleted after the restore is successful.

## Install

With Go version 1.22.1 or higher, run `make install`

## Usage

```sh
# Initialize the PG project in the current directory
gho init my_local_pg "postgresql://admin:admin@localhost/main?sslmode=disable"

# Create a snapshot
gho snapshot before_user_migration

# List snapshots
gho ls

# Enable fast restore
gho set fastRestore true

# Restore snapshot
gho restore before_user_migration

# Remove snapshot
gho rm before_user_migration

# Initialize a MongoDB project in the current directory
gho init local_mongo "mongodb://admin:admin@localhost:27017/primary?tls=false"

# View all projects (databases) in the current directory
gho status

# From here, the "local_mongo" DB is selected and all subsequent operations will affect that DB

# Select the PG database
gho select my_local_pg
```

## Faster Restore
By default, restoring a snapshot will first create a backup of the original database. Then only upon successfully restoring the snapshot will the backup be deleted.

But the initial backup step can be skipped to speed up the restore operation.

```sh
# Select the project
gho select myDb

# Enable fast restore
gho set fastRestore true

# All subsequent restores will be faster 
gho restore v1
```

## Supporting other databases

If you want to add support for other databases, just implement interfaces:
```go
type IDBOperator interface {
  Snapshot(snapshotName string) error
  Restore(snapshotName string, fast bool) error
  Delete(snapshotName string) error
  List() (List, error)
}

type IDBOperatorBuilder interface {
  ID() string
  BuildOperator(dbURL string) (IDBOperator, error)
}
```
