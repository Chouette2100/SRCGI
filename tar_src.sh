#! /bin/sh
cp thpoint.txt thpoint.txt.ref
filename=`date +%Y%m%d-%H%M`
#tar zcvf SRCGI_$filename.tar.gz SRCGI SRCGI_maint run.sh DBConfig.yml ServerConfig.yml excl.txt rvl.txt DenyIp.txt templates public/index.html 
#tar zcvf SRCGI_$filename.tar.gz SRCGI run.sh my_script.env DBConfig.yml ServerConfig.yml bots.yml Env.yml excl.txt rvl.txt DenyIp.txt templates public/index.html 
# tar zcvf SRCGI_$filename.tar.gz SRCGI run.sh my_script.env DBConfig.yml ServerConfig.yml bots.yml nontargetentry.yml Env.yml excl.txt rvl.txt DenyIp.txt templates public/index.html 
# tar zcvf SRCGI_$filename.tar.gz SRCGI run.sh chgLM.sh my_script.env DBConfig.yml ServerConfig.yml bots.yml nontargetentry.yml Env.yml excl.txt rvl.txt DenyIp.txt templates public/index.html 
tar zcvf SRCGI_$filename.tar.gz \
*.go \
ShowroomCGIlib/*.go \
.vscode \
go.mod \
go.sum \
run.sh \
chgLM.sh \
*.enc.yaml \
bots.yml \
nontargetentry.yml \
Env.yml \
excl.txt \
thpoint.txt.ref \
rvl.txt \
DenyIp.txt \
cidr.txt \
templates \
public/index.html \
tar_src.sh \
tar.sh
