#!/bin/bash
source .gitlab/common/guard.sh

check_env_variable "$GITLAB_TOKEN" "GITLAB_TOKEN"
check_env_variable "$GO_VERSION" "GO_VERSION"
check_env_variable "$TASK_VERSION" "TASK_VERSION"

/kaniko/executor \
    --context "${CI_PROJECT_DIR}" \
    --dockerfile "${CI_PROJECT_DIR}"/Dockerfile \
    --no-push \
    --build-arg GITLAB_TOKEN="${GITLAB_TOKEN}" \
    --build-arg GO_VERSION="${GO_VERSION}" \
    --build-arg TASK_VERSION="${TASK_VERSION}"
