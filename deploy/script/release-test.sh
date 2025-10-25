#!/usr/bin/env bash
set -e
need_start_server_shell=(
  # rpc
  im-ws-test.sh
  user-rpc-test.sh
  social-rpc-test.sh
  im-rpc-test.sh

  # api
  user-api-test.sh
  social-api-test.sh
  im-api-test.sh

  # task
  task-mq-test.sh
)

for i in ${need_start_server_shell[*]} ; do
    chmod +x $i
    ./$i
done


docker ps

docker exec -it etcd etcdctl get --prefix ""