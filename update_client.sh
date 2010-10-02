#!/bin/sh

lunch_path=$1
replacement_path=$2

echo Copying updated version...
sleep 3

cp $replacement_path $lunch_path
echo Updated version installed.
