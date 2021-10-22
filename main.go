package main

import (
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func main() {
	// ...
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	tracer, closer, err := cfg.New(
		"Polygon",
		config.Logger(jaeger.StdLogger),
	)

	if err != nil {
		panic(err)
		// log.Fatal("Error returned from filepath.Walk:", err)
	}

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	allFunctions()
}

func allFunctions() {

	var n uint

	fmt.Print("Enter Number To Compute : \n")
	fmt.Scan(&n)

	parent := opentracing.GlobalTracer().StartSpan("All Functions")
	defer parent.Finish()

	child := opentracing.GlobalTracer().StartSpan("Fibonacci Function", opentracing.ChildOf(parent.Context()))

	ansFib, fibErr := Fibonacci(n)

	if fibErr != nil {
		panic((fibErr))
	}

	child.Finish()

	child2 := opentracing.GlobalTracer().StartSpan("Factorial Function", opentracing.ChildOf(parent.Context()))

	ansFact := Factorial(int(n))
	time.Sleep(1 * time.Second)

	child2.Finish()

	fmt.Print("Fibonacci Answer : ", ansFib, "\n")
	fmt.Print("Factorial Answer : ", ansFact, "\n")
}
