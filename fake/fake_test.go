package fake

import (
	"fmt"
	"regexp"
	"sync"
	"testing"
	"time"
)

func TestGenerateUUID(t *testing.T) {
	var uuidSet sync.Map
	for i := range 1000 {
		t.Run(fmt.Sprintf("UUIDTest-%d", i), func(t *testing.T) {
			t.Parallel()

			uuid, err := RandomUUID()
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			// Test if UUID is of correct length
			if len(uuid) != 36 {
				t.Errorf("Expected length 36, got %d", len(uuid))
			}

			// Test if UUID matches the correct format
			isValidUUID := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`).MatchString
			if !isValidUUID(uuid) {
				t.Errorf("UUID %s does not match the required format", uuid)
			}

			// Test for uniqueness
			if _, exists := uuidSet.Load(uuid); exists {
				t.Errorf("Duplicate UUID found: %s", uuid)
			}
			uuidSet.Store(uuid, struct{}{})
		})
	}
}

func TestRandomDate(t *testing.T) {
	start := time.Date(EpochYear, time.Month(EpochMonth), EpochDay, EpochHour, EpochMinute, EpochSecond, EpochNano, time.UTC)
	end := time.Now()

	randomDate, err := RandomDate()
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	if randomDate.Before(start) || randomDate.After(end) {
		t.Fatalf("Random date %v is outside the expected range [%v, %v]", randomDate, start, end)
	}
}

func TestGenerateRandomPhoneNumber(t *testing.T) {
	phoneNumber, err := RandomPhoneNumber()
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	re := regexp.MustCompile(`^\+1 \(\d{3}\) \d{3}-\d{4}$`)

	if !re.MatchString(phoneNumber) {
		t.Errorf("Generated phone number %v does not match the expected format", phoneNumber)
	}
}

func TestRandomAddress(t *testing.T) {
	address, err := RandomAddress()
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	re := regexp.MustCompile(`^\d+ [A-Za-z ]+, [A-Za-z ]+, [A-Z]{2} \d{5}, USA$`)

	if !re.MatchString(address) {
		t.Errorf("Generated address %v does not match the expected format", address)
	}
}

// ================================================================================
// ### BENCHMARKS
// ================================================================================

func BenchmarkGenerateUUID(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_, _ = RandomUUID()
	}
}

func BenchmarkRandomDate(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_, _ = RandomDate()
	}
}

func BenchmarkRandomPhoneNumber(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_, _ = RandomPhoneNumber()
	}
}

func BenchmarkRandomAddress(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_, _ = RandomAddress()
	}
}
