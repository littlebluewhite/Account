package convert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

// Updated Embedded struct remains the same
type EmbeddedSrc struct {
	Num         *int
	StringField *string
}

type EmbeddedDst struct {
	Num         int
	StringField *string
}

// Updated Src struct now includes a TimeField of type time.Time
type Src struct {
	IntField      int
	StringField   *string
	PointerField  *string
	SliceField    []int
	EmbeddedField []EmbeddedSrc
	TimeField     time.Time
}

// Updated Dst struct also includes a TimeField of type time.Time
type Dst struct {
	IntField      int
	StringField   string
	PointerField  *string
	SliceField    []int
	EmbeddedField []EmbeddedDst
	TimeField     time.Time
}

// Updated Model struct includes an UpdatedAt field
type Model struct {
	ID          int32
	Name        string
	Age         int32
	UpdatedAt   time.Time
	Rule        json.RawMessage
	Embedded    []ModelEmbedded
	CommonSlice []int
}

type ModelEmbedded struct {
	ID    int
	Check bool
}

// Updated Update struct includes an UpdatedAt field (pointer to time.Time)
type Update struct {
	ID          int32
	Name        *string
	Age         *int32
	UpdatedAt   *time.Time
	Rule        *json.RawMessage
	Embedded    []UpdateEmbedded
	CommonSlice []int
}

type UpdateEmbedded struct {
	ID    int
	Check *bool
}

// TestUpdateConvert now tests updating of time.Time fields
func TestUpdateConvert(t *testing.T) {
	// Step 1: Initialize models with initial UpdatedAt and CommonSlice values
	initialTime := time.Now().Add(-24 * time.Hour)
	models := []Model{
		{
			ID:          1,
			Name:        "Alice",
			Age:         30,
			UpdatedAt:   initialTime,
			Rule:        []byte("{}"),
			Embedded:    []ModelEmbedded{{ID: 1, Check: true}, {ID: 2, Check: false}},
			CommonSlice: []int{1, 2, 3},
		},
		{
			ID:          2,
			Name:        "Bob",
			Age:         25,
			UpdatedAt:   initialTime,
			Rule:        nil,
			Embedded:    []ModelEmbedded{{ID: 1, Check: true}, {ID: 2, Check: true}},
			CommonSlice: []int{4, 5},
		},
	}

	// Step 2: Define new values for updates, including CommonSlice
	aName := "Alice Updated"
	aAge := int32(31)
	bName := "Bob Updated"
	newTime := time.Now()
	newCheck := false
	var bRule json.RawMessage = []byte(`{"test":1}`)
	aCommonSlice := []int{10, 20}
	bCommonSlice := []int{30, 40}

	// Step 3: Define updates, including updates to Embedded and CommonSlice
	updates := []*Update{
		{
			ID:        1,
			Name:      &aName,
			Age:       &aAge,
			UpdatedAt: &newTime,
			Rule:      nil, // Should not overwrite existing Rule
			Embedded: []UpdateEmbedded{
				{ID: 1, Check: &newCheck}, // Update existing embedded item
				{ID: 3, Check: nil},       // New embedded item with nil Check (should not be added if all fields nil except key)
			},
			CommonSlice: aCommonSlice, // Update CommonSlice for Alice
		},
		{
			ID:        2,
			Name:      &bName,
			Age:       nil, // Should not overwrite existing Age
			UpdatedAt: &newTime,
			Rule:      &bRule,
			Embedded: []UpdateEmbedded{
				{ID: 1, Check: nil},       // Should not overwrite existing Check
				{ID: 2, Check: &newCheck}, // Update existing embedded item
			},
			CommonSlice: bCommonSlice, // Update CommonSlice for Bob
		},
		{
			ID:        3,
			Name:      &bName,
			Age:       nil,
			UpdatedAt: &newTime,
			Rule:      &bRule,
			Embedded: []UpdateEmbedded{
				{ID: 1, Check: nil},
				{ID: 2, Check: &newCheck},
			},
			CommonSlice: bCommonSlice,
		},
	}

	// Step 4: Call UpdateConvert with the correct key "ID" (case-sensitive)
	updatedModels, err := UpdateConvert[Model, *Update](models, updates, "ID")
	if err != nil {
		t.Fatal(err)
	}

	// Step 5: Create a map for easy lookup of updated models by ID
	updatedModelsMap := make(map[int]Model)
	for _, m := range updatedModels {
		updatedModelsMap[int(m.ID)] = m
	}

	// Step 6: Check results for Alice
	alice, exists := updatedModelsMap[1]
	if !exists {
		t.Errorf("Model with ID=1 (Alice) not found in updated models")
	} else {
		// Check if Alice's data has been updated correctly
		if alice.Name != aName {
			t.Errorf("Update failed for Alice's Name; expected: %s, got: %s", aName, alice.Name)
		}
		if alice.Age != aAge {
			t.Errorf("Update failed for Alice's Age; expected: %d, got: %d", aAge, alice.Age)
		}
		if !alice.UpdatedAt.Equal(newTime) {
			t.Errorf("Expected UpdatedAt to be %v, got %v", newTime, alice.UpdatedAt)
		}
		if !bytes.Equal(alice.Rule, []byte("{}")) {
			t.Errorf("Expected Rule to remain unchanged, got %s", alice.Rule)
		}
		if !reflect.DeepEqual(alice.CommonSlice, aCommonSlice) {
			t.Errorf("Expected CommonSlice to be %v, got %v", aCommonSlice, alice.CommonSlice)
		}

		// Check embedded updates for Alice
		if len(alice.Embedded) != 2 {
			t.Errorf("Expected 2 embedded items for Alice, got %d", len(alice.Embedded))
		}
		embeddedMap := make(map[int]ModelEmbedded)
		for _, e := range alice.Embedded {
			embeddedMap[e.ID] = e
		}
		// Check that Embedded item with ID=1 was updated
		if e, ok := embeddedMap[1]; ok {
			if e.Check != newCheck {
				t.Errorf("Expected Embedded[ID=1].Check to be %v, got %v", newCheck, e.Check)
			}
		} else {
			t.Errorf("Embedded item with ID=1 not found for Alice")
		}
		// Check that Embedded item with ID=3 was added
		if e, ok := embeddedMap[3]; ok {
			// Since Check was nil in the update, it should default to false
			if e.Check != false {
				t.Errorf("Expected Embedded[ID=3].Check to be false, got %v", e.Check)
			}
		} else {
			t.Errorf("Embedded item with ID=3 not added for Alice")
		}
	}

	// Step 7: Check results for Bob
	bob, exists := updatedModelsMap[2]
	if !exists {
		t.Errorf("Model with ID=2 (Bob) not found in updated models")
	} else {
		// Check if Bob's data has been updated correctly
		if bob.Name != bName {
			t.Errorf("Update failed for Bob's Name; expected: %s, got: %s", bName, bob.Name)
		}
		if bob.Age != 25 {
			t.Errorf("Bob's Age should remain unchanged; expected: %d, got: %d", 25, bob.Age)
		}
		if !bob.UpdatedAt.Equal(newTime) {
			t.Errorf("Expected UpdatedAt to be %v, got %v", newTime, bob.UpdatedAt)
		}
		if !bytes.Equal(bob.Rule, bRule) {
			t.Errorf("Expected Rule to be %s, got %s", bRule, bob.Rule)
		}
		if !reflect.DeepEqual(bob.CommonSlice, bCommonSlice) {
			t.Errorf("Expected CommonSlice to be %v, got %v", bCommonSlice, bob.CommonSlice)
		}

		// Check embedded updates for Bob
		if len(bob.Embedded) != 2 {
			t.Errorf("Expected 2 embedded items for Bob, got %d", len(bob.Embedded))
		}
		embeddedMap := make(map[int]ModelEmbedded)
		for _, e := range bob.Embedded {
			embeddedMap[e.ID] = e
		}
		// Check that Embedded item with ID=1 remains unchanged
		if e, ok := embeddedMap[1]; ok {
			if e.Check != true {
				t.Errorf("Expected Embedded[ID=1].Check to remain unchanged, got %v", e.Check)
			}
		} else {
			t.Errorf("Embedded item with ID=1 not found for Bob")
		}
		// Check that Embedded item with ID=2 was updated
		if e, ok := embeddedMap[2]; ok {
			if e.Check != newCheck {
				t.Errorf("Expected Embedded[ID=2].Check to be %v, got %v", newCheck, e.Check)
			}
		} else {
			t.Errorf("Embedded item with ID=2 not found for Bob")
		}
	}

	// Step 8: Check results for new Model ID=3
	model3, exists := updatedModelsMap[3]
	if !exists {
		t.Errorf("Model with ID=3 not found in updated models")
	} else {
		// Check Model ID=3's data
		if model3.Name != bName {
			t.Errorf("Update failed for Model ID=3's Name; expected: %s, got: %s", bName, model3.Name)
		}
		if model3.Age != 0 { // Age was nil in update, so should remain zero
			t.Errorf("Expected Model ID=3's Age to remain zero, got: %d", model3.Age)
		}
		if !model3.UpdatedAt.Equal(newTime) {
			t.Errorf("Expected Model ID=3's UpdatedAt to be %v, got %v", newTime, model3.UpdatedAt)
		}
		if !bytes.Equal(model3.Rule, bRule) {
			t.Errorf("Expected Model ID=3's Rule to be %s, got %s", bRule, model3.Rule)
		}
		if !reflect.DeepEqual(model3.CommonSlice, bCommonSlice) {
			t.Errorf("Expected Model ID=3's CommonSlice to be %v, got %v", bCommonSlice, model3.CommonSlice)
		}

		// Check embedded updates for Model ID=3
		if len(model3.Embedded) != 2 {
			t.Errorf("Expected 2 embedded items for Model ID=3, got %d", len(model3.Embedded))
		}
		embeddedMap := make(map[int]ModelEmbedded)
		for _, e := range model3.Embedded {
			embeddedMap[e.ID] = e
		}
		// Check that Embedded item with ID=1 was not added because Check was nil
		if e, ok := embeddedMap[1]; ok {
			if e.Check != false { // Default false
				t.Errorf("Expected Embedded[ID=1].Check to be false, got %v", e.Check)
			}
		} else {
			t.Errorf("Embedded item with ID=1 not added for Model ID=3")
		}
		// Check that Embedded item with ID=2 was updated
		if e, ok := embeddedMap[2]; ok {
			if e.Check != newCheck {
				t.Errorf("Expected Embedded[ID=2].Check to be %v, got %v", newCheck, e.Check)
			}
		} else {
			t.Errorf("Embedded item with ID=2 not found for Model ID=3")
		}
	}

	// Optionally, print updated models for debugging
	fmt.Printf("Updated Models:\n")
	for _, m := range updatedModels {
		mJSON, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			t.Errorf("Error marshalling model: %v", err)
			continue
		}
		fmt.Printf("%s\n", mJSON)
	}
}

func TestUpdateFieldConvert(t *testing.T) {
	hello := "hello"
	currentTime := time.Now()
	numValue := 10
	testStr := "embedded_string"

	// Source struct with TimeField initialized
	src := Src{
		IntField:     42,
		StringField:  nil,
		PointerField: &hello,
		SliceField:   []int{1, 2, 3},
		EmbeddedField: []EmbeddedSrc{
			{Num: &numValue, StringField: nil},
			{Num: &numValue, StringField: &testStr},
		},
		TimeField: currentTime,
	}

	dst := Dst{
		IntField:      30,
		StringField:   "test",
		PointerField:  &hello,
		SliceField:    []int{1, 2, 3},
		EmbeddedField: []EmbeddedDst{},
		TimeField:     currentTime}

	srcVal := reflect.ValueOf(&src).Elem()
	dstVal := reflect.ValueOf(&dst).Elem()

	// Call UpdateFieldConvert
	UpdateFieldConvert(srcVal, dstVal)

	fmt.Printf("srcVal: %+v\n", srcVal)
	fmt.Printf("dstVal: %+v\n", dstVal)

	// Check results
	if dst.IntField != src.IntField {
		t.Errorf("Expected IntField %d, got %d", src.IntField, dst.IntField)
	}
	if dst.StringField != "test" {
		t.Errorf("Expected StringField '%s', got '%s'", "test", dst.StringField)
	}
	if dst.PointerField == nil || *dst.PointerField != *src.PointerField {
		t.Errorf("Expected PointerField '%s', got '%v'", *src.PointerField, dst.PointerField)
	}
	if len(dst.SliceField) != len(src.SliceField) {
		t.Errorf("Expected SliceField length %d, got %d", len(src.SliceField), len(dst.SliceField))
	} else {
		for i := range dst.SliceField {
			if dst.SliceField[i] != src.SliceField[i] {
				t.Errorf("Expected SliceField[%d] %d, got %d", i, src.SliceField[i], dst.SliceField[i])
			}
		}
	}
	if len(dst.EmbeddedField) != len(src.EmbeddedField) {
		t.Errorf("Expected EmbeddedField length %d, got %d", len(src.EmbeddedField), len(dst.EmbeddedField))
	} else {
		for i := range dst.EmbeddedField {
			if dst.EmbeddedField[i].Num != *src.EmbeddedField[i].Num {
				t.Errorf("Expected EmbeddedField[%d].Num %d, got %d", i, *src.EmbeddedField[i].Num, dst.EmbeddedField[i].Num)
			}
			if dst.EmbeddedField[i].StringField == nil || *dst.EmbeddedField[i].StringField != *src.EmbeddedField[i].StringField {
				t.Errorf("Expected EmbeddedField[%d].StringField '%s', got '%v'", i, *src.EmbeddedField[i].StringField, dst.EmbeddedField[i].StringField)
			}
		}
	}
	if !dst.TimeField.Equal(src.TimeField) {
		t.Errorf("Expected TimeField %v, got %v", src.TimeField, dst.TimeField)
	}
}
