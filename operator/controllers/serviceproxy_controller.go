/*
Copyright 2021.

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

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	stablev1alpha1 "github.com/dfds/crossplane-sandbox/dfds-serviceproxy/operator-go/api/v1alpha1"
)

// ServiceProxyReconciler reconciles a ServiceProxy object
type ServiceProxyReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=stable.dfds.cloud,resources=serviceproxies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stable.dfds.cloud,resources=serviceproxies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=stable.dfds.cloud,resources=serviceproxies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ServiceProxy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile

//
// When ServiceProxy is created, create proxy deployment(s) to the specified service(s)
// When ServiceProxy is deleted, tear down associated proxy deployment(s)
// When ServiceProxy is edited, reconcile as necessary
//
func (r *ServiceProxyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("serviceproxy", req.NamespacedName)

	// your logic here
	serviceProxy := &stablev1alpha1.ServiceProxy{}
	err := r.Get(ctx, req.NamespacedName, serviceProxy)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			fmt.Println("ServiceProxy resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		fmt.Println(err, "Failed to get ServiceProxy")
		return ctrl.Result{}, err
	}

	// Check if Deployment used for proxy already exists, if not create it.

	for _, svc := range serviceProxy.Spec.Services {
		found := &v1.Deployment{}
		name := fmt.Sprintf("%s-%s", serviceProxy.Name, svc.Name)
		err = r.Get(ctx, types.NamespacedName{
			Namespace: serviceProxy.Namespace,
			Name:      name,
		}, found)

		if err != nil && errors.IsNotFound(err) {
			// Define a new deployment
			dep := r.deploymentForServiceProxy(serviceProxy, name, svc)
			fmt.Println("Creating a new Deployment", "Deployment.Namespace", dep, "Deployment.Name", dep.Name)
			err = r.Create(ctx, dep)
			if err != nil {
				fmt.Println(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
				return ctrl.Result{}, err
			}
			// Deployment created successfully - return and requeue

			// Create Service for Deployment
			currentSvc := r.serviceForDeployment(dep, serviceProxy, svc)
			fmt.Println("Creating a new Service", "Service.Namespace", currentSvc, "Service.Name", currentSvc.Name)
			err = r.Create(ctx, currentSvc)
			if err != nil {
				fmt.Println(err, "Failed to create new Service", "Service.Namespace", currentSvc.Namespace, "Service.Name", currentSvc.Name)
				return ctrl.Result{}, err
			}

		} else if err != nil {
			fmt.Println(err, "Failed to get Deployment")
			return ctrl.Result{}, err
		}
	}

	// 			return ctrl.Result{Requeue: true}, nil

	return ctrl.Result{}, nil
}

func (r *ServiceProxyReconciler) deploymentForServiceProxy(s *stablev1alpha1.ServiceProxy, name string, svc stablev1alpha1.ServiceProxyService) *v1.Deployment {
	replicas := int32(1)

	ls := labelsForServiceProxy(s.Name, svc.Name)
	addr := fmt.Sprintf("http://%s.%s.svc.cluster.local:%v", svc.LookupServiceName, svc.LookupServiceNamespace, svc.LookupServicePort)

	dep := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: s.Namespace,
		},
		Spec: v1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: v12.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: v12.PodSpec{
					Containers: []v12.Container{
						{
							Image: "dfdsdk/serviceproxy-agent:latest",
							Name:  "nginx-proxy",
							Env: []v12.EnvVar{
								{
									Name:  "ADDR",
									Value: addr,
								},
							},
							Ports: []v12.ContainerPort{
								{
									ContainerPort: 80,
									Name:          "http",
								},
							},
						},
					},
				},
			},
		},
	}

	ctrl.SetControllerReference(s, dep, r.Scheme)

	return dep
}

func (r *ServiceProxyReconciler) serviceForDeployment(d *v1.Deployment, s *stablev1alpha1.ServiceProxy, svcx stablev1alpha1.ServiceProxyService) *v12.Service {
	selector := map[string]string{}
	selector["serviceproxy_svc"] = d.Name

	svc := &v12.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: d.Namespace,
			Name:      d.Name,
		},
		Spec: v12.ServiceSpec{
			Ports: []v12.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(80),
				},
			},
			Selector: selector,
		},
	}

	ctrl.SetControllerReference(s, svc, r.Scheme)

	return svc
}

// labelsForServiceProxy returns the labels for selecting the resources
// belonging to the given memcached CR name.
func labelsForServiceProxy(name string, svc string) map[string]string {
	return map[string]string{"app": "serviceproxy", "serviceproxy_cr": name, "serviceproxy_svc": fmt.Sprintf("%s-%s", name, svc)}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceProxyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1alpha1.ServiceProxy{}).
		Complete(r)
}
