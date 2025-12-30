package inmemory

import (
	"context"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/x/storage/inmemory/simplequery"
)

func (s *Dataset) QuerySimple(ctx context.Context, q simplequery.Query, opts simplequery.QueryOptions) (simplequery.QueryResult, error) {
	var selectBindings []simplequery.QueryResultBinding

	first, remaining := q.Where.SortStaticEfficiency().Shift()

	err := s.querySimpleWhere(
		ctx,
		q,
		&selectBindings,
		simplequery.NewQueryResultBinding(map[string]rdf.Term{}),
		first,
		remaining,
	)
	if err != nil {
		return nil, err
	}

	return simplequery.NewQueryResult(selectBindings), nil
}

func (s *Dataset) querySimpleWhere(ctx context.Context, q simplequery.Query, selectBindings *[]simplequery.QueryResultBinding, currentBindings simplequery.QueryResultBinding, whereStatement simplequery.WhereTriple, remaining simplequery.WhereTripleList) error {
	iter, err := s.NewQuadIterator(ctx, whereStatement)
	if err != nil {
		return err
	}

	var oValues rdf.TermList
	var oValuesFound bool

	if oVar, ok := whereStatement.Object.(simplequery.Var); ok {
		vs := q.Values.GetByVar(oVar)

		oValues = vs.Terms
		oValuesFound = len(vs.Var) > 0
	}

	var found int

	for _, edge := range iter.(*QuadIterator).edges {
		if oValuesFound {
			var matched bool

			for _, oValueTerm := range oValues {
				if oValueTerm.TermEquals(edge.GetQuad().Triple.Object) {
					matched = true

					break
				}
			}

			if !matched {
				continue
			}
		}

		found++

		nextBindings := whereStatement.UpdateBindings(currentBindings, edge.GetQuad())

		if len(remaining) == 0 {
			finalTermsByVar := map[string]rdf.Term{}

			for _, v := range q.Select {
				finalTermsByVar[string(v)] = nextBindings.Get(string(v))
			}

			*selectBindings = append(*selectBindings, simplequery.NewQueryResultBinding(finalTermsByVar))

			continue
		} else {
			nextStatement, nextRemaining := remaining.ResolveBindings(nextBindings).SortStaticEfficiency().Shift()

			err := s.querySimpleWhere(ctx, q, selectBindings, nextBindings, nextStatement, nextRemaining)
			if err != nil {
				return err
			}
		}
	}

	if found == 0 && whereStatement.Optional {
		if len(remaining) == 0 {
			finalTermsByVar := map[string]rdf.Term{}

			for _, v := range q.Select {
				finalTermsByVar[string(v)] = currentBindings.Get(string(v))
			}

			*selectBindings = append(*selectBindings, simplequery.NewQueryResultBinding(finalTermsByVar))
		} else {
			nextStatement, nextRemaining := remaining.Shift()

			err := s.querySimpleWhere(ctx, q, selectBindings, currentBindings, nextStatement, nextRemaining)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
