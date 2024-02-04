#!/bin/bash
.gitlab/common/git.sh
source .gitlab/common/guard.sh
check_env_variable "$GOPATH" "GOPATH"

go test ./... -coverprofile=coverage.txt -covermode count
go get github.com/boumenot/gocover-cobertura
go install github.com/boumenot/gocover-cobertura
"$GOPATH"/bin/gocover-cobertura <coverage.txt >coverage.xml
