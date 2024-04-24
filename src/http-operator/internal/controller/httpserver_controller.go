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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	httpv1alpha1 "github.com/example/http-server-operator/api/v1alpha1"
	"github.com/go-logr/logr"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HTTPServerReconciler reconciles a HTTPServer object
type HTTPServerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=http.example.com,resources=httpservers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=http.example.com,resources=httpservers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=http.example.com,resources=httpservers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HTTPServer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *HTTPServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx = context.Background()
	log := r.Log.WithValues("httpserver", req.NamespacedName)

	// TODO(user): your logic here

	// Fetch the HTTPServer instance
	httpServer := &httpv1alpha1.HTTPServer{}
	if err := r.Get(ctx, req.NamespacedName, httpServer); err != nil {
		log.Error(err, "Failed to get HTTPServer")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Logic to manage HTTP server lifecycle based on httpServer.Spec
	// For example, create a Deployment with Pods using the provided hosts and ports

	// Check if the desired number of replicas has changed
	desiredReplicas := httpServer.Spec.Replicas
	deploymentName := "http-server-deployment"
	currentDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: deploymentName, Namespace: req.Namespace}, currentDeployment)
	if err != nil && errors.IsNotFound(err) {
		// Deployment not found, create a new one with the desired replicas
		deployment := r.newDeploymentForHTTPServer(httpServer)
		log.Info("Creating a new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		err = r.Create(ctx, deployment)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// Check if the number of replicas needs to be updated
	if *currentDeployment.Spec.Replicas != desiredReplicas {
		log.Info("Updating Deployment replicas", "Deployment.Namespace", req.Namespace, "Deployment.Name", deploymentName, "DesiredReplicas", desiredReplicas)
		currentDeployment.Spec.Replicas = &desiredReplicas
		err = r.Update(ctx, currentDeployment)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", req.Namespace, "Deployment.Name", deploymentName)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Define Pod template
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "http-server-pod",
			Namespace: req.Namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "http-server-container",
					Image: "your-http-server-image",
					Ports: []corev1.ContainerPort{
						{ContainerPort: 80}, // Assuming your HTTP server listens on port 80
					},
				},
			},
		},
	}

	// Define Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "http-server-service",
			Namespace: req.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "http-server",
			},
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					Port:     80,
				},
			},
		},
	}

	// Set HTTPServer instance as the owner and controller
	if err := controllerutil.SetControllerReference(httpServer, pod, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// Check if the Pod already exists, if not create a new one
	foundPod := &corev1.Pod{}
	err = r.Get(ctx, types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, foundPod)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.Create(ctx, pod)
		if err != nil {
			log.Error(err, "Failed to create new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
			return ctrl.Result{}, err
		}
		// Pod created successfully
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Pod")
		return ctrl.Result{}, err
	}

	// Check if the Service already exists, if not create a new one
	foundService := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		err = r.Create(ctx, service)
		if err != nil {
			log.Error(err, "Failed to create new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
			return ctrl.Result{}, err
		}
		// Service created successfully
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Service")
		return ctrl.Result{}, err
	}

	// Pod and Service already exist - do nothing
	return ctrl.Result{}, nil

}

func (r *HTTPServerReconciler) newDeploymentForHTTPServer(httpServer *httpv1alpha1.HTTPServer) *appsv1.Deployment {
	// Define Pod template
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "http-server-pod",
			Namespace: httpServer.Namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "http-server-container",
					Image: "your-http-server-image",
					Ports: []corev1.ContainerPort{
						{ContainerPort: 80}, // Assuming your HTTP server listens on port 80
					},
				},
			},
		},
	}

	// Define Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "http-server-deployment",
			Namespace: httpServer.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &httpServer.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "http-server",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "http-server",
					},
				},
				Spec: pod.Spec,
			},
		},
	}

	// Set HTTPServer instance as the owner and controller
	controllerutil.SetControllerReference(httpServer, deployment, r.Scheme)
	return deployment
}

func (r *HTTPServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&httpv1alpha1.HTTPServer{}).
		Owns(&corev1.Pod{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
