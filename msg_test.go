package main

import (
	"fmt"
	"testing"

	SampleV1 "github.com/rah-0/hyperion/entities/Sample/v1"
)

func TestMessageInsert(t *testing.T) {
	c, err := ConnectToNode(GlobalNode)
	if err != nil {
		t.Fatal(err)
	}

	entity := SampleV1.Sample{
		Name:    "Something",
		Surname: "Else",
	}

	_, err = entity.DbInsert(c)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMessageInsert1000(t *testing.T) {
	c, err := ConnectToNode(GlobalNode)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		entity := SampleV1.Sample{
			Name:    fmt.Sprintf("Something%d", i),
			Surname: fmt.Sprintf("Else%d", i),
		}
		if _, err := entity.DbInsert(c); err != nil {
			t.Fatal(err)
		}
	}
}

func BenchmarkMessageInsert(b *testing.B) {
	c, err := ConnectToNode(GlobalNode)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		entity := SampleV1.Sample{
			Name:    fmt.Sprintf("Something%d", i),
			Surname: fmt.Sprintf("Else%d", i),
		}
		if _, err := entity.DbInsert(c); err != nil {
			b.Fatal(err)
		}
	}
}
