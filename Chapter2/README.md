# Schedule a pod and notice the cluster
## List and filter resource we are interesting in
Now we can communicate with the cluster by the help of client SDK. Let us move on, to complete our scheduler. For the simplest case, the pod and node are the minimum information we need to schedule. Let us add some code to query them from the cluster. However, other than the node, not all the pod we are interested in, right? We can use a filter to filter out these that need to schedule and these need schedule by our scheduler. Thanks to client SDK again, they had already provided such interface.
```
...
// interestedPod selects pods that are assigned (scheduled and running).
func interestedPod(pod *v1.Pod) bool {
	return pod.Spec.SchedulerName == SchedulerName && len(pod.Spec.NodeName) == 0
}
...
	informerfactory.Core().V1().Pods().Informer().AddEventHandler(
		cache.FilteringResourceEventHandler{
			FilterFunc: func(obj interface{}) bool {
				switch t := obj.(type) {
				case *v1.Pod:
					return !interestedPod(t)
				case cache.DeletedFinalStateUnknown:
					if pod, ok := t.Obj.(*v1.Pod); ok {
						return !interestedPod(pod)
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
```
Question: 

Could you list resource as possible as you can that a scheduler may interest in?

Extending reading: addAllEventHandlers is the function that kube-scheduler adds events to watch on resouces. The source code is here: [The resource that kube scheduler list and watch](https://github.com/kubernetes/kubernetes/blob/master/pkg/scheduler/eventhandlers.go)
```
...
// addAllEventHandlers is a helper function used in tests and in Scheduler
// to add event handlers for various informers.
func addAllEventHandlers(
	sched *Scheduler,
	informerFactory informers.SharedInformerFactory,
) {
	// scheduled pod cache
	informerFactory.Core().V1().Pods().Informer().AddEventHandler(
		cache.FilteringResourceEventHandler{
...
```
Questions:
1. For a great number of pod and nodes, how should we store them, and which form of data struct to use?
[Scheduler cache](https://github.com/kubernetes/kubernetes/blob/master/pkg/scheduler/internal/cache/interface.go)
2. Should we schedule the pod based on first in first out? Could you suggest a more reasonable way?
[Priority queue](https://github.com/kubernetes/kubernetes/blob/master/pkg/scheduler/internal/queue/scheduling_queue.go)
 
## Bind node and notice the cluster
Now we have enough information to schedule a pod. We should choose the node that the pod run on. This is the so-called filter and score procedure in the kube-scheduler. [filter and score](https://kubernetes.io/docs/concepts/scheduling-eviction/kube-scheduler/#kube-scheduler-implementation)
For our test environment there is only one node and we just want to see the schedule result so let us skip the filter and score. Move on to the next phase. We should tell the cluster our schedule result. This phase called binding and of course there is an api in client SDK to help you.
```
	binding := &v1.Binding{
		ObjectMeta: metav1.ObjectMeta{Namespace: pod.Namespace, Name: pod.Name, UID: pod.UID},
		Target:     v1.ObjectReference{Kind: "Node", Name: pod.Spec.NodeName},
	}
	err := client.CoreV1().Pods(binding.Namespace).Bind(ctx, binding, metav1.CreateOptions{})
```
Let us test our scheduler. Deploy a nginx which specify the schedule name to our scheduler.
```
apiVersion: v1
kind: Pod
metadata:
  name: webserver
spec:
  containers:
  - name: webserver  # The name that this container will have.
    image: nginx:1.14.0 # The image on which it is based.
  hostNetwork: true
  schedulerName: "StupidScheduler"
```
Deploy before start our scheduler.
```
$ kubectl apply -f nginx_deployment.yaml
pod/webserver created
$ kubectl get pods
NAME        READY   STATUS    RESTARTS   AGE
webserver   0/1     Pending   0          4s
```
We can see the pod is in pending status.
Now start our scheduler
```
# go run scheduler.go ~/scheduler.conf
add event for pod default/webserver
add event for node "test"
Success bind pod default/webserver
```
We can see our pod is running now.
```
$ kubectl get pods
NAME        READY   STATUS    RESTARTS   AGE
webserver   1/1     Running   0          2m3s
```
Test wherthe it works
```
$ curl localhost:80
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```
If you see this. Congratulation! You just finish a your own kubernetes scheduler.
This basic tutorial just help you understand the kubernetes scheduler. Here we skip many problems occurs in the production environment. You can continue to read the advance topic or find the answer in the kube-scheduler.
More questions:
1. Can you split the schedule procedure to different phases? And point out each phase main task?
[Kube scheduler framework](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/)
2. Can a single scheduler to schedule pod with different policy according to its configuration?
[Kube scheduler profile]()


 
