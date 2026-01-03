# Infra Spec: Docker

## Purpose
Defines the containerized environment for development and debugging.

## Requirements

### Requirement: Development Environment
The system MUST provide a standard development environment via `docker-compose.yml`.

#### Scenario: Service Orchestration
When the composition is started
Then it brings up the `tinyurl` service, `redis`, and any other dependencies defined in `builders/docker-compose.yml`.

### Requirement: Debug Environment
The system MUST provide a debug-ready environment via `docker-compose.debug.yml`.

#### Scenario: Debug Attachment
When the debug composition is started
Then the `tinyurl` service runs with a configuration (e.g., Delve) allowing remote debugger attachment.
