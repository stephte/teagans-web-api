### go-chi Users App

# Setup

- requires go 1.18 and a running postgres server

- copy ".template.env" into ".env" and update it with appropriate data

# Run App:

    go run app/app.go

# Run tests:
- first, set up new postgres DB for testing

- copy .template.env into .testenv, fill with appropriate DB data (omit email SMTP data)

## Then run:

    cd tests
    go test -v

  To run individual or group of tests:

    go test -v -run <testname>
