package Sample

// NOTE: this file is generated only once, if you want to update it you can delete it and run the generator again!

import (
	"strings"

	v1 "github.com/rah-0/hyperion/entities/Sample/v1"
)

func Upgrade(previous v1.Sample) (current Sample) {
	current.FullName = previous.Name + " " + previous.Surname
	return
}

func Downgrade(current Sample) (previous v1.Sample) {
	parts := strings.Split(current.FullName, " ")
	previous.Name = parts[0]
	previous.Surname = parts[1]
	return
}
