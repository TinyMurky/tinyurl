# Pkg Spec: Database

## Description
Provides a wrapper around `modernc.org/sqlite` for persistent storage.

## Requirements

### Requirement: Configuration
The package must be configurable via environment variables to construct the DSN.

#### Scenario: Initialization
Given valid file path/DSN in the environment
When `NewFromEnv` is called
Then a connection pool to the SQLite database is established.

### Requirement: Connection Management
#### Scenario: Close
When `Close` is called
Then the underlying SQL connection pool is closed.
