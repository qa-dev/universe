#!/bin/bash

echo "mode: set" > acc.out
for Dir in $(go list ./...);
do
    if [[ ${Dir} != *"/vendor/"* ]]
    then
    	Cwd=`pwd`
        returnval=`go test -coverprofile=profile.out $Dir`
        echo ${returnval}
        if [[ ${returnval} != *FAIL* ]]
        then
            if [ -f profile.out ]
            then
                cat profile.out | grep -v "mode: set" >> acc.out
            fi
        else
            exit 1
        fi
    fi

done
if [ -n "$COVERALLS_TOKEN" ]
then
    goveralls -coverprofile=acc.out -repotoken=$COVERALLS_TOKEN -service=travis-ci
fi
