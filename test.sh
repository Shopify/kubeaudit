#!/usr/bin/env bash

set -e
touch coverage.txt

for d in $(go list ./... | grep -v vendor); do
    go test -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out > coverage.txt
        rm profile.out
    fi
done

cp -r fixtures/ fixtures_temp

for f in fixtures_temp/*;
do
    ./kubeaudit autofix -f $f -v ERROR
done;

for f in fixtures_temp/*;
do
    ./kubeaudit all -f $f -v ERROR 2>> result_file
done;

if (( $(grep -c 'level=error' result_file) > 0 ));
    then
    echo ERROR: Regression Test Failed
    exit 1 # terminate and indicate error
fi

rm -rf result_file
rm -rf fixtures_temp
