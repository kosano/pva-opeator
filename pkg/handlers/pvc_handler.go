package handlers

import (
	"context"
	"time"

	pvav1 "github.com/kosano/pva-operator/api/v1"
	"github.com/kosano/pva-operator/pkg/resources"
	apicorev1 "k8s.io/api/core/v1"
	statsapi "k8s.io/kubelet/pkg/apis/stats/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PVCHandler struct {
	client.Client
}

func NewPVCHander(client client.Client) *PVCHandler {
	return &PVCHandler{
		client,
	}
}

func (p *PVCHandler) IncreaseCapacity(ctx context.Context, vs *statsapi.VolumeStats, pva pvav1.PersistentVolumeAutocapacity) error {
	usageRate := float64((*vs.CapacityBytes - *vs.AvailableBytes) / *vs.CapacityBytes)
	if usageRate >= float64(pva.Spec.UsageRateOver/100) {
		pvc := apicorev1.PersistentVolumeClaim{}
		p.Get(ctx, client.ObjectKey{Namespace: vs.PVCRef.Namespace, Name: vs.PVCRef.Name}, &pvc)
		if resources.IsIncludeStorageProvisioner(pvc, pva.Spec.StorageProvisioners) {
			err := resources.IncreaseCapacity(ctx, p.Client, pvc, pva.Spec.RequestTo)
			return err
		}
		for {
			queryNum := 0
			err := resources.StatusCheck(ctx, p.Client, pvc)
			if err != nil {
				// errors
				queryNum += 1
				time.Sleep(time.Second * 60)
			} else {
				break
			}
			if queryNum >= 3 {
				return err
			}
		}
	}
	return nil
}
