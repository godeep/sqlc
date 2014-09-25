package sqlc

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

var predicateTypes = map[PredicateType]string{
	EqPredicate: "=",
	GtPredicate: ">",
	GePredicate: ">=",
	LtPredicate: "<",
	LePredicate: "<=",
}

func (s *selection) String() string {
	var buf bytes.Buffer
	s.Render(&buf)
	return buf.String()
}

func (s *selection) Render(w io.Writer) (placeholders []interface{}) {

	var alias string

	// TODO This type switch is used twice, consider refactoring
	switch sub := s.selection.(type) {
	case table:
		alias = sub.name
	}

	fmt.Fprint(w, "SELECT ")

	if len(s.projection) == 0 {
		fmt.Fprint(w, "*")
	} else {
		colClause := columnClause(alias, s.projection)
		fmt.Fprint(w, colClause)
	}

	fmt.Fprintf(w, " FROM ")

	switch sub := s.selection.(type) {
	case table:
		fmt.Fprint(w, sub.name)
	case *selection:
		fmt.Fprint(w, "(")
		sub.Render(w)
		fmt.Fprint(w, ")")
	}

	// TODO Support more than one join
	if len(s.joins) == 1 {
		join := s.joins[0]
		fmt.Fprintf(w, " JOIN %s ON %s = %s", join.target.Name(), join.conds[0].Binding.Field.Name(), join.conds[0].Binding.Value)
	}

	if len(s.predicate) > 0 {
		fmt.Fprint(w, " ")
		placeholders = renderWhereClause(alias, s.predicate, w)
	} else {
		placeholders = []interface{}{}
	}

	if (len(s.groups)) > 0 {
		fmt.Fprint(w, " GROUP BY ")
		colClause := columnClause(alias, s.groups)
		fmt.Fprint(w, colClause)
	}

	// TODO eliminate copy and paste
	if (len(s.ordering)) > 0 {
		fmt.Fprint(w, " ORDER BY ")
		colClause := columnClause(alias, s.ordering)
		fmt.Fprint(w, colClause)
	}

	return placeholders
}

func columnClause(alias string, cols []Field) string {
	colFragments := make([]string, len(cols))
	for i, col := range cols {
		colFragments[i] = fmt.Sprintf("%s.%s", alias, col.Name())
	}
	return strings.Join(colFragments, ", ")
}

func renderWhereClause(alias string, conds []Condition, w io.Writer) []interface{} {
	fmt.Fprint(w, "WHERE ")

	whereFragments := make([]string, len(conds))
	values := make([]interface{}, len(conds))

	for i, condition := range conds {
		col := condition.Binding.Field.Name()
		pred := condition.Predicate
		whereFragments[i] = fmt.Sprintf("%s.%s %s ?", alias, col, predicateTypes[pred])
		values[i] = condition.Binding.Value
	}

	whereClause := strings.Join(whereFragments, " AND ")
	fmt.Fprint(w, whereClause)

	return values
}
