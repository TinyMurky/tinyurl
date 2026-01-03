# urlshortener-singleflight Specification

## Purpose
TBD - created by archiving change add-urlshortener-singleflight-wrapper. Update Purpose after archive.
## Requirements
### Requirement: Domain Key Namespacing
The internal singleflight wrapper MUST automatically namespace keys to prevent collisions with other domain entities.

#### Scenario: ShortURL Prefix
Given a wrapper instance initialized with a "shorturl:" prefix (or hardcoded)
When `Do` is called with ID "abc"
Then the underlying singleflight group receives the key "shorturl:abc".

### Requirement: Typed Execution (Optional)
The wrapper MUST delegate execution to the underlying group.

#### Scenario: Execution Delegation
Given a wrapper instance
When `Do` is called with a function
Then the function is executed via the underlying group
And the result is returned to the caller.

