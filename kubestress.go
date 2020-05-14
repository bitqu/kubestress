package kubestress

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

    // pod parameters for stress test
    perPodProcess, err := strconv.Atoi(os.Getenv("PER_POD_PROCESS"))
    totalCpuLoad, err := strconv.Atoi(os.Getenv("TOTAL_CPU_PERCENT"))
    testPeriodSecond, err := strconv.Atoi(os.Getenv("TEST_PERIOD_SECOND"))
    perPodCpuLoad := totalCpuLoad / numReplicas
    perProcessCpuLoad := perPodCpuLoad / perPodProcess

    fmt.Printf("Total CPU load [%v]m. Total replicas [%v]. Per-pod CPU load [%v]. Per-process load [%v].\n",
      totalCpuLoad, numReplicas, perPodCpuLoad, perProcessCpuLoad)
    
    // run stress test
    out, err := exec.Command("/usr/bin/stress-ng",
      "--cpu=" + strconv.Itoa(perPodProcess),
      "--cpu-load=" + strconv.Itoa(perProcessCpuLoad),
      "--timeout="+ strconv.Itoa(testPeriodSecond) +"s").CombinedOutput()

    if err != nil {
      panic(fmt.Errorf("Failed to run stress test: %v", err))
    }
    fmt.Printf("Output from stress-ng: %s\n", out)

  }
}
