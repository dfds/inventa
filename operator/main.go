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

package main

import (
	"flag"
	"fmt"
	"github.com/dfds/inventa/operator/misc"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	stablev1alpha1 "github.com/dfds/inventa/operator/api/v1alpha1"
	"github.com/dfds/inventa/operator/controllers"
	// +kubebuilder:scaffold:imports
)

var (
	scheme             = runtime.NewScheme()
	setupLog           = ctrl.Log.WithName("setup")
	enableServiceProxy = misc.GetEnvBool("INVENTA_OPERATOR_ENABLE_SERVICEPROXY_CONTROLLER", true)
	enableHttpApi      = misc.GetEnvBool("INVENTA_OPERATOR_ENABLE_HTTP_API", true)
	enableApiAuth      = misc.GetEnvBool("INVENTA_OPERATOR_API_ENABLE_AUTH", false)
	enablePublisher    = misc.GetEnvBool(fmt.Sprintf("%s_ENABLE_PUBLISHER", misc.CONF_PREFIX), false)

	enableIngressProxyAnnotationController = misc.GetEnvBool("INVENTA_OPERATOR_ENABLE_INGRESSPROXY_ANNOTATION_CONTROLLER", true)
	enableServiceProxyAnnotationController = misc.GetEnvBool("INVENTA_OPERATOR_ENABLE_SERVICEPROXY_ANNOTATION_CONTROLLER", true)
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(stablev1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "ab15f4e4.dfds.cloud",
		Namespace:              "",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	store := misc.NewInMemoryStore()

	if enablePublisher {
		messageChannel := make(chan interface{}, 99)
		go misc.RunPublisherService(messageChannel)

		if err = (&controllers.PublishEventsReconciler{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("PublishEventsReconciler"),
			Scheme: mgr.GetScheme(),
			Store:  store,
		}).SetupWithManager(mgr, messageChannel); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "PublishEventsReconciler")
			os.Exit(1)
		}
	}

	if enableHttpApi {
		// Start separate Goroutine that runs the API server
		fmt.Println("HTTP api enabled")

		go misc.InitApi(store, enableApiAuth)
	}

	if enableServiceProxyAnnotationController {
		fmt.Println("ServiceProxyAnnotationController enabled")

		if err = (&controllers.ServiceProxyAnnotationReconciler{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("ServiceProxyReconciler"),
			Scheme: mgr.GetScheme(),
			Store:  store,
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "ServiceProxyReconciler")
			os.Exit(1)
		}

	}

	if enableIngressProxyAnnotationController {
		fmt.Println("IngressProxyAnnotationController enabled")

		if err = (&controllers.IngressProxyAnnotationReconciler{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("IngressProxyReconciler"),
			Scheme: mgr.GetScheme(),
			Store:  store,
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "IngressProxyReconciler")
			os.Exit(1)
		}
	}

	if enableServiceProxy {
		fmt.Println("ServiceProxy CRD controller enabled")
		if err = (&controllers.ServiceProxyReconciler{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("ServiceProxy"),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "ServiceProxy")
			os.Exit(1)
		}
		// +kubebuilder:scaffold:builder
	}

	if err := mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("check", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}