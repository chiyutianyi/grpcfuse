/*
 * Copyright 2022 Han Xin, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test utility functions for the utils package
func TestStringUtils(t *testing.T) {
	t.Run("IsEmpty", func(t *testing.T) {
		assert.True(t, IsEmpty(""), "Empty string should be empty")
		assert.True(t, IsEmpty("   "), "Whitespace string should be empty")
		assert.False(t, IsEmpty("hello"), "Non-empty string should not be empty")
		assert.False(t, IsEmpty("  hello  "), "String with content should not be empty")
	})

	t.Run("TruncateString", func(t *testing.T) {
		assert.Equal(t, "hello", TruncateString("hello world", 5), "Should truncate to 5 characters")
		assert.Equal(t, "hello world", TruncateString("hello world", 20), "Should not truncate if length is sufficient")
		assert.Equal(t, "", TruncateString("hello world", 0), "Should return empty string for length 0")
		assert.Equal(t, "hello world", TruncateString("hello world", -1), "Should not truncate for negative length")
	})

	t.Run("SafeString", func(t *testing.T) {
		assert.Equal(t, "hello", SafeString("hello"), "Should return original string if not nil")
		assert.Equal(t, "", SafeString(nil), "Should return empty string for nil")
	})
}

func TestTimeUtils(t *testing.T) {
	t.Run("FormatDuration", func(t *testing.T) {
		duration := 2*time.Hour + 30*time.Minute + 45*time.Second
		formatted := FormatDuration(duration)
		
		assert.Contains(t, formatted, "2h", "Should contain hours")
		assert.Contains(t, formatted, "30m", "Should contain minutes")
		assert.Contains(t, formatted, "45s", "Should contain seconds")
	})

	t.Run("ParseDuration", func(t *testing.T) {
		duration, err := ParseDuration("2h30m45s")
		require.NoError(t, err, "Should parse valid duration")
		
		expected := 2*time.Hour + 30*time.Minute + 45*time.Second
		assert.Equal(t, expected, duration, "Should parse duration correctly")
	})

	t.Run("IsExpired", func(t *testing.T) {
		past := time.Now().Add(-1 * time.Hour)
		future := time.Now().Add(1 * time.Hour)
		
		assert.True(t, IsExpired(past), "Past time should be expired")
		assert.False(t, IsExpired(future), "Future time should not be expired")
	})
}

func TestSliceUtils(t *testing.T) {
	t.Run("ContainsString", func(t *testing.T) {
		slice := []string{"apple", "banana", "cherry"}
		
		assert.True(t, ContainsString(slice, "apple"), "Should contain apple")
		assert.True(t, ContainsString(slice, "banana"), "Should contain banana")
		assert.False(t, ContainsString(slice, "orange"), "Should not contain orange")
		assert.False(t, ContainsString(slice, ""), "Should not contain empty string")
	})

	t.Run("RemoveString", func(t *testing.T) {
		slice := []string{"apple", "banana", "cherry", "banana"}
		
		result := RemoveString(slice, "banana")
		expected := []string{"apple", "cherry"}
		
		assert.Equal(t, expected, result, "Should remove all occurrences of banana")
	})

	t.Run("UniqueStrings", func(t *testing.T) {
		slice := []string{"apple", "banana", "apple", "cherry", "banana"}
		
		result := UniqueStrings(slice)
		expected := []string{"apple", "banana", "cherry"}
		
		assert.ElementsMatch(t, expected, result, "Should return unique strings")
	})
}

func TestMapUtils(t *testing.T) {
	t.Run("MergeMaps", func(t *testing.T) {
		map1 := map[string]string{"a": "1", "b": "2"}
		map2 := map[string]string{"b": "3", "c": "4"}
		
		result := MergeMaps(map1, map2)
		expected := map[string]string{"a": "1", "b": "3", "c": "4"}
		
		assert.Equal(t, expected, result, "Should merge maps correctly")
	})

	t.Run("MapKeys", func(t *testing.T) {
		testMap := map[string]int{"a": 1, "b": 2, "c": 3}
		
		keys := MapKeys(testMap)
		expected := []string{"a", "b", "c"}
		
		assert.ElementsMatch(t, expected, keys, "Should return all map keys")
	})

	t.Run("MapValues", func(t *testing.T) {
		testMap := map[string]int{"a": 1, "b": 2, "c": 3}
		
		values := MapValues(testMap)
		expected := []int{1, 2, 3}
		
		assert.ElementsMatch(t, expected, values, "Should return all map values")
	})
}

func TestErrorUtils(t *testing.T) {
	t.Run("WrapError", func(t *testing.T) {
		originalErr := assert.AnError
		wrappedErr := WrapError(originalErr, "context")
		
		assert.Error(t, wrappedErr, "Should return an error")
		assert.Contains(t, wrappedErr.Error(), "context", "Should contain context")
	})

	t.Run("IsTimeoutError", func(t *testing.T) {
		// This would test actual timeout errors
		// For now, we'll just test the function exists
		assert.False(t, IsTimeoutError(assert.AnError), "Should handle non-timeout errors")
	})
}

func TestValidationUtils(t *testing.T) {
	t.Run("IsValidPort", func(t *testing.T) {
		assert.True(t, IsValidPort(80), "Port 80 should be valid")
		assert.True(t, IsValidPort(8080), "Port 8080 should be valid")
		assert.True(t, IsValidPort(65535), "Port 65535 should be valid")
		assert.False(t, IsValidPort(0), "Port 0 should not be valid")
		assert.False(t, IsValidPort(65536), "Port 65536 should not be valid")
		assert.False(t, IsValidPort(-1), "Negative port should not be valid")
	})

	t.Run("IsValidIP", func(t *testing.T) {
		assert.True(t, IsValidIP("127.0.0.1"), "127.0.0.1 should be valid")
		assert.True(t, IsValidIP("192.168.1.1"), "192.168.1.1 should be valid")
		assert.True(t, IsValidIP("::1"), "::1 should be valid")
		assert.False(t, IsValidIP("256.256.256.256"), "Invalid IP should not be valid")
		assert.False(t, IsValidIP("not-an-ip"), "Non-IP string should not be valid")
		assert.False(t, IsValidIP(""), "Empty string should not be valid")
	})

	t.Run("IsValidHostname", func(t *testing.T) {
		assert.True(t, IsValidHostname("localhost"), "localhost should be valid")
		assert.True(t, IsValidHostname("example.com"), "example.com should be valid")
		assert.True(t, IsValidHostname("sub.example.com"), "sub.example.com should be valid")
		assert.False(t, IsValidHostname(""), "Empty hostname should not be valid")
		assert.False(t, IsValidHostname("invalid hostname"), "Hostname with spaces should not be valid")
	})
}

// Benchmark tests
func BenchmarkStringUtils(b *testing.B) {
	b.Run("IsEmpty", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = IsEmpty("test string")
		}
	})

	b.Run("TruncateString", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = TruncateString("this is a long test string", 10)
		}
	})
}

func BenchmarkSliceUtils(b *testing.B) {
	slice := []string{"apple", "banana", "cherry", "date", "elderberry"}
	
	b.Run("ContainsString", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ContainsString(slice, "cherry")
		}
	})

	b.Run("UniqueStrings", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = UniqueStrings(slice)
		}
	})
}

func BenchmarkValidationUtils(b *testing.B) {
	b.Run("IsValidPort", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = IsValidPort(8080)
		}
	})

	b.Run("IsValidIP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = IsValidIP("192.168.1.1")
		}
	})
}

// Test utility functions that don't exist yet (to show what could be added)
func TestFutureUtils(t *testing.T) {
	t.Skip("These tests show potential future utility functions")
	
	/*
	// These functions could be implemented in the future:
	
	t.Run("JSONUtils", func(t *testing.T) {
		// Safe JSON marshaling/unmarshaling
		// JSON validation
		// JSON schema validation
	})
	
	t.Run("CryptoUtils", func(t *testing.T) {
		// Hash functions
		// Encryption/decryption
		// Random string generation
	})
	
	t.Run("FileUtils", func(t *testing.T) {
		// Safe file operations
		// Path validation
		// File type detection
	})
	
	t.Run("NetworkUtils", func(t *testing.T) {
		// URL validation
		// Network connectivity checks
		// Port scanning utilities
	})
	*/
}
