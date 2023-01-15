### go-chi Users App

# Setup

- requires go 1.18 and a running postgres server

- copy .env/dev_template.env into ./env/dev.env, and fill in template values with appropriate data

then run command:

    go install

# Run App:

    go run app/app.go

# Run tests:

    cd tests
    go test -v

  To run individual or group of tests:

    go test -v -run <testname>
