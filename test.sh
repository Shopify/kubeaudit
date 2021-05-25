#!/usr/bin/env bash
    
set -e

echo "Starting tests. This may take a while..."

function finish {
    if [ "$USE_KIND" != "false" ] ; then
        make test-teardown
    fi
}
trap finish EXIT

if [ "$USE_KIND" != "false" ] ; then
    make test-setup
fi

touch coverage.txt

for d in $(go list ./...); do
    go test -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done
