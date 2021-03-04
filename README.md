# A tutorial to crack kubenetes scheduler
This is a tutorial to help you understand the kubernetes scheduler and the philosophy behind it.

This tutorial will guide you to write your own scheduler step by step. Besides that, there are also some questions to help you understand the difficulties in the production environment and the practice way to solve them by kube-scheduler.

I would try my best to make the tutorial as simple as possible and focus on the scheduler trunk. So that we would be hindered by robust, performance, error-handling etc. This does not mean these are not important, on the contrary, they cost the engineer enoumous time to tune it. After this tutorial, I believe you can find out these answers in kube-scheduler source code.

Before starting this tutorial, you still need some prerequisites for it:
- You need some basic knowledge about kubernetes and what the role is the scheduler. Here two links to help you understand it. 
[What-is-kubernetes](https://kubernetes.io/docs/concepts/overview/what-is-kubernetes/) and [Kubernetes Components](https://kubernetes.io/docs/concepts/overview/components/)
- The Golang programming language, the tutorial chooses golang as the programming language so you need some knowledge to read the code and how to compile and run it. The kubernetes provides many programming language SDKs, you could choose one as you like to write a C++ , python, or Jave scheduler for kubernetes. You can download and install the golang SDK from [here](https://golang.org/dl/) and I highly recommend this [book](https://www.gopl.io/) for the beginner. 
- You also need a computer to setup the development environment and setup a kubernetes cluster to test it. 
- The other prerequisites I would point out in the following chapters.


## 1. [Communicate with the kubernetes cluster](Chapter1/README.md)
# Advance topics
## 1. [Metrics and logs]()
## 2. [Scheduler Cache and priority queue]()
## 3. [Framework and profile]()
