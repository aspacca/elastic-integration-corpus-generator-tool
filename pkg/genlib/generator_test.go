package genlib

import (
	"bytes"
	"context"
	"github.com/elastic/elastic-integration-corpus-generator-tool/pkg/genlib/config"
	"github.com/elastic/elastic-integration-corpus-generator-tool/pkg/genlib/fields"
	"testing"
)

func Benchmark_GeneratorCustomTemplateJSONContent(b *testing.B) {
	ctx := context.Background()
	flds, err := fields.LoadFields(ctx, fields.ProductionBaseURL, "endpoint", "process", "8.2.0")

	template := generateCustomTemplateFromField(Config{}, flds)
	g, err := NewGeneratorWithCustomTemplate(template, Config{}, flds)
	if err != nil {
		b.Fatal(err)
	}

	var buf bytes.Buffer

	state := NewGenState()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := g.Emit(state, &buf)
		if err != nil {
			b.Fatal(err)
		}
		buf.Reset()
	}
}

func Benchmark_GeneratorJetHTMLJSONContent(b *testing.B) {
	ctx := context.Background()
	flds, err := fields.LoadFields(ctx, fields.ProductionBaseURL, "endpoint", "process", "8.2.0")

	template := generateJetTemplateFromField(Config{}, flds)
	g, err := NewGeneratorWithJetHTML(template, Config{}, flds)
	if err != nil {
		b.Fatal(err)
	}

	var buf bytes.Buffer

	state := NewGenState()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := g.Emit(state, &buf)
		if err != nil {
			b.Fatal(err)
		}
		buf.Reset()
	}
}

func Benchmark_GeneratorTextTemplateJSONContent(b *testing.B) {
	ctx := context.Background()
	flds, err := fields.LoadFields(ctx, fields.ProductionBaseURL, "endpoint", "process", "8.2.0")

	template := generateTextTemplateFromField(Config{}, flds)
	g, err := NewGeneratorWithTemplate(template, Config{}, flds)
	if err != nil {
		b.Fatal(err)
	}

	var buf bytes.Buffer

	state := NewGenState()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := g.Emit(state, &buf)
		if err != nil {
			b.Fatal(err)
		}
		buf.Reset()
	}
}

func Benchmark_GeneratorCustomTemplate(b *testing.B) {
	flds := Fields{
		{
			Name: "Version",
			Type: FieldTypeLong,
		},
		{
			Name: "AccountID",
			Type: FieldTypeLong,
		},
		{
			Name:    "InterfaceID",
			Type:    FieldTypeKeyword,
			Example: "eni-1235b8ca123456789",
		},
		{
			Name: "SrcAddr",
			Type: FieldTypeIP,
		},
		{
			Name: "DstAddr",
			Type: FieldTypeIP,
		},
		{
			Name: "SrcPort",
			Type: FieldTypeLong,
		},
		{
			Name: "DstPort",
			Type: FieldTypeLong,
		},
		{
			Name: "Protocol",
			Type: FieldTypeLong,
		},
		{
			Name: "Packets",
			Type: FieldTypeLong,
		},
		{
			Name: "Bytes",
			Type: FieldTypeLong,
		},
		{
			Name: "Start",
			Type: FieldTypeDate,
		},
		{
			Name: "End",
			Type: FieldTypeDate,
		},
		{
			Name: "Action",
			Type: FieldTypeKeyword,
		},
		{
			Name: "LogStatus",
			Type: FieldTypeKeyword,
		},
	}

	configYaml := `- name: Version
  value: 2
- name: AccountID
  value: 627286350134
- name: InterfaceID
  cardinality: 10
- name: SrcAddr
  cardinality: 1
- name: DstAddr
  cardinality: 100
- name: SrcPort
  range: 65535
- name: DstPort
  range: 65535
  cardinality: 100
- name: Protocol
  range: 256
- name: Packets
  range: 1048576
- name: Bytes
  range: 15728640
- name: Action
  enum: ["ACCEPT", "REJECT"]
- name: LogStatus
  enum: ["NODATA", OK", "SKIPDATA"]
`
	cfg, err := config.LoadConfigFromYaml([]byte(configYaml))

	if err != nil {
		b.Fatal(err)
	}

	template := []byte(`{{.Version}} {{.AccountID}} {{.InterfaceID}} {{.SrcAddr}} {{.DstAddr}} {{.SrcPort}} {{.DstPort}} {{.Protocol}} {{.Packets}} {{.Bytes}} {{.Start}} {{.End}} {{.Action}} {{.LogStatus}}`)
	g, err := NewGeneratorWithCustomTemplate(template, cfg, flds)
	if err != nil {
		b.Fatal(err)
	}

	var buf bytes.Buffer

	state := NewGenState()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := g.Emit(state, &buf)
		if err != nil {
			b.Fatal(err)
		}
		buf.Reset()
	}
}

func Benchmark_GeneratorJetHTML(b *testing.B) {
	flds := Fields{
		{
			Name: "Version",
			Type: FieldTypeLong,
		},
		{
			Name: "AccountID",
			Type: FieldTypeLong,
		},
		{
			Name:    "InterfaceID",
			Type:    FieldTypeKeyword,
			Example: "eni-1235b8ca123456789",
		},
		{
			Name: "SrcAddr",
			Type: FieldTypeIP,
		},
		{
			Name: "DstAddr",
			Type: FieldTypeIP,
		},
		{
			Name: "SrcPort",
			Type: FieldTypeLong,
		},
		{
			Name: "DstPort",
			Type: FieldTypeLong,
		},
		{
			Name: "Protocol",
			Type: FieldTypeLong,
		},
		{
			Name: "Packets",
			Type: FieldTypeLong,
		},
		{
			Name: "StartOffset",
			Type: FieldTypeLong,
		},
		{
			Name: "End",
			Type: FieldTypeDate,
		},
		{
			Name: "Action",
			Type: FieldTypeKeyword,
		},
		{
			Name: "LogStatus",
			Type: FieldTypeKeyword,
		},
	}

	configYaml := `- name: Version
  value: 2
- name: AccountID
  value: 627286350134
- name: InterfaceID
  cardinality: 10
- name: SrcAddr
  cardinality: 1
- name: DstAddr
  cardinality: 100
- name: SrcPort
  range: 65535
- name: DstPort
  range: 65535
  cardinality: 100
- name: Protocol
  range: 256
- name: Packets
  range: 1048576
- name: StartOffset
  range: 60
- name: Action
  enum: ["ACCEPT", "REJECT"]
- name: LogStatus
  enum: ["OK", "SKIPDATA"]
`
	cfg, err := config.LoadConfigFromYaml([]byte(configYaml))

	if err != nil {
		b.Fatal(err)
	}

	template := []byte(`{{generate: "Version"}} {{generate: "AccountID"}} {{generate: "InterfaceID"}} {{generate: "SrcAddr"}} {{generate: "DstAddr"}} {{generate: "SrcPort"}} {{generate: "DstPort"}} {{generate: "Protocol"}}{{ packets := generate("Packets")}} {{ packets }} {{ packets * 15}} {{endTime := generate("End")}}{{endTime.Add(generate("StartOffset")*-1000000000).Format:"2006-01-02T15:04:05.999999Z07:00" }} {{endTime.Format:"2006-01-02T15:04:05.999999Z07:00"}} {{generate: "Action"}}{{ if packets == 0 }} NODATA {{ else }} {{generate: "LogStatus"}} {{ end }}`)
	g, err := NewGeneratorWithJetHTML(template, cfg, flds)
	if err != nil {
		b.Fatal(err)
	}

	var buf bytes.Buffer

	state := NewGenState()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := g.Emit(state, &buf)
		if err != nil {
			b.Fatal(err)
		}
		buf.Reset()
	}
}

func Benchmark_GeneratorTextTemplate(b *testing.B) {
	flds := Fields{
		{
			Name: "Version",
			Type: FieldTypeLong,
		},
		{
			Name: "AccountID",
			Type: FieldTypeLong,
		},
		{
			Name:    "InterfaceID",
			Type:    FieldTypeKeyword,
			Example: "eni-1235b8ca123456789",
		},
		{
			Name: "SrcAddr",
			Type: FieldTypeIP,
		},
		{
			Name: "DstAddr",
			Type: FieldTypeIP,
		},
		{
			Name: "SrcPort",
			Type: FieldTypeLong,
		},
		{
			Name: "DstPort",
			Type: FieldTypeLong,
		},
		{
			Name: "Protocol",
			Type: FieldTypeLong,
		},
		{
			Name: "Packets",
			Type: FieldTypeLong,
		},
		{
			Name: "Bytes",
			Type: FieldTypeLong,
		},
		{
			Name: "Start",
			Type: FieldTypeDate,
		},
		{
			Name: "End",
			Type: FieldTypeDate,
		},
		{
			Name: "Action",
			Type: FieldTypeKeyword,
		},
		{
			Name: "LogStatus",
			Type: FieldTypeKeyword,
		},
	}

	configYaml := `- name: Version
  value: 2
- name: AccountID
  value: 627286350134
- name: InterfaceID
  cardinality: 10
- name: SrcAddr
  cardinality: 1
- name: DstAddr
  cardinality: 100
- name: SrcPort
  range: 65535
- name: DstPort
  range: 65535
  cardinality: 100
- name: Protocol
  range: 256
- name: Packets
  range: 1048576
- name: Bytes
  range: 15728640
- name: Action
  enum: ["ACCEPT", "REJECT"]
- name: LogStatus
  enum: ["OK", "SKIPDATA"]
`
	cfg, err := config.LoadConfigFromYaml([]byte(configYaml))

	if err != nil {
		b.Fatal(err)
	}

	template := []byte(`{{generate "Version"}} {{generate "AccountID"}} {{generate "InterfaceID"}} {{generate "SrcAddr"}} {{generate "DstAddr"}} {{generate "SrcPort"}} {{generate "DstPort"}} {{generate "Protocol"}}{{ $packets := generate "Packets" }} {{ $packets }} {{generate "Bytes"}} {{$start := generate "Start" }}{{$start.Format "2006-01-02T15:04:05.999999Z07:00" }} {{$end := generate "End" }}{{$end.Format "2006-01-02T15:04:05.999999Z07:00"}} {{generate "Action"}}{{ if eq $packets 0 }} NODATA {{ else }} {{generate "LogStatus"}} {{ end }}`)
	g, err := NewGeneratorWithTemplate(template, cfg, flds)
	if err != nil {
		b.Fatal(err)
	}

	var buf bytes.Buffer

	state := NewGenState()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := g.Emit(state, &buf)
		if err != nil {
			b.Fatal(err)
		}
		buf.Reset()
	}
}
