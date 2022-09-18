package resources

import (
	"context"
	"fmt"

	"github.com/kosano/pva-operator/pkg/utils"
	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	StorageProvisionerAnnotation = "volume.beta.kubernetes.io/storage-provisioner"
)

func IsIncludeStorageProvisioner(pvc apicorev1.PersistentVolumeClaim, storageProvisioner []string) bool {
	for k, v := range pvc.Annotations {
		if k == StorageProvisionerAnnotation && utils.IsInclude(v, storageProvisioner) {
			return true
		}
	}
	return false
}

func IncreaseCapacity(ctx context.Context, c client.Client, pvc apicorev1.PersistentVolumeClaim, requestStorage string) error {
	q := resource.MustParse(requestStorage)
	rq, _ := q.AsInt64()
	pvc.Spec.Resources.Requests.Storage().Set(rq)
	return c.Update(ctx, &pvc, nil)
}

func StatusCheck(ctx context.Context, c client.Client, pvc apicorev1.PersistentVolumeClaim) error {
	p := apicorev1.PersistentVolumeClaim{}
	c.Get(ctx, client.ObjectKey{Namespace: pvc.Namespace, Name: pvc.Name}, &p)
	if *p.Status.ResizeStatus != apicorev1.PersistentVolumeClaimNoExpansionInProgress {
		return fmt.Errorf("resize failed, status: %s", *p.Status.ResizeStatus)
	}
	if p.Status.Phase != apicorev1.ClaimBound {
		return fmt.Errorf("%s status is abnormal, status: %s", p.Name, p.Status.Phase)
	}
	return nil
}
