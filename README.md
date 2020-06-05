# kubestress
Kubestress is a tool designed for testing auto-scaling of Kubernetes. It can generate a fixed amount of workload (currently only `CPU`) to simulate a certain amount of stress from external sources. The goal of `kubestress` is to test the scaling behavior for Kubernetes cluster without having to use stress testing tool like Locust or JMeter that generate traffic externally. 

## Environmental Variables

The following ENVs are required for running the stress testing.

| ENV Name           | Description                                                  |
| ------------------ | ------------------------------------------------------------ |
| DEPLOYMENT_NAME    | Name of the deployment. This will be automatically maintained by the Helm chart. |
| TOTAL_CPU_LOAD     | Total desired CPU load for the entire deployment. The unit is *milli cores*. |
| PER_POD_PROCESS    | Number of `stress-ng` processes per pod. By default, there is only `1` process per pod. |
| PER_POD_CPU_LIMIT  | CPU limit of each pod. Same value as `spec.template.spec.containers[0].resources.limits.cpu`. |
| TEST_PERIOD_SECOND | Length of each run of `stress-ng`. The unit is *seconds*.    |

