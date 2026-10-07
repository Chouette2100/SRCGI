#! /bin/sh
filename=`date +%Y%m%d-%H%M`
tar zcvf SRCGIenv_$filename.tar.gz \
DBConfig_K10.enc.yaml \
ServerConfig_HTTPPORT.enc.yaml \
bots.yml \
nontargetentry.yml \
Env.yml \
excl.txt \
thpoint.txt.ref \
rvl.txt \
DenyIp.txt \
cidr.txt \
templates \
public/index.html 