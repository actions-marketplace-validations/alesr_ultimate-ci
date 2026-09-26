package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFoo verifies the Foo function returns the correct greeting.
func TestFoo(t *testing.T) {
	result := Foo()
	assert.Equal(t, "Hello, World!", result)
}

// TestFoo_ContainsGreeting verifies the greeting is present.
func TestFoo_ContainsGreeting(t *testing.T) {
	result := Foo()
	assert.Contains(t, result, "Hello")
}

// TestBar verifies the Bar addition function.
func TestBar(t *testing.T) {
	tests := []struct {
		name     string
		x, y     int
		expected int
	}{
		{"positive numbers", 1, 2, 3},
		{"zero", 0, 5, 5},
		{"negative", -1, 5, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Bar(tt.x, tt.y)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestBar_EdgeCase verifies Bar with large numbers.
func TestBar_EdgeCase(t *testing.T) {
	result := Bar(1000000, 2000000)
	assert.Equal(t, 3000000, result)
}

// TestBaz_Success verifies Baz with valid input.
func TestBaz_Success(t *testing.T) {
	result, err := Baz(5)
	require.NoError(t, err)
	assert.Equal(t, 10, result)
}

// TestBaz_Error verifies Baz error handling with negative input.
func TestBaz_Error(t *testing.T) {
	result, err := Baz(-1)
	assert.Error(t, err)
	assert.Equal(t, 0, result)
	assert.Contains(t, err.Error(), "negative")
}

// TestBaz_Zero verifies Baz with zero input.
func TestBaz_Zero(t *testing.T) {
	result, err := Baz(0)
	require.NoError(t, err)
	assert.Equal(t, 0, result)
}
