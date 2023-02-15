#!/usr/bin/bash
for i in {1..26}
do
    cd $i/
    mv main.go $i.go
    mv $i.go ../
    cd ..
    rmdir $i
done
