package main

import (
  "context"
  "fmt"
  "strconv"
  "os"
  "os/exec"
  apiv1 "k8s.io/api/core/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/rest"
)

func main() {
  // creates the in-cluster config
  config, err := rest.InClusterConfig()
  if err != nil {
    panic(err.Error())
  }
  // creates the clientset
  clientset, err := kubernetes.NewForConfig(config)
  if err != nil {
    panic(err.Error())
  }

  for {

    deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
    result, err := deploymentsClient.Get(context.TODO(), os.Getenv("DEPLOYMENT_NAME"), metav1.GetOptions{})
    if err != nil {
      panic(fmt.Errorf("Failed to get latest version of Deployment: %v", err))
    }
    numReplicas := int(*(result.Spec.Replicas))

    testPeriodSecond, err := strconv.Atoi(os.Getenv("TEST_PERIOD_SECOND"))
    // pod parameters for stress test
    perPodProcess, err := strconv.Atoi(os.Getenv("PER_POD_PROCESS"))
    totalCpuLoad, err := strconv.Atoi(os.Getenv("TOTAL_CPU_LOAD"))
    
    perPodCpuLimit, err := strconv.Atoi(os.Getenv("PER_POD_CPU_LIMIT"))    
    perPodCpuLoad := totalCpuLoad / numReplicas
    
    perProcessCpuLoadPercent := 100 * perPodCpuLoad / perPodProcess / perPodCpuLimit
    if(perProcessCpuLoadPercent > 100) {
      // only set per-process CPU load to less than 100 CPU limit is plenty
      perProcessCpuLoadPercent = 100
    }

    fmt.Printf("Total CPU load [%v]m. Total replicas [%v]. Per-pod CPU load [%v]m. Per-pod processes [%v]. Per-process load [%v] percent.\n",
      totalCpuLoad, numReplicas, perPodCpuLoad, perPodProcess, perProcessCpuLoadPercent)
    
    // run stress test
    out, err := exec.Command("/usr/bin/stress-ng",
      "--cpu=" + strconv.Itoa(perPodProcess),
      "--cpu-load=" + strconv.Itoa(perProcessCpuLoadPercent),
      "--timeout="+ strconv.Itoa(testPeriodSecond) +"s").CombinedOutput()

    if err != nil {
      panic(fmt.Errorf("Failed to run stress test: %v", err))
    }
    fmt.Printf("Output from stress-ng: %s\n", out)

  }
}
