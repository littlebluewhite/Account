package convert

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

// Define test structs
type NestedStructSrc struct {
	Value       int
	ExpiredTime time.Time
}

type NestedStructDst struct {
	Value       int
	ExpiredTime time.Time
}

type TestSrc struct {
	PtrEmptyValue *string
	IntValue      int
	OtherValue    *int // Changed to *int to match TestDst
	StringValue   string
	NestedPtr     *NestedStructSrc
	Nested        NestedStructSrc
	PtrValue      *string
	SliceValue    []int
}

type TestDst struct {
	ID            int
	PtrEmptyValue *string
	IntValue      int
	OtherValue    *int
	StringValue   string
	NestedPtr     *NestedStructDst
	Nested        NestedStructDst
	PtrValue      *string
	SliceValue    []int
	Empty         string
}

func TestCreateFieldConvert(t *testing.T) {
	t.Run("Complete Field Transfer", func(t *testing.T) {
		// Initialize time for testing
		expiredTime := time.Date(2025, 10, 30, 9, 32, 59, 205999000, time.UTC)

		// Initialize nested source struct
		nestedSrc := &NestedStructSrc{
			Value:       42,
			ExpiredTime: expiredTime,
		}

		// Initialize source struct
		str := "Hello"
		slice := []int{1, 2, 3}
		other := 5
		src := TestSrc{
			PtrEmptyValue: nil,        // Should remain nil in destination
			IntValue:      10,         // Direct copy
			OtherValue:    &other,     // Pointer copy
			StringValue:   "test",     // Direct copy
			Nested:        *nestedSrc, // Struct copy
			NestedPtr:     nestedSrc,  // Pointer to struct copy
			PtrValue:      &str,       // Pointer copy
			SliceValue:    slice,      // Slice copy
		}

		// Initialize destination struct
		dst := TestDst{}

		// Perform the field conversion
		srcVal := reflect.ValueOf(&src).Elem()
		dstVal := reflect.ValueOf(&dst).Elem()

		CreateFieldConvert(srcVal, dstVal)

		// Debug output
		fmt.Printf("srcVal: %+v, dstVal: %+v\n", srcVal, dstVal)

		// Check direct field values
		if dst.IntValue != src.IntValue {
			t.Errorf("Expected IntValue %d, got %d", src.IntValue, dst.IntValue)
		}
		if dst.StringValue != src.StringValue {
			t.Errorf("Expected StringValue '%s', got '%s'", src.StringValue, dst.StringValue)
		}

		// Check nested struct (non-pointer)
		if dst.Nested.Value != src.Nested.Value {
			t.Errorf("Expected Nested.Value %d, got %d", src.Nested.Value, dst.Nested.Value)
		}
		if !dst.Nested.ExpiredTime.Equal(src.Nested.ExpiredTime) {
			t.Errorf("Expected Nested.ExpiredTime %v, got %v", src.Nested.ExpiredTime, dst.Nested.ExpiredTime)
		}

		// Check nested pointer
		if dst.NestedPtr == nil {
			t.Errorf("Expected NestedPtr to be non-nil")
		} else {
			if dst.NestedPtr.Value != src.NestedPtr.Value {
				t.Errorf("Expected NestedPtr.Value %d, got %d", src.NestedPtr.Value, dst.NestedPtr.Value)
			}
			if !dst.NestedPtr.ExpiredTime.Equal(src.NestedPtr.ExpiredTime) {
				t.Errorf("Expected NestedPtr.ExpiredTime %v, got %v", src.NestedPtr.ExpiredTime, dst.NestedPtr.ExpiredTime)
			}
		}

		// Check pointer values
		if dst.PtrValue == nil || *dst.PtrValue != *src.PtrValue {
			t.Errorf("Expected PtrValue '%s', got '%v'", *src.PtrValue, dst.PtrValue)
		}

		// Check OtherValue
		if dst.OtherValue == nil || *dst.OtherValue != *src.OtherValue {
			t.Errorf("Expected OtherValue %d, got %v", *src.OtherValue, dst.OtherValue)
		}

		// Check slice values
		if len(dst.SliceValue) != len(src.SliceValue) {
			t.Errorf("Expected Slice length %d, got %d", len(src.SliceValue), len(dst.SliceValue))
		}
		for i, v := range dst.SliceValue {
			if v != src.SliceValue[i] {
				t.Errorf("Expected SliceValue[%d] %d, got %d", i, src.SliceValue[i], v)
			}
		}

		// Check Empty field (should remain default)
		if dst.Empty != "" {
			t.Errorf("Expected Empty to be '', got '%s'", dst.Empty)
		}

		// Check PtrEmptyValue (should remain nil)
		if dst.PtrEmptyValue != nil {
			t.Errorf("Expected PtrEmptyValue to be nil, got '%v'", *dst.PtrEmptyValue)
		}
	})
}

func TestCreateConvert(t *testing.T) {
	t.Run("Complete Slice Field Transfer", func(t *testing.T) {
		// Initialize specific times for testing time.Time fields
		expiredTime1 := time.Date(2025, 10, 30, 9, 32, 59, 205999000, time.UTC)
		expiredTime2 := time.Date(2026, 5, 15, 14, 0, 0, 0, time.UTC)

		// Initialize nested source structs
		nestedSrc1 := &NestedStructSrc{
			Value:       42,
			ExpiredTime: expiredTime1,
		}

		nestedSrc2 := &NestedStructSrc{
			Value:       84,
			ExpiredTime: expiredTime2,
		}

		// Initialize source slice with various field types and values
		str1 := "Hello"
		str2 := "World"
		other1 := 100
		other2 := 200
		sources := []*TestSrc{
			{
				PtrEmptyValue: nil, // Should remain nil in destination
				IntValue:      10,
				OtherValue:    &other1,
				StringValue:   "test1",
				Nested:        *nestedSrc1,
				NestedPtr:     nestedSrc1,
				PtrValue:      &str1,
				SliceValue:    []int{1, 2, 3},
			},
			{
				PtrEmptyValue: &str2,
				IntValue:      20,
				OtherValue:    &other2,
				StringValue:   "test2",
				Nested:        *nestedSrc2,
				NestedPtr:     nestedSrc2,
				PtrValue:      nil, // Should remain nil in destination
				SliceValue:    []int{4, 5, 6},
			},
			{
				PtrEmptyValue: nil,
				IntValue:      30,
				OtherValue:    nil, // Should remain nil in destination
				StringValue:   "test3",
				Nested:        NestedStructSrc{}, // Empty struct
				NestedPtr:     nil,               // Should remain nil in destination
				PtrValue:      nil,               // Should remain nil in destination
				SliceValue:    nil,               // Should remain nil in destination
			},
		}

		// Perform the conversion using CreateConvert
		destinations := CreateConvert[TestDst, TestSrc](sources)

		// Verify that the number of destinations matches the number of sources
		if len(destinations) != len(sources) {
			t.Fatalf("Expected %d destinations, got %d", len(sources), len(destinations))
		}

		// Iterate over each source and corresponding destination to verify field values
		for i, src := range sources {
			dst := destinations[i]
			if dst == nil {
				t.Fatalf("Destination at index %d is nil", i)
			}

			// Debug output (optional)
			fmt.Printf("Source[%d]: %+v\nDestination[%d]: %+v\n", i, src, i, dst)

			// Check PtrEmptyValue
			if src.PtrEmptyValue == nil {
				if dst.PtrEmptyValue != nil {
					t.Errorf("Destination PtrEmptyValue at index %d expected nil, got %v", i, *dst.PtrEmptyValue)
				}
			} else {
				if dst.PtrEmptyValue == nil {
					t.Errorf("Destination PtrEmptyValue at index %d expected %v, got nil", i, *src.PtrEmptyValue)
				} else if *dst.PtrEmptyValue != *src.PtrEmptyValue {
					t.Errorf("Destination PtrEmptyValue at index %d expected %v, got %v", i, *src.PtrEmptyValue, *dst.PtrEmptyValue)
				}
			}

			// Check IntValue
			if dst.IntValue != src.IntValue {
				t.Errorf("Destination IntValue at index %d expected %d, got %d", i, src.IntValue, dst.IntValue)
			}

			// Check OtherValue
			if src.OtherValue == nil {
				if dst.OtherValue != nil {
					t.Errorf("Destination OtherValue at index %d expected nil, got %v", i, *dst.OtherValue)
				}
			} else {
				if dst.OtherValue == nil {
					t.Errorf("Destination OtherValue at index %d expected %v, got nil", i, *src.OtherValue)
				} else if *dst.OtherValue != *src.OtherValue {
					t.Errorf("Destination OtherValue at index %d expected %d, got %d", i, *src.OtherValue, *dst.OtherValue)
				}
			}

			// Check StringValue
			if dst.StringValue != src.StringValue {
				t.Errorf("Destination StringValue at index %d expected %s, got %s", i, src.StringValue, dst.StringValue)
			}

			// Check Nested (struct)
			if dst.Nested.Value != src.Nested.Value {
				t.Errorf("Destination Nested.Value at index %d expected %d, got %d", i, src.Nested.Value, dst.Nested.Value)
			}
			if !dst.Nested.ExpiredTime.Equal(src.Nested.ExpiredTime) {
				t.Errorf("Destination Nested.ExpiredTime at index %d expected %v, got %v", i, src.Nested.ExpiredTime, dst.Nested.ExpiredTime)
			}

			// Check NestedPtr
			if src.NestedPtr == nil {
				if dst.NestedPtr != nil {
					t.Errorf("Destination NestedPtr at index %d expected nil, got %+v", i, *dst.NestedPtr)
				}
			} else {
				if dst.NestedPtr == nil {
					t.Errorf("Destination NestedPtr at index %d expected %+v, got nil", i, *src.NestedPtr)
				} else {
					if dst.NestedPtr.Value != src.NestedPtr.Value {
						t.Errorf("Destination NestedPtr.Value at index %d expected %d, got %d", i, src.NestedPtr.Value, dst.NestedPtr.Value)
					}
					if !dst.NestedPtr.ExpiredTime.Equal(src.NestedPtr.ExpiredTime) {
						t.Errorf("Destination NestedPtr.ExpiredTime at index %d expected %v, got %v", i, src.NestedPtr.ExpiredTime, dst.NestedPtr.ExpiredTime)
					}
				}
			}

			// Check PtrValue
			if src.PtrValue == nil {
				if dst.PtrValue != nil {
					t.Errorf("Destination PtrValue at index %d expected nil, got %v", i, *dst.PtrValue)
				}
			} else {
				if dst.PtrValue == nil {
					t.Errorf("Destination PtrValue at index %d expected %v, got nil", i, *src.PtrValue)
				} else if *dst.PtrValue != *src.PtrValue {
					t.Errorf("Destination PtrValue at index %d expected %v, got %v", i, *src.PtrValue, *dst.PtrValue)
				}
			}

			// Check SliceValue
			if src.SliceValue == nil {
				if dst.SliceValue != nil {
					t.Errorf("Destination SliceValue at index %d expected nil, got %v", i, dst.SliceValue)
				}
			} else {
				if dst.SliceValue == nil {
					t.Errorf("Destination SliceValue at index %d expected %v, got nil", i, src.SliceValue)
				} else if !reflect.DeepEqual(dst.SliceValue, src.SliceValue) {
					t.Errorf("Destination SliceValue at index %d expected %v, got %v", i, src.SliceValue, dst.SliceValue)
				}
			}

			// Check Empty field (should remain zero value)
			if dst.Empty != "" {
				t.Errorf("Destination Empty at index %d expected '', got '%s'", i, dst.Empty)
			}

			// Check ID field (should remain zero value)
			if dst.ID != 0 {
				t.Errorf("Destination ID at index %d expected 0, got %d", i, dst.ID)
			}
		}
	})
}
