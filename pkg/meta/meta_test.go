/*
Copyright 2018 The Crossplane Authors.

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

package meta

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	groupVersion = "coolstuff/v1"
	kind         = "coolresource"
	namespace    = "coolns"
	name         = "cool"
	uid          = types.UID("definitely-a-uuid")
)

func TestReferenceTo(t *testing.T) {
	tests := map[string]struct {
		o    TypedObject
		want *corev1.ObjectReference
	}{
		"WithTypeMeta": {
			o: &corev1.Pod{
				TypeMeta: metav1.TypeMeta{
					APIVersion: groupVersion,
					Kind:       kind,
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      name,
					UID:       uid,
				},
			},
			want: &corev1.ObjectReference{
				APIVersion: groupVersion,
				Kind:       kind,
				Namespace:  namespace,
				Name:       name,
				UID:        uid,
			},
		},
		"WithoutTypeMeta": {
			o: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      name,
					UID:       uid,
				},
			},
			want: &corev1.ObjectReference{
				Namespace: namespace,
				Name:      name,
				UID:       uid,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := ReferenceTo(tc.o)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("ReferenceTo(): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestAsOwner(t *testing.T) {
	tests := map[string]struct {
		r    *corev1.ObjectReference
		want metav1.OwnerReference
	}{
		"Successful": {
			r: &corev1.ObjectReference{
				APIVersion: groupVersion,
				Kind:       kind,
				Namespace:  name,
				Name:       name,
				UID:        uid,
			},
			want: metav1.OwnerReference{
				APIVersion: groupVersion,
				Kind:       kind,
				Name:       name,
				UID:        uid,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := AsOwner(tc.r)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("AsOwner(): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestAsController(t *testing.T) {
	controller := true

	tests := map[string]struct {
		r    *corev1.ObjectReference
		want metav1.OwnerReference
	}{
		"Successful": {
			r: &corev1.ObjectReference{
				APIVersion: groupVersion,
				Kind:       kind,
				Namespace:  name,
				Name:       name,
				UID:        uid,
			},
			want: metav1.OwnerReference{
				APIVersion: groupVersion,
				Kind:       kind,
				Name:       name,
				UID:        uid,
				Controller: &controller,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := AsController(tc.r)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("AsController(): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestHaveSameController(t *testing.T) {
	controller := true

	controllerA := metav1.OwnerReference{
		UID:        uid,
		Controller: &controller,
	}

	controllerB := metav1.OwnerReference{
		UID:        types.UID("a-different-uuid"),
		Controller: &controller,
	}

	cases := map[string]struct {
		a    metav1.Object
		b    metav1.Object
		want bool
	}{
		"SameController": {
			a: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{controllerA},
				},
			},
			b: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{controllerA},
				},
			},
			want: true,
		},
		"AHasNoController": {
			a: &corev1.Pod{},
			b: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{controllerB},
				},
			},
			want: false,
		},
		"BHasNoController": {
			a: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{controllerA},
				},
			},
			b:    &corev1.Pod{},
			want: false,
		},
		"ControllersDiffer": {
			a: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{controllerA},
				},
			},
			b: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{controllerB},
				},
			},
			want: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := HaveSameController(tc.a, tc.b)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("HaveSameController(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestNamespacedNameOf(t *testing.T) {
	cases := map[string]struct {
		r    *corev1.ObjectReference
		want types.NamespacedName
	}{
		"Success": {
			r:    &corev1.ObjectReference{Namespace: namespace, Name: name},
			want: types.NamespacedName{Namespace: namespace, Name: name},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := NamespacedNameOf(tc.r)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("NamespacedNameOf(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestAddOwnerReference(t *testing.T) {
	owner := metav1.OwnerReference{UID: uid}
	other := metav1.OwnerReference{UID: "a-different-uuid"}

	type args struct {
		o metav1.Object
		r metav1.OwnerReference
	}

	cases := map[string]struct {
		args args
		want []metav1.OwnerReference
	}{
		"NoExistingOwners": {
			args: args{
				o: &corev1.Pod{},
				r: owner,
			},
			want: []metav1.OwnerReference{owner},
		},
		"OwnerAlreadyExists": {
			args: args{
				o: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						OwnerReferences: []metav1.OwnerReference{owner},
					},
				},
				r: owner,
			},
			want: []metav1.OwnerReference{owner},
		},
		"OwnedByAnotherObject": {
			args: args{
				o: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						OwnerReferences: []metav1.OwnerReference{other},
					},
				},
				r: owner,
			},
			want: []metav1.OwnerReference{other, owner},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			AddOwnerReference(tc.args.o, tc.args.r)

			got := tc.args.o.GetOwnerReferences()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("tc.args.o.GetOwnerReferences(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestAddFinalizer(t *testing.T) {
	finalizer := "fin"
	funalizer := "fun"

	type args struct {
		o         metav1.Object
		finalizer string
	}

	cases := map[string]struct {
		args args
		want []string
	}{
		"NoExistingFinalizers": {
			args: args{
				o:         &corev1.Pod{},
				finalizer: finalizer,
			},
			want: []string{finalizer},
		},
		"FinalizerAlreadyExists": {
			args: args{
				o: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Finalizers: []string{finalizer},
					},
				},
				finalizer: finalizer,
			},
			want: []string{finalizer},
		},
		"AnotherFinalizerExists": {
			args: args{
				o: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Finalizers: []string{funalizer},
					},
				},
				finalizer: finalizer,
			},
			want: []string{funalizer, finalizer},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			AddFinalizer(tc.args.o, tc.args.finalizer)

			got := tc.args.o.GetFinalizers()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("tc.args.o.GetFinalizers(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestRemoveFinalizer(t *testing.T) {
	finalizer := "fin"
	funalizer := "fun"

	type args struct {
		o         metav1.Object
		finalizer string
	}

	cases := map[string]struct {
		args args
		want []string
	}{
		"NoExistingFinalizers": {
			args: args{
				o:         &corev1.Pod{},
				finalizer: finalizer,
			},
			want: nil,
		},
		"FinalizerExists": {
			args: args{
				o: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Finalizers: []string{finalizer},
					},
				},
				finalizer: finalizer,
			},
			want: []string{},
		},
		"AnotherFinalizerExists": {
			args: args{
				o: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Finalizers: []string{finalizer, funalizer},
					},
				},
				finalizer: finalizer,
			},
			want: []string{funalizer},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			RemoveFinalizer(tc.args.o, tc.args.finalizer)

			got := tc.args.o.GetFinalizers()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("tc.args.o.GetFinalizers(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestAddLabels(t *testing.T) {
	key, value := "key", "value"
	existingKey, existingValue := "ekey", "evalue"

	type args struct {
		o      metav1.Object
		labels map[string]string
	}

	cases := map[string]struct {
		args args
		want map[string]string
	}{
		"ExistingLabels": {
			args: args{
				o: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							existingKey: existingValue,
						},
					},
				},
				labels: map[string]string{key: value},
			},
			want: map[string]string{
				existingKey: existingValue,
				key:         value,
			},
		},
		"NoExistingLabels": {
			args: args{
				o:      &corev1.Pod{},
				labels: map[string]string{key: value},
			},
			want: map[string]string{key: value},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			AddLabels(tc.args.o, tc.args.labels)

			got := tc.args.o.GetLabels()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("tc.args.o.GetLabels(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestRemoveLabels(t *testing.T) {
	keyA, valueA := "keyA", "valueA"
	keyB, valueB := "keyB", "valueB"

	type args struct {
		o      metav1.Object
		labels []string
	}

	cases := map[string]struct {
		args args
		want map[string]string
	}{
		"ExistingLabels": {
			args: args{
				o: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							keyA: valueA,
							keyB: valueB,
						},
					},
				},
				labels: []string{keyA},
			},
			want: map[string]string{keyB: valueB},
		},
		"NoExistingLabels": {
			args: args{
				o:      &corev1.Pod{},
				labels: []string{keyA},
			},
			want: nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			RemoveLabels(tc.args.o, tc.args.labels...)

			got := tc.args.o.GetLabels()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("tc.args.o.GetLabels(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestAddAnnotations(t *testing.T) {
	key, value := "key", "value"
	existingKey, existingValue := "ekey", "evalue"

	type args struct {
		o           metav1.Object
		annotations map[string]string
	}

	cases := map[string]struct {
		args args
		want map[string]string
	}{
		"ExistingAnnotations": {
			args: args{
				o: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							existingKey: existingValue,
						},
					},
				},
				annotations: map[string]string{key: value},
			},
			want: map[string]string{
				existingKey: existingValue,
				key:         value,
			},
		},
		"NoExistingAnnotations": {
			args: args{
				o:           &corev1.Pod{},
				annotations: map[string]string{key: value},
			},
			want: map[string]string{key: value},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			AddAnnotations(tc.args.o, tc.args.annotations)

			got := tc.args.o.GetAnnotations()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("tc.args.o.GetAnnotations(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestRemoveAnnotations(t *testing.T) {
	keyA, valueA := "keyA", "valueA"
	keyB, valueB := "keyB", "valueB"

	type args struct {
		o           metav1.Object
		annotations []string
	}

	cases := map[string]struct {
		args args
		want map[string]string
	}{
		"ExistingAnnotations": {
			args: args{
				o: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							keyA: valueA,
							keyB: valueB,
						},
					},
				},
				annotations: []string{keyA},
			},
			want: map[string]string{keyB: valueB},
		},
		"NoExistingAnnotations": {
			args: args{
				o:           &corev1.Pod{},
				annotations: []string{keyA},
			},
			want: nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			RemoveAnnotations(tc.args.o, tc.args.annotations...)

			got := tc.args.o.GetAnnotations()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("tc.args.o.GetAnnotations(...): -want, +got:\n%s", diff)
			}
		})
	}
}
