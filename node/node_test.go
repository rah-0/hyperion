package node

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/rah-0/hyperion/config"
	SampleV1 "github.com/rah-0/hyperion/entities/Sample/v1"
	"github.com/rah-0/hyperion/hconn"
	"github.com/rah-0/hyperion/model"
	"github.com/rah-0/hyperion/query"
	"github.com/rah-0/hyperion/util"
)

const (
	FieldName    = 3
	FieldSurname = 4
)

var (
	connection *hconn.HConn
)

func TestMain(m *testing.M) {
	util.TestMainWrapper(util.TestConfig{
		M: m,
		LoadResources: func() error {
			p, err := filepath.Abs(filepath.Join("..", "config", "config.json"))
			if err != nil {
				return err
			}

			content, err := util.FileRead(p)
			if err != nil {
				return err
			}

			err = json.Unmarshal(content, &config.Loaded)
			if err != nil {
				return err
			}

			var wg sync.WaitGroup
			for _, n := range config.Loaded.Nodes {
				wg.Add(1)

				x := NewNode().
					WithHost(n.Host.Name, n.Host.Port).
					WithPath(n.Path.Data)

				for _, e := range n.Entities {
					x.AddEntity(e.Name)
				}

				go func(node *Node) {
					defer wg.Done()
					_ = util.BuildBinary(node.Host.Name)

					pathConfigForNode := filepath.Join(os.TempDir(), "hyperion_test_"+node.Host.Name+".config")
					_ = util.FileCopy(p, pathConfigForNode)

					logFilePath := filepath.Join(os.TempDir(), "hyperion_test_"+node.Host.Name+".log")
					_ = util.FileDelete(logFilePath)
					logFile, _ := os.Create(logFilePath)
					defer logFile.Close()

					cmd := exec.Command(filepath.Join(os.TempDir(), "hyperion_test_"+node.Host.Name),
						"-pathConfig", pathConfigForNode,
						"-forceHost", node.Host.Name)
					cmd.Stdout = logFile
					cmd.Stderr = logFile

					if err := cmd.Start(); err != nil {
						fmt.Printf("Error running instance for host %s: %v\n", node.Host.Name, err)
					}
				}(x)
			}
			wg.Wait()

			connection, err = ConnectToNodeWithHostAndPort("127.0.0.1", "5000")
			return err
		},
		UnloadResources: func() error {
			if err := connection.Close(); err != nil {
				return err
			}

			err := util.Pkill("hyperion_test_*")
			for _, n := range config.Loaded.Nodes {
				util.DirectoryRemove(n.Path.Data)
			}
			return err
		},
	})
}

func TestNodeDirectConnectionIpAndPort(t *testing.T) {
	msg := model.Message{
		Type:   model.MessageTypeTest,
		String: "Test",
	}

	err := connection.Send(msg)
	if err != nil {
		t.Fatal(err)
	}
	msg, err = connection.Receive()
	if err != nil {
		t.Fatal(err)
	}

	if msg.String != "TestReceived" {
		t.Fatalf("Unexpected message: %s", msg.String)
	}
}

func TestNodesDirectConnection(t *testing.T) {
	var totalExpectedMessages = len(config.Loaded.Nodes)
	var totalSuccessfulMessages int

	for _, nodeConfig := range config.Loaded.Nodes {
		n := NewNode().
			WithHost(nodeConfig.Host.Name, nodeConfig.Host.Port).
			WithPath(nodeConfig.Path.Data)
		for _, e := range nodeConfig.Entities {
			n.AddEntity(e.Name)
		}

		c, err := ConnectToNode(n)
		if err != nil {
			t.Errorf("Failed to connect to node [%s:%d]: %v", n.Host.Name, n.Host.Port, err)
			continue
		}

		msg := model.Message{
			Type:   model.MessageTypeTest,
			String: "Test",
		}

		err = c.Send(msg)
		if err != nil {
			t.Errorf("Failed to send message to node [%s:%d]: %v", n.Host.Name, n.Host.Port, err)
			continue
		}

		msg, err = c.Receive()
		if err != nil {
			t.Errorf("Failed to receive message from node [%s:%d]: %v", n.Host.Name, n.Host.Port, err)
			continue
		}

		if msg.String != "TestReceived" {
			t.Errorf("Unexpected response from node [%s:%d]: got %q, want %q",
				n.Host.Name, n.Host.Port, msg.String, "TestReceived")
			continue
		}

		totalSuccessfulMessages++
	}

	if totalSuccessfulMessages != totalExpectedMessages {
		t.Fatalf("Mismatch in total successful messages: got %d, expected %d",
			totalSuccessfulMessages, totalExpectedMessages)
	}
}

func TestMessageInsert(t *testing.T) {
	entity := SampleV1.Sample{
		Name:    "Something",
		Surname: "Else",
	}

	err := entity.DbInsert(connection)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMessageInsertAndDelete(t *testing.T) {
	entity := SampleV1.Sample{
		Name:    "Something",
		Surname: "Else",
	}

	err := entity.DbInsert(connection)
	if err != nil {
		t.Fatal(err)
	}

	err = entity.DbDelete(connection)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMessageInsert1000(t *testing.T) {
	var expected []SampleV1.Sample
	for i := 0; i < 1000; i++ {
		entity := SampleV1.Sample{
			Name:    fmt.Sprintf("Something%d", i),
			Surname: fmt.Sprintf("Else%d", i),
		}

		err := entity.DbInsert(connection)
		if err != nil {
			t.Fatal(err)
		}

		expected = append(expected, entity)
	}

	all, err := SampleV1.DbGetAll(connection)
	if err != nil {
		t.Fatal(err)
	}

	for _, expectedEntity := range expected {
		found := false
		for _, actual := range all {
			if expectedEntity.GetUuid() == actual.GetUuid() {
				found = true
				if expectedEntity.GetFieldValue(FieldName) != actual.GetFieldValue(FieldName) {
					t.Fatalf("Name mismatch for UUID %v", actual.GetUuid())
				}
				if expectedEntity.GetFieldValue(FieldSurname) != actual.GetFieldValue(FieldSurname) {
					t.Fatalf("Surname mismatch for UUID %v", actual.GetUuid())
				}
				break
			}
		}
		if !found {
			t.Fatalf("Expected entity UUID %v not found in received list", expectedEntity.GetUuid())
		}
	}
}

func TestMessageUpdate(t *testing.T) {
	entity := &SampleV1.Sample{
		Name:    "Initial",
		Surname: "User",
	}
	if err := entity.DbInsert(connection); err != nil {
		t.Fatal(err)
	}

	// Modify values
	entity.Name = "Updated"
	entity.Surname = "User"

	err := entity.DbUpdate(connection)
	if err != nil {
		t.Fatal(err)
	}

	all, err := SampleV1.DbGetAll(connection)
	if err != nil {
		t.Fatal(err)
	}

	var match *SampleV1.Sample
	for _, e := range all {
		if e.GetUuid() == entity.GetUuid() {
			match = e
			break
		}
	}

	if match == nil {
		t.Fatalf("Updated entity not found")
	}
	if match.Name != "Updated" || match.Surname != "User" {
		t.Fatalf("Update not applied correctly. Got Name: %s, Surname: %s", match.Name, match.Surname)
	}
}

func TestMessageGetAll(t *testing.T) {
	var inserted []*SampleV1.Sample

	// Insert 3 entities and track them
	for i := 0; i < 3; i++ {
		entity := &SampleV1.Sample{
			Name:    fmt.Sprintf("Name%d", i),
			Surname: fmt.Sprintf("Surname%d", i),
		}
		if err := entity.DbInsert(connection); err != nil {
			t.Fatal(err)
		}
		inserted = append(inserted, entity)
	}

	// Get all from remote
	entities, err := SampleV1.DbGetAll(connection)
	if err != nil {
		t.Fatal(err)
	}

	if len(entities) < len(inserted) {
		t.Fatalf("Expected at least %d entities, got %d", len(inserted), len(entities))
	}

	// Check UUIDs exist in received list
	for _, ins := range inserted {
		found := false
		for _, got := range entities {
			if ins.Uuid == got.Uuid {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Expected UUID %v not found in response", ins.Uuid)
		}
	}
}

func TestQueryStringFilter(t *testing.T) {
	entities := []*SampleV1.Sample{
		{Name: "Alice", Surname: "Smith"},
		{Name: "Bob", Surname: "Jones"},
		{Name: "Charlie", Surname: "Smith"},
		{Name: "Diana", Surname: "Brown"},
	}

	for _, e := range entities {
		if err := e.DbInsert(connection); err != nil {
			t.Fatalf("insert failed: %v", err)
		}
	}

	all, err := SampleV1.DbGetAll(connection)
	if err != nil {
		t.Fatal(err)
	}

	for _, inserted := range entities {
		found := false
		for _, fetched := range all {
			if inserted.GetUuid() == fetched.GetUuid() {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("entity with UUID %v not found in DbGetAll results", inserted.GetUuid())
		}
	}

	q := SampleV1.NewQuery().SetFilters(query.Filters{
		Type: query.FilterTypeOr,
		Filters: []query.Filter{
			{Field: SampleV1.FieldSurname, Op: query.OpTypeEqual, Value: "Smith"},
		},
	})

	results, err := q.Execute(connection)
	if err != nil {
		t.Fatal(err)
	}

	for _, r := range results {
		if r.Surname != "Smith" {
			t.Fatalf("unexpected result: %+v", r)
		}
	}

	expectedCount := 0
	for _, e := range entities {
		if e.Surname == "Smith" {
			expectedCount++
		}
	}
	if len(results) != expectedCount {
		t.Fatalf("expected %d results for surname=Smith, got %d", expectedCount, len(results))
	}
}

func TestQueryStringFilterAnd(t *testing.T) {
	entities := []*SampleV1.Sample{
		{Name: "Alice", Surname: "Smith"},
		{Name: "Alice", Surname: "Brown"},
		{Name: "Bob", Surname: "Smith"},
		{Name: "Diana", Surname: "Jones"},
	}

	for _, e := range entities {
		if err := e.DbInsert(connection); err != nil {
			t.Fatalf("insert failed: %v", err)
		}
	}

	all, err := SampleV1.DbGetAll(connection)
	if err != nil {
		t.Fatal(err)
	}

	for _, inserted := range entities {
		found := false
		for _, fetched := range all {
			if inserted.GetUuid() == fetched.GetUuid() {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("entity with UUID %v not found in DbGetAll results", inserted.GetUuid())
		}
	}

	q := SampleV1.NewQuery().SetFilters(query.Filters{
		Type: query.FilterTypeAnd,
		Filters: []query.Filter{
			{Field: SampleV1.FieldName, Op: query.OpTypeEqual, Value: "Alice"},
			{Field: SampleV1.FieldSurname, Op: query.OpTypeEqual, Value: "Smith"},
		},
	})

	results, err := q.Execute(connection)
	if err != nil {
		t.Fatal(err)
	}

	for _, r := range results {
		if r.Name != "Alice" || r.Surname != "Smith" {
			t.Fatalf("unexpected result: %+v", r)
		}
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result for Name='Alice' AND Surname='Smith', got %d", len(results))
	}
}

var messageSizes = map[string]int{
	"2KB":   util.Size2KB,
	"4KB":   util.Size4KB,
	"8KB":   util.Size8KB,
	"16KB":  util.Size16KB,
	"32KB":  util.Size32KB,
	"64KB":  util.Size64KB,
	"128KB": util.Size128KB,
	"256KB": util.Size256KB,
	"512KB": util.Size512KB,
	"1MB":   util.Size1MB,
	"10MB":  util.Size10MB,
	"100MB": util.Size100MB,
}

func BenchmarkListenPortForStatus(b *testing.B) {
	c, err := ConnectToNodeWithHostAndPort("127.0.0.1", "5000")
	if err != nil {
		b.Fatal(err)
	}

	for sizeLabel, size := range messageSizes {
		b.Run("Size: "+sizeLabel, func(b *testing.B) {
			// Generate a random message of the given size
			testStr, err := util.GenerateRandomStringMessage(size)
			if err != nil {
				b.Fatalf("Failed to generate message of size %d: %v", size, err)
			}

			msg := model.Message{
				Type:   model.MessageTypeTest,
				String: testStr,
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Send message
				if err := c.Send(msg); err != nil {
					b.Fatal(err)
				}

				receivedMsg, err := c.Receive()
				if err != nil {
					b.Fatal(err)
				}

				// Validate response
				if receivedMsg.String != testStr+"Received" {
					b.Fatalf("Unexpected response: got %q, want %q", receivedMsg.String, testStr+"Received")
				}
			}
		})
	}
}

func BenchmarkConnectionEstablishment(b *testing.B) {
	// Ensure the node is reachable once before benchmarking
	for {
		warmupConn, err := ConnectToNodeWithHostAndPort("127.0.0.1", "5000")
		if err == nil {
			warmupConn.Close()
			break
		}
		time.Sleep(1 * time.Second)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		conn, err := ConnectToNodeWithHostAndPort("127.0.0.1", "5000")
		if err != nil {
			b.Fatalf("Connection failed on iteration %d: %v", i, err)
		}
		conn.Close()
	}
}

func BenchmarkMessageInsert(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		entity := SampleV1.Sample{
			Name:    fmt.Sprintf("Something%d", i),
			Surname: fmt.Sprintf("Else%d", i),
		}

		err := entity.DbInsert(connection)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDbGetAll(b *testing.B) {
	counts := []int{1, 10, 100, 1000}

	for _, count := range counts {
		b.Run(fmt.Sprintf("Count%d", count), func(b *testing.B) {
			// Insert 'count' entities before benchmarking
			for i := 0; i < count; i++ {
				entity := SampleV1.Sample{
					Name:    fmt.Sprintf("Name%d", i),
					Surname: fmt.Sprintf("Surname%d", i),
				}
				err := entity.DbInsert(connection)
				if err != nil {
					b.Fatalf("Insert failed for count %d at i=%d: %v", count, i, err)
				}
			}

			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				entities, err := SampleV1.DbGetAll(connection)
				if err != nil {
					b.Fatal(err)
				}
				if len(entities) < count {
					b.Fatalf("Expected at least %d entities, got %d", count, len(entities))
				}
			}
		})
	}
}
