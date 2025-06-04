package stage

import (
	"path/filepath"

	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"

	"github.com/mss-boot-io/cd-template/chart"
	"github.com/mss-boot-io/cd-template/pkg/config"
)

// Synth chart generate
func Synth(stage string, paths ...string) {
	path := filepath.Join("dist", stage, config.Cfg.Service)
	if len(paths) > 0 {
		path = filepath.Join("dist", stage, filepath.Join(paths...))
	}
	app := cdk8s.NewApp(&cdk8s.AppProps{Outdir: jsii.String(path)})
	chartProps := &cdk8s.ChartProps{
		Labels: &map[string]*string{
			"app":     &config.Cfg.Service,
			"version": &config.Cfg.Version,
		},
	}
	if &config.Cfg.Project != nil && *&config.Cfg.Project != "" {
		chartProps = &cdk8s.ChartProps{
			Labels: &map[string]*string{
				"app":     &config.Cfg.Service,
				"version": &config.Cfg.Version,
				"project": &config.Cfg.Project,
			},
		}
	}
	if len(config.Cfg.Ports) > 0 {
		chart.NewServiceChart(app, config.Cfg.GetName()+"-service", chartProps)
	}
	needConfigmap := false
	if len(config.Cfg.Config) > 0 {
		for i := range config.Cfg.Config {
			if len(config.Cfg.Config[i].Data) > 0 {
				needConfigmap = true
			}
		}
	}
	if needConfigmap {
		chart.NewConfigmapChart(app, config.Cfg.GetName()+"-configmap", chartProps)
	}
	chart.NewWorkloadChart(app, config.Cfg.GetName()+"-workload", chartProps)
	if config.Cfg.Hpa.Enabled {
		chart.NewHpaChart(app, config.Cfg.GetName()+"-hpa", chartProps)
	}

	if len(config.Cfg.Storages) > 0 && config.Cfg.WorkloadType != "statefulset" {
		storageName := make(map[string]struct{}, 0)
		for i := range config.Cfg.Storages {
			_, ok := storageName[config.Cfg.Storages[i].Name]
			if ok || config.Cfg.Storages[i].Size == "" {
				continue
			}
			storageName[config.Cfg.Storages[i].Name] = struct{}{}
			chart.NewPvcChart(app, config.Cfg.GetName()+"-pvc-"+config.Cfg.Storages[i].Name, chartProps, config.Cfg.Storages[i])
		}
	}
	app.Synth()
}
