#!/bin/bash

source .github/env.sh

[ $rel ] || sed -i "s/buildDate .*/buildDate := \"`date +'%Y%m%d'`\"/g" date.go

BUILD="build"

rm -rf $BUILD/android \
  $BUILD/java \
  $BUILD/javac-output \
  $BUILD/src
