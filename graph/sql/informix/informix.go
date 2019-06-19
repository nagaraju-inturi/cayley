package informix

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/cayleygraph/cayley/clog"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/graph/log"
	csql "github.com/cayleygraph/cayley/graph/sql"
	"github.com/cayleygraph/cayley/quad"
       _ "github.com/alexbrainman/odbc"
)

const Type = "informix"

var QueryDialect = csql.QueryDialect{
	RegexpOp: "REGEXP",
	FieldQuote: func(name string) string {
		return "`" + name + "`"
	},
	Placeholder: func(n int) string { return "?" },
}

func init() {
	csql.Register(Type, csql.Registration{
		Driver:               "odbc",
		HashType:             fmt.Sprintf(`CHAR(%d)`, quad.HashSize*2),
		BytesType:            `CHAR(2048)`,
		HorizonType:          `BIGSERIAL`,
		TimeType:             `DATETIME YEAR TO SECOND`,
		QueryDialect:         QueryDialect,
		NoOffsetWithoutLimit: true,
		Error: func(err error) error {
			return err
		},
		Estimated: nil,
		RunTx:     runTxInformix,
	})
}

func runTxInformix(tx *sql.Tx, nodes []graphlog.NodeUpdate, quads []graphlog.QuadUpdate, opts graph.IgnoreOpts) error {
	// update node ref counts and insert nodes
	var (
		// prepared statements for each value type
		insertValue = make(map[csql.ValueType]*sql.Stmt)
		updateValue *sql.Stmt
	)
	for _, n := range nodes {
		if n.RefInc >= 0 {
			nodeKey, values, err := csql.NodeValues(csql.NodeHash{n.Hash}, n.Val)
			if err != nil {
				return err
			}
			values = append([]interface{}{n.RefInc}, values...)
			//nag values = append(values, n.RefInc) // one more time for UPDATE
			stmt, ok := insertValue[nodeKey]
			if !ok {
				// nag var ph = make([]string, len(values)-1) // excluding last increment
				var ph = make([]string, len(values)) // excluding last increment
				for i := range ph {
					ph[i] = "?"
				}
				stmt, err = tx.Prepare(`INSERT INTO nodes(refs, hash, ` +
					strings.Join(nodeKey.Columns(), ", ") +
					`) VALUES (` + strings.Join(ph, ", ") +
                                        `) ;`)
                                        //nag `) ON DUPLICATE KEY UPDATE refs = refs + ?;`)
				if err != nil {
					return err
				}
				insertValue[nodeKey] = stmt
			}
			_, err = stmt.Exec(values...)
			err = convInsertError(err)
			// if err != nil {
			if err != nil && (strings.Contains(strings.ToUpper(err.Error()), "UNIQUE") == false) {
				clog.Errorf("couldn't exec INSERT statement: %v %v", err, values)
				return err
			}
		} else {
			panic("unexpected node update")
		}
	}
	for _, s := range insertValue {
		s.Close()
	}
	if s := updateValue; s != nil {
		s.Close()
	}
	insertValue = nil
	updateValue = nil

	// now we can deal with quads
	ignore := ""
	if opts.IgnoreDup {
		// nag ignore = " IGNORE"
		 ignore = ""
	}

	var (
		insertQuad *sql.Stmt
		err        error
	)
	for _, d := range quads {
		dirs := make([]interface{}, 0, len(quad.Directions))
		for _, h := range d.Quad.Dirs() {
			dirs = append(dirs, csql.NodeHash{h}.SQLValue())
		}
		if !d.Del {
			if insertQuad == nil {
				insertQuad, err = tx.Prepare(`INSERT` + ignore + ` INTO quads(subject_hash, predicate_hash, object_hash, label_hash, ts) VALUES (?, ?, ?, ?, current);`)
				if err != nil {
					return err
				}
				insertValue = make(map[csql.ValueType]*sql.Stmt)
			}
			_, err := insertQuad.Exec(dirs...)
			err = convInsertError(err)
			if err != nil {
				if _, ok := err.(*graph.DeltaError); !ok {
					clog.Errorf("couldn't exec INSERT statement: %v", err)
				}
				return err
			}
		} else {
			panic("unexpected quad delete")
		}
	}
	return nil
}

func convInsertError(err error) error {
	if err == nil {
		return nil
	}
//	if e, ok := err.(*odbc.SQLError); ok {
//		if e.Number == 1062 {
//			// TODO: reference to delta
//			return &graph.DeltaError{Err: graph.ErrQuadExists}
//		}
//	}
	return err
}
