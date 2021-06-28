#!/bin/sh
TAG=$(git describe --tags --abbrev=0 2>/dev/null)

if [ -n "${TAG}" ]; then
  COMMITS_SINCE_TAG=$(git rev-list "${TAG}".. --count)
  if [ "${COMMITS_SINCE_TAG}" -gt 0 ]; then
    COMMIT_SHA=$(git show -s --format=%ct.%h HEAD)
    SUFFIX="+git.dev${COMMITS_SINCE_TAG}.${COMMIT_SHA}"
  fi
else
  TAG="0"
fi

echo "${TAG}${SUFFIX}"
