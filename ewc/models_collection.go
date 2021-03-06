// generated by collection-wrapper -- DO NOT EDIT
package ewc

import (
	"encoding/json"
	"fmt"
)

type MessageCollection struct {
	s []*Message
}

func NewMessageCollection() *MessageCollection {
	return &MessageCollection{}
}

func (v *MessageCollection) Clear() {
	v.s = v.s[:0]
}

func (v *MessageCollection) Equal(rhs *MessageCollection) bool {
	if rhs == nil {
		return false
	}

	if len(v.s) != len(rhs.s) {
		return false
	}

	for i := range v.s {
		if !v.s[i].Equal(rhs.s[i]) {
			return false
		}
	}

	return true
}

func (v *MessageCollection) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.s)
}

func (v *MessageCollection) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &v.s)
}

func (v *MessageCollection) Copy(rhs *MessageCollection) {
	v.s = make([]*Message, len(rhs.s))
	copy(v.s, rhs.s)
}

func (v *MessageCollection) Clone() *MessageCollection {
	return &MessageCollection{
		s: v.s[:],
	}
}

func (v *MessageCollection) Index(rhs *Message) int {
	for i, lhs := range v.s {
		if lhs == rhs {
			return i
		}
	}
	return -1
}

func (v *MessageCollection) Insert(i int, n *Message) {
	if i < 0 || i > len(v.s) {
		fmt.Printf("Vapi::MessageCollection field_values.go error trying to insert at index %d\n", i)
		return
	}
	v.s = append(v.s, nil)
	copy(v.s[i+1:], v.s[i:])
	v.s[i] = n
}

func (v *MessageCollection) Remove(i int) {
	if i < 0 || i >= len(v.s) {
		fmt.Printf("Vapi::MessageCollection field_values.go error trying to remove bad index %d\n", i)
		return
	}
	copy(v.s[i:], v.s[i+1:])
	v.s[len(v.s)-1] = nil
	v.s = v.s[:len(v.s)-1]
}

func (v *MessageCollection) Count() int {
	return len(v.s)
}

func (v *MessageCollection) At(i int) *Message {
	if i < 0 || i >= len(v.s) {
		fmt.Printf("Vapi::MessageCollection field_values.go invalid index %d\n", i)
	}
	return v.s[i]
}

type FriendCollection struct {
	s []*Friend
}

func NewFriendCollection() *FriendCollection {
	return &FriendCollection{}
}

func (v *FriendCollection) Clear() {
	v.s = v.s[:0]
}

func (v *FriendCollection) Equal(rhs *FriendCollection) bool {
	if rhs == nil {
		return false
	}

	if len(v.s) != len(rhs.s) {
		return false
	}

	for i := range v.s {
		if !v.s[i].Equal(rhs.s[i]) {
			return false
		}
	}

	return true
}

func (v *FriendCollection) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.s)
}

func (v *FriendCollection) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &v.s)
}

func (v *FriendCollection) Copy(rhs *FriendCollection) {
	v.s = make([]*Friend, len(rhs.s))
	copy(v.s, rhs.s)
}

func (v *FriendCollection) Clone() *FriendCollection {
	return &FriendCollection{
		s: v.s[:],
	}
}

func (v *FriendCollection) Index(rhs *Friend) int {
	for i, lhs := range v.s {
		if lhs == rhs {
			return i
		}
	}
	return -1
}

func (v *FriendCollection) Insert(i int, n *Friend) {
	if i < 0 || i > len(v.s) {
		fmt.Printf("Vapi::FriendCollection field_values.go error trying to insert at index %d\n", i)
		return
	}
	v.s = append(v.s, nil)
	copy(v.s[i+1:], v.s[i:])
	v.s[i] = n
}

func (v *FriendCollection) Remove(i int) {
	if i < 0 || i >= len(v.s) {
		fmt.Printf("Vapi::FriendCollection field_values.go error trying to remove bad index %d\n", i)
		return
	}
	copy(v.s[i:], v.s[i+1:])
	v.s[len(v.s)-1] = nil
	v.s = v.s[:len(v.s)-1]
}

func (v *FriendCollection) Count() int {
	return len(v.s)
}

func (v *FriendCollection) At(i int) *Friend {
	if i < 0 || i >= len(v.s) {
		fmt.Printf("Vapi::FriendCollection field_values.go invalid index %d\n", i)
	}
	return v.s[i]
}

type ChatCollection struct {
	s []*Chat
}

func NewChatCollection() *ChatCollection {
	return &ChatCollection{}
}

func (v *ChatCollection) Clear() {
	v.s = v.s[:0]
}

func (v *ChatCollection) Equal(rhs *ChatCollection) bool {
	if rhs == nil {
		return false
	}

	if len(v.s) != len(rhs.s) {
		return false
	}

	for i := range v.s {
		if !v.s[i].Equal(rhs.s[i]) {
			return false
		}
	}

	return true
}

func (v *ChatCollection) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.s)
}

func (v *ChatCollection) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &v.s)
}

func (v *ChatCollection) Copy(rhs *ChatCollection) {
	v.s = make([]*Chat, len(rhs.s))
	copy(v.s, rhs.s)
}

func (v *ChatCollection) Clone() *ChatCollection {
	return &ChatCollection{
		s: v.s[:],
	}
}

func (v *ChatCollection) Index(rhs *Chat) int {
	for i, lhs := range v.s {
		if lhs == rhs {
			return i
		}
	}
	return -1
}

func (v *ChatCollection) Insert(i int, n *Chat) {
	if i < 0 || i > len(v.s) {
		fmt.Printf("Vapi::ChatCollection field_values.go error trying to insert at index %d\n", i)
		return
	}
	v.s = append(v.s, nil)
	copy(v.s[i+1:], v.s[i:])
	v.s[i] = n
}

func (v *ChatCollection) Remove(i int) {
	if i < 0 || i >= len(v.s) {
		fmt.Printf("Vapi::ChatCollection field_values.go error trying to remove bad index %d\n", i)
		return
	}
	copy(v.s[i:], v.s[i+1:])
	v.s[len(v.s)-1] = nil
	v.s = v.s[:len(v.s)-1]
}

func (v *ChatCollection) Count() int {
	return len(v.s)
}

func (v *ChatCollection) At(i int) *Chat {
	if i < 0 || i >= len(v.s) {
		fmt.Printf("Vapi::ChatCollection field_values.go invalid index %d\n", i)
	}
	return v.s[i]
}
