package log

import (
	"errors"
	"testing"
	"time"
)

type user struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Age       int
	Created   time.Time
}

var (
	err          = errors.New("something happened")
	testDuration = time.Second * 5
	testTime     = time.Date(2025, 9, 25, 12, 30, 45, 0, time.UTC)
	testUser     = user{
		ID:        "u123456",
		FirstName: "John",
		LastName:  "Dow",
		Email:     "john.doe@example.com",
		Age:       30,
		Created:   testTime,
	}
)

type discardWriter struct{}

func (d *discardWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func BenchmarkStructuredLogging(b *testing.B) {
	SetOutput(&discardWriter{})
	SetLevel(InfoLevel)

	b.Run("WithFileInfo", func(b *testing.B) {
		SetIncludeFileInfo(true)
		defer SetIncludeFileInfo(false)
		b.ResetTimer()
		for b.Loop() {
			InfoS(
				WithString("action", "login"),
				WithError("error", err),
				WithAny("user", testUser),
			)
		}
	})

	b.Run("WithoutFileInfo", func(b *testing.B) {
		SetIncludeFileInfo(false)
		b.ResetTimer()
		for b.Loop() {
			InfoS(
				WithString("action", "login"),
				WithError("error", err),
				WithAny("user", testUser),
			)
		}
	})
}

func BenchmarkFormattedLogging(b *testing.B) {
	SetOutput(&discardWriter{})
	SetLevel(InfoLevel)

	b.Run("WithFileInfo", func(b *testing.B) {
		SetIncludeFileInfo(true)
		defer SetIncludeFileInfo(false)
		b.ResetTimer()
		for b.Loop() {
			Infof("user %q failed to login: %v", testUser.ID, err)
		}
	})

	b.Run("WithoutFileInfo", func(b *testing.B) {
		SetIncludeFileInfo(false)
		b.ResetTimer()
		for b.Loop() {
			Infof("user %q failed to login: %v", testUser.ID, err)
		}
	})
}

func BenchmarkFieldTypes(b *testing.B) {
	SetOutput(&discardWriter{})
	SetLevel(InfoLevel)

	b.Run("String", func(b *testing.B) {
		for b.Loop() {
			InfoS(WithString("key", "value"))
		}
	})

	b.Run("Int", func(b *testing.B) {
		for b.Loop() {
			InfoS(WithInt("count", 42))
		}
	})

	b.Run("Float", func(b *testing.B) {
		for b.Loop() {
			InfoS(WithFloat("value", 3.14159))
		}
	})

	b.Run("Bool", func(b *testing.B) {
		for b.Loop() {
			InfoS(WithBool("success", true))
		}
	})

	b.Run("Error", func(b *testing.B) {
		for b.Loop() {
			InfoS(WithError("error", err))
		}
	})

	b.Run("Duration", func(b *testing.B) {
		for b.Loop() {
			InfoS(WithDuration("elapsed", testDuration))
		}
	})
}

func BenchmarkFiltered(b *testing.B) {
	SetLevel(InfoLevel)
	SetOutput(&discardWriter{})
	b.ResetTimer()

	for b.Loop() {
		Debug("debug message that should be filtered")
	}
}

func BenchmarkMixedFields(b *testing.B) {
	SetOutput(&discardWriter{})
	SetLevel(InfoLevel)
	b.ResetTimer()

	for b.Loop() {
		InfoS(
			WithString("user", "john"),
			WithInt("age", 30),
			WithFloat("score", 95.5),
			WithBool("active", true),
			WithError("error", err),
		)
	}
}

func BenchmarkLargePayload(b *testing.B) {
	SetOutput(&discardWriter{})
	SetLevel(InfoLevel)
	b.ResetTimer()

	for b.Loop() {
		InfoS(
			WithString("str1", "value1"),
			WithString("str2", "value2"),
			WithString("str3", "value3"),
			WithInt("int1", 100),
			WithInt("int2", 200),
			WithInt("int3", 300),
			WithFloat("float1", 1.1),
			WithFloat("float2", 2.2),
			WithFloat("float3", 3.3),
			WithBool("bool1", true),
			WithBool("bool2", false),
			WithError("error", err),
		)
	}
}

func BenchmarkConcurrent(b *testing.B) {
	SetOutput(&discardWriter{})
	SetLevel(InfoLevel)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			InfoS(
				WithString("user", "john"),
				WithInt("count", 42),
				WithFloat("score", 95.5),
			)
		}
	})
}

func BenchmarkAllocations(b *testing.B) {
	SetOutput(&discardWriter{})
	SetLevel(InfoLevel)
	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		InfoS(
			WithString("string", "test"),
			WithInt("int", 123),
			WithFloat("float", 1.23),
			WithBool("bool", true),
		)
	}
}
