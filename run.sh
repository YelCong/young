#!/bin/bash

echo hello
exit

#etcd

### pull etcd
docker pull bitnami/etcd

### start etcd
docker run -d --name etcd-server \
    --publish 2379:2379 \
    --publish 2380:2380 \
    --env ALLOW_NONE_AUTHENTICATION=yes \
    --env ETCD_ADVERTISE_CLIENT_URLS=http://etcd-server:2379 \
    bitnami/etcd:latest


### pull mongoDB
docker pull mongo

### start mongo
docker run --name=mongodb \
 -v /usr/local/var/mongo/data:/data/db \
 -v /usr/local/var/mongo/backup:/data/backup \
 -p 27017:27017 \
 -d mongo:latest --auth

docker exec -it mongodb /bin/bash
#mongo admin
#db.createUser({user:'admin',pwd:'123456',roles:[{role:'root',db:'admin'}],})
#db.auth('admin','123456')
#mongodb://admin:123456@127.0.0.1:27017
#mongo --host 127.0.0.1 --port 27017 -u admin -p 123456


curl -X POST 'http://127.0.0.1:8070/job/save' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'job={"name":"job1","command":"echo hello;","cronExpr":"5/* * * * * *"}'


//etcd
etcdctl watch "/cron/killer" --prefix
etcdctl watch --prefix  /cron/killer/