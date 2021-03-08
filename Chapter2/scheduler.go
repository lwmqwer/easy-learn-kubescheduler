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
	"sync"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	schedulerName = "StupidScheduler"
)

var (
	podslock  sync.Mutex
	nodeslock sync.Mutex
	pods      map[string]*v1.Pod
	nodes     map[string]*v1.Node
)

func addPod(obj interface{}) {
	pod := obj.(*v1.Pod)
	fmt.Printf("add event for pod %s/%s\n", pod.Namespace, pod.Name)
	podslock.Lock()
	defer podslock.Unlock()
	pods[pod.Namespace+pod.Name] = pod
}

func addNode(obj interface{}) {
	node := obj.(*v1.Node)
	fmt.Printf("add event for node %q\n", node.Name)
	nodeslock.Lock()
	defer nodeslock.Unlock()
	nodes[node.Name] = node
}

// interestedPod selects pods that are assigned (scheduled and running).
func interestedPod(pod *v1.Pod) bool {
	return pod.Spec.SchedulerName == schedulerName && len(pod.Spec.NodeName) == 0
}

func scheduleOne(ctx context.Context, client clientset.Interface) {
	var pod *v1.Pod

	podslock.Lock()
	for k, v := range pods {
		pod = v
		delete(pods, k)
		break
	}
	podslock.Unlock()
	if pod == nil {
		return
	}
	nodeslock.Lock()
	for _, v := range nodes {
		pod.Spec.NodeName = v.Name
		break
	}
	nodeslock.Unlock()
	binding := &v1.Binding{
		ObjectMeta: metav1.ObjectMeta{Namespace: pod.Namespace, Name: pod.Name, UID: pod.UID},
		Target:     v1.ObjectReference{Kind: "Node", Name: pod.Spec.NodeName},
	}
	err := client.CoreV1().Pods(binding.Namespace).Bind(ctx, binding, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("failed bind pod %s/%s\n", pod.Namespace, pod.Name)
		pod.Spec.NodeName = ""
		podslock.Lock()
		defer podslock.Unlock()
		pods[pod.Namespace+pod.Name] = pod
	}
	fmt.Printf("Success bind pod %s/%s\n", pod.Namespace, pod.Name)
}

func main() {
	pods = make(map[string]*v1.Pod)
	nodes = make(map[string]*v1.Node)
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
	// Here we only care about pods and nodes.
	informerfactory.Core().V1().Pods().Lister()
	informerfactory.Core().V1().Nodes().Lister()
	// Start all informers.
	informerfactory.Start(ctx.Done())

	// Wait for all caches to sync before scheduling.
	informerfactory.WaitForCacheSync(ctx.Done())

	informerfactory.Core().V1().Pods().Informer().AddEventHandler(
		cache.FilteringResourceEventHandler{
			FilterFunc: func(obj interface{}) bool {
				switch t := obj.(type) {
				case *v1.Pod:
					return interestedPod(t)
				case cache.DeletedFinalStateUnknown:
					if pod, ok := t.Obj.(*v1.Pod); ok {
						return interestedPod(pod)
					}
					fmt.Errorf("unable to convert object %T to *v1.Pod", obj)
					return false
				default:
					fmt.Errorf("unable to handle object in %T", obj)
					return false
				}
			},
			Handler: cache.ResourceEventHandlerFuncs{
				AddFunc: addPod,
			},
		},
	)

	informerfactory.Core().V1().Nodes().Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: addNode,
		},
	)

	for {
		scheduleOne(ctx, client)
	}
}
