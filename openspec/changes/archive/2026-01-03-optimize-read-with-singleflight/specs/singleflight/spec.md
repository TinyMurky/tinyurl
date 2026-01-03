## ADDED Requirements

### Requirement: Request Coalescing
The package MUST provide a mechanism to coalesce duplicate function calls associated with a key into a single execution.

#### Scenario: Deduplication
Given a `Group` instance
And a function `fn` that returns a value
When `Do` is called concurrently 10 times with the key "key1"
Then `fn` executes only once
And all 10 calls return the same result.

#### Scenario: Isolation
Given a `Group` instance
When `Do` is called with key "A"
And `Do` is called with key "B"
Then the two function executions are independent.

### Requirement: Configuration
The package MUST be configurable via environment variables.

#### Scenario: Initialization
Given valid environment variables for `singleflight` (if any are applicable, e.g., timeouts or disabled state)
When `NewFromEnv` is called
Then the returned `Group` reflects those settings.