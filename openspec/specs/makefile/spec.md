# Infra Spec: Makefile

## Purpose
Provides a standard interface for common development tasks.

## Requirements

### Requirement: Lifecycle Management
The Makefile MUST provide commands to start and stop the application.

#### Scenario: Start Dev
When `make up` is executed
Then `docker compose` starts the standard environment.

#### Scenario: Start Debug
When `make up-debug` is executed
Then `docker compose` starts the debug environment.

#### Scenario: Stop
When `make down` is executed
Then `docker compose` stops and removes containers for the active environment.
