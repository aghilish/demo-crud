package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
)

func main() {
	home, _ := os.UserHomeDir()
	kubeConfigPath := filepath.Join(home, ".kube/config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	createPod(client)
	updatePod(client, "jack", "httpd:latest")
	listPods(client)
	deletePod(client, "jack")
}

func listPods(client *kubernetes.Clientset) {
	pods, err := client.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}
}

func createPod(client *kubernetes.Clientset) {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "shahrooz-",
			Namespace:    "default",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "my-container",
					Image: "nginx",
				},
			},
		},
	}
	createdPod, err := client.CoreV1().Pods("default").Create(context.Background(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Created pod:", createdPod.Name)
}
func updatePod(client *kubernetes.Clientset, podName string, newImageName string) {

	retryError := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		pod, err := client.CoreV1().Pods("default").Get(context.Background(), podName, metav1.GetOptions{})
		if err != nil {
			panic(err.Error())
		}

		// Modify the pod object as needed
		pod.Spec.Containers[0].Image = newImageName

		updatedPod, err := client.CoreV1().Pods("default").Update(context.Background(), pod, metav1.UpdateOptions{})
		if err != nil {
			panic(err.Error())
		}

		fmt.Println("Updated pod:", updatedPod.Name)
		return nil
	})

	if retryError != nil {
		fmt.Println("retrying...")
		panic(retryError.Error())
	}
}
func deletePod(client *kubernetes.Clientset, podName string) {
	err := client.CoreV1().Pods("default").Delete(context.Background(), podName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
}
