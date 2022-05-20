package controller

import (
	"github.com/hwameistor/local-storage/pkg/apis"
	"github.com/hwameistor/local-storage/pkg/apis/hwameistor/v1alpha1"
	apisv1alpha1 "github.com/hwameistor/local-storage/pkg/apis/hwameistor/v1alpha1"
	"github.com/hwameistor/local-storage/pkg/common"
	"github.com/hwameistor/local-storage/pkg/member/controller/scheduler"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		name           string
		namespace      string
		cli            client.Client
		scheme         *runtime.Scheme
		informersCache cache.Cache
		systemConfig   v1alpha1.SystemConfig
	}
	tests := []struct {
		name    string
		args    args
		want    apis.ControllerManager
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.name, tt.args.namespace, tt.args.cli, tt.args.scheme, tt.args.informersCache, tt.args.systemConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_manager_ReconcileNode(t *testing.T) {
	type fields struct {
		name                   string
		namespace              string
		apiClient              client.Client
		informersCache         cache.Cache
		scheme                 *runtime.Scheme
		volumeScheduler        scheduler.Scheduler
		nodeTaskQueue          *common.TaskQueue
		k8sNodeTaskQueue       *common.TaskQueue
		volumeTaskQueue        *common.TaskQueue
		volumeExpandTaskQueue  *common.TaskQueue
		volumeMigrateTaskQueue *common.TaskQueue
		volumeConvertTaskQueue *common.TaskQueue
		localNodes             map[string]apisv1alpha1.State
		logger                 *log.Entry
	}
	type args struct {
		node *apisv1alpha1.LocalStorageNode
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &manager{
				name:                   tt.fields.name,
				namespace:              tt.fields.namespace,
				apiClient:              tt.fields.apiClient,
				informersCache:         tt.fields.informersCache,
				scheme:                 tt.fields.scheme,
				volumeScheduler:        tt.fields.volumeScheduler,
				nodeTaskQueue:          tt.fields.nodeTaskQueue,
				k8sNodeTaskQueue:       tt.fields.k8sNodeTaskQueue,
				volumeTaskQueue:        tt.fields.volumeTaskQueue,
				volumeExpandTaskQueue:  tt.fields.volumeExpandTaskQueue,
				volumeMigrateTaskQueue: tt.fields.volumeMigrateTaskQueue,
				volumeConvertTaskQueue: tt.fields.volumeConvertTaskQueue,
				localNodes:             tt.fields.localNodes,
				logger:                 tt.fields.logger,
			}
			m.ReconcileNode(tt.args.node)
		})
	}
}

func Test_manager_ReconcileVolume(t *testing.T) {
	type fields struct {
		name                   string
		namespace              string
		apiClient              client.Client
		informersCache         cache.Cache
		scheme                 *runtime.Scheme
		volumeScheduler        scheduler.Scheduler
		nodeTaskQueue          *common.TaskQueue
		k8sNodeTaskQueue       *common.TaskQueue
		volumeTaskQueue        *common.TaskQueue
		volumeExpandTaskQueue  *common.TaskQueue
		volumeMigrateTaskQueue *common.TaskQueue
		volumeConvertTaskQueue *common.TaskQueue
		localNodes             map[string]apisv1alpha1.State
		logger                 *log.Entry
	}
	type args struct {
		vol *apisv1alpha1.LocalVolume
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &manager{
				name:                   tt.fields.name,
				namespace:              tt.fields.namespace,
				apiClient:              tt.fields.apiClient,
				informersCache:         tt.fields.informersCache,
				scheme:                 tt.fields.scheme,
				volumeScheduler:        tt.fields.volumeScheduler,
				nodeTaskQueue:          tt.fields.nodeTaskQueue,
				k8sNodeTaskQueue:       tt.fields.k8sNodeTaskQueue,
				volumeTaskQueue:        tt.fields.volumeTaskQueue,
				volumeExpandTaskQueue:  tt.fields.volumeExpandTaskQueue,
				volumeMigrateTaskQueue: tt.fields.volumeMigrateTaskQueue,
				volumeConvertTaskQueue: tt.fields.volumeConvertTaskQueue,
				localNodes:             tt.fields.localNodes,
				logger:                 tt.fields.logger,
			}
			m.ReconcileVolume(tt.args.vol)
		})
	}
}

func Test_manager_ReconcileVolumeConvert(t *testing.T) {
	type fields struct {
		name                   string
		namespace              string
		apiClient              client.Client
		informersCache         cache.Cache
		scheme                 *runtime.Scheme
		volumeScheduler        scheduler.Scheduler
		nodeTaskQueue          *common.TaskQueue
		k8sNodeTaskQueue       *common.TaskQueue
		volumeTaskQueue        *common.TaskQueue
		volumeExpandTaskQueue  *common.TaskQueue
		volumeMigrateTaskQueue *common.TaskQueue
		volumeConvertTaskQueue *common.TaskQueue
		localNodes             map[string]apisv1alpha1.State
		logger                 *log.Entry
	}
	type args struct {
		convert *apisv1alpha1.LocalVolumeConvert
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &manager{
				name:                   tt.fields.name,
				namespace:              tt.fields.namespace,
				apiClient:              tt.fields.apiClient,
				informersCache:         tt.fields.informersCache,
				scheme:                 tt.fields.scheme,
				volumeScheduler:        tt.fields.volumeScheduler,
				nodeTaskQueue:          tt.fields.nodeTaskQueue,
				k8sNodeTaskQueue:       tt.fields.k8sNodeTaskQueue,
				volumeTaskQueue:        tt.fields.volumeTaskQueue,
				volumeExpandTaskQueue:  tt.fields.volumeExpandTaskQueue,
				volumeMigrateTaskQueue: tt.fields.volumeMigrateTaskQueue,
				volumeConvertTaskQueue: tt.fields.volumeConvertTaskQueue,
				localNodes:             tt.fields.localNodes,
				logger:                 tt.fields.logger,
			}
			m.ReconcileVolumeConvert(tt.args.convert)
		})
	}
}

func Test_manager_ReconcileVolumeExpand(t *testing.T) {
	type fields struct {
		name                   string
		namespace              string
		apiClient              client.Client
		informersCache         cache.Cache
		scheme                 *runtime.Scheme
		volumeScheduler        scheduler.Scheduler
		nodeTaskQueue          *common.TaskQueue
		k8sNodeTaskQueue       *common.TaskQueue
		volumeTaskQueue        *common.TaskQueue
		volumeExpandTaskQueue  *common.TaskQueue
		volumeMigrateTaskQueue *common.TaskQueue
		volumeConvertTaskQueue *common.TaskQueue
		localNodes             map[string]apisv1alpha1.State
		logger                 *log.Entry
	}
	type args struct {
		expand *apisv1alpha1.LocalVolumeExpand
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &manager{
				name:                   tt.fields.name,
				namespace:              tt.fields.namespace,
				apiClient:              tt.fields.apiClient,
				informersCache:         tt.fields.informersCache,
				scheme:                 tt.fields.scheme,
				volumeScheduler:        tt.fields.volumeScheduler,
				nodeTaskQueue:          tt.fields.nodeTaskQueue,
				k8sNodeTaskQueue:       tt.fields.k8sNodeTaskQueue,
				volumeTaskQueue:        tt.fields.volumeTaskQueue,
				volumeExpandTaskQueue:  tt.fields.volumeExpandTaskQueue,
				volumeMigrateTaskQueue: tt.fields.volumeMigrateTaskQueue,
				volumeConvertTaskQueue: tt.fields.volumeConvertTaskQueue,
				localNodes:             tt.fields.localNodes,
				logger:                 tt.fields.logger,
			}
			m.ReconcileVolumeExpand(tt.args.expand)
		})
	}
}

func Test_manager_ReconcileVolumeGroup(t *testing.T) {
	type fields struct {
		name                   string
		namespace              string
		apiClient              client.Client
		informersCache         cache.Cache
		scheme                 *runtime.Scheme
		volumeScheduler        scheduler.Scheduler
		nodeTaskQueue          *common.TaskQueue
		k8sNodeTaskQueue       *common.TaskQueue
		volumeTaskQueue        *common.TaskQueue
		volumeExpandTaskQueue  *common.TaskQueue
		volumeMigrateTaskQueue *common.TaskQueue
		volumeConvertTaskQueue *common.TaskQueue
		localNodes             map[string]apisv1alpha1.State
		logger                 *log.Entry
	}
	type args struct {
		convert *apisv1alpha1.LocalVolumeGroup
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &manager{
				name:                   tt.fields.name,
				namespace:              tt.fields.namespace,
				apiClient:              tt.fields.apiClient,
				informersCache:         tt.fields.informersCache,
				scheme:                 tt.fields.scheme,
				volumeScheduler:        tt.fields.volumeScheduler,
				nodeTaskQueue:          tt.fields.nodeTaskQueue,
				k8sNodeTaskQueue:       tt.fields.k8sNodeTaskQueue,
				volumeTaskQueue:        tt.fields.volumeTaskQueue,
				volumeExpandTaskQueue:  tt.fields.volumeExpandTaskQueue,
				volumeMigrateTaskQueue: tt.fields.volumeMigrateTaskQueue,
				volumeConvertTaskQueue: tt.fields.volumeConvertTaskQueue,
				localNodes:             tt.fields.localNodes,
				logger:                 tt.fields.logger,
			}
			m.ReconcileVolumeGroup(tt.args.convert)
		})
	}
}

func Test_manager_ReconcileVolumeMigrate(t *testing.T) {
	type fields struct {
		name                   string
		namespace              string
		apiClient              client.Client
		informersCache         cache.Cache
		scheme                 *runtime.Scheme
		volumeScheduler        scheduler.Scheduler
		nodeTaskQueue          *common.TaskQueue
		k8sNodeTaskQueue       *common.TaskQueue
		volumeTaskQueue        *common.TaskQueue
		volumeExpandTaskQueue  *common.TaskQueue
		volumeMigrateTaskQueue *common.TaskQueue
		volumeConvertTaskQueue *common.TaskQueue
		localNodes             map[string]apisv1alpha1.State
		logger                 *log.Entry
	}
	type args struct {
		expand *apisv1alpha1.LocalVolumeMigrate
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &manager{
				name:                   tt.fields.name,
				namespace:              tt.fields.namespace,
				apiClient:              tt.fields.apiClient,
				informersCache:         tt.fields.informersCache,
				scheme:                 tt.fields.scheme,
				volumeScheduler:        tt.fields.volumeScheduler,
				nodeTaskQueue:          tt.fields.nodeTaskQueue,
				k8sNodeTaskQueue:       tt.fields.k8sNodeTaskQueue,
				volumeTaskQueue:        tt.fields.volumeTaskQueue,
				volumeExpandTaskQueue:  tt.fields.volumeExpandTaskQueue,
				volumeMigrateTaskQueue: tt.fields.volumeMigrateTaskQueue,
				volumeConvertTaskQueue: tt.fields.volumeConvertTaskQueue,
				localNodes:             tt.fields.localNodes,
				logger:                 tt.fields.logger,
			}
			m.ReconcileVolumeMigrate(tt.args.expand)
		})
	}
}

func Test_manager_Run(t *testing.T) {
	type fields struct {
		name                   string
		namespace              string
		apiClient              client.Client
		informersCache         cache.Cache
		scheme                 *runtime.Scheme
		volumeScheduler        scheduler.Scheduler
		nodeTaskQueue          *common.TaskQueue
		k8sNodeTaskQueue       *common.TaskQueue
		volumeTaskQueue        *common.TaskQueue
		volumeExpandTaskQueue  *common.TaskQueue
		volumeMigrateTaskQueue *common.TaskQueue
		volumeConvertTaskQueue *common.TaskQueue
		localNodes             map[string]apisv1alpha1.State
		logger                 *log.Entry
	}
	type args struct {
		stopCh <-chan struct{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &manager{
				name:                   tt.fields.name,
				namespace:              tt.fields.namespace,
				apiClient:              tt.fields.apiClient,
				informersCache:         tt.fields.informersCache,
				scheme:                 tt.fields.scheme,
				volumeScheduler:        tt.fields.volumeScheduler,
				nodeTaskQueue:          tt.fields.nodeTaskQueue,
				k8sNodeTaskQueue:       tt.fields.k8sNodeTaskQueue,
				volumeTaskQueue:        tt.fields.volumeTaskQueue,
				volumeExpandTaskQueue:  tt.fields.volumeExpandTaskQueue,
				volumeMigrateTaskQueue: tt.fields.volumeMigrateTaskQueue,
				volumeConvertTaskQueue: tt.fields.volumeConvertTaskQueue,
				localNodes:             tt.fields.localNodes,
				logger:                 tt.fields.logger,
			}
			m.Run(tt.args.stopCh)
		})
	}
}

func Test_manager_handleK8sNodeUpdatedEvent(t *testing.T) {
	type fields struct {
		name                   string
		namespace              string
		apiClient              client.Client
		informersCache         cache.Cache
		scheme                 *runtime.Scheme
		volumeScheduler        scheduler.Scheduler
		nodeTaskQueue          *common.TaskQueue
		k8sNodeTaskQueue       *common.TaskQueue
		volumeTaskQueue        *common.TaskQueue
		volumeExpandTaskQueue  *common.TaskQueue
		volumeMigrateTaskQueue *common.TaskQueue
		volumeConvertTaskQueue *common.TaskQueue
		localNodes             map[string]apisv1alpha1.State
		logger                 *log.Entry
	}
	type args struct {
		oldObj interface{}
		newObj interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &manager{
				name:                   tt.fields.name,
				namespace:              tt.fields.namespace,
				apiClient:              tt.fields.apiClient,
				informersCache:         tt.fields.informersCache,
				scheme:                 tt.fields.scheme,
				volumeScheduler:        tt.fields.volumeScheduler,
				nodeTaskQueue:          tt.fields.nodeTaskQueue,
				k8sNodeTaskQueue:       tt.fields.k8sNodeTaskQueue,
				volumeTaskQueue:        tt.fields.volumeTaskQueue,
				volumeExpandTaskQueue:  tt.fields.volumeExpandTaskQueue,
				volumeMigrateTaskQueue: tt.fields.volumeMigrateTaskQueue,
				volumeConvertTaskQueue: tt.fields.volumeConvertTaskQueue,
				localNodes:             tt.fields.localNodes,
				logger:                 tt.fields.logger,
			}
			m.handleK8sNodeUpdatedEvent(tt.args.oldObj, tt.args.newObj)
		})
	}
}

func Test_manager_handleVolumeCRDDeletedEvent(t *testing.T) {
	type fields struct {
		name                   string
		namespace              string
		apiClient              client.Client
		informersCache         cache.Cache
		scheme                 *runtime.Scheme
		volumeScheduler        scheduler.Scheduler
		nodeTaskQueue          *common.TaskQueue
		k8sNodeTaskQueue       *common.TaskQueue
		volumeTaskQueue        *common.TaskQueue
		volumeExpandTaskQueue  *common.TaskQueue
		volumeMigrateTaskQueue *common.TaskQueue
		volumeConvertTaskQueue *common.TaskQueue
		localNodes             map[string]apisv1alpha1.State
		logger                 *log.Entry
	}
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &manager{
				name:                   tt.fields.name,
				namespace:              tt.fields.namespace,
				apiClient:              tt.fields.apiClient,
				informersCache:         tt.fields.informersCache,
				scheme:                 tt.fields.scheme,
				volumeScheduler:        tt.fields.volumeScheduler,
				nodeTaskQueue:          tt.fields.nodeTaskQueue,
				k8sNodeTaskQueue:       tt.fields.k8sNodeTaskQueue,
				volumeTaskQueue:        tt.fields.volumeTaskQueue,
				volumeExpandTaskQueue:  tt.fields.volumeExpandTaskQueue,
				volumeMigrateTaskQueue: tt.fields.volumeMigrateTaskQueue,
				volumeConvertTaskQueue: tt.fields.volumeConvertTaskQueue,
				localNodes:             tt.fields.localNodes,
				logger:                 tt.fields.logger,
			}
			m.handleVolumeCRDDeletedEvent(tt.args.obj)
		})
	}
}

func Test_manager_handleVolumeExpandCRDDeletedEvent(t *testing.T) {
	type fields struct {
		name                   string
		namespace              string
		apiClient              client.Client
		informersCache         cache.Cache
		scheme                 *runtime.Scheme
		volumeScheduler        scheduler.Scheduler
		nodeTaskQueue          *common.TaskQueue
		k8sNodeTaskQueue       *common.TaskQueue
		volumeTaskQueue        *common.TaskQueue
		volumeExpandTaskQueue  *common.TaskQueue
		volumeMigrateTaskQueue *common.TaskQueue
		volumeConvertTaskQueue *common.TaskQueue
		localNodes             map[string]apisv1alpha1.State
		logger                 *log.Entry
	}
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &manager{
				name:                   tt.fields.name,
				namespace:              tt.fields.namespace,
				apiClient:              tt.fields.apiClient,
				informersCache:         tt.fields.informersCache,
				scheme:                 tt.fields.scheme,
				volumeScheduler:        tt.fields.volumeScheduler,
				nodeTaskQueue:          tt.fields.nodeTaskQueue,
				k8sNodeTaskQueue:       tt.fields.k8sNodeTaskQueue,
				volumeTaskQueue:        tt.fields.volumeTaskQueue,
				volumeExpandTaskQueue:  tt.fields.volumeExpandTaskQueue,
				volumeMigrateTaskQueue: tt.fields.volumeMigrateTaskQueue,
				volumeConvertTaskQueue: tt.fields.volumeConvertTaskQueue,
				localNodes:             tt.fields.localNodes,
				logger:                 tt.fields.logger,
			}
			m.handleVolumeExpandCRDDeletedEvent(tt.args.obj)
		})
	}
}

func Test_manager_handleVolumeMigrateCRDDeletedEvent(t *testing.T) {
	type fields struct {
		name                   string
		namespace              string
		apiClient              client.Client
		informersCache         cache.Cache
		scheme                 *runtime.Scheme
		volumeScheduler        scheduler.Scheduler
		nodeTaskQueue          *common.TaskQueue
		k8sNodeTaskQueue       *common.TaskQueue
		volumeTaskQueue        *common.TaskQueue
		volumeExpandTaskQueue  *common.TaskQueue
		volumeMigrateTaskQueue *common.TaskQueue
		volumeConvertTaskQueue *common.TaskQueue
		localNodes             map[string]apisv1alpha1.State
		logger                 *log.Entry
	}
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &manager{
				name:                   tt.fields.name,
				namespace:              tt.fields.namespace,
				apiClient:              tt.fields.apiClient,
				informersCache:         tt.fields.informersCache,
				scheme:                 tt.fields.scheme,
				volumeScheduler:        tt.fields.volumeScheduler,
				nodeTaskQueue:          tt.fields.nodeTaskQueue,
				k8sNodeTaskQueue:       tt.fields.k8sNodeTaskQueue,
				volumeTaskQueue:        tt.fields.volumeTaskQueue,
				volumeExpandTaskQueue:  tt.fields.volumeExpandTaskQueue,
				volumeMigrateTaskQueue: tt.fields.volumeMigrateTaskQueue,
				volumeConvertTaskQueue: tt.fields.volumeConvertTaskQueue,
				localNodes:             tt.fields.localNodes,
				logger:                 tt.fields.logger,
			}
			m.handleVolumeMigrateCRDDeletedEvent(tt.args.obj)
		})
	}
}

func Test_manager_setupInformers(t *testing.T) {
	type fields struct {
		name                   string
		namespace              string
		apiClient              client.Client
		informersCache         cache.Cache
		scheme                 *runtime.Scheme
		volumeScheduler        scheduler.Scheduler
		nodeTaskQueue          *common.TaskQueue
		k8sNodeTaskQueue       *common.TaskQueue
		volumeTaskQueue        *common.TaskQueue
		volumeExpandTaskQueue  *common.TaskQueue
		volumeMigrateTaskQueue *common.TaskQueue
		volumeConvertTaskQueue *common.TaskQueue
		localNodes             map[string]apisv1alpha1.State
		logger                 *log.Entry
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &manager{
				name:                   tt.fields.name,
				namespace:              tt.fields.namespace,
				apiClient:              tt.fields.apiClient,
				informersCache:         tt.fields.informersCache,
				scheme:                 tt.fields.scheme,
				volumeScheduler:        tt.fields.volumeScheduler,
				nodeTaskQueue:          tt.fields.nodeTaskQueue,
				k8sNodeTaskQueue:       tt.fields.k8sNodeTaskQueue,
				volumeTaskQueue:        tt.fields.volumeTaskQueue,
				volumeExpandTaskQueue:  tt.fields.volumeExpandTaskQueue,
				volumeMigrateTaskQueue: tt.fields.volumeMigrateTaskQueue,
				volumeConvertTaskQueue: tt.fields.volumeConvertTaskQueue,
				localNodes:             tt.fields.localNodes,
				logger:                 tt.fields.logger,
			}
			m.setupInformers()
		})
	}
}

func Test_manager_start(t *testing.T) {
	type fields struct {
		name                   string
		namespace              string
		apiClient              client.Client
		informersCache         cache.Cache
		scheme                 *runtime.Scheme
		volumeScheduler        scheduler.Scheduler
		nodeTaskQueue          *common.TaskQueue
		k8sNodeTaskQueue       *common.TaskQueue
		volumeTaskQueue        *common.TaskQueue
		volumeExpandTaskQueue  *common.TaskQueue
		volumeMigrateTaskQueue *common.TaskQueue
		volumeConvertTaskQueue *common.TaskQueue
		localNodes             map[string]apisv1alpha1.State
		logger                 *log.Entry
	}
	type args struct {
		stopCh <-chan struct{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &manager{
				name:                   tt.fields.name,
				namespace:              tt.fields.namespace,
				apiClient:              tt.fields.apiClient,
				informersCache:         tt.fields.informersCache,
				scheme:                 tt.fields.scheme,
				volumeScheduler:        tt.fields.volumeScheduler,
				nodeTaskQueue:          tt.fields.nodeTaskQueue,
				k8sNodeTaskQueue:       tt.fields.k8sNodeTaskQueue,
				volumeTaskQueue:        tt.fields.volumeTaskQueue,
				volumeExpandTaskQueue:  tt.fields.volumeExpandTaskQueue,
				volumeMigrateTaskQueue: tt.fields.volumeMigrateTaskQueue,
				volumeConvertTaskQueue: tt.fields.volumeConvertTaskQueue,
				localNodes:             tt.fields.localNodes,
				logger:                 tt.fields.logger,
			}
			m.start(tt.args.stopCh)
		})
	}
}
