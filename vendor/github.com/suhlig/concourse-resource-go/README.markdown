# Concourse Resource Interface in Go

This repo specifies an opinionated interface for Concourse resources in Go. It simplifies resource implementations, by providing

- Unmarshaling and validating the request from JSON,
- Calling the implementation of the interface, and
- Validating and marshaling the response to JSON.

The interface is generic and expects concrete types for `Source`, `Version` and `Params`. You can add additional validation rules by annotating your types with the ones implemented by [go-playground/validator](https://pkg.go.dev/github.com/go-playground/validator).

For an example that uses this interface, check the [Euro Exchange Rates Resource](https://github.com/suhlig/euro-exchange-rates-resource).
