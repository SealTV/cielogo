// Package cielogo provides typed HTTP and WebSocket clients for the Cielo API.
//
// The client accepts an API key and supports concurrent HTTP operations. Contexts
// control individual calls. No paid request is retried automatically. Errors can
// be inspected with errors.Is and errors.As; response and request models live in
// api/apiv1. Portfolio V2 requires explicit private-endpoint configuration.
//
// See API_STABILITY.md for the 1.x contract and MIGRATION.md when upgrading.
package cielogo
