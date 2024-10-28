package tracing

import (
	"ascale/pkg/conf/env"
	"log"

	gtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Config struct {
	// samples a given fraction of traces. Fractions >= 1 will
	// always sample. If the parent span is sampled, then it's child spans will
	// automatically be sampled. Fractions < 0 are treated as zero, but spans may
	// still be sampled if their parent is.
	Probability float64 `dsn:"-"`
}

func serviceNameFromEnv() string {
	return env.AppID
}

func Init(c *Config) {
	exporter, err := gtrace.NewExporter(gtrace.WithProjectID(env.ProjectID))
	if err != nil {
		log.Fatalf("texporter.NewExporter: %v", err)
	}

	// Create trace provider with the exporter.
	//
	// By default it uses AlwaysSample() which samples all traces.
	// In a production environment or high QPS setup please use
	// ProbabilitySampler set at the desired probability.
	// Example:
	//   config := sdktrace.Config{DefaultSampler:sdktrace.ProbabilitySampler(0.0001)}
	//   tp, err := sdktrace.NewProvider(sdktrace.WithConfig(config), ...)

	tp, err := sdktrace.NewProvider(
		sdktrace.WithSyncer(exporter),
		sdktrace.WithConfig(sdktrace.Config{
			DefaultSampler: sdktrace.ProbabilitySampler(c.Probability),
			Resource: resource.New(
				label.String("serivce.name", env.AppID),
				label.String("serivce.env", env.DeployEnv),
				label.String("serivce.host", env.Hostname),
				label.String("serivce.ip", env.IP),
				label.String("serivce.region", env.Region),
				label.String("serivce.zone", env.Zone),
				label.String("serivce.projectID", env.ProjectID),
			),
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	global.SetTraceProvider(tp)
}
