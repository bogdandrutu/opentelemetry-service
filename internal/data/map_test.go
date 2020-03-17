// Copyright 2020 OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package data

import (
	"strconv"
	"testing"

	otlpcommon "github.com/open-telemetry/opentelemetry-proto/gen/go/common/v1"
)

type AttributeV2 struct {
	orig *otlpcommon.AttributeKeyValue
}

func newAttributeV2(k string, av AttributeValue) AttributeValue {
	return AttributeValue{orig: &otlpcommon.AttributeKeyValue{
		Key:         k,
		Type:        av.orig.Type,
		StringValue: av.orig.StringValue,
		IntValue:    av.orig.IntValue,
		DoubleValue: av.orig.DoubleValue,
		BoolValue:   av.orig.BoolValue,
	}}
}

func newAttributeV2String(k, v string) AttributeValue {
	return AttributeValue{orig: &otlpcommon.AttributeKeyValue{Key: k, Type: otlpcommon.AttributeKeyValue_STRING, StringValue: v}}
}

func (a AttributeV2) Type() AttributeValueType {
	return AttributeValueType(a.orig.Type)
}

func (a AttributeV2) StringVal() string {
	return a.orig.StringValue
}

func (a AttributeV2) SetStringVal(v string) {
	a.orig.Type = otlpcommon.AttributeKeyValue_STRING
	a.orig.StringValue = v
}

func (a AttributeV2) Key() string {
	return a.orig.Key
}

func (a AttributeV2) SetValue(av AttributeValue) {
	a.orig.Type = av.orig.Type
	a.orig.StringValue = av.orig.StringValue
	a.orig.IntValue = av.orig.IntValue
	a.orig.DoubleValue = av.orig.DoubleValue
	a.orig.BoolValue = av.orig.BoolValue
}

type AttributeMapV2 struct {
	origs *[]*otlpcommon.AttributeKeyValue
}

func (am AttributeMapV2) GetValue(k string) (AttributeValue, bool) {
	for _, a := range *am.origs {
		if a.Key == k {
			return AttributeValue{a}, true
		}
	}
	return AttributeValue{nil}, false
}

func (am AttributeMapV2) Get(k string) (AttributeV2, bool) {
	for _, a := range *am.origs {
		if a.Key == k {
			return AttributeV2{a}, true
		}
	}
	return AttributeV2{nil}, false
}

func (am AttributeMapV2) Delete(k string) bool {
	for i, a := range *am.origs {
		if a.Key == k {
			(*am.origs)[i] = (*am.origs)[len(*am.origs)-1]
			*am.origs = (*am.origs)[:len(*am.origs)-1]
			return true
		}
	}
	return false
}

func (am AttributeMapV2) InsertValue(k string, av AttributeValue) {
	if _, existing := am.Get(k); !existing {
		*am.origs = append(*am.origs, newAttributeV2(k, av).orig)
	}
}

func (am AttributeMapV2) UpdateValue(k string, av AttributeValue) {
	if attr, existing := am.Get(k); existing {
		attr.SetValue(av)
	}
}

func (am AttributeMapV2) UpsertValue(k string, av AttributeValue) {
	if attr, existing := am.Get(k); existing {
		attr.SetValue(av)
	} else {
		*am.origs = append(*am.origs, newAttributeV2(k, av).orig)
	}
}

func (am AttributeMapV2) AttributesCount() int {
	return len(*am.origs)
}

func (am AttributeMapV2) GetAttribute(ix int) AttributeV2 {
	return AttributeV2{(*am.origs)[ix]}
}

type AttributeMapV1 struct {
	// Cannot track changes in the map, so if this is initialized we
	// always reconstruct the labels.
	attributesMap map[string]AttributeValue
	// True if the pimpl was initialized.
	initialized bool
}

// NewLabels creates a new Labels.
func NewAttributeMapV1() *AttributeMapV1 {
	return &AttributeMapV1{nil, false}
}

func (am *AttributeMapV1) initAndGet(orig []*otlpcommon.AttributeKeyValue) map[string]AttributeValue {
	if !am.initialized {
		if len(orig) == 0 {
			am.attributesMap = map[string]AttributeValue{}
			am.initialized = true
			return am.attributesMap
		}
		// Extra overhead here if we decode the orig attributes
		// then immediately overwrite them in set.
		labels := make(map[string]AttributeValue, len(orig))
		for i := range orig {
			labels[orig[i].Key] = AttributeValue{orig[i]}
		}
		am.attributesMap = labels
		am.initialized = true
	}
	return am.attributesMap
}

func (am *AttributeMapV1) toOrig(orig []*otlpcommon.AttributeKeyValue) []*otlpcommon.AttributeKeyValue {
	if !am.initialized {
		// Guaranteed no changes via internal fields.
		return orig
	}

	if len(orig) != len(am.attributesMap) {
		skvs := make([]otlpcommon.AttributeKeyValue, len(am.attributesMap))
		orig = make([]*otlpcommon.AttributeKeyValue, len(am.attributesMap))
		for i := range orig {
			orig[i] = &skvs[i]
		}
	}

	i := 0
	for k, v := range am.attributesMap {
		orig[i].Key = k
		orig[i].Type = v.orig.Type
		orig[i].StringValue = v.orig.StringValue
		orig[i].BoolValue = v.orig.BoolValue
		orig[i].IntValue = v.orig.IntValue
		orig[i].DoubleValue = v.orig.DoubleValue
		i++
	}

	return orig
}

const numAttributes = 200

func BenchmarkMapV1(b *testing.B) {
	ap := generateAttributesProcessor()
	akvs := generateAttributes(numAttributes)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		am := NewAttributeMapV1()
		ap.ConsumeAttributeMapV1(am.initAndGet(akvs))
		// Sync to start from the same map
		akvs = am.toOrig(akvs)
	}
}

func BenchmarkMapV1_NoAlloc(b *testing.B) {
	ap := generateAttributesProcessor()
	akvs := generateAttributes(numAttributes)
	ams := NewAttributeMapV1()
	if _, exists := ams.initAndGet(akvs)["key_1"]; !exists {
		b.Fatal("Cannot happen")
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		// akvs is not touched here because all the maps are initialized
		ap.ConsumeAttributeMapV1(ams.initAndGet(akvs))
	}
}

func BenchmarkMapV2(b *testing.B) {
	ap := generateAttributesProcessor()
	akvs := generateAttributes(numAttributes)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ap.ConsumeAttributeMapV2(AttributeMapV2{&akvs})
	}
}

func generateAttributes(len int) []*otlpcommon.AttributeKeyValue {
	akvs := make([]*otlpcommon.AttributeKeyValue, len)
	for i := range akvs {
		akvs[i] = generateAttribute(i)
	}
	return akvs
}

func generateAttribute(ix int) *otlpcommon.AttributeKeyValue {
	return &otlpcommon.AttributeKeyValue{Key: "key_" + strconv.Itoa(ix), Type: otlpcommon.AttributeKeyValue_STRING, StringValue: "value_" + strconv.Itoa(ix)}
}

type attributesProcessor struct {
	actions []attributeAction
}

func generateAttributesProcessor() attributesProcessor {
	return attributesProcessor{[]attributeAction{
		attributeAction{
			Key:    "key_10",
			Action: DELETE,
		},
		attributeAction{
			Key:            "key_7",
			AttributeValue: NewAttributeValueBool(true),
			Action:         UPSERT,
		},
		attributeAction{
			Key:           "key_15",
			FromAttribute: "key_13",
			Action:        UPDATE,
		},
		attributeAction{
			Key:            "key_15",
			AttributeValue: NewAttributeValueString("value_15"),
			Action:         UPSERT,
		},
		attributeAction{
			Key:    "key_24",
			Action: DELETE,
		},
		attributeAction{
			Key:            "key_24",
			AttributeValue: NewAttributeValueString("value_24"),
			Action:         INSERT,
		},
		attributeAction{
			Key:            "key_7",
			AttributeValue: NewAttributeValueString("value_7"),
			Action:         UPSERT,
		},
		attributeAction{
			Key:            "key_10",
			AttributeValue: NewAttributeValueString("value_10"),
			Action:         INSERT,
		},
	}}
}

type attributeAction struct {
	Key           string
	FromAttribute string
	// TODO https://github.com/open-telemetry/opentelemetry-collector/issues/296
	// Do benchmark testing between having action be of type string vs integer.
	// The reason is attributes processor will most likely be commonly used
	// and could impact performance.
	Action         Action
	AttributeValue AttributeValue
}

func (aa *attributeAction) getAttributeValue(attributesMap map[string]AttributeValue) (AttributeValue, bool) {
	// Set the key with a value from the configuration.
	if aa.AttributeValue.orig != nil {
		return aa.AttributeValue, true
	}
	// Don't know why return attributesMap[aa.FromAttribute] does not work
	value, existing := attributesMap[aa.FromAttribute]
	return value, existing
}

func (aa *attributeAction) getAttributeValueV2(am AttributeMapV2) (AttributeValue, bool) {
	// Set the key with a value from the configuration.
	if aa.AttributeValue.orig != nil {
		return aa.AttributeValue, true
	}
	return am.GetValue(aa.FromAttribute)
}

func (ac *attributesProcessor) ConsumeAttributeMapV2(am AttributeMapV2) {
	for _, action := range ac.actions {
		if av, existing := action.getAttributeValueV2(am); existing {
			switch action.Action {
			case DELETE:
				am.Delete(action.Key)
			case INSERT:
				am.InsertValue(action.Key, av)
			case UPDATE:
				am.UpdateValue(action.Key, av)
			case UPSERT:
				// There is no need to check if the target key exists in the attribute map
				// because the value is to be set regardless.
				am.UpsertValue(action.Key, av)
			case HASH:
				av.SetString("hashed_value")
			}
		}
	}
}

func (ac *attributesProcessor) ConsumeAttributeMapV1(attributesMap map[string]AttributeValue) {
	for _, action := range ac.actions {
		switch action.Action {
		case DELETE:
			delete(attributesMap, action.Key)
		case INSERT:
			insertAttribute(action, attributesMap)
		case UPDATE:
			updateAttribute(action, attributesMap)
		case UPSERT:
			// There is no need to check if the target key exists in the attribute map
			// because the value is to be set regardless.
			setAttribute(action, attributesMap)
		case HASH:
			hashAttribute(action, attributesMap)
		}
	}
}

func insertAttribute(action attributeAction, attributesMap map[string]AttributeValue) {
	// Insert is only performed when the target key does not already exist
	// in the attribute map.
	if _, exists := attributesMap[action.Key]; exists {
		return
	}

	setAttribute(action, attributesMap)
}

func updateAttribute(action attributeAction, attributesMap map[string]AttributeValue) {
	// Update is only performed when the target key already exists in
	// the attribute map.
	if _, exists := attributesMap[action.Key]; !exists {
		return
	}

	setAttribute(action, attributesMap)
}

func setAttribute(action attributeAction, attributesMap map[string]AttributeValue) {
	// Set the key with a value from the configuration.
	if action.AttributeValue.orig != nil {
		attributesMap[action.Key] = action.AttributeValue
	} else if value, fromAttributeExists := attributesMap[action.FromAttribute]; fromAttributeExists {
		// Set the key with a value from another attribute, if it exists.
		attributesMap[action.Key] = value
	}
}

func hashAttribute(action attributeAction, attributesMap map[string]AttributeValue) {
	if value, exists := attributesMap[action.Key]; exists {
		value.SetString("hashed_value")
		attributesMap[action.Key] = value
	}
}

// Action is the enum to capture the four types of actions to perform on an
// attribute.
type Action string

const (
	// INSERT adds the key/value to spans when the key does not exist.
	// No action is applied to spans where the key already exists.
	INSERT Action = "insert"

	// UPDATE updates an existing key with a value. No action is applied
	// to spans where the key does not exist.
	UPDATE Action = "update"

	// UPSERT performs the INSERT or UPDATE action. The key/value is
	// insert to spans that did not originally have the key. The key/value is
	// updated for spans where the key already existed.
	UPSERT Action = "upsert"

	// DELETE deletes the attribute from the span. If the key doesn't exist,
	//no action is performed.
	DELETE Action = "delete"

	// HASH calculates the SHA-1 hash of an existing value and overwrites the value
	// with it's SHA-1 hash result.
	HASH Action = "hash"
)
