#!/bin/sh
# Copyright 2020 Wuming Liu (lwmqwer@163.com)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

mirror="registry.aliyuncs.com/google_containers"
kubernetes_version=v1.20.2
pause_version=3.2
etcd_version=3.4.13-0
coredns_version=1.7.0

sudo docker pull $mirror/kube-apiserver:$kubernetes_version
sudo docker pull $mirror/kube-controller-manager:$kubernetes_version
sudo docker pull $mirror/kube-scheduler:$kubernetes_version
sudo docker pull $mirror/kube-proxy:$kubernetes_version
sudo docker pull $mirror/pause:$pause_version
sudo docker pull $mirror/etcd:$etcd_version
sudo docker pull $mirror/coredns:$coredns_version

sudo docker tag $mirror/kube-apiserver:$kubernetes_version k8s.gcr.io/kube-apiserver:$kubernetes_version
sudo docker tag $mirror/kube-controller-manager:$kubernetes_version k8s.gcr.io/kube-controller-manager:$kubernetes_version
sudo docker tag $mirror/kube-scheduler:$kubernetes_version k8s.gcr.io/kube-scheduler:$kubernetes_version
sudo docker tag $mirror/kube-proxy:$kubernetes_version k8s.gcr.io/kube-proxy:$kubernetes_version
sudo docker tag $mirror/pause:$pause_version k8s.gcr.io/pause:$pause_version
sudo docker tag $mirror/etcd:$etcd_version k8s.gcr.io/etcd:$etcd_version
sudo docker tag $mirror/coredns:$coredns_version k8s.gcr.io/coredns:$coredns_version

sudo docker rmi $mirror/kube-apiserver:$kubernetes_version
sudo docker rmi $mirror/kube-controller-manager:$kubernetes_version
sudo docker rmi $mirror/kube-scheduler:$kubernetes_version
sudo docker rmi $mirror/kube-proxy:$kubernetes_version
sudo docker rmi $mirror/pause:$pause_version
sudo docker rmi $mirror/etcd:$etcd_version
sudo docker rmi $mirror/coredns:$coredns_version