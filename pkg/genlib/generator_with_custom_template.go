// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package genlib

import (
	"bytes"
	"regexp"
)

type emitter struct {
	fieldName string
	fieldType string
	emitFunc  emitFNotReturn
	state     *GenState
	prefix    []byte
}

// GeneratorWithCustomTemplate is resolved at construction to a slice of emit functions
type GeneratorWithCustomTemplate struct {
	emitters         []emitter
	trailingTemplate []byte
}

func parseCustomTemplate(template []byte) ([]string, map[string][]byte, []byte) {
	if len(template) == 0 {
		return nil, nil, nil
	}

	tokenizer := regexp.MustCompile(`([^{]*)({{\.[^}]+}})*`)
	allIndexes := tokenizer.FindAllSubmatchIndex(template, -1)

	orderedFields := make([]string, 0, len(allIndexes))
	templateFieldsMap := make(map[string][]byte, len(allIndexes))

	var fieldPrefixBuffer []byte
	var fieldPrefixPreviousN int
	var trimTrailingTemplateN int

	for i, loc := range allIndexes {
		var fieldName []byte
		var fieldPrefix []byte

		if loc[4] > -1 && loc[5] > -1 {
			fieldName = template[loc[4]+3 : loc[5]-2]
		}

		if loc[2] > -1 && loc[3] > -1 {
			fieldPrefix = template[loc[2]:loc[3]]
		}

		if len(fieldName) == 0 {
			if template[fieldPrefixPreviousN] == byte(123) {
				fieldPrefixBuffer = append(fieldPrefixBuffer, byte(123))
			} else {
				if i == len(allIndexes)-1 {
					fieldPrefixBuffer = template[trimTrailingTemplateN:]
				} else {
					fieldPrefixBuffer = append(fieldPrefixBuffer, fieldPrefix...)
					fieldPrefixBufferIdx := bytes.Index(template[trimTrailingTemplateN:], fieldPrefixBuffer)
					if fieldPrefixBufferIdx > 0 {
						trimTrailingTemplateN += fieldPrefixBufferIdx
					}

				}
			}
		} else {
			fieldPrefixBuffer = append(fieldPrefixBuffer, fieldPrefix...)
			trimTrailingTemplateN = loc[5]
			templateFieldsMap[string(fieldName)] = fieldPrefixBuffer
			orderedFields = append(orderedFields, string(fieldName))
			fieldPrefixBuffer = nil
		}

		fieldPrefixPreviousN = loc[2]
	}

	return orderedFields, templateFieldsMap, fieldPrefixBuffer

}

func NewGeneratorWithCustomTemplate(template []byte, cfg Config, fields Fields) (*GeneratorWithCustomTemplate, error) {
	// Parse the template and extract relevant information
	orderedFields, templateFieldsMap, trailingTemplate := parseCustomTemplate(template)

	// Preprocess the fields, generating appropriate emit functions
	fieldMap := make(map[string]any)
	fieldTypes := make(map[string]string)
	for _, field := range fields {
		if err := bindField(cfg, field, fieldMap, false); err != nil {
			return nil, err
		}

		fieldTypes[field.Name] = field.Type
	}

	// Roll into slice of emit functions
	emitters := make([]emitter, 0, len(fieldMap))
	for _, fieldName := range orderedFields {
		emitters = append(emitters, emitter{
			fieldName: fieldName,
			fieldType: fieldTypes[fieldName],
			emitFunc:  fieldMap[fieldName].(emitFNotReturn),
			state:     NewGenState(),
			prefix:    templateFieldsMap[fieldName],
		})
	}

	return &GeneratorWithCustomTemplate{emitters: emitters, trailingTemplate: trailingTemplate}, nil
}

func (gen GeneratorWithCustomTemplate) Close() error {
	return nil
}

func (gen GeneratorWithCustomTemplate) Emit(buf *bytes.Buffer) error {
	if err := gen.emit(buf); err != nil {
		return err
	}

	return nil
}

func (gen GeneratorWithCustomTemplate) emit(buf *bytes.Buffer) error {
	for _, e := range gen.emitters {
		buf.Write(e.prefix)
		if err := e.emitFunc(e.state, buf); err != nil {
			return err
		}
	}

	buf.Write(gen.trailingTemplate)
	return nil
}
