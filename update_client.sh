#!/bin/sh

lunch_path=$1
replacement_path=$2

echo Copying updated version...

cp $replacement_path $lunch_path
chmod u+x $lunch_path
echo Updated version installed.
