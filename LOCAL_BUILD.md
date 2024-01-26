# Local Build

Use this set of commands to perform a local build for tesing.

```bash
SEMVER=v0.0.999; echo ${SEMVER}
BUILD_DATE=$(gdate --utc +%FT%T.%3NZ); echo ${BUILD_DATE}
GIT_COMMIT=$(git rev-parse HEAD); echo ${GIT_COMMIT}

MODULE_NAME=daily-notes
go build -ldflags "-X ${MODULE_NAME}/cmd.semVer=${SEMVER} -X ${MODULE_NAME}/cmd.buildDate=${BUILD_DATE} -X ${MODULE_NAME}/cmd.gitCommit=${GIT_COMMIT} -X ${MODULE_NAME}/cmd.gitRef=/refs/tags/${SEMVER}" && \
./daily-notes version | jq .

if [[ -d ~/tbin ]]; then
  cp ./daily-notes ~/tbin/daily-notes
fi
```

