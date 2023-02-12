#!/bin/bash

# val is the current version
val=$(echo | awk '{print $2}' FS='=' version.go)
v=`expr $val`
newnum=`expr $val + 1`
sed -i 's/'$v'/'$newnum'/g' version.go
# $newnum is the next version
echo $newnum