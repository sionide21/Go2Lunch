#!/bin/sh

lunch_path=$1
replacement_path=$2

echo Copying updated version...
sleep 3

echo cp $replacement_path $lunch_path
cp $replacement_path $lunch_path
chmod u+x $replacement_path
echo Updated version installed.
