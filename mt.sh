#!/bin/sh
f1=`ls -rt SRCGI*.txt | tail -n 1`
# /MyProject/Showroom/

multitail $f1
# multitail $f1  $f2 $f3 $f4

#if [ $? -eq 0 ]; then
#       less +F $flast
#fi

