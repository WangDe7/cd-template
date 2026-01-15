package chart

import (
	"encoding/base64"
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
	for _, secretData := range config.Cfg.SecretResource.SecretData {
		data[secretData.Key] = jsii.String(secretData.Value)
	}
	stage := os.Getenv("config_stage")
	if len(config.Cfg.SecretResource.StageSecrets) > 0 {
		for _, stageSecret := range config.Cfg.SecretResource.StageSecrets {
			if stageSecret.Stage == stage {
				for _, secretData := range stageSecret.SecretData {
					base64Str := base64Encode(secretData.Value)
					data[secretData.Key] = jsii.String(base64Str)
				}
			}
		}
	}
	if len(data) > 0 {
		k8s.NewKubeSecret(chart, jsii.String("secret"), &k8s.KubeSecretProps{
			Data: &data,
			Metadata: &k8s.ObjectMeta{
				Name: jsii.String(config.Cfg.SecretResource.Name),
			},
		})
	}

	return chart
}

func base64Encode(plainStr string) string {
	encodeStr := base64.StdEncoding.EncodeToString([]byte(plainStr))
	return encodeStr
}
