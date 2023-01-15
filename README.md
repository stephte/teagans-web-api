### go-chi Users App

# Setup

requires go 1.18

and a running postgres server

fill in .env/dev.env template with appropriate data

then run command:

    go install

# Run App:

    go run app/app.go

# Run tests:

    cd tests
    go test -v

  To run individual or group of tests:

    go test -v -run <testname>
