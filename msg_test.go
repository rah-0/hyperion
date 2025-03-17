package main

import (
	"testing"

	"github.com/google/uuid"

	SampleV1 "github.com/rah-0/hyperion/entities/Sample/v1"
)

func TestMessageInsert(t *testing.T) {
	c, err := ConnectToNode(GlobalNode)
	if err != nil {
		t.Fatal(err)
	}

	entity := SampleV1.Sample{
		Uuid:    uuid.New(),
		Name:    "Something",
		Surname: "Else",
	}

	if err = entity.Encode(); err != nil {
		t.Fatal(err)
	}

	msg := Message{
		Type:   MessageTypeInsert,
		Mode:   ModeAsync,
		Entity: entity.GetBufferData(),
	}

	entity.BufferReset()

	err = c.Send(msg)
	if err != nil {
		t.Fatal(err)
	}
}
