# API Spec: Shorten URL

## Endpoint
`POST /api/v1/data/shorten`

## Purpose
Accepts a long URL and returns a shortened URL containing a unique Base62 ID.

## Requirements

### Requirement: Create Short URL
The system MUST generate or retrieve a unique identifier for the provided URL.

#### Scenario: New URL
Given a valid long URL that has not been shortened before
When a POST request is made with `long_url` in the form data
Then the system generates a new Snowflake ID
And persists the mapping to SQLite
And updates the Redis cache and Bloom Filter
And returns the new Short URL.

#### Scenario: Existing URL
Given a long URL that already exists in the database
When a POST request is made
Then the system returns the existing Short URL without creating a new ID.

#### Scenario: Invalid URL
Given a malformed URL string
When a POST request is made
Then the system returns 400 Bad Request.
