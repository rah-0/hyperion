package Sample

// NOTE: this file is generated only once, if you want to update it you can delete it and run the generator again!

import (
	"testing"

	v1 "github.com/rah-0/hyperion/entities/Sample/v1"
)

func TestUpgrade(t *testing.T) {
	previous := v1.Sample{
		Name:    "Name",
		Surname: "Surname",
	}

	current := Upgrade(previous)

	if current.FullName != "Name Surname" {
		t.Fatal()
	}
}

func TestDowngrade(t *testing.T) {
	current := Sample{
		FullName: "Name Surname",
	}

	previous := Downgrade(current)

	if previous.Name != "Name" {
		t.Fatal()
	}
	if previous.Surname != "Surname" {
		t.Fatal()
	}
}
