package stmts

import (
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
)

// Logf interface for logging formatted messages such as log.Fatalf, log.Panicf, or log.Printf
type Logf func(format string, v ...interface{})

// Named map of queries and their *sqlx.NamedStmt counterparts.
type Named struct {
	sync.RWMutex

	statements map[string]*sqlx.NamedStmt
}

// Add queries to the Named stmts map. Statements for these queries can be retrieved using the same query after the map
// has been prepared.
func (n *Named) Add(queries ...string) {
	if len(queries) == 0 {
		return
	}
	n.Lock()

	if n.statements == nil {
		n.statements = map[string]*sqlx.NamedStmt{}
	}

	for _, q := range queries {
		if _, exists := n.statements[q]; exists {
			continue
		}
		// Set the value to nil, it will be set again later when preparing statements
		n.statements[q] = nil
	}

	n.Unlock()
	return
}

// Prepare all the queries in the map
func (n Named) Prepare(db *sqlx.DB) error {
	if n.statements == nil {
		return fmt.Errorf("No statements to prepare")
	}
	var err error
	for q := range n.statements {
		n.Lock()
		n.statements[q], err = db.PrepareNamed(q)
		if err != nil {
			return fmt.Errorf("%v in query %s", err, q)
		}
		n.Unlock()
	}
	return nil
}

// MustPrepare prepares all the queries in the map and triggers the provided logf function if something goes wrong
func (n Named) MustPrepare(db *sqlx.DB, logf Logf) {
	err := n.Prepare(db)
	if err != nil {
		logf("Unable to prepare statement with error: %v", err)
	}
}

// Stmt gets a prepared *sqlx.NamedStmt by query
func (n Named) Stmt(query string) *sqlx.NamedStmt {
	if n.statements == nil {
		return nil
	}

	n.RLock()
	s := n.statements[query]
	n.RUnlock()
	return s
}
