package bk

import (
	"fmt"
)

type (
	KeyName      string
	KeyCondition string //KeyCondition represents a requirement for using an item. A KeyCondition is either can_act, or the name of another key

	KeyPhrase string //KeyPhrase is a temporary typename used to indicate a conditional string which requires a parser to pluck conditional logic from

	ItemDistributionSetting struct {
		SlotType string `json:"slot_type,omitempty"`

	}

	// Key represents game state, or player save file state. Anything that can be used to indicate progression, really.
	Key struct {
		Name      KeyName      `json:"name"`             // Name is the human-readable ID of this key.
		Type      string       `json:"type"`             // Type is an extra descriptor for a key that can be added in lieu of listing all required items at once
		Pinned    bool         `json:"pinned,omitempty"` // Pinned indicates if this item can be locked to its location
		Condition KeyCondition `json:"use_condition"`    // Condition is a parseable representation of requirements FOR USING an item or flag

		//TODO: Add Ability to represent item USAGE, kappa
		Settings *ItemDistributionSetting `json:"settings,omitempty"`
	}
)

//ParseRequirements should implement the recursive-descent scanner.
//It should return a KeyNodeList, a tree of items which reflect something similar to lisp syntax for conditionals
func (kp KeyPhrase) ParseRequirements() {
	panic("not yet implemented")
}

//Basic sanity checks
func (k *Key) Validate() error {
	if len(k.Name) == 0 {
		return fmt.Errorf("key missing name")
	}

	return nil
}
