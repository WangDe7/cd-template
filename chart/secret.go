package chart

import (
	"github.com/WangDe7/cd-template/imports/k8s"
	"github.com/WangDe7/cd-template/pkg/config"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"os"
)

func NewSecretChart(scope constructs.Construct, id string, props *cdk8s.ChartProps) cdk8s.Chart {
	if len(config.Cfg.Secret) == 0 {
		return nil
	}

	chart := cdk8s.NewChart(scope, jsii.String(id), props)

	data := make(map[string]*string)
	for key, value := range config.Cfg.SecretResource.SecretData {
		data[key] = jsii.String(value)
	}
	stage := os.Getenv("config_stage")
	if len(config.Cfg.SecretResource.StageSecrets) > 0 {
		for _, stageSecret := range config.Cfg.SecretResource.StageSecrets {
			if stageSecret.Stage == stage {
				for key, value := range stageSecret.SecretData {
					data[key] = jsii.String(value)
				}
			}
		}
	}
	if len(data) > 0 {
		k8s.NewKubeSecret(chart, jsii.String("secret"), &k8s.KubeSecretProps{
			Data: &data,
		})
	}

	return chart
}
