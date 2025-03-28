/*
Copyright 2025.

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

// +kubebuilder:validation:Enum=basic-auth;jwt
// SecretType defines the type of secret to generate.
type SecretType string

// CustomSecretSpec defines the desired state of CustomSecret.
type CustomSecretSpec struct {
	SecretType     SecretType `json:"secretType"`
	Username       string     `json:"username"`
	PasswordLength int        `json:"passwordLength"`
	RotationPeriod string     `json:"rotationPeriod"`
}

// CustomSecretStatus defines the observed state of CustomSecret.
type CustomSecretStatus struct {
	LastUpdated string `json:"lastUpdated"`
	SecretName  string `json:"secretName"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// CustomSecret is the Schema for the customsecrets API.
type CustomSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CustomSecretSpec   `json:"spec,omitempty"`
	Status CustomSecretStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CustomSecretList contains a list of CustomSecret.
type CustomSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CustomSecret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CustomSecret{}, &CustomSecretList{})
}
