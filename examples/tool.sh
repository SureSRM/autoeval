#!/bin/sh

echo -n "Name: "
read name
echo "Hello, ${name}"

echo -n "Age: "
read age
echo "You are ${age} years old"

echo "Your arg was: $1"
