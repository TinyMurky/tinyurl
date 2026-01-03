# Pkg Spec: Bloom Filter

## Purpose
Provides a wrapper around Redis for Bloom Filter operations to check for the probable existence of keys.

## Requirements

### Requirement: Configuration
The package MUST be configurable via environment variables.

#### Scenario: Initialization
Given valid Redis connection details in the environment
When `NewFromEnv` is called
Then a connection to the specific Redis DB for Bloom Filters is established.

### Requirement: Connection Management
The package MUST provide a way to close the connection.

#### Scenario: Close
When `Close` is called
Then the underlying Redis client connection is closed.
