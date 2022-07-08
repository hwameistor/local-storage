### 介绍
Convertible 功能是HwameiStor中重要的运维管理功能。数据高可用是存储系统重要能力之一，Hwameistor不仅能创建高可用的 HA 数据卷，还提供非 HA 卷的高可用转化能力。通过该能力，可以将副本数为1的卷，扩展为副本数为2，进而使用 HA 的数据卷能力。

### 基本概念

**LocalVolumeGroup(LVG)**: LocalVolumeGroup（数据卷组）管理是HwameiStor中重要的一个功能。当应用Pod申请多个数据卷PVCs时，为了保证Pod能正确运行，这些数据卷必须具有某些相同属性，例如：数据卷的副本数量，副本所在的节点。通过数据卷组管理功能正确地管理这些相关联的数据卷，是HwameiStor中非常重要的能力。

### 前提条件

LocalVolumeConvert需要部署在Kubernetes系统中，需要部署应用满足下列条件：

* 支持 lvm 类型的卷
* convertible类型卷（需要在sc中增加配置项convertible: true）
  * 应用pod申请多个数据卷PVCs时，对应数据卷需要使用相同配置sc
  * 基于LocalVolume粒度转化时，所属相同LocalVolumeGroup的数据卷一并转化

### 步骤 1: 创建 convertible StorageClass

```console
$ cd ../../deploy/
$ kubectl apply -f storageclass-convertible-lvm.yaml
```

### 步骤 2: 创建 multiple PVC

```console
$ kubectl apply -f pvc-multiple-lvm.yaml
```

### 步骤 3: 部署 多数据卷 Pod

```console
$ kubectl apply -f nginx-multiple-lvm.yaml
```

### 步骤 4: 创建 转化任务

```console
cat > ./convert_lv.yaml <<- EOF
apiVersion: hwameistor.io/v1alpha1
kind: LocalVolumeConvert
metadata:
  namespace: hwameistor
  name: <localVolumeConvertName>
spec:
  volumeName: <volName>
  replicaNumber: 2
EOF
```

```console
$ kubectl apply -f ./convert_lv.yaml
```

### 步骤 5: 查看 转化状态

```console
[root@172-30-45-222 deploy]# kubectl  get LocalVolumeConvert -o yaml
apiVersion: v1
items:
- apiVersion: hwameistor.io/v1alpha1
  kind: LocalVolumeConvert
  metadata:
    annotations:
      kubectl.kubernetes.io/last-applied-configuration: |
        {"apiVersion":"hwameistor.io/v1alpha1","kind":"LocalVolumeConvert","metadata":{"annotations":{},"name":"localvolumeconvert-1","namespace":"hwameistor"},"spec":{"replicaNumber":2,"volumeName":"pvc-1a0913ac-32b9-46fe-8258-39b4e3b696a4"}}
    creationTimestamp: "2022-07-07T12:41:40Z"
    generation: 1
    name: localvolumeconvert-1
    namespace: hwameistor
    resourceVersion: "12830469"
    uid: fa333754-ec38-49bf-8211-9e6346f188e6
  spec:
    abort: false
    replicaNumber: 2
    volumeName: pvc-1a0913ac-32b9-46fe-8258-39b4e3b696a4
  status:
    state: InProgress
kind: List
metadata:
  resourceVersion: ""
  selfLink: ""
```

### 步骤 6: 查看迁移成功状态

```console
[root@172-30-45-222 deploy]# kubectl  get lvr
NAME                                              CAPACITY     NODE               STATE   SYNCED   DEVICE                                                                  AGE
pvc-1a0913ac-32b9-46fe-8258-39b4e3b696a4-9cdkkn   1073741824   172-30-45-223      Ready   true     /dev/LocalStorage_PoolHDD-HA/pvc-1a0913ac-32b9-46fe-8258-39b4e3b696a4   7m57s
pvc-1a0913ac-32b9-46fe-8258-39b4e3b696a4-jnt7km   1073741824   dce-172-30-40-61   Ready   true     /dev/LocalStorage_PoolHDD-HA/pvc-1a0913ac-32b9-46fe-8258-39b4e3b696a4   48s
pvc-d9d3ae9f-64af-44de-baad-4c69b9e0744a-7ppmrx   1073741824   172-30-45-223      Ready   true     /dev/LocalStorage_PoolHDD-HA/pvc-d9d3ae9f-64af-44de-baad-4c69b9e0744a   7m57s
pvc-d9d3ae9f-64af-44de-baad-4c69b9e0744a-sx2j4v   1073741824   dce-172-30-40-61   Ready   true     /dev/LocalStorage_PoolHDD-HA/pvc-d9d3ae9f-64af-44de-baad-4c69b9e0744a   48s
```