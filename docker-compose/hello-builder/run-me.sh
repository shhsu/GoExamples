#!/bin/bash

for i in 1 2 3 4 5
do
    echo $i
	/bin/sleep 10
done

>&2 echo "Some errors, no problem"
