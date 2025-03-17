package main

import (
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
	if err = entity.DbInsertAsync(c); err != nil {
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
			Name:    "Something",
			Surname: "Else",
		}
		if err = entity.DbInsertAsync(c); err != nil {
			t.Fatal(err)
		}
	}
}

func BenchmarkMessageInsert(b *testing.B) {
	c, err := ConnectToNode(GlobalNode)
	if err != nil {
		b.Fatal(err)
	}

	entity := SampleV1.Sample{
		Name:    "Something",
		Surname: "Else",
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err = entity.DbInsertAsync(c); err != nil {
			b.Fatal(err)
		}
	}
}
