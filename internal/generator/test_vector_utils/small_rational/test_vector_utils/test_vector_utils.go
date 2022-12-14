// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package test_vector_utils

import (
	"encoding/json"
	"fmt"
	"github.com/consensys/gnark-crypto/internal/generator/test_vector_utils/small_rational"
	"github.com/consensys/gnark-crypto/internal/generator/test_vector_utils/small_rational/polynomial"
	"hash"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type ElementTriplet struct {
	key1        small_rational.SmallRational
	key2        small_rational.SmallRational
	key2Present bool
	value       small_rational.SmallRational
	used        bool
}

func (t *ElementTriplet) CmpKey(o *ElementTriplet) int {
	if cmp1 := t.key1.Cmp(&o.key1); cmp1 != 0 {
		return cmp1
	}

	if t.key2Present {
		if o.key2Present {
			return t.key2.Cmp(&o.key2)
		}
		return 1
	} else {
		if o.key2Present {
			return -1
		}
		return 0
	}
}

var MapCache = make(map[string]*ElementMap)

func ElementMapFromFile(path string) (*ElementMap, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	if h, ok := MapCache[path]; ok {
		return h, nil
	}
	var bytes []byte
	if bytes, err = os.ReadFile(path); err == nil {
		var asMap map[string]interface{}
		if err = json.Unmarshal(bytes, &asMap); err != nil {
			return nil, err
		}

		var h ElementMap
		if h, err = CreateElementMap(asMap); err == nil {
			MapCache[path] = &h
		}

		return &h, err

	} else {
		return nil, err
	}
}

func CreateElementMap(rawMap map[string]interface{}) (ElementMap, error) {
	res := make(ElementMap, 0, len(rawMap))

	for k, v := range rawMap {
		var entry ElementTriplet
		if _, err := entry.value.SetInterface(v); err != nil {
			return nil, err
		}

		key := strings.Split(k, ",")
		switch len(key) {
		case 1:
			entry.key2Present = false
		case 2:
			entry.key2Present = true
			if _, err := entry.key2.SetInterface(key[1]); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("cannot parse %T as one or two field elements", v)
		}
		if _, err := entry.key1.SetInterface(key[0]); err != nil {
			return nil, err
		}

		res = append(res, &entry)
	}

	res.sort()
	return res, nil
}

type ElementMap []*ElementTriplet

type MapHash struct {
	Map        *ElementMap
	state      small_rational.SmallRational
	stateValid bool
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m *MapHash) Write(p []byte) (n int, err error) {
	var x small_rational.SmallRational
	for i := 0; i < len(p); i += small_rational.Bytes {
		x.SetBytes(p[i:min(len(p), i+small_rational.Bytes)])
		if err = m.write(x); err != nil {
			return
		}
	}
	n = len(p)
	return
}

func (m *MapHash) Sum(b []byte) []byte {
	mP := *m
	if _, err := mP.Write(b); err != nil {
		panic(err)
	}
	bytes := mP.state.Bytes()
	return bytes[:]
}

func (m *MapHash) Reset() {
	m.stateValid = false
}

func (m *MapHash) Size() int {
	return small_rational.Bytes
}

func (m *MapHash) BlockSize() int {
	return small_rational.Bytes
}

func (m *MapHash) write(x small_rational.SmallRational) error {
	X := &x
	Y := &m.state
	if !m.stateValid {
		Y = nil
	}
	var err error
	if m.state, err = m.Map.FindPair(X, Y); err == nil {
		m.stateValid = true
	}
	return err
}

func (t *ElementTriplet) writeKey(sb *strings.Builder) {
	sb.WriteRune('"')
	sb.WriteString(t.key1.String())
	if t.key2Present {
		sb.WriteRune(',')
		sb.WriteString(t.key2.String())
	}
	sb.WriteRune('"')
}

func SaveUsedHashEntries() error {
	for path, hash := range MapCache {
		if err := hash.SaveUsedEntries(path); err != nil {
			return err
		}
	}
	return nil
}

func (t *ElementTriplet) writeKeyValue(sb *strings.Builder) error {
	t.writeKey(sb)
	sb.WriteRune(':')

	if valueBytes, err := json.Marshal(ElementToInterface(&t.value)); err == nil {
		sb.WriteString(string(valueBytes))
		return nil
	} else {
		return err
	}
}

func (m *ElementMap) serializedUsedEntries() (string, error) {
	var sb strings.Builder
	sb.WriteRune('{')

	first := true

	for _, element := range *m {
		if !element.used {
			continue
		}
		if !first {
			sb.WriteRune(',')
		}
		first = false
		sb.WriteString("\n\t")
		if err := element.writeKeyValue(&sb); err != nil {
			return "", err
		}
	}

	sb.WriteString("\n}")

	return sb.String(), nil
}

func (m *ElementMap) SaveUsedEntries(path string) error {

	if s, err := m.serializedUsedEntries(); err != nil {
		return err
	} else {
		return os.WriteFile(path, []byte(s), 0)
	}
}

func (m *ElementMap) sort() {
	sort.Slice(*m, func(i, j int) bool {
		return (*m)[i].CmpKey((*m)[j]) <= 0
	})
}

func (m *ElementMap) find(toFind *ElementTriplet) (small_rational.SmallRational, error) {
	i := sort.Search(len(*m), func(i int) bool { return (*m)[i].CmpKey(toFind) >= 0 })

	if i < len(*m) && (*m)[i].CmpKey(toFind) == 0 {
		(*m)[i].used = true
		return (*m)[i].value, nil
	}
	// if not found, add it:
	if _, err := toFind.value.SetInterface(rand.Int63n(11) - 5); err != nil {
		panic(err.Error())
	}
	toFind.used = true
	*m = append(*m, toFind)
	m.sort() //Inefficient, but it's okay. This is only run when a new test case is introduced

	return toFind.value, nil
}

func (m *ElementMap) FindPair(x *small_rational.SmallRational, y *small_rational.SmallRational) (small_rational.SmallRational, error) {

	toFind := ElementTriplet{
		key1:        *x,
		key2Present: y != nil,
	}

	if y != nil {
		toFind.key2 = *y
	}

	return m.find(&toFind)
}

func ToElement(i int64) *small_rational.SmallRational {
	var res small_rational.SmallRational
	res.SetInt64(i)
	return &res
}

type MessageCounter struct {
	startState uint64
	state      uint64
	step       uint64
}

func (m *MessageCounter) Write(p []byte) (n int, err error) {
	inputBlockSize := (len(p)-1)/small_rational.Bytes + 1
	m.state += uint64(inputBlockSize) * m.step
	return len(p), nil
}

func (m *MessageCounter) Sum(b []byte) []byte {
	inputBlockSize := (len(b)-1)/small_rational.Bytes + 1
	resI := m.state + uint64(inputBlockSize)*m.step
	var res small_rational.SmallRational
	res.SetInt64(int64(resI))
	resBytes := res.Bytes()
	return resBytes[:]
}

func (m *MessageCounter) Reset() {
	m.state = m.startState
}

func (m *MessageCounter) Size() int {
	return small_rational.Bytes
}

func (m *MessageCounter) BlockSize() int {
	return small_rational.Bytes
}

func NewMessageCounter(startState, step int) hash.Hash {
	transcript := &MessageCounter{startState: uint64(startState), step: uint64(step)}
	return transcript
}

func NewMessageCounterGenerator(startState, step int) func() hash.Hash {
	return func() hash.Hash {
		return NewMessageCounter(startState, step)
	}
}

func SliceToElementSlice[T any](slice []T) ([]small_rational.SmallRational, error) {
	elementSlice := make([]small_rational.SmallRational, len(slice))
	for i, v := range slice {
		if _, err := elementSlice[i].SetInterface(v); err != nil {
			return nil, err
		}
	}
	return elementSlice, nil
}

func SliceEquals(a []small_rational.SmallRational, b []small_rational.SmallRational) error {
	if len(a) != len(b) {
		return fmt.Errorf("length mismatch %d≠%d", len(a), len(b))
	}
	for i := range a {
		if !a[i].Equal(&b[i]) {
			return fmt.Errorf("at index %d: %s ≠ %s", i, a[i].String(), b[i].String())
		}
	}
	return nil
}

func SliceSliceEquals(a [][]small_rational.SmallRational, b [][]small_rational.SmallRational) error {
	if len(a) != len(b) {
		return fmt.Errorf("length mismatch %d≠%d", len(a), len(b))
	}
	for i := range a {
		if err := SliceEquals(a[i], b[i]); err != nil {
			return fmt.Errorf("at index %d: %w", i, err)
		}
	}
	return nil
}

func PolynomialSliceEquals(a []polynomial.Polynomial, b []polynomial.Polynomial) error {
	if len(a) != len(b) {
		return fmt.Errorf("length mismatch %d≠%d", len(a), len(b))
	}
	for i := range a {
		if err := SliceEquals(a[i], b[i]); err != nil {
			return fmt.Errorf("at index %d: %w", i, err)
		}
	}
	return nil
}

func ElementToInterface(x *small_rational.SmallRational) interface{} {
	text := x.Text(10)
	if len(text) < 10 && !strings.Contains(text, "/") {
		if i, err := strconv.Atoi(text); err != nil {
			panic(err.Error())
		} else {
			return i
		}
	}
	return text
}

func ElementSliceToInterfaceSlice(x interface{}) []interface{} {
	if x == nil {
		return nil
	}

	X := reflect.ValueOf(x)

	res := make([]interface{}, X.Len())
	for i := range res {
		xI := X.Index(i).Interface().(small_rational.SmallRational)
		res[i] = ElementToInterface(&xI)
	}
	return res
}

func ElementSliceSliceToInterfaceSliceSlice(x interface{}) [][]interface{} {
	if x == nil {
		return nil
	}

	X := reflect.ValueOf(x)

	res := make([][]interface{}, X.Len())
	for i := range res {
		res[i] = ElementSliceToInterfaceSlice(X.Index(i).Interface())
	}

	return res
}
