package rzap

import (
	"errors"
	"testing"

	"go.uber.org/zap"
)

func BenchmarkLogger(b *testing.B) {
	object := map[string]interface{}{
		"string": "string",
		"map": map[string]interface{}{
			"1": 1,
			"2": 2,
		},
		"map1": map[string]interface{}{
			"1": 1,
			"2": 2,
		},
		"map2": map[string]interface{}{
			"1": 1,
			"2": 2,
		},
		"map3": map[string]interface{}{
			"1": 1,
			"2": 2,
		},
	}
	err := errors.New("error")

	rzap := NewLogger(WithIODiscard())
	zapl := rzap.Zap()
	zaplsu := zapl.Sugar()

	table := []struct {
		scenario string
		logging  func()
	}{
		{
			scenario: "Zap",
			logging: func() {
				zapl.Info("bench", zap.Any("object", object), zap.NamedError("error", err))
			},
		},
		{
			scenario: "ZapSugar",
			logging: func() {
				zaplsu.Infow("bench", "object", object, "error", err)
			},
		},
		{
			scenario: "RZap",
			logging: func() {
				rzap.WithFields("object", object, "error", err).Info("bench")
			},
		},
		{
			scenario: "RZapf",
			logging: func() {
				rzap.WithFields("object", object, "error", err).Infof("bench %s %s", "with", "formating")
			},
		},
	}
	for _, d := range table {
		b.Run(d.scenario, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					d.logging()
				}
			})
		})
	}
}
