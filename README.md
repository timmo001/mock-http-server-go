# Mock HTTP Server

This is a simple mock HTTP server written in Go that can be used to test HTTP requests and responses.

## Features

- `/echo` endpoint that returns the request body as the response body
- `/echo/details` endpoint that returns the request details as the response body in JSON format
- `/write?path=abc.json` endpoint that writes the request body to a file defined. The file path is defined in the query parameter `path`
