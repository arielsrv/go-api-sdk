#!/bin/bash
.gitlab/common/git.sh
source .gitlab/common/guard.sh
check_env_variable "$GOPATH" "GOPATH"

go install gotest.tools/gotestsum@latest
"$GOPATH"/bin/gotestsum --junitfile report.xml --format testname
