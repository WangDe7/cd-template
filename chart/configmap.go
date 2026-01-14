/*
 * @Author: lwnmengjing
 * @Date: 2021/11/1 11:28 上午
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2021/11/1 11:28 上午
 */

package chart

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"os"

	"github.com/WangDe7/cd-template/imports/k8s"
	"github.com/WangDe7/cd-template/pkg/config"
)

func NewConfigmapChart(scope constructs.Construct, id string, props *cdk8s.ChartProps) cdk8s.Chart {
	if len(config.Cfg.Config) == 0 {
		return nil
	}

	chart := cdk8s.NewChart(scope, jsii.String(id), props)

	data := make(map[string]*string)
	for _, configData := range config.Cfg.ConfigmapResource.ConfigData {
		data[configData.Key] = jsii.String(configData.Value)
	}
	stage := os.Getenv("config_stage")
	if len(config.Cfg.ConfigmapResource.StageConfigs) > 0 {
		for _, stageConfig := range config.Cfg.ConfigmapResource.StageConfigs {
			if stageConfig.Stage == stage {
				for _, configData := range stageConfig.ConfigData {
					data[configData.Key] = jsii.String(configData.Value)
				}
			}
		}
	}
	if len(data) > 0 {
		k8s.NewKubeConfigMap(chart, jsii.String("configmap"), &k8s.KubeConfigMapProps{
			Data: &data,
			Metadata: &k8s.ObjectMeta{
				Name: jsii.String(config.Cfg.ConfigmapResource.Name),
			},
		})
	}

	//for i := range config.Cfg.Config {
	//	if len(config.Cfg.Config[i].Data) == 0 {
	//		continue
	//	}
	//	data := make(map[string]*string)
	//	for path := range config.Cfg.Config[i].Data {
	//		if strings.Index(config.Cfg.Config[i].Data[path], string(os.PathSeparator)) > -1 &&
	//			strings.Index(config.Cfg.Config[i].Data[path], "\n") == -1 {
	//			//path
	//			rb, err := os.ReadFile(config.Cfg.Config[i].Data[path])
	//			if err != nil {
	//				log.Fatalf("read %s error, %s", config.Cfg.Config[i].Data[path], err.Error())
	//				return nil
	//			}
	//			data[path] = jsii.String(string(rb))
	//			continue
	//		}
	//		data[path] = jsii.String(config.Cfg.Config[i].Data[path])
	//	}
	//	if len(data) > 0 {
	//		cm := k8s.NewKubeConfigMap(chart, jsii.String(fmt.Sprintf("%d", i)), &k8s.KubeConfigMapProps{
	//			Data: &data,
	//		})
	//		config.Cfg.Config[i].Name = *cm.Name()
	//
	//	}
	//}
	return chart
}
