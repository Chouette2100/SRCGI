#!/bin/sh
f1=`ls -rt SRCGI_*.txt | tail -n 1`
grep "can't evaluate" $f1
egrep "not defined|panic" $f1
