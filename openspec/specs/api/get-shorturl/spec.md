# API Spec: Get Short URL

## Endpoint
`GET /api/v1/shortUrl/{id}`

## Description
Redirects a client to the original long URL associated with the provided Base62 ID.

## Requirements

### Requirement: Redirect to Original URL
The system must retrieve the long URL and redirect the user.

#### Scenario: ID in Cache
Given a valid Base62 ID that exists in the Redis cache
When a GET request is made
Then the system returns a 301/302 redirect to the long URL immediately.

#### Scenario: ID in Database (Cache Miss)
Given a valid Base62 ID that exists in SQLite but not in Redis
When a GET request is made
Then the system retrieves the URL from the DB
And updates the Redis cache
And returns a 301/302 redirect to the long URL.

#### Scenario: ID Not Found (Bloom Filter)
Given a Base62 ID that is not in the Bloom Filter
When a GET request is made
Then the system returns 404 Not Found immediately without checking the DB.

#### Scenario: ID Not Found (Database)
Given a Base62 ID that passes the Bloom Filter but is not in the DB
When a GET request is made
Then the system returns 404 Not Found.
