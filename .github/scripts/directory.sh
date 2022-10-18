#!/bin/bash
dir_path=$1
release_version=$2
# bash check if directory exists
if [ -d $dir_path/$release_version ] ; then
	echo "Directory exists"
	echo "new_version_available=true" >> $GITHUB_OUTPUT
else 
	echo "Directory does not exists"
	echo "new_version_available=false" >> $GITHUB_OUTPUT
fi