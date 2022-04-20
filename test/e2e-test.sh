#! /usr/bin/env bash
# simple scripts mng machine
# link hosts
export GOVC_INSECURE=1
export GOVC_RESOURCE_POOL="e2e"
export hosts="fupan-e2e-k8s-master fupan-e2e-k8s-node1 fupan-e2e-k8s-node2"
export snapshot="e2etest"
# for h in hosts; do govc vm.power -off -force $h; done
# for h in hosts; do govc snapshot.revert -vm $h "机器配置2"; done
# for h in hosts; do govc vm.power -on -force $h; done

# govc vm.info $hosts[0].Power state
# govc find . -type m -runtime.powerState poweredOn
# govc find . -type m -runtime.powerState poweredOn | xargs govc vm.info
# govc vm.info $hosts

git clone https://github.com/hwameistor/helm-charts.git test/helm-charts
echo "it is a test"
cat test/helm-charts/charts/hwameistor/values.yaml | while read line
##
do
result=$(echo $line | grep "imageRepository")
if [[ "$result" != "" ]]
then
img=${line:17:50}
fi
result=$(echo $line | grep "tag")
if [[ "$result" != "" ]]
then
hwamei=$(echo $img | grep "hwameistor")
if [[ "$hwamei" != "" ]]
then
image=$img:${line:5:50}
echo "docker pull ghcr.io/$image"
#         docker pull $image
echo "docker tag ghcr.io/$image 10.6.170.180/hwamei-e2e/$image"
echo "docker push 10.6.170.180/hwamei-e2e/$image"
fi
fi
done
##
#sed -i '/local-storage/{n;d}' test/helm-charts/charts/hwameistor/values.yaml
#sed -i '/local-storage/a \ \ \ \ tag: test' test/helm-charts/charts/hwameistor/values.yaml
# ginkgo --fail-fast --label-filter=${E2E_TESTING_LEVEL} test/e2e