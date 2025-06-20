#!/bin/bash

nfn=`ls -tr SRCGI.2025*|tail -1`; echo $nfn
nbe=$(basename "$nfn");echo $nbe

systemctl --user stop SRCGI
rm SRCGI
ln -s ./$nbe SRCGI
ls -l
systemctl --user start SRCGI
systemctl --user status SRCGI
