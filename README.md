# Mock HTTP Server

This is a mock HTTP server written in [Go](https://go.dev) that can be used to test HTTP requests.

## Features

- `/echo` endpoint that returns the request body as the response body
- `/echo/details` endpoint that returns the request details as the response body in JSON format
- `/write?path=temp/abc.json` endpoint that writes the request body to a file defined. The file path is defined in the query parameter `path`
  - Supports `multipart/form-data` and any other content type
