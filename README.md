### Teagan's YouTube downloader App

# Setup

- requires go 1.18 and a running postgres server

- copy ".template.env" into ".env" and update it with appropriate data

# Run App:

    go run app/app.go

# Run tests:

    cd tests
    go test -v

  To run individual or group of tests:

    go test -v -run <testname>
