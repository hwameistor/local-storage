#! /usr/bin/env bash

export GOVC_INSECURE=1
export GOVC_RESOURCE_POOL="e2e"
export hosts="fupan-e2e-k8s-master fupan-e2e-k8s-node1 fupan-e2e-k8s-node2"
export snapshot="e2etest"

for h in $hosts; do
  if [[ `govc vm.info $h | grep poweredOn | wc -l` -eq 1 ]]; then
    govc vm.power -off -force $h
    echo -e "\033[35m === $h has been down === \033[0m"
  fi

  govc snapshot.revert -vm $h $snapshot
  echo -e "\033[35m === $h reverted to snapshot: `govc snapshot.tree -vm $h -C -D -i -d` === \033[0m"

  govc vm.power -on $h
  echo -e "\033[35m === $h: power turned on === \033[0m"
done

echo -e "\033[35m === task will end in 1m 30s === \033[0m"
for i in `seq 1 15`; do
  echo -e "\033[35m === `date  '+%Y-%m-%d %H:%M:%S'` === \033[0m"
  sleep 6s
done
