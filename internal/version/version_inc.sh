#!/bin/bash

val=$(echo | awk '{print $2}' FS='=' version.go)
v=`expr $val`
newnum=`expr $val + 1`
sed -i 's/'$v'/'$newnum'/g' version.go
echo $newnum