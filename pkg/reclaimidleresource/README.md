# Overview

This folder holds the reclaim idle resource plugin implemented as discussed

## Maturity Level

<!-- Check one of the values: Sample, Alpha, Beta, GA -->

- [x] 💡 Sample (for demonstrating and inspiring purpose)
- [ ] 👶 Alpha (used in companies for pilot projects)
- [ ] 👦 Beta (used in companies and developed actively)
- [ ] 👨 Stable (used in companies for production workloads)

## Example scheduler config:

```yaml
apiVersion: kubescheduler.config.k8s.io/v1beta2
kind: KubeSchedulerConfiguration
leaderElection:
  leaderElect: false
clientConnection:
  kubeconfig: "REPLACE_ME_WITH_KUBE_CONFIG_PATH"
profiles:
- schedulerName: default-scheduler
  plugins:
    postFilter:
      enabled:
      - name: ReclaimIdleResource
      disabled:
      - name: DefaultPreemption
```

## How to define ReclaimIdleResource policy on PriorityClass resource

Reclaim Idle Resource policy can be defined on each `PriorityClass` resource by annotations like below:

```yaml
# PriorityClass with ReclaimIdleResource policy:
# Any pod P in this priority class can not be preempted (can tolerate preemption)
# - by preemptor pods with priority < 10000 
# - and if P is within 1h since being scheduled
# - and if resource type is gpu and has been below usage threshold of 0.0 for the last 1h
kind: PriorityClass
metadata:
  name: toleration-policy-sample
  annotations:
    reclaim-idle-resource.scheduling.x-k8s.io/minimum-preemptable-priority: "10000"
    reclaim-idle-resource.scheduling.x-k8s.io/toleration-seconds: "3600"
    reclaim-idle-resource.scheduling.x-k8s.io/resource-type: "gpu"
    reclaim-idle-resource.scheduling.x-k8s.io/resource-idle-seconds: "3600"
    reclaim-idle-resource.scheduling.x-k8s.io/resource-idle-usage-threshold: "0.0"
value: 8000
```
