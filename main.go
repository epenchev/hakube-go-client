package main

import (
	"context"
	"fmt"
	"log"

	v1beta1 "hakube-go-client/v1beta1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func createOperatorScheme() (*runtime.Scheme, error) {
	// create a new scheme specifically for this manager
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

func main() {
	cl := getClient()

	objectName := "saphana-scaleup"
	namespace := "hakube-manager"

	// Get node attributes storage object.
	// hkbnattr := &v1beta1.HAKubeNodeAttributes{}
	//if err := cl.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: objectName}, hkbnattr); err != nil {
	//    log.Fatalln("failed to fetch node attributes store object, error:", err)
	//}

	// Get resource controller object.
	hactrl := &v1beta1.HAKubeController{}
	if err := cl.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: objectName}, hactrl); err != nil {
		log.Fatalln("failed to fetch resource controller object, error:", err)
	}
	// Dump objects
	fmt.Printf("resource configuration %+v\n", hactrl.Spec.Resources)
	//fmt.Printf("node-attributes %+v\n", hkbnattr.NodeAttributes)
}
