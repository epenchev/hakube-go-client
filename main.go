package main

import (
	"context"
	"fmt"
	"log"

	v1beta1 "hakube-go-client/v1beta1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func createOperatorScheme() (*runtime.Scheme, error) {
	hakubeScheme := runtime.NewScheme()
	// add standard resource types to the scheme
	if err := scheme.AddToScheme(hakubeScheme); err != nil {
		return nil, err
	}
	// add custom resource types to the default scheme
	if err := v1beta1.AddToScheme(hakubeScheme); err != nil {
		return nil, err
	}
	return hakubeScheme, nil
}

func getClient() client.Client {
	cfg, err := ctrl.GetConfig()
	if err != nil {
		log.Fatalln("Error getting config")
	}
	hakubeScheme, _ := createOperatorScheme()
	cl, err := client.New(cfg, client.Options{Scheme: hakubeScheme})
	if err != nil {
		log.Fatalln("failed to create client")
	}
	return cl
}

func getClientNoScheme() client.Client {
	cfg, err := ctrl.GetConfig()
	if err != nil {
		log.Fatalln("Error getting config")
	}
	cl, err := client.New(cfg, client.Options{})
	if err != nil {
		log.Fatalln("failed to create client")
	}
	return cl
}

func getObjectsWithScheme() {
	cl := getClient()

	objectName := "saphana-scaleup"
	namespace := "hakube-manager"

	// Get node attributes storage object.
	hkbnattr := &v1beta1.HAKubeNodeAttributes{}
	if err := cl.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: objectName}, hkbnattr); err != nil {
		log.Fatalln("failed to fetch node attributes store object, error:", err)
	}

	// Get resource controller object.
	hactrl := &v1beta1.HAKubeController{}
	if err := cl.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: objectName}, hactrl); err != nil {
		log.Fatalln("failed to fetch resource controller object, error:", err)
	}
	// Dump objects
	fmt.Printf("resource configuration %+v\n", hactrl.Spec.Resources)
	fmt.Printf("node-attributes %+v\n", hkbnattr.NodeAttributes)
}

// see https://mittachaitu.medium.com/dealing-with-unstructured-objects-f190629e7d60
// see https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1/unstructured
// see https://github.com/kubernetes/apimachinery/blob/master/pkg/apis/meta/v1/unstructured/helpers.go
// for dealing with unstructured objects
func getObjectsUnstructured() {
	cl := getClientNoScheme()

	objectName := "saphana-scaleup"
	namespace := "hakube-manager"

	// Get resource controller object.
	kind := "HAKubeController"
	apiversion := "hakube.io/v1beta1"

	obj := new(unstructured.Unstructured)
	obj.SetName(objectName)
	obj.SetAPIVersion(apiversion)
	obj.SetGroupVersionKind(schema.FromAPIVersionAndKind(apiversion, kind))
	key := client.ObjectKey{Name: obj.GetName(), Namespace: namespace}
	if err := cl.Get(context.Background(), key, obj); err != nil {
		log.Fatalln("failed to retrieve object", obj.GetKind(), key.Namespace, key.Name)
	}
	fmt.Printf("object contents %+v\n\n", obj)

	resourceConfig, _, err := unstructured.NestedFieldNoCopy(obj.Object, "spec", "resources")
	if err == nil {
		fmt.Printf("resource configuration %+v\n", resourceConfig)
	}
}

func main() {
	getObjectsWithScheme()
	//getObjectsUnstructured()
}
