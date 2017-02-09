#!/dgr/bin/busybox sh
set -e
. /dgr/bin/functions.sh
isLevelEnabled "debug" && set -x

cat > /etc/apt/sources.list.d/cassandra.list<<EOF
deb http://ftp.fr.debian.org/debian/ jessie main non-free contrib # needed for java8
deb http://www.apache.org/dist/cassandra/debian 30x main
deb-src http://www.apache.org/dist/cassandra/debian 30x main
deb http://ftp.debian.org/debian jessie-backports main
EOF

apt-cache policy openjdk-8-jre-headless

apt-get install -y --force-yes curl
curl https://www.apache.org/dist/cassandra/KEYS | apt-key add -

apt-get update
apt-get install -y --force-yes cassandra cassandra-tools

chown -R cassandra: /etc/cassandra
mkdir /data
chown cassandra: /data
