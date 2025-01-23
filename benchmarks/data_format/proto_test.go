package data_format

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rah-0/hyperion/benchmarks/data_format/pb"
	"github.com/rah-0/hyperion/utils/testutil"
)

func TestEncodeDecodePerson(t *testing.T) {
	persons := generateRandomPersons(1)
	original := persons[0]

	data, err := EncodeProto(&pb.Person{
		Name:    original.Name,
		Age:     int32(original.Age),
		Surname: original.Surname,
	})
	assert.NoError(t, err, "Serialization of Person should not fail")

	decoded := &pb.Person{}
	err = DecodeProto(data, decoded)
	assert.NoError(t, err, "Deserialization of Person should not fail")

	assert.Equal(t, original.Name, decoded.Name, "Names should match")
	assert.Equal(t, original.Age, int(decoded.Age), "Ages should match")
	assert.Equal(t, original.Surname, decoded.Surname, "Surnames should match")
}

func BenchmarkEncodeProto1Small(b *testing.B) {
	runEncodeProtoBenchmark(b, 1)
}

func BenchmarkEncodeProto100Small(b *testing.B) {
	runEncodeProtoBenchmark(b, 100)
}

func BenchmarkEncodeProto10000Small(b *testing.B) {
	runEncodeProtoBenchmark(b, 10000)
}

func BenchmarkEncodeProto1000000Small(b *testing.B) {
	runEncodeProtoBenchmark(b, 1000000)
}

func runEncodeProtoBenchmark(b *testing.B, count int) {
	defer testutil.RecoverBenchHandler(b)

	data := generateRandomPersons(count)
	protoData := make([]*pb.Person, len(data))
	for i, p := range data {
		protoData[i] = &pb.Person{
			Name:    p.Name,
			Age:     int32(p.Age),
			Surname: p.Surname,
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, person := range protoData {
			_, err := EncodeProto(person)
			if err != nil {
				b.Fatalf("Failed to Encode Protobuf: %v", err)
			}
		}
	}
}

func BenchmarkEncodeProto1Unreal(b *testing.B) {
	runEncodeProtoUnrealBenchmark(b, 1)
}

func BenchmarkEncodeProto10Unreals(b *testing.B) {
	runEncodeProtoUnrealBenchmark(b, 10)
}

func BenchmarkEncodeProto100Unreals(b *testing.B) {
	runEncodeProtoUnrealBenchmark(b, 100)
}

func BenchmarkEncodeProto1000Unreals(b *testing.B) {
	runEncodeProtoUnrealBenchmark(b, 1000)
}

func runEncodeProtoUnrealBenchmark(b *testing.B, count int) {
	data := generateRandomUnreals(count)
	protoData := make([]*pb.Unreal, len(data))
	for i, u := range data {
		protoData[i] = &pb.Unreal{
			Prop1:   u.Prop1,
			Prop2:   u.Prop2,
			Prop3:   u.Prop3,
			Prop4:   u.Prop4,
			Prop5:   u.Prop5,
			Prop6:   u.Prop6,
			Prop7:   u.Prop7,
			Prop8:   u.Prop8,
			Prop9:   u.Prop9,
			Prop10:  u.Prop10,
			Prop11:  u.Prop11,
			Prop12:  u.Prop12,
			Prop13:  u.Prop13,
			Prop14:  u.Prop14,
			Prop15:  u.Prop15,
			Prop16:  u.Prop16,
			Prop17:  u.Prop17,
			Prop18:  u.Prop18,
			Prop19:  u.Prop19,
			Prop20:  u.Prop20,
			Prop21:  u.Prop21,
			Prop22:  u.Prop22,
			Prop23:  u.Prop23,
			Prop24:  u.Prop24,
			Prop25:  u.Prop25,
			Prop26:  u.Prop26,
			Prop27:  u.Prop27,
			Prop28:  u.Prop28,
			Prop29:  u.Prop29,
			Prop30:  u.Prop30,
			Prop31:  u.Prop31,
			Prop32:  u.Prop32,
			Prop33:  u.Prop33,
			Prop34:  u.Prop34,
			Prop35:  u.Prop35,
			Prop36:  u.Prop36,
			Prop37:  u.Prop37,
			Prop38:  u.Prop38,
			Prop39:  u.Prop39,
			Prop40:  u.Prop40,
			Prop41:  u.Prop41,
			Prop42:  u.Prop42,
			Prop43:  u.Prop43,
			Prop44:  u.Prop44,
			Prop45:  u.Prop45,
			Prop46:  u.Prop46,
			Prop47:  u.Prop47,
			Prop48:  u.Prop48,
			Prop49:  u.Prop49,
			Prop50:  u.Prop50,
			Prop51:  u.Prop51,
			Prop52:  u.Prop52,
			Prop53:  u.Prop53,
			Prop54:  u.Prop54,
			Prop55:  u.Prop55,
			Prop56:  u.Prop56,
			Prop57:  u.Prop57,
			Prop58:  u.Prop58,
			Prop59:  u.Prop59,
			Prop60:  u.Prop60,
			Prop61:  u.Prop61,
			Prop62:  u.Prop62,
			Prop63:  u.Prop63,
			Prop64:  u.Prop64,
			Prop65:  u.Prop65,
			Prop66:  u.Prop66,
			Prop67:  u.Prop67,
			Prop68:  u.Prop68,
			Prop69:  u.Prop69,
			Prop70:  u.Prop70,
			Prop71:  u.Prop71,
			Prop72:  u.Prop72,
			Prop73:  u.Prop73,
			Prop74:  u.Prop74,
			Prop75:  u.Prop75,
			Prop76:  u.Prop76,
			Prop77:  u.Prop77,
			Prop78:  u.Prop78,
			Prop79:  u.Prop79,
			Prop80:  u.Prop80,
			Prop81:  u.Prop81,
			Prop82:  u.Prop82,
			Prop83:  u.Prop83,
			Prop84:  u.Prop84,
			Prop85:  u.Prop85,
			Prop86:  u.Prop86,
			Prop87:  u.Prop87,
			Prop88:  u.Prop88,
			Prop89:  u.Prop89,
			Prop90:  u.Prop90,
			Prop91:  u.Prop91,
			Prop92:  u.Prop92,
			Prop93:  u.Prop93,
			Prop94:  u.Prop94,
			Prop95:  u.Prop95,
			Prop96:  u.Prop96,
			Prop97:  u.Prop97,
			Prop98:  u.Prop98,
			Prop99:  u.Prop99,
			Prop100: u.Prop100,
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, unreal := range protoData {
			_, err := EncodeProto(unreal)
			if err != nil {
				b.Fatalf("Failed to Encode Protobuf: %v", err)
			}
		}
	}
}
