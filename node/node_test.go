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

	"github.com/google/uuid"
	"github.com/rah-0/testmark/testutil"

	"github.com/rah-0/hyperion/config"
	"github.com/rah-0/hyperion/disk"
	SampleV1 "github.com/rah-0/hyperion/entities/Sample/v1"
	"github.com/rah-0/hyperion/hconn"
	"github.com/rah-0/hyperion/model"
	"github.com/rah-0/hyperion/query"
	"github.com/rah-0/hyperion/register"
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
	testutil.TestMainWrapper(testutil.TestConfig{
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
					WithHost(n.Host.Name, n.Host.IP, n.Host.Port).
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

// TestNodeShutdown verifies that the Shutdown method properly cleans up resources
func TestNodeShutdown(t *testing.T) {
	// Create a new node
	n := NewNode()
	n.WithHost("test-shutdown", "127.0.0.1", util.GetAvailablePort())
	tempDir := t.TempDir()
	n.WithPath(tempDir)

	// Add entity storage with disk
	e := &register.Entity{
		EntityBase: &register.EntityBase{
			Name:    "Sample",
			Version: "v1",
		},
	}

	// Use the entity setup pattern
	d := disk.NewDisk()
	d.WithNewRandomPath()
	d.WithEntity(e)
	if err := d.OpenFile(); err != nil {
		t.Fatalf("Failed to open disk: %v", err)
	}

	es := &EntityStorage{
		Disk:   d,
		Memory: e,
	}
	n.EntitiesStorage = append(n.EntitiesStorage, es)

	// Start the node and create a connection
	errChan := make(chan error, 1)
	n.ErrCh = errChan

	// Node must be started for connection to work
	go func() {
		if err := n.Start(); err != nil {
			errChan <- err
		}
	}()

	// Give node time to start
	time.Sleep(100 * time.Millisecond)

	// Execute shutdown
	err := n.Shutdown()
	if err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	// Verify that resources are properly cleaned up:

	// Check if file handle is closed by trying to use it
	_, err = d.File.Stat()
	if err == nil {
		t.Error("File handle should be closed after Shutdown")
	}

	// Check if disk resources are properly cleaned
	if d.SyncTicker != nil {
		t.Error("SyncTicker should be nil after Shutdown")
	}

	if d.StopChan != nil {
		t.Error("StopChan should be nil after Shutdown")
	}

	// Check if error channel was properly handled
	if n.ErrCh != nil {
		t.Error("Error channel should be nil after Shutdown")
	}
}

// TestNodeShutdownMultiple ensures calling Shutdown multiple times is safe
// TestNodeMessageDuringShutdown verifies that clients receive the proper shutdown error
// when sending messages to a node that is in the process of shutting down
func TestNodeMessageDuringShutdown(t *testing.T) {
	// Create a node that will be shut down
	shutdownNode := NewNode()
	shutdownPort := util.GetAvailablePort()
	shutdownNode.WithHost("shutdown-test-node", "127.0.0.1", shutdownPort)
	shutdownNode.WithPath(t.TempDir())

	// Add entity storage with disk
	e := &register.Entity{
		EntityBase: &register.EntityBase{
			Name:    "Sample",
			Version: "v1",
		},
	}

	// Use the entity setup pattern
	d := disk.NewDisk()
	d.WithNewRandomPath()
	d.WithEntity(e)
	if err := d.OpenFile(); err != nil {
		t.Fatalf("Failed to open disk: %v", err)
	}

	es := &EntityStorage{
		Disk:   d,
		Memory: e,
	}
	shutdownNode.EntitiesStorage = append(shutdownNode.EntitiesStorage, es)

	// Create error channel
	shutdownErrChan := make(chan error, 10)
	shutdownNode.ErrCh = shutdownErrChan

	// Start the node that will be shutdown
	go func() {
		if err := shutdownNode.Start(); err != nil {
			t.Errorf("Failed to start shutdownNode: %v", err)
		}
	}()

	// Give node time to start
	time.Sleep(100 * time.Millisecond)

	// Create a direct connection to the node
	conn, err := ConnectToNodeWithHostAndPort("127.0.0.1", fmt.Sprintf("%d", shutdownPort))
	if err != nil {
		t.Fatalf("Failed to connect to node: %v", err)
	}

	// Start shutting down the node in a goroutine
	go func() {
		// Small sleep to ensure the connection is established before shutdown starts
		time.Sleep(10 * time.Millisecond)
		if err := shutdownNode.Shutdown(); err != nil {
			t.Errorf("Shutdown failed: %v", err)
		}
	}()

	// Give it a moment to register the shutdown status
	time.Sleep(20 * time.Millisecond)

	// Try to send a message during shutdown and verify we get the correct error
	msg := model.Message{
		Type:   model.MessageTypeTest,
		String: "test message",
	}

	err = conn.Send(msg)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Receive the response
	response, err := conn.Receive()
	if err != nil {
		t.Fatalf("Failed to receive message: %v", err)
	}

	// Check that we received the shutdown error
	if response.Status != model.StatusError {
		t.Errorf("Expected status error, got %d", response.Status)
	}

	// The error message should be in the String field when Status is StatusError
	expectedErrMsg := model.ErrNodeShutdown.Error()
	if response.String != expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, response.String)
	}

	// Clean up
	conn.Close()
}

func TestNodeShutdownMultiple(t *testing.T) {
	n := NewNode()
	n.WithHost("test-multi-shutdown", "127.0.0.1", util.GetAvailablePort())
	tempDir := t.TempDir()
	n.WithPath(tempDir)

	// Start briefly to initialize resources
	errChan := make(chan error, 1)
	n.ErrCh = errChan

	go func() {
		if err := n.Start(); err != nil {
			errChan <- err
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// First shutdown
	err := n.Shutdown()
	if err != nil {
		t.Fatalf("First shutdown failed: %v", err)
	}

	// Second shutdown - should not error
	err = n.Shutdown()
	if err != nil {
		t.Fatalf("Second shutdown should be safe but got error: %v", err)
	}
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
			WithHost(nodeConfig.Host.Name, nodeConfig.Host.IP, nodeConfig.Host.Port).
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

func TestMessageInsertAndDeleteAll(t *testing.T) {
	// Insert multiple entities
	entityCount := 5
	for i := 0; i < entityCount; i++ {
		entity := SampleV1.Sample{
			Name:    fmt.Sprintf("Test%d", i),
			Surname: fmt.Sprintf("User%d", i),
		}

		err := entity.DbInsert(connection)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Verify entities were inserted
	entities, err := SampleV1.DbGetAll(connection)
	if err != nil {
		t.Fatal(err)
	}

	// We might have existing entities from other tests, so just check we have at least our new ones
	if len(entities) < entityCount {
		t.Fatalf("Expected at least %d entities, but got %d", entityCount, len(entities))
	}

	// Delete all entities
	err = SampleV1.DbDeleteAll(connection)
	if err != nil {
		t.Fatal(err)
	}

	// Verify all entities are deleted
	entities, err = SampleV1.DbGetAll(connection)
	if err != nil {
		t.Fatal(err)
	}

	// Check
	if len(entities) != 0 {
		t.Errorf("Entities were not deleted")
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

	q := query.NewQuery().SetFilters(query.FilterTypeOr, []query.Filter{
		{Field: SampleV1.FieldSurname, Op: query.OperatorTypeEqual, Value: "Smith"},
	})

	results, err := SampleV1.DbQuery(connection, q)
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
		{Name: "Alice1", Surname: "Smith1"},
		{Name: "Alice1", Surname: "Brown"},
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

	q := query.NewQuery().SetFilters(query.FilterTypeAnd, []query.Filter{
		{Field: SampleV1.FieldName, Op: query.OperatorTypeEqual, Value: "Alice1"},
		{Field: SampleV1.FieldSurname, Op: query.OperatorTypeEqual, Value: "Smith1"},
	})

	results, err := SampleV1.DbQuery(connection, q)
	if err != nil {
		t.Fatal(err)
	}

	for _, r := range results {
		if r.Name != "Alice1" || r.Surname != "Smith1" {
			t.Fatalf("unexpected result: %+v", r)
		}
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result for Name='Alice' AND Surname='Smith', got %d", len(results))
	}
}

func TestQueryOrderAscSingleField(t *testing.T) {
	entities := []*SampleV1.Sample{
		{Name: "Asc1_C"},
		{Name: "Asc1_A"},
		{Name: "Asc1_B"},
	}
	insertAndQueryCheck(t, entities, []query.Order{
		{Type: query.OrderTypeAsc, Field: SampleV1.FieldName},
	}, []string{"Asc1_A", "Asc1_B", "Asc1_C"})
}

func TestQueryOrderDescSingleField(t *testing.T) {
	entities := []*SampleV1.Sample{
		{Name: "Desc1_B"},
		{Name: "Desc1_C"},
		{Name: "Desc1_A"},
	}
	insertAndQueryCheck(t, entities, []query.Order{
		{Type: query.OrderTypeDesc, Field: SampleV1.FieldName},
	}, []string{"Desc1_C", "Desc1_B", "Desc1_A"})
}

func TestQueryOrderMultiFieldAsc(t *testing.T) {
	entities := []*SampleV1.Sample{
		{Name: "Asc2_A", Surname: "C"},
		{Name: "Asc2_A", Surname: "A"},
		{Name: "Asc2_B", Surname: "B"},
	}
	insertAndQueryCheck(t, entities, []query.Order{
		{Type: query.OrderTypeAsc, Field: SampleV1.FieldName},
		{Type: query.OrderTypeAsc, Field: SampleV1.FieldSurname},
	}, []string{"Asc2_A:A", "Asc2_A:C", "Asc2_B:B"})
}

func TestQueryOrderMultiFieldDesc(t *testing.T) {
	entities := []*SampleV1.Sample{
		{Name: "Desc2_A", Surname: "B"},
		{Name: "Desc2_A", Surname: "C"},
		{Name: "Desc2_B", Surname: "A"},
	}
	insertAndQueryCheck(t, entities, []query.Order{
		{Type: query.OrderTypeDesc, Field: SampleV1.FieldName},
		{Type: query.OrderTypeDesc, Field: SampleV1.FieldSurname},
	}, []string{"Desc2_B:A", "Desc2_A:C", "Desc2_A:B"})
}

func TestQueryOrderMultiFieldMixed(t *testing.T) {
	entities := []*SampleV1.Sample{
		{Name: "Mix1_A", Surname: "B"},
		{Name: "Mix1_A", Surname: "C"},
		{Name: "Mix1_B", Surname: "A"},
	}
	insertAndQueryCheck(t, entities, []query.Order{
		{Type: query.OrderTypeAsc, Field: SampleV1.FieldName},
		{Type: query.OrderTypeDesc, Field: SampleV1.FieldSurname},
	}, []string{"Mix1_A:C", "Mix1_A:B", "Mix1_B:A"})
}

func insertAndQueryCheck(t *testing.T, entities []*SampleV1.Sample, orders []query.Order, expectedOrdered []string) {
	t.Helper()

	// Insert each entity
	for _, e := range entities {
		if err := e.DbInsert(connection); err != nil {
			t.Fatalf("insert failed: %v", err)
		}
	}

	// Run the ordered query
	q := query.NewQuery().SetOrders(orders)
	results, err := SampleV1.DbQuery(connection, q)
	if err != nil {
		t.Fatal(err)
	}

	// Extract actual ordered values
	var actual []string
	for _, r := range results {
		n := r.GetFieldValue(SampleV1.FieldName).(string)
		s := r.GetFieldValue(SampleV1.FieldSurname).(string)
		if s != "" {
			actual = append(actual, n+":"+s)
		} else {
			actual = append(actual, n)
		}
	}

	// Find expected values in correct order inside actual list
	idx := 0
	for _, val := range actual {
		if idx < len(expectedOrdered) && val == expectedOrdered[idx] {
			idx++
		}
	}
	if idx != len(expectedOrdered) {
		t.Fatalf("expected sequence %v not matched in result %v", expectedOrdered, actual)
	}
}

func TestQueryFilterMultipleValuesOr(t *testing.T) {
	entities := []*SampleV1.Sample{
		{Name: "FilterOr_Unique_X1", Surname: "A"},
		{Name: "FilterOr_Unique_X2", Surname: "B"},
		{Name: "FilterOr_Unique_X3", Surname: "C"},
	}
	var target []*SampleV1.Sample
	for _, e := range entities {
		if err := e.DbInsert(connection); err != nil {
			t.Fatal(err)
		}
		if e.Surname == "A" || e.Surname == "C" {
			target = append(target, e)
		}
	}

	q := query.NewQuery().SetFilters(query.FilterTypeOr, []query.Filter{
		{Field: SampleV1.FieldSurname, Op: query.OperatorTypeEqual, Value: "A"},
		{Field: SampleV1.FieldSurname, Op: query.OperatorTypeEqual, Value: "C"},
	})

	results, err := SampleV1.DbQuery(connection, q)
	if err != nil {
		t.Fatal(err)
	}

	for _, expected := range target {
		found := false
		for _, r := range results {
			if r.GetUuid() == expected.GetUuid() {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("UUID %v with surname %s not found", expected.GetUuid(), expected.Surname)
		}
	}
}

func TestQueryFilterContains(t *testing.T) {
	e := &SampleV1.Sample{Name: "FilterContains_Unique", Surname: "XYZ_ABC_123"}
	if err := e.DbInsert(connection); err != nil {
		t.Fatal(err)
	}

	q := query.NewQuery().SetFilters(query.FilterTypeAnd, []query.Filter{
		{Field: SampleV1.FieldSurname, Op: query.OperatorTypeContains, Value: "ABC"},
	})

	results, err := SampleV1.DbQuery(connection, q)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, r := range results {
		if r.GetUuid() == e.GetUuid() {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Expected UUID %v not found in results", e.GetUuid())
	}
}

func TestQueryLimitResults(t *testing.T) {
	var inserted []*SampleV1.Sample
	for i := 0; i < 5; i++ {
		e := &SampleV1.Sample{Name: fmt.Sprintf("LimitTest_Unique_%d", i)}
		if err := e.DbInsert(connection); err != nil {
			t.Fatal(err)
		}
		inserted = append(inserted, e)
	}

	q := query.NewQuery().
		SetFilters(query.FilterTypeOr, []query.Filter{
			{Field: SampleV1.FieldName, Op: query.OperatorTypeContains, Value: "LimitTest_Unique_"},
		}).
		SetLimit(2)

	results, err := SampleV1.DbQuery(connection, q)
	if err != nil {
		t.Fatal(err)
	}

	found := 0
	for _, expected := range inserted {
		for _, r := range results {
			if r.GetUuid() == expected.GetUuid() {
				found++
			}
		}
	}
	if found == 0 {
		t.Fatalf("None of the inserted UUIDs found in limited query result")
	}
}

func TestQueryOrderBySurnameDesc(t *testing.T) {
	entities := []*SampleV1.Sample{
		{Name: "OrderDesc_Unique", Surname: "C"},
		{Name: "OrderDesc_Unique", Surname: "A"},
		{Name: "OrderDesc_Unique", Surname: "B"},
	}
	insertAndQueryCheck(t, entities, []query.Order{
		{Type: query.OrderTypeDesc, Field: SampleV1.FieldSurname},
	}, []string{"OrderDesc_Unique:C", "OrderDesc_Unique:B", "OrderDesc_Unique:A"})
}

func TestQueryFilterOrderLimitCombo(t *testing.T) {
	entities := []*SampleV1.Sample{
		{Name: "Combo_Unique", Surname: "Z"},
		{Name: "Combo_Unique", Surname: "X"},
		{Name: "Combo_Unique", Surname: "Y"},
	}
	var uuids []uuid.UUID
	for _, e := range entities {
		if err := e.DbInsert(connection); err != nil {
			t.Fatal(err)
		}
		uuids = append(uuids, e.GetUuid())
	}

	q := query.NewQuery().
		SetFilters(query.FilterTypeAnd, []query.Filter{
			{Field: SampleV1.FieldName, Op: query.OperatorTypeEqual, Value: "Combo_Unique"},
		}).
		SetOrders([]query.Order{
			{Type: query.OrderTypeAsc, Field: SampleV1.FieldSurname},
		}).
		SetLimit(2)

	results, err := SampleV1.DbQuery(connection, q)
	if err != nil {
		t.Fatal(err)
	}

	found := 0
	for _, r := range results {
		for _, u := range uuids {
			if r.GetUuid() == u {
				found++
			}
		}
	}
	if found == 0 {
		t.Fatalf("Expected at least 1 match from filtered + ordered + limited results")
	}
}

func TestQueryBirthLastYear(t *testing.T) {
	now := time.Now()
	lastYear := now.AddDate(-1, 0, 0)

	// Create 10 recent entries within the last year
	var recent []*SampleV1.Sample
	for i := 0; i < 10; i++ {
		e := &SampleV1.Sample{
			Name:    fmt.Sprintf("Recent%d", i),
			Surname: "User",
			Birth:   now.AddDate(0, 0, -i*10), // Spread out over the last ~100 days
		}
		recent = append(recent, e)
	}

	// Create 10 old entries more than a year old
	var old []*SampleV1.Sample
	for i := 0; i < 10; i++ {
		e := &SampleV1.Sample{
			Name:    fmt.Sprintf("Old%d", i),
			Surname: "User",
			Birth:   now.AddDate(-2, 0, -i*30), // ~2 years ago
		}
		old = append(old, e)
	}

	// Insert all records
	for _, e := range append(recent, old...) {
		if err := e.DbInsert(connection); err != nil {
			t.Fatalf("failed to insert entity: %v", err)
		}
	}

	// Query only those born in the last year, sorted ascending
	q := query.NewQuery().
		SetFilters(query.FilterTypeAnd, []query.Filter{
			{Field: SampleV1.FieldBirth, Op: query.OperatorTypeGreaterThanEqual, Value: lastYear},
		}).
		SetOrders([]query.Order{
			{Type: query.OrderTypeAsc, Field: SampleV1.FieldBirth},
		})

	results, err := SampleV1.DbQuery(connection, q)
	if err != nil {
		t.Fatal(err)
	}

	// Check all returned records are from 'recent' and in ascending order
	if len(results) != len(recent) {
		t.Fatalf("expected %d recent records, got %d", len(recent), len(results))
	}

	prev := time.Time{}
	for i, r := range results {
		if r.Birth.Before(lastYear) {
			t.Fatalf("record %d has birth before last year: %+v", i, r)
		}
		if !prev.IsZero() && r.Birth.Before(prev) {
			t.Fatalf("records not in ascending birth order at index %d", i)
		}
		prev = r.Birth
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
	count := 100000

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

func BenchmarkQueryExecution(b *testing.B) {
	count := 1000000

	b.Run(fmt.Sprintf("Count%d", count), func(b *testing.B) {
		var inserted []*SampleV1.Sample
		for i := 0; i < count; i++ {
			entity := &SampleV1.Sample{
				Name:    fmt.Sprintf("Name%d", i),
				Surname: fmt.Sprintf("Surname%d", i),
			}
			if err := entity.DbInsert(connection); err != nil {
				b.Fatalf("insert failed at i=%d: %v", i, err)
			}
			inserted = append(inserted, entity)
		}

		// Pick one known UUID to query for
		target := inserted[count/2]

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			q := query.NewQuery().
				SetFilters(query.FilterTypeAnd, []query.Filter{
					{Field: SampleV1.FieldUuid, Op: query.OperatorTypeEqual, Value: target.Uuid},
				}).
				//AddOrder(query.OrderTypeAsc, SampleV1.FieldName).
				SetLimit(1)

			results, err := SampleV1.DbQuery(connection, q)
			if err != nil {
				b.Fatal(err)
			}
			if len(results) != 1 || results[0].GetUuid() != target.GetUuid() {
				b.Fatalf("expected match for UUID %v not found", target.GetUuid())
			}
		}
	})
}

func BenchmarkHyperionInsert100kAndSort(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	const totalRows = 1000000
	var inserted []*SampleV1.Sample

	// Insert fresh unique records every run
	for i := 0; i < totalRows; i++ {
		entity := &SampleV1.Sample{
			Name: fmt.Sprintf("SortBench_Unique_%03d", totalRows-i), // reverse order for obvious sorting
		}
		if err := entity.DbInsert(connection); err != nil {
			b.Fatalf("Insert failed: %v", err)
		}
		inserted = append(inserted, entity)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q := query.NewQuery().
			SetFilters(query.FilterTypeOr, func() []query.Filter {
				var filters []query.Filter
				for _, e := range inserted {
					filters = append(filters, query.Filter{
						Field: SampleV1.FieldUuid,
						Op:    query.OperatorTypeEqual,
						Value: e.Uuid,
					})
				}
				return filters
			}()).
			SetOrders([]query.Order{
				{Type: query.OrderTypeAsc, Field: SampleV1.FieldName},
			})

		results, err := SampleV1.DbQuery(connection, q)
		if err != nil {
			b.Fatal(err)
		}

		if len(results) != totalRows {
			b.Fatalf("Expected %d results, got %d", totalRows, len(results))
		}

		for j := 1; j < len(results); j++ {
			if results[j].Name < results[j-1].Name {
				b.Fatalf("Sort order incorrect at index %d: %s < %s", j, results[j].Name, results[j-1].Name)
			}
		}
	}
}
