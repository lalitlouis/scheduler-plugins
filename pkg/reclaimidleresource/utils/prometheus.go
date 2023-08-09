package utils

import (
	"context"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"sigs.k8s.io/scheduler-plugins/pkg/reclaimidleresource/constants"
)

// GetPromServiceIP gets the ip to the prometheus service of the cluster
func GetPromServiceIP(cs kubernetes.Interface) (string, error) {
	if os.Getenv("ENV") == constants.DevENVFlag {
		return "localhost", nil
	}

	promService, err := cs.CoreV1().Services(constants.PrometheusNamespace).Get(context.TODO(), constants.PrometheusServiceName, meta_v1.GetOptions{})
	if err != nil {
		klog.Error("Couldn't fetch prometheus service IP")
		return "", err
	}

	return promService.Spec.ClusterIP, nil
}

// RunPromQuery runs prom query against the prometheus service
func RunPromQuery(query string, ip string) (model.Value, error) {

	promAddress := constants.HTTPPrefix + ip + ":" + constants.PrometheusPort
	client, err := api.NewClient(api.Config{
		Address: promAddress,
	})
	if err != nil {
		klog.Error("Error creating client: %v\n", err)
		return nil, err
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := v1api.Query(ctx, query, time.Now(), v1.WithTimeout(5*time.Second))
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return nil, err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
		return nil, err
	}

	return result, nil
}

func getAverageUsageCalculationMetricString(resourceType, podName, podNamespace, period string) string {

	resourceMetricCalculationStr := ""
	switch resourceType {
	case constants.ResourceTypeGPU:
		gpuUsageString := `scalar(avg_over_time(DCGM_FI_PROF_GR_ENGINE_ACTIVE{exported_pod="%s", exported_namespace="%s"}[%s]))`
		resourceMetricCalculationStr = fmt.Sprintf(gpuUsageString, podName, podNamespace, period)

	case constants.ResourceTypeCPU:
		cpuUsageString := `scalar(sum(rate(container_cpu_usage_seconds_total{pod="%s",namespace="%s",container!=""}[%s])) by (pod_name))`
		resourceMetricCalculationStr = fmt.Sprintf(cpuUsageString, podName, podNamespace, period)

	}
	klog.Info("Average Resource usage calculation metric string  - ", resourceMetricCalculationStr)
	return resourceMetricCalculationStr
}

// CalculateAverageGPUUsage retrurns the average gpu usage for the period mentioned
func CalculateAverageUsage(resourceType, podName, podNamespace, ip string, period string) float64 {

	measureMetric := getAverageUsageCalculationMetricString(resourceType, podName, podNamespace, period+"s")
	promOutput, _ := GetScalarMetricValue(measureMetric, ip)
	return promOutput.Value
}

// GetScalarMetricValue for prometheus scalar metrics value
func GetScalarMetricValue(query, ip string) (PromScalar, error) {
	promOutput := PromScalar{}
	output, err := RunPromQuery(query, ip)
	if err != nil {
		return promOutput, err
	}
	scalarVal := output.(*model.Scalar)
	promOutput = PromScalar{
		Value:     float64(scalarVal.Value),
		TimeStamp: float64(scalarVal.Timestamp),
	}
	if math.IsNaN(promOutput.Value) {
		promOutput.Value = 0
	}
	return promOutput, nil
}
