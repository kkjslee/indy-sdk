/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package json

import (
	"encoding/json"
	"fmt"
)

// Message is a raw JSON message
type Message json.RawMessage

// Map is a map of messages keyed by string
type Map map[string]*Message

// M converts the map into a message
func (m Map) M() *Message {
	return New(marshal(m))
}

// JSON converts the map into a JSON string
func (m Map) JSON() string {
	return string(marshal(m))
}

// Slice is a slice of messages
type Slice []*Message

// M converts the slice into a message
func (s Slice) M() *Message {
	jsnew := "["
	for i, msg := range s {
		if i > 0 {
			jsnew += ","
		}
		jsnew += fmt.Sprintf("%s", *msg)
	}
	jsnew += "]"
	return New(jsnew)
}

// New creates a new Message from the given string
func New(str string) *Message {
	var msg *json.RawMessage
	err := json.Unmarshal([]byte(str), &msg)
	if err != nil {
		panic(fmt.Sprintf("Error unmarshalling message: %s", err))
	}
	var jmsg Message
	jmsg = Message(*msg)
	return &jmsg
}

// AsMap returns a message map from the given JSON string.
func AsMap(json string) Map {
	return New(json).AsMap()
}

// Str converts the given string to a Message
func Str(str string) *Message {
	return New(fmt.Sprintf(`"%s"`, str))
}

// Int converts the given int to a Message
func Int(i int64) *Message {
	return New(fmt.Sprintf("%d", i))
}

// Bool converts the given bool to a Message
func Bool(b bool) *Message {
	return New(fmt.Sprintf("%t", b))
}

// Val returns the map value for the given key.
// If the message is not a JSON object then a panic occurs.
func (msg *Message) Val(key string) *Message {
	return msg.AsMap()[key]
}

// AsMap returns a map of Messages.
// A panic occurs if the message is not a JSON object.
func (msg *Message) AsMap() Map {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(*msg, &m)
	if err != nil {
		panic(fmt.Sprintf("Error unmarshaling map: %s", err))
	}

	jmap := make(Map)
	for key, val := range m {
		var jmsg Message
		if val != nil {
			jmsg = Message(*val)
		}
		jmap[key] = &jmsg
	}
	return jmap
}

// Idx returns the message at the given index
// A panic occurs if the message is not a JSON array.
func (msg *Message) Idx(index int) *Message {
	return msg.Slice()[index]
}

// Slice returns a Slice of messages
// A panic occurs if the message is not a JSON array.
func (msg *Message) Slice() Slice {
	var l []*json.RawMessage
	err := json.Unmarshal(*msg, &l)
	if err != nil {
		panic(fmt.Sprintf("Error unmarshalling slice: %s", err))
	}

	var s Slice
	for _, val := range l {
		var jmsg Message
		jmsg = Message(*val)
		s = append(s, &jmsg)
	}
	return s
}

// String returns the string value of the message.
// If the message is not a string then a panic occurs.
func (msg *Message) String() string {
	var str string
	err := json.Unmarshal(*msg, &str)
	if err != nil {
		panic(fmt.Sprintf("Error unmarshalling string: %s", err))
	}
	return str
}

// Int64 returns the int64 value of the JSON message
// If the message is not an int then a panic occurs.
func (msg *Message) Int64() int64 {
	var i int64
	err := json.Unmarshal(*msg, &i)
	if err != nil {
		panic(fmt.Sprintf("Error unmarshalling int: %s", err))
	}
	return i
}

// Int returns the int value of the JSON message
// If the message is not an int then a panic occurs.
func (msg *Message) Int() int {
	return int(msg.Int64())
}

func marshal(msgMap Map) string {
	i := 0
	jsnew := "{"
	for key, msg := range msgMap {
		if i > 0 {
			jsnew += ","
		}
		jsnew += fmt.Sprintf(`"%s":%s`, key, *msg)
		i++
	}
	jsnew += "}"
	return jsnew
}
