package main

import (
	"testing"
)

func TestSuccess(t *testing.T) {
    // Just a simple test to confirm everything is set up correctly
    t.Run("Success", func(t *testing.T) {
        t.Log("This test case is a simple success test")
        // No actual logic here; the test should just pass
    })
}

func TestSuccessOne(t *testing.T) {
    // Just a simple test to confirm everything is set up correctly
    t.Run("Success", func(t *testing.T) {
        t.Log("This test case is a simple success test")
        // No actual logic here; the test should just pass
    })
}

func TestSuccessSecond(t *testing.T) {
    // Just a simple test to confirm everything is set up correctly
    t.Run("Success", func(t *testing.T) {
        t.Log("This test case is a simple success test")
        // No actual logic here; the test should just pass
    })
}

// func TestSimpleFail(t *testing.T) {
//     // Simple assertion that is intended to fail
//     expected := 1
//     actual := 2

//     if expected != actual {
//         t.Errorf("Expected %d, but got %d", expected, actual)
//     }
// }