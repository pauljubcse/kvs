package kvs

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	//"sync"

	"github.com/gorilla/websocket"
)

type Request struct {
	Action        string      `json:"action"`
	Domain        string      `json:"domain,omitempty"`
	Key           string      `json:"key,omitempty"`
	SLKey         string      `json:"slkey,omitempty"`
	Value         string      `json:"value,omitempty"`
	MinKey        string      `json:"min_key,omitempty"`
	MaxKey        string      `json:"max_key,omitempty"`
}

type Response struct {
	Status        string      `json:"status"`
	Message       string      `json:"message,omitempty"`
	Value         string      `json:"value,omitempty"`
	Values        []string    `json:"values,omitempty"`
}

type Store struct {
	domains map[string]*Domain
	//mu      sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		domains: make(map[string]*Domain),
	}
}

func (s *Store) CreateDomain(name string) {
	//s.mu.Lock()
	//defer s.mu.Unlock()
	s.domains[name] = NewDomain()
}

func (s *Store) SetString(domain, key, value string) error {
	//s.mu.RLock()
	d, ok := s.domains[domain]
	//s.mu.RUnlock()
	if !ok {
		return fmt.Errorf("domain not found")
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	d.stringStore[key] = value
	return nil
}

func (s *Store) GetString(domain, key string) (string, error) {
	//s.mu.RLock()
	d, ok := s.domains[domain]
	//s.mu.RUnlock()
	if !ok {
		return "", fmt.Errorf("domain not found")
	}

	d.mu.RLock()
	defer d.mu.RUnlock()
	value, ok := d.stringStore[key]
	if !ok {
		return "", fmt.Errorf("key not found")
	}
	return value, nil
}

func (s *Store) Increment(domain, key string) error {
	d, ok := s.domains[domain]
	if (!ok) {
		return fmt.Errorf("domain not found")
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	val, err := strconv.Atoi(d.stringStore[key])
	if (err != nil) {
		return fmt.Errorf("value is not an integer")
	}
	d.stringStore[key] = strconv.Itoa(val + 1)
	return nil
}

func (s *Store) Decrement(domain, key string) error {
	d, ok := s.domains[domain]
	if (!ok) {
		return fmt.Errorf("domain not found")
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	val, err := strconv.Atoi(d.stringStore[key])
	if (err != nil) {
		return fmt.Errorf("value is not an integer")
	}
	d.stringStore[key] = strconv.Itoa(val - 1)
	return nil
}

func (s *Store) InsertToSkipList(domain, slkey, key, value string) error {
	intKey, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return fmt.Errorf("key must be integer")
	}

	//s.mu.RLock()
	d, ok := s.domains[domain]
	//s.mu.RUnlock()
	if !ok {
		return fmt.Errorf("domain not found")
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	sl, ok := d.skipListStore[slkey]
	if !ok {
		sl = NewSkipList()
		d.skipListStore[slkey] = sl
	}

	sl.Insert(int(intKey), value)
	return nil
}

func (s *Store) DeleteFromSkipList(domain, slkey, key string) error {
	intKey, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return fmt.Errorf("key must be integer")
	}

	//s.mu.RLock()
	d, ok := s.domains[domain]
	//s.mu.RUnlock()
	if !ok {
		return fmt.Errorf("domain not found")
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	sl, ok := d.skipListStore[slkey]
	if !ok {
		return fmt.Errorf("skip list not found")
	}

	sl.Delete(int(intKey))
	return nil
}

func (s *Store) DeleteRangeFromSkipList(domain, slkey, minKey, maxKey string) error {
	intMinKey, err := strconv.ParseInt(minKey, 10, 64)
	if err != nil {
		return fmt.Errorf("minKey must be integer")
	}
	intMaxKey, err := strconv.ParseInt(maxKey, 10, 64)
	if err != nil {
		return fmt.Errorf("maxKey must be integer")
	}

	//s.mu.RLock()
	d, ok := s.domains[domain]
	//s.mu.RUnlock()
	if !ok {
		return fmt.Errorf("domain not found")
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	sl, ok := d.skipListStore[slkey]
	if !ok {
		return fmt.Errorf("skip list not found")
	}

	sl.DeleteRange(int(intMinKey), int(intMaxKey))
	return nil
}

// func (s *Store) GetAllValuesFromSkipList(domain, slkey string) ([]string, error) {
// 	//s.mu.RLock()
// 	d, ok := s.domains[domain]
// 	//s.mu.RUnlock()
// 	if !ok {
// 		return nil, fmt.Errorf("domain not found")
// 	}

// 	d.mu.RLock()
// 	defer d.mu.RUnlock()
// 	sl, ok := d.skipListStore[slkey]
// 	if !ok {
// 		return nil, fmt.Errorf("skip list not found")
// 	}

// 	return sl.GetAllValues(), nil
// }

func (s *Store) SearchInSkipList(domain, slkey, key string) (string, error) {
	intKey, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return "", fmt.Errorf("key must be integer")
	}

	//s.mu.RLock()
	d, ok := s.domains[domain]
	//s.mu.RUnlock()
	if !ok {
		return "", fmt.Errorf("domain not found")
	}

	d.mu.RLock()
	defer d.mu.RUnlock()
	sl, ok := d.skipListStore[slkey]
	if !ok {
		return "", fmt.Errorf("skip list not found")
	}

	value, found := sl.Search(int(intKey))
	if !found {
		return "", fmt.Errorf("key not found")
	}
	return value, nil
}

// WebSocket connection upgrade
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Store) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	for {
		var req Request
		err := conn.ReadJSON(&req)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error: %v", err)
			}
			break
		}

		var resp Response
		switch req.Action {
		case "create_domain":
			s.CreateDomain(req.Domain)
			resp = Response{Status: "success"}
		case "set_string":
			err := s.SetString(req.Domain, req.Key, req.Value)
			if err != nil {
				resp = Response{Status: "error", Message: err.Error()}
			} else {
				resp = Response{Status: "success"}
			}
		case "get_string":
			value, err := s.GetString(req.Domain, req.Key)
			if err != nil {
				resp = Response{Status: "error", Message: err.Error()}
			} else {
				resp = Response{Status: "success", Value: value}
			}
		case "insert_skiplist":
			err := s.InsertToSkipList(req.Domain, req.SLKey, req.Key, req.Value)
			if err != nil {
				resp = Response{Status: "error", Message: err.Error()}
			} else {
				resp = Response{Status: "success"}
			}
		case "delete_skiplist":
			err := s.DeleteFromSkipList(req.Domain, req.SLKey, req.Key)
			if err != nil {
				resp = Response{Status: "error", Message: err.Error()}
			} else {
				resp = Response{Status: "success"}
			}
		case "delete_range_skiplist":
			err := s.DeleteRangeFromSkipList(req.Domain, req.SLKey, req.MinKey, req.MaxKey)
			if err != nil {
				resp = Response{Status: "error", Message: err.Error()}
			} else {
				resp = Response{Status: "success"}
			}
		// case "get_all_skiplist":
		// 	values, err := s.GetAllValuesFromSkipList(req.Domain, req.SLKey)
		// 	if err != nil {
		// 		resp = Response{Status: "error", Message: err.Error()}
		// 	} else {
		// 		resp = Response{Status: "success", Values: values}
		// 	}
		case "increment":
			err := s.Increment(req.Domain, req.Key)
			if (err != nil) {
				resp = Response{Status: "error", Message: err.Error()}
			} else {
				resp = Response{Status: "success"}
			}
		case "decrement":
			err := s.Decrement(req.Domain, req.Key)
			if (err != nil) {
				resp = Response{Status: "error", Message: err.Error()}
			} else {
				resp = Response{Status: "success"}
			}
		case "search_skiplist":
			value, err := s.SearchInSkipList(req.Domain, req.SLKey, req.Key)
			if err != nil {
				resp = Response{Status: "error", Message: err.Error()}
			} else {
				resp = Response{Status: "success", Value: value}
			}
		default:
			resp = Response{Status: "error", Message: "unknown action"}
		}

		err = conn.WriteJSON(resp)
		if err != nil {
			fmt.Printf("error: %v", err)
			break
		}
	}
}


type Server struct {
	httpServer *http.Server
}

func StartServer(urlStr string) (*Server, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	store := NewStore()
	http.HandleFunc(u.Path, func(w http.ResponseWriter, r *http.Request) {
		store.HandleWebSocket(w, r)
	})

	server := &http.Server{Addr: u.Host, Handler: nil}

	go func() {
		fmt.Printf("Starting server on %s...\n", u.Host)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	return &Server{httpServer: server}, nil
}

func (s *Server) CloseServer() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}