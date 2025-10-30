# CI/CD Configuration for CGO-Disabled Testing

## GitHub Actions Example (.github/workflows/test.yml)

```yaml
name: Test Suite

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.25.x]

    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ matrix.go-version }}

    - name: Verify Go version
      run: go version

    - name: Install dependencies
      run: go mod tidy

    - name: Test with CGO disabled (platform-independent)
      run: |
        export CGO_ENABLED=0
        go test ./... -v

    - name: Test with CGO enabled (for comparison)
      run: |
        export CGO_ENABLED=1
        go test ./... -v

    - name: Build for multiple platforms
      run: |
        echo "Building for Linux..."
        GOOS=linux CGO_ENABLED=0 go build -o free2free-linux .
        echo "Building for Windows..."
        GOOS=windows CGO_ENABLED=0 go build -o free2free.exe .
        echo "Building for macOS..."
        GOOS=darwin CGO_ENABLED=0 go build -o free2free-mac .
```

## GitLab CI Example (.gitlab-ci.yml)

```yaml
stages:
  - test
  - build

variables:
  GO_VERSION: "1.25"

test_cgo_disabled:
  stage: test
  image: golang:${GO_VERSION}
  script:
    - go mod tidy
    - export CGO_ENABLED=0
    - go test ./... -v
  artifacts:
    reports:
      junit: junit-report.xml

test_cgo_enabled:
  stage: test
  image: golang:${GO_VERSION}
  script:
    - export CGO_ENABLED=1
    - go test ./... -v

build:
  stage: build
  image: golang:${GO_VERSION}
  script:
    - export CGO_ENABLED=0
    - go build -o free2free .
  artifacts:
    paths:
      - free2free
```

## Docker-based Testing

For containerized builds and tests:

```Dockerfile
FROM golang:1.25-alpine AS builder

# Install ca-certificates for HTTPS requests
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY . .

# Build the application with CGO disabled
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o free2free .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/free2free .

# Run tests in container (optional step)
RUN CGO_ENABLED=0 go test ./tests/... -v || echo "Tests completed with exit code $?"

CMD ["./free2free"]
```

## Jenkins Pipeline Example

```groovy
pipeline {
    agent any

    environment {
        GOVERSION = '1.25'
        GOPATH = "${WORKSPACE}"
    }

    stages {
        stage('Setup') {
            steps {
                sh 'go mod tidy'
            }
        }
        
        stage('Test with CGO Disabled') {
            steps {
                sh '''
                    export CGO_ENABLED=0
                    go test ./tests/... -v
                '''
            }
        }
        
        stage('Test with CGO Enabled') {
            steps {
                sh '''
                    export CGO_ENABLED=1
                    go test ./tests/... -v
                '''
            }
        }
        
        stage('Build') {
            steps {
                sh '''
                    export CGO_ENABLED=0
                    go build -o free2free .
                '''
            }
        }
    }
}
```

## Key Configuration Points

1. Set `CGO_ENABLED=0` in the build/test environment
2. Ensure the `modernc.org/sqlite` import is present in your codebase
3. Test across multiple platforms to verify portability
4. Monitor performance differences between CGO-enabled and disabled builds
5. Ensure all existing functionality remains intact