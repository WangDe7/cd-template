package chart

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"

	"github.com/WangDe7/cd-template/imports/k8s"
	"github.com/WangDe7/cd-template/pkg/config"
)

func NewHpaChart(scope constructs.Construct, id string, props *cdk8s.ChartProps) cdk8s.Chart {
	chart := cdk8s.NewChart(scope, jsii.String(id), props)

	k8s.NewKubeHorizontalPodAutoscalerV2(chart, jsii.String("hpa"), &k8s.KubeHorizontalPodAutoscalerV2Props{
		Metadata: &k8s.ObjectMeta{
			Labels: props.Labels,
			Name:   &config.Cfg.Service,
		},
		Spec: &k8s.HorizontalPodAutoscalerSpecV2{
			MinReplicas: jsii.Number(float64(config.Cfg.Hpa.MinReplicas)),
			MaxReplicas: jsii.Number(float64(config.Cfg.Hpa.MaxReplicas)),
			ScaleTargetRef: &k8s.CrossVersionObjectReferenceV2{
				Kind:       jsii.String(config.Cfg.WorkloadType),
				Name:       jsii.String(config.Cfg.Service),
				ApiVersion: jsii.String("apps/v1"),
			},
			Metrics: &[]*k8s.MetricSpecV2{
				{
					Type: jsii.String("Resource"),
					Resource: &k8s.ResourceMetricSourceV2{
						Name: jsii.String("memory"),
						Target: &k8s.MetricTargetV2{
							Type:               jsii.String("Utilization"),
							AverageUtilization: jsii.Number(80),
						},
					},
				},
				{
					Type: jsii.String("Resource"),
					Resource: &k8s.ResourceMetricSourceV2{
						Name: jsii.String("cpu"),
						Target: &k8s.MetricTargetV2{
							Type:               jsii.String("Utilization"),
							AverageUtilization: jsii.Number(80),
						},
					},
				},
			},
		},
	})

	return chart
}
