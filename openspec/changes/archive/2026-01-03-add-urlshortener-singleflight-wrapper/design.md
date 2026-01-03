# Design: Domain-Specific Singleflight

## Context
The `pkg/singleflight` package is a generic infrastructure utility. The `urlshortener` domain has specific requirements for key namespacing ("shorturl:{id}") to prevent collisions if the global singleflight group is shared.

## Solution
Introduce an adapter/wrapper in `internal/urlshortener/singleflight`.

### Interface
```go
package singleflight

type Group interface {
    Do(id string, fn func() (any, error)) (model.URL, error, bool)
}
```

*Note: The return type can be strongly typed (`model.URL`) or keep it `any` depending on strictness. Since this is a domain wrapper, strong typing is preferred if possible, but `any` is flexible.* -> *Decision: Keep it `any` for now to match the underlying signature, or make it specific. Let's keep it consistent with the underlying `Do` for now but handle the key.*

### Key Prefixing
The wrapper will automatically prepend `shorturl:` to the `id` passed in.

`wrapper.Do(id, fn)` -> `underlying.Do("shorturl:" + id, fn)`
