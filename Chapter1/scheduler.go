/*
Copyright 2021 Wuming Liu (lwmqwer@163.com)

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
	"context"
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func addPod(obj interface{}) {
	pod := obj.(*v1.Pod)
	fmt.Printf("add event for pod %s/%s\n", pod.Namespace, pod.Name)
}

func main() {
	// Prepare for informerfactory
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(os.Args) <= 1 || len(os.Args[1]) == 0 {
		panic("No --kubeconfig was specified. Using default API client. This might not work")
	}

	// This creates a client, load kubeconfig
	kubeConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: os.Args[1]}, nil).ClientConfig()
	if err != nil {
		os.Exit(-1)
	}
	client, err := clientset.NewForConfig(restclient.AddUserAgent(kubeConfig, "scheduler"))
	if err != nil {
		os.Exit(-1)
	}
	informerfactory := informers.NewSharedInformerFactory(client, 0)
	// Here we only care about pods.
	informerfactory.Core().V1().Pods().Lister()
	// Start all informers.
	informerfactory.Start(ctx.Done())

	// Wait for all caches to sync before scheduling.
	informerfactory.WaitForCacheSync(ctx.Done())
	// Add event handle for add pod
	informerfactory.Core().V1().Pods().Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: addPod,
		},
	)

	for {

	}
}
