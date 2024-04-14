/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctfv1alpha1 "securinetes.com/ctf/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// CloudReconciler reconciles a Cloud object
type CloudReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ctf.securinetes.com,resources=clouds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ctf.securinetes.com,resources=clouds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ctf.securinetes.com,resources=clouds/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Cloud object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *CloudReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	ctf := ctfv1alpha1.Cloud{}
	if err := r.Get(ctx, req.NamespacedName, &ctf); err != nil {
		log.Error(err, "you are not qualified yet...")

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if ctf.Spec.Escape {
		fmt.Printf("You need to escape...\n")
		time.Sleep(30 * time.Second)
		if err := r.createescapepod(ctx); err != nil {
			return ctrl.Result{}, err
		}
		time.Sleep(3 * time.Minute)
		if err := r.deletepod(ctx); err != nil {
			return ctrl.Result{}, err
		}

	} else {
		fmt.Printf("You are not ready to escape...\n")
		if err := r.deletepod(ctx); err != nil {
			return ctrl.Result{}, err
		}

	}

	return ctrl.Result{}, nil
}

func (r *CloudReconciler) escapepod() corev1.Pod {
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "escape",
			Namespace: "chunin",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "escape",
					Image: "youffes/escape",
				},
			},
		},
	}
	return pod
}

func (r *CloudReconciler) createescapepod(ctx context.Context) error {
	pod := corev1.Pod{}
	podname := types.NamespacedName{
		Name:      "escape",
		Namespace: "chunin",
	}
	if err := r.Get(ctx, podname, &pod); err != nil {
		if errors.IsNotFound(err) {
			pod := r.escapepod()
			if err := r.Create(ctx, &pod); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *CloudReconciler) deletepod(ctx context.Context) error {
	pod := corev1.Pod{}
	podname := types.NamespacedName{
		Name:      "escape",
		Namespace: "chunin",
	}
	if err := r.Get(ctx, podname, &pod); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		if err := r.Delete(ctx, &pod); err != nil {
			return err
		}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CloudReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ctfv1alpha1.Cloud{}).
		Complete(r)
}
