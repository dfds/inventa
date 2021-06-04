package controllers

import (
	"context"
	"fmt"
	"github.com/dfds/inventa/operator/misc"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"
	v13 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"log"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type PublishEventsReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	Store  *misc.InMemoryStore
}

func (r *PublishEventsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func (r *PublishEventsReconciler) SetupWithManager(mgr ctrl.Manager, msgChannel chan interface{}) error {
	fmt.Println("Launching PublishEventsReconciler")
	conf, err := ctrl.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	clientSet, err := kubernetes.NewForConfig(conf)
	if err != nil {
		log.Fatal(err)
	}
	sharedInformers := informers.NewSharedInformerFactory(clientSet, 5 * time.Second)
	SetupInformerv13(sharedInformers.Core().V1().Events(), msgChannel)

	sharedInformers.Start(nil)

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Pod{}).
		Complete(r)
}

func SetupInformerv13(informer v13.EventInformer, msgChannel chan interface{}) {
	handler := &PublishEventsHandler{msgChannel: msgChannel}

	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: handler.Add,
		UpdateFunc: handler.Update,
		DeleteFunc: handler.Delete,
	})
}

type PublishEventsHandler struct {
	msgChannel chan interface{}
}

func (p *PublishEventsHandler) Add(obj interface{}) {
	casted := obj.(*v1.Event)

	log.Printf("Add event :: %v :: %v (%v)\n", casted.Action, casted.Name, casted.Kind)
	p.msgChannel <- casted
}

func (p *PublishEventsHandler) Update(oldObj, newObj interface{}) {
	//casted := oldObj.(*v1.Event)

	//log.Printf("Update event :: %v :: %v (%v)\n", casted.Action, casted.Name, casted.Kind)
	//p.msgChannel <- casted
}

func (p *PublishEventsHandler) Delete(obj interface{}) {
	casted := obj.(*v1.Event)

	log.Printf("Delete event :: %v :: %v (%v)\n", casted.Action, casted.Name, casted.Kind)
	p.msgChannel <- casted
}

