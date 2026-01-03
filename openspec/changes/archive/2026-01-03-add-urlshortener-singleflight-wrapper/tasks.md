# Tasks

-   [ ] Create `internal/urlshortener/singleflight` package.
-   [ ] Implement `Group` struct that wraps `pkg/singleflight.Group`.
-   [ ] Implement `Do(id string, fn ...)` method that handles key prefixing.
-   [ ] Update `Handler` in `handle_get_shorturl` to use `internal/urlshortener/singleflight.Group`.
-   [ ] Update `New` in `handle_get_shorturl` to initialize the wrapper.
-   [ ] Remove manual "shorturl:" string concatenation from `ServeHTTP`.
-   [ ] Verify with tests.
