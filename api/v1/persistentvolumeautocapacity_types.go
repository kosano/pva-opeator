/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PersistentVolumeAutocapacitySpec defines the desired state of PersistentVolumeAutocapacity
type PersistentVolumeAutocapacitySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of PersistentVolumeAutocapacity. Edit persistentvolumeautocapacity_types.go to remove/update
	PVCNames            []string `json:"pvcName,omitempty"`
	Namespaces          []string `json:"namespace,omitempty"`
	StorageProvisioners []string `json:"storageProvisioners,omitempty"`
	UsageRateOver       uint16   `json:"usageRateOver,omitempty"`
	Threshold           []string `json:"metrics,omitempty"`
	RequestTo           string   `json:"requestTo,omitempty"` //[#ORIGINAL_SIZE]
}

// PersistentVolumeAutocapacityStatus defines the observed state of PersistentVolumeAutocapacity
type PersistentVolumeAutocapacityStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PersistentVolumeAutocapacity is the Schema for the persistentvolumeautocapacities API
type PersistentVolumeAutocapacity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PersistentVolumeAutocapacitySpec   `json:"spec,omitempty"`
	Status PersistentVolumeAutocapacityStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PersistentVolumeAutocapacityList contains a list of PersistentVolumeAutocapacity
type PersistentVolumeAutocapacityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PersistentVolumeAutocapacity `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PersistentVolumeAutocapacity{}, &PersistentVolumeAutocapacityList{})
}
