// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package genlib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/lithammer/shortuuid/v3"
)

// GeneratorJson is resolved at construction to a slice of emit functions
type GeneratorJson struct {
	emitFuncs []emitF
}

func NewGenerator(cfg Config, fields Fields) (Generator, error) {

	// Preprocess the fields, generating appropriate emit functions
	fieldMap := make(map[string]emitF)
	for _, field := range fields {
		if err := bindField(cfg, field, fieldMap, nil, nil); err != nil {
			return nil, err
		}
	}

	// Roll into slice of emit functions
	emitFuncs := make([]emitF, 0, len(fieldMap))
	for _, f := range fieldMap {
		emitFuncs = append(emitFuncs, f)
	}

	return &GeneratorJson{emitFuncs: emitFuncs}, nil

}

func bindConstantKeyword(field Field, fieldMap map[string]emitF) error {

	prefix := fmt.Sprintf("\"%s\":\"", field.Name)

	fieldMap[field.Name] = func(state *GenState, dupes map[string]struct{}, buf *bytes.Buffer) error {
		value, ok := state.prevCache[field.Name].(string)
		if !ok {
			value = randomdata.Noun()
			state.prevCache[field.Name] = value
		}
		buf.WriteString(prefix)
		buf.WriteString(value)
		buf.WriteByte('"')
		return nil
	}

	return nil
}

func bindKeyword(fieldCfg ConfigField, field Field, fieldMap map[string]emitF) error {
	if len(fieldCfg.Enum) > 0 {
		prefix := fmt.Sprintf("\"%s\":\"", field.Name)

		fieldMap[field.Name] = func(state *GenState, dupes map[string]struct{}, buf *bytes.Buffer) error {
			idx := rand.Intn(len(fieldCfg.Enum) - 1)
			value := fieldCfg.Enum[idx]
			buf.WriteString(prefix)
			buf.WriteString(value)
			buf.WriteByte('"')
			return nil
		}
	} else if len(field.Example) > 0 {

		totWords := len(keywordRegex.Split(field.Example, -1))

		var joiner string
		if strings.Contains(field.Example, "\\.") {
			joiner = "\\."
		} else if strings.Contains(field.Example, "-") {
			joiner = "-"
		} else if strings.Contains(field.Example, "_") {
			joiner = "_"
		} else if strings.Contains(field.Example, " ") {
			joiner = " "
		}

		return bindJoinRand(field, totWords, joiner, fieldMap)
	} else {
		prefix := fmt.Sprintf("\"%s\":\"", field.Name)

		fieldMap[field.Name] = func(state *GenState, dupes map[string]struct{}, buf *bytes.Buffer) error {
			value := randomdata.Noun()
			buf.WriteString(prefix)
			buf.WriteString(value)
			buf.WriteByte('"')
			return nil
		}
	}
	return nil
}

func bindJoinRand(field Field, N int, joiner string, fieldMap map[string]emitF) error {

	prefix := fmt.Sprintf("\"%s\":\"", field.Name)

	fieldMap[field.Name] = func(state *GenState, dupes map[string]struct{}, buf *bytes.Buffer) error {

		buf.WriteString(prefix)

		for i := 0; i < N-1; i++ {
			buf.WriteString(randomdata.Noun())
			buf.WriteString(joiner)
		}
		buf.WriteString(randomdata.Noun())
		buf.WriteByte('"')
		return nil
	}

	return nil
}

func bindStatic(field Field, v interface{}, fieldMap map[string]emitF) error {

	vstr, err := json.Marshal(v)
	if err != nil {
		return err
	}

	payload := fmt.Sprintf("\"%s\":%s", field.Name, vstr)

	fieldMap[field.Name] = func(state *GenState, dupes map[string]struct{}, buf *bytes.Buffer) error {
		buf.WriteString(payload)
		return nil
	}

	return nil
}

func bindBool(field Field, fieldMap map[string]emitF) error {

	prefix := fmt.Sprintf("\"%s\":", field.Name)

	fieldMap[field.Name] = func(state *GenState, dupes map[string]struct{}, buf *bytes.Buffer) error {
		buf.WriteString(prefix)
		switch rand.Int() % 2 {
		case 0:
			buf.WriteString("false")
		case 1:
			buf.WriteString("true")
		}
		return nil
	}

	return nil
}

func bindGeoPoint(field Field, fieldMap map[string]emitF) error {

	prefix := fmt.Sprintf("\"%s\":\"", field.Name)

	fieldMap[field.Name] = func(state *GenState, dupes map[string]struct{}, buf *bytes.Buffer) error {
		buf.WriteString(prefix)
		err := randGeoPoint(buf)
		buf.WriteByte('"')
		return err
	}

	return nil
}

func bindWordN(field Field, n int, fieldMap map[string]emitF) error {
	prefix := fmt.Sprintf("\"%s\":\"", field.Name)

	fieldMap[field.Name] = func(state *GenState, dupes map[string]struct{}, buf *bytes.Buffer) error {
		buf.WriteString(prefix)
		genNounsN(rand.Intn(n), buf)
		buf.WriteByte('"')
		return nil
	}

	return nil
}

func bindNearTime(field Field, fieldMap map[string]emitF) error {
	prefix := fmt.Sprintf("\"%s\":\"", field.Name)

	fieldMap[field.Name] = func(state *GenState, dupes map[string]struct{}, buf *bytes.Buffer) error {
		offset := time.Duration(rand.Intn(FieldTypeTimeRange)*-1) * time.Second
		newTime := time.Now().Add(offset)

		buf.WriteString(prefix)
		buf.WriteString(newTime.Format(FieldTypeTimeLayout))
		buf.WriteByte('"')
		return nil
	}

	return nil
}

func bindIP(field Field, fieldMap map[string]emitF) error {
	prefix := fmt.Sprintf("\"%s\":", field.Name)

	fieldMap[field.Name] = func(state *GenState, dupes map[string]struct{}, buf *bytes.Buffer) error {

		buf.WriteString(prefix)

		i0 := rand.Intn(255)
		i1 := rand.Intn(255)
		i2 := rand.Intn(255)
		i3 := rand.Intn(255)

		_, err := fmt.Fprintf(buf, "\"%d.%d.%d.%d\"", i0, i1, i2, i3)
		return err
	}

	return nil
}

func bindLong(fieldCfg ConfigField, field Field, fieldMap map[string]emitF) error {

	dummyFunc := makeIntFunc(fieldCfg, field)

	fuzziness := fieldCfg.Fuzziness

	prefix := fmt.Sprintf("\"%s\":", field.Name)

	if fuzziness <= 0 {
		fieldMap[field.Name] = func(state *GenState, dupes map[string]struct{}, buf *bytes.Buffer) error {
			buf.WriteString(prefix)
			v := make([]byte, 0, 32)
			v = strconv.AppendInt(v, int64(dummyFunc()), 10)
			buf.Write(v)
			return nil
		}

		return nil
	}

	fieldMap[field.Name] = func(state *GenState, dupes map[string]struct{}, buf *bytes.Buffer) error {
		dummyInt := dummyFunc()
		if previousDummyInt, ok := state.prevCache[field.Name].(int); ok {
			adjustedRatio := 1. - float64(rand.Intn(fuzziness))/100.
			if rand.Int()%2 == 0 {
				adjustedRatio = 1. + float64(rand.Intn(fuzziness))/100.
			}
			dummyInt = int(math.Ceil(float64(previousDummyInt) * adjustedRatio))
		}
		state.prevCache[field.Name] = dummyInt
		buf.WriteString(prefix)
		v := make([]byte, 0, 32)
		v = strconv.AppendInt(v, int64(dummyInt), 10)
		buf.Write(v)
		return nil
	}

	return nil
}

func bindDouble(fieldCfg ConfigField, field Field, fieldMap map[string]emitF) error {

	dummyFunc := makeIntFunc(fieldCfg, field)

	fuzziness := fieldCfg.Fuzziness

	prefix := fmt.Sprintf("\"%s\":", field.Name)

	if fuzziness <= 0 {
		fieldMap[field.Name] = func(state *GenState, dupes map[string]struct{}, buf *bytes.Buffer) error {
			dummyFloat := float64(dummyFunc()) / rand.Float64()
			buf.WriteString(prefix)
			_, err := fmt.Fprintf(buf, "%f", dummyFloat)
			return err
		}

		return nil
	}

	fieldMap[field.Name] = func(state *GenState, dupes map[string]struct{}, buf *bytes.Buffer) error {
		dummyFloat := float64(dummyFunc()) / rand.Float64()
		if previousDummyFloat, ok := state.prevCache[field.Name].(float64); ok {
			adjustedRatio := 1. - float64(rand.Intn(fuzziness))/100.
			if rand.Int()%2 == 0 {
				adjustedRatio = 1. + float64(rand.Intn(fuzziness))/100.
			}
			dummyFloat = previousDummyFloat * adjustedRatio
		}
		state.prevCache[field.Name] = dummyFloat
		buf.WriteString(prefix)
		_, err := fmt.Fprintf(buf, "%f", dummyFloat)
		return err
	}

	return nil
}

func bindCardinality(cfg Config, field Field, fieldMap map[string]emitF) error {

	fieldCfg, _ := cfg.GetField(field.Name)
	cardinality := int(math.Ceil((1000. / float64(fieldCfg.Cardinality))))

	if strings.HasSuffix(field.Name, ".*") {
		field.Name = replacer.Replace(field.Name)
	}

	// Go ahead and bind the original field
	if err := bindByType(cfg, field, fieldMap, nil, nil); err != nil {
		return err
	}

	// We will wrap the function we just generated
	boundF := fieldMap[field.Name]

	fieldMap[field.Name] = func(state *GenState, dupes map[string]struct{}, buf *bytes.Buffer) error {
		var va []bytes.Buffer

		if v, ok := state.prevCache[field.Name]; ok {
			va = v.([]bytes.Buffer)
		}

		// Have we rolled over once?  If not, generate a value and cache it.
		if len(va) < cardinality {

			// Do college try dupe detection on value;
			// Allow dupe if no unique value in nTries.
			nTries := 11 // "These go to 11."
			var tmp bytes.Buffer
			for i := 0; i < nTries; i++ {

				tmp.Reset()
				if err := boundF(state, dupes, &tmp); err != nil {
					return err
				}

				if !isDupe(va, tmp.Bytes()) {
					break
				}
			}

			va = append(va, tmp)
			state.prevCache[field.Name] = va
		}

		idx := int(state.counter % uint64(cardinality))

		// Safety check; should be a noop
		if idx >= len(va) {
			idx = len(va) - 1
		}

		choice := va[idx]
		buf.Write(choice.Bytes())
		return nil
	}

	return nil

}

func makeDynamicStub(root, key string, boundF emitF) emitF {
	target := fmt.Sprintf("\"%s\":", key)

	return func(state *GenState, dupes map[string]struct{}, buf *bytes.Buffer) error {
		// Fire or skip
		if rand.Int()%2 == 0 {
			return nil
		}

		v := state.pool.Get()
		tmp := v.(*bytes.Buffer)
		tmp.Reset()
		defer state.pool.Put(tmp)

		// Fire the bound function, write into temp buffer
		if err := boundF(state, dupes, tmp); err != nil {
			return err
		}

		// If bound function did not write for some reason; abort
		if tmp.Len() <= len(target) {
			return nil
		}

		if !bytes.HasPrefix(tmp.Bytes(), []byte(target)) {
			return fmt.Errorf("Malformed dynamic function payload %s", tmp.String())
		}

		var try int
		const maxTries = 10
		rNoun := randomdata.Noun()
		_, ok := dupes[rNoun]
		for ; ok && try < maxTries; try++ {
			rNoun = randomdata.Noun()
			_, ok = dupes[rNoun]
		}

		// If all else fails, use a shortuuid.
		// Try to avoid this as it is alloc expensive
		if try >= maxTries {
			rNoun = shortuuid.New()
		}

		dupes[rNoun] = struct{}{}

		// ok, formatted as expected, swap it out the payload
		buf.WriteByte('"')
		buf.WriteString(root)
		buf.WriteByte('.')
		buf.WriteString(rNoun)
		buf.WriteString("\":")
		buf.Write(tmp.Bytes()[len(target):])
		return nil
	}
}

func (gen GeneratorJson) Emit(state *GenState, buf *bytes.Buffer) error {

	buf.WriteByte('{')

	if err := gen.emit(state, buf); err != nil {
		return err
	}

	buf.WriteByte('}')

	state.counter += 1

	return nil
}

func (gen GeneratorJson) emit(state *GenState, buf *bytes.Buffer) error {

	dupes := make(map[string]struct{})

	lastComma := -1
	for _, f := range gen.emitFuncs {
		pos := buf.Len()
		if err := f(state, dupes, buf); err != nil {
			return err
		}

		// If we emitted something, write the comma, otherwise skip.
		if buf.Len() > pos {
			buf.WriteByte(',')
			lastComma = buf.Len()
		}
	}

	// Strip dangling comma
	if lastComma == buf.Len() {
		buf.Truncate(buf.Len() - 1)
	}

	return nil
}
