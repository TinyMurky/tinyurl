# Pkg Spec: Database

## Purpose
Provides a wrapper around `modernc.org/sqlite` for persistent storage.

## Requirements

### Requirement: Configuration
The package MUST be configurable via environment variables to construct the DSN.

#### Scenario: Initialization
Given valid file path/DSN in the environment
When `NewFromEnv` is called
Then a connection pool to the SQLite database is established.

### Requirement: Connection Management
The package MUST provide a way to close the connection pool.

#### Scenario: Close
When `Close` is called
Then the underlying SQL connection pool is closed.
