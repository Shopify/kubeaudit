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
    ./kubeaudit autofix -f $f 2>> /dev/null
done;

for f in fixtures_temp/*;
do   
    echo $f >> result_file
    ./kubeaudit all -f $f 2>> result_file
done;

if (( $(grep -c 'level=error' result_file) > 0 ));
    then
    echo ERROR: Regression Test Failed
    sed -n '/^fixtures_temp/{ x; /level=error/p; d; }; /level=error/H; ${ x; /level=error/p; }' result_file
    rm result_file
    rm -r fixtures_temp
    exit 1 # terminate and indicate error
fi

rm result_file
rm -r fixtures_tem
