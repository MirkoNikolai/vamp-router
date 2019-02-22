#!/bin/bash
set -e

# Set the package name

PACKAGE='vamp-router'

# Set the version
VERSION=$1
if [ -z $VERSION ]; then
    echo "Enter a version"
    exit 1
fi

#clear the target/dist dir and recreate it
rm -rf ./target/dist
mkdir -p ./target/dist

# build the app for linux/i386 and create a zip with necessary artifacts
for GOOS in darwin linux windows; do
  for GOARCH in 386 amd64; do
    
    DISTRIBUTABLE=${PACKAGE}_${VERSION}_${GOOS}_${GOARCH}

    echo "Building $DISTRIBUTABLE"
    
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    
    # remove and create a tmp dir for collecting and zipping per package
    rm -rf ./target/dist/tmp
    mkdir -p ./target/dist/tmp

    go build

    if [ "${GOOS}" == "windows" ]; then
        mv ${PACKAGE}.exe ./target/dist/tmp
    else
        mv ${PACKAGE} ./target/dist/tmp
        chmod +x ./target/dist/tmp/${PACKAGE}
    fi

    cp -r ./configuration ./target/dist/tmp
    cp -r ./examples ./target/dist/tmp
    cd ./target/dist/tmp
    zip -r ${DISTRIBUTABLE}.zip *

    # move the package into the target/dist dir
    mv ${DISTRIBUTABLE}.zip ../
    cd ../../../
  done
done

# remove the last tmp dir
rm -rf ./target/dist/tmp

# change to dist files

cd target/dist

