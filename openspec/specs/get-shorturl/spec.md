# API Spec: Get Short URL

## Endpoint
`GET /api/v1/shortUrl/{id}`

## Purpose
Redirects a client to the original long URL associated with the provided Base62 ID.
## Requirements
### Requirement: Redirect to Original URL
The system MUST retrieve the long URL and redirect the user.

#### Scenario: ID in Database (Cache Miss with Singleflight)
Given a valid Base62 ID (ex:"quZWvVVg") that exists in SQLite but not in Redis
And multiple concurrent GET requests are made for "quZWvVVg"
When the requests are processed
Then the system executes only one database query for "quZWvVVg"
And updates the Redis cache once
And returns a 301/302 redirect to the long URL for all requests.

