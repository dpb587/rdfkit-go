package simplequery

import (
	"sort"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdfio"
)

type QueryScope int

const (
	QueryScopeGlobal QueryScope = iota
	QueryScopeSubjectEntity
)

type QueryOptions struct {
	Scope QueryScope
}

type Query struct {
	Select []Var
	Where  WhereTripleList
	Values ValueSetList
}

type ValueSet struct {
	Var   Var
	Terms rdf.TermList
}

type ValueSetList []ValueSet

func (vsl ValueSetList) GetByVar(v Var) ValueSet {
	for _, vs := range vsl {
		if vs.Var == v {
			return vs
		}
	}

	return ValueSet{}
}

type WhereTriple struct {
	Subject   VarOrTerm
	Predicate VarOrTerm
	Object    VarOrTerm
	Optional  bool
}

var _ rdfio.StatementMatcher = WhereTriple{}

func (wt WhereTriple) staticEfficiency() int {
	var offset int

	if wt.Optional {
		offset = -100
	}

	return offset + wt.Subject.staticEfficiency('s') + wt.Predicate.staticEfficiency('p') + wt.Object.staticEfficiency('o')
}

func (wt WhereTriple) ResolveBindings(res QueryResultBinding) WhereTriple {
	if s, ok := wt.Subject.(Var); ok {
		if bound, ok := res.termsByVar[string(s)]; ok {
			wt.Subject = Term{Term: bound}
		}
	}

	if p, ok := wt.Predicate.(Var); ok {
		if bound, ok := res.termsByVar[string(p)]; ok {
			wt.Predicate = Term{Term: bound}
		}
	}

	if o, ok := wt.Object.(Var); ok {
		if bound, ok := res.termsByVar[string(o)]; ok {
			wt.Object = Term{Term: bound}
		}
	}

	return wt
}

func (wt WhereTriple) MatchStatement(e rdfio.Statement) bool {
	et := e.GetTriple()

	if s, ok := wt.Subject.(Term); ok {
		if !s.Term.TermEquals(et.Subject) {
			return false
		}
	}

	if p, ok := wt.Predicate.(Term); ok {
		if !p.Term.TermEquals(et.Predicate) {
			return false
		}
	}

	if o, ok := wt.Object.(Term); ok {
		if !o.Term.TermEquals(et.Object) {
			return false
		}
	}

	return true
}

func (wt WhereTriple) UpdateBindings(res QueryResultBinding, e rdfio.Statement) QueryResultBinding {
	res = res.Clone()

	et := e.GetTriple()

	if s, ok := wt.Subject.(Var); ok {
		res.termsByVar[string(s)] = et.Subject
		res.tripleBindingsByVar[string(s)] = e
	}

	if p, ok := wt.Predicate.(Var); ok {
		res.termsByVar[string(p)] = et.Predicate
		res.tripleBindingsByVar[string(p)] = e
	}

	if o, ok := wt.Object.(Var); ok {
		res.termsByVar[string(o)] = et.Object
		res.tripleBindingsByVar[string(o)] = e
	}

	return res
}

//

type VarOrTerm interface {
	staticEfficiency(role byte) int
}

type Var string

var _ VarOrTerm = Var("")

func (v Var) staticEfficiency(role byte) int {
	return 0
}

//

type Term struct {
	Term rdf.Term
}

var _ VarOrTerm = Term{}

func (t Term) staticEfficiency(role byte) int {
	switch t.Term.(type) {
	case rdf.BlankNode, rdf.IRI:
		return 1
	}

	return 0
}

//

type WhereTripleList []WhereTriple

func (wtl WhereTripleList) ResolveBindings(res QueryResultBinding) WhereTripleList {
	var resolved = make(WhereTripleList, len(wtl))

	for wtIdx, wt := range wtl {
		resolved[wtIdx] = wt.ResolveBindings(res)
	}

	return resolved
}

func (wtl WhereTripleList) SortStaticEfficiency() WhereTripleList {
	sorted := wtl[:]

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].staticEfficiency() > sorted[j].staticEfficiency()
	})

	return sorted
}

func (wtl WhereTripleList) Shift() (WhereTriple, WhereTripleList) {
	if len(wtl) == 0 {
		return WhereTriple{}, nil
	}

	return wtl[0], wtl[1:]
}
