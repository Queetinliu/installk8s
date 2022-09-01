#!/bin/bash
set -uxo pipefail
SCRIPT=$(readlink -f "$0")
SCRIPTPATH=$(dirname "$SCRIPT")
ping -c 1 -q baidu.com >& /dev/null
if [[ $? != 0 ]];then
set -e
rm -f /etc/yum.repos.d/*.repo
cd $SCRIPTPATH/packages/baserpms
yum install -y *  # install httpd vim sshpass ansible createrepo ntfs-g3 nfs-utils pciutils
cp $SCRIPTPATH/packages/globalrpms/* /var/www/html  
set -e
sed -i 's/Listen 80/Listen 18080/' /etc/httpd/conf/httpd.conf
if [ -d /var/www/html/repodata ];then
rm -rf /var/www/html/repodata
fi
createrepo /var/www/html
selinux=$(getenforce)
if [[ $selinux != "Disabled" ]];then
setenforce 0
fi
systemctl restart httpd
systemctl enable httpd
 
fi



