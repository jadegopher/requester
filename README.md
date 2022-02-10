# Requester

CLI tool that has ability to send requests to http and print MD5 hashes of responses.

## How to build

To build this tool you should execute command sequence in root directory of project:

````
cd cmd

go build
````
## Example of usage

```
./cmd facebook.com amazon.com apple.com google.com
```

You could specify count of workers by using -parallel flag. By default, application runs 10 workers

```
./cmd -parallel 2 facebook.com amazon.com apple.com google.com
```
## Testing

To run tests you could execute ``go test ./...`` in root directory of this project

The most part of business logic covered with unit tests. 
All unit tests has written by using classical school manner except one.

I've written unit test for worker package in London school manner. 
London school of unit testing assumes mocking for all dependencies. 
It leads to tests that are not resistant to refactoring because we know about details of implementation.

To re-generate mocks execute ``go generate`` command in root directory of this project. 
