package kvs

import (
	"net/http"
	"net/http/httptest"
	"testing"

	//"your_module_path/kvs" // replace with the actual module path

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestWebSocketStore(t *testing.T) {
	store := NewStore()
	server := httptest.NewServer(http.HandlerFunc(store.HandleWebSocket))
	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)
	defer conn.Close()

	// Create Domain
	createDomainRequest := Request{
		Action: "create_domain",
		Domain: "test_domain",
	}
	conn.WriteJSON(createDomainRequest)
	var createDomainResponse Response
	conn.ReadJSON(&createDomainResponse)
	assert.Equal(t, "success", createDomainResponse.Status)

	// Set String
	setStringRequest := Request{
		Action: "set_string",
		Domain: "test_domain",
		Key:    "test_key",
		Value:  "test_value",
	}
	conn.WriteJSON(setStringRequest)
	var setStringResponse Response
	conn.ReadJSON(&setStringResponse)
	assert.Equal(t, "success", setStringResponse.Status)

	// Get String
	getStringRequest := Request{
		Action: "get_string",
		Domain: "test_domain",
		Key:    "test_key",
	}
	conn.WriteJSON(getStringRequest)
	var getStringResponse Response
	conn.ReadJSON(&getStringResponse)
	assert.Equal(t, "success", getStringResponse.Status)
	assert.Equal(t, "test_value", getStringResponse.Value)

	// Insert to SkipList
	insertSkipListRequest := Request{
		Action: "insert_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		Key:    "1",
		Value:  "value1",
	}
	conn.WriteJSON(insertSkipListRequest)
	var insertSkipListResponse Response
	conn.ReadJSON(&insertSkipListResponse)
	assert.Equal(t, "success", insertSkipListResponse.Status)

	// Search in SkipList
	searchSkipListRequest := Request{
		Action: "search_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		Key:    "1",
	}
	conn.WriteJSON(searchSkipListRequest)
	var searchSkipListResponse Response
	conn.ReadJSON(&searchSkipListResponse)
	assert.Equal(t, "success", searchSkipListResponse.Status)
	assert.Equal(t, "value1", searchSkipListResponse.Value)

	// Delete from SkipList
	deleteSkipListRequest := Request{
		Action: "delete_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		Key:    "1",
	}
	conn.WriteJSON(deleteSkipListRequest)
	var deleteSkipListResponse Response
	conn.ReadJSON(&deleteSkipListResponse)
	assert.Equal(t, "success", deleteSkipListResponse.Status)

	// Confirm Deletion from SkipList
	searchSkipListRequest2 := Request{
		Action: "search_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		Key:    "1",
	}
	conn.WriteJSON(searchSkipListRequest2)
	var searchSkipListResponse2 Response
	conn.ReadJSON(&searchSkipListResponse2)
	assert.Equal(t, "error", searchSkipListResponse2.Status)
	assert.Equal(t, "key not found", searchSkipListResponse2.Message)

	// Insert Range in SkipList
	insertSkipListRequest1 := Request{
		Action: "insert_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		Key:    "1",
		Value:  "value1",
	}
	insertSkipListRequest2 := Request{
		Action: "insert_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		Key:    "2",
		Value:  "value2",
	}
	insertSkipListRequest3 := Request{
		Action: "insert_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		Key:    "3",
		Value:  "value3",
	}
	conn.WriteJSON(insertSkipListRequest1)
	conn.ReadJSON(&insertSkipListResponse)
	conn.WriteJSON(insertSkipListRequest2)
	conn.ReadJSON(&insertSkipListResponse)
	conn.WriteJSON(insertSkipListRequest3)
	conn.ReadJSON(&insertSkipListResponse)

	// Delete Range from SkipList
	deleteRangeSkipListRequest := Request{
		Action: "delete_range_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		MinKey: "1",
		MaxKey: "2",
	}
	conn.WriteJSON(deleteRangeSkipListRequest)
	var deleteRangeSkipListResponse Response
	conn.ReadJSON(&deleteRangeSkipListResponse)
	assert.Equal(t, "success", deleteRangeSkipListResponse.Status)

	// Confirm Range Deletion from SkipList
	searchSkipListRequest3 := Request{
		Action: "search_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		Key:    "1",
	}
	conn.WriteJSON(searchSkipListRequest3)
	var searchSkipListResponse3 Response
	conn.ReadJSON(&searchSkipListResponse3)
	assert.Equal(t, "error", searchSkipListResponse3.Status)
	assert.Equal(t, "key not found", searchSkipListResponse3.Message)

	searchSkipListRequest4 := Request{
		Action: "search_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		Key:    "2",
	}
	conn.WriteJSON(searchSkipListRequest4)
	var searchSkipListResponse4 Response
	conn.ReadJSON(&searchSkipListResponse4)
	assert.Equal(t, "error", searchSkipListResponse4.Status)
	assert.Equal(t, "key not found", searchSkipListResponse4.Message)

	searchSkipListRequest5 := Request{
		Action: "search_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		Key:    "3",
	}
	conn.WriteJSON(searchSkipListRequest5)
	var searchSkipListResponse5 Response
	conn.ReadJSON(&searchSkipListResponse5)
	assert.Equal(t, "success", searchSkipListResponse5.Status)
	assert.Equal(t, "value3", searchSkipListResponse5.Value)
}
func TestRank(t *testing.T) {
	// Initialize the store and create a domain
	store := NewStore()
	domain := "test_domain"
	slKey := "test_skiplist"
	store.CreateDomain(domain)

	// Insert elements into the skip list
	elements := []struct {
		key   string
		value string
	}{
		{"1", "one"},
		{"2", "two"},
		{"3", "three"},
		{"5", "five"},
		{"7", "seven"},
	}
	for _, elem := range elements {
		err := store.InsertToSkipList(domain, slKey, elem.key, elem.value)
		if err != nil {
			t.Fatalf("failed to insert %s: %v", elem.key, err)
		}
	}

	tests := []struct {
		key        string
		expectedRank string
	}{
		{"1", "0"},
		{"2", "1"},
		{"3", "2"},
		{"4", "3"},  // Non-existent element, should return rank as if it was present
		{"5", "3"},
		{"6", "4"},  // Non-existent element, should return rank as if it was present
		{"7", "4"},
		{"8", "5"},  // Non-existent element, should return rank as if it was present
	}

	for _, tt := range tests {
		t.Run("rank_of_"+tt.key, func(t *testing.T) {
			rank, err := store.RankInSkipList(domain, slKey, tt.key)
			if err != nil {
				t.Fatalf("failed to get rank for %s: %v", tt.key, err)
			}
			if rank != tt.expectedRank {
				t.Errorf("rank of %s = %s; want %s", tt.key, rank, tt.expectedRank)
			}
		})
	}
}