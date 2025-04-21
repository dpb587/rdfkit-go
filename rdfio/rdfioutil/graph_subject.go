package rdfioutil

import (
	"context"
	"fmt"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/termutil"
	"github.com/dpb587/rdfkit-go/rdfio"
)

type GraphSubjectDropOptions struct {
	RecurseDetachedBlankNodes bool
}

func GraphSubjectDrop(ctx context.Context, g rdfio.Graph, s rdf.SubjectValue, opts GraphSubjectDropOptions) error {
	iter, err := g.NewStatementIterator(ctx, rdfio.SubjectStatementMatcher{
		Matcher: termutil.Equals{
			Expected: s,
		},
	})
	if err != nil {
		return err
	}

	defer iter.Close()

	for iter.Next() {
		triple := iter.GetStatement().GetTriple()

		if err := g.DeleteTriple(ctx, triple); err != nil {
			return err
		}

		if !opts.RecurseDetachedBlankNodes {
			continue
		} else if _, ok := triple.Object.(rdf.BlankNode); !ok {
			continue
		}

		err = func() error {
			recurseIter, err := g.NewStatementIterator(ctx, rdfio.ObjectStatementMatcher{
				Matcher: termutil.Equals{
					Expected: triple.Object,
				},
			})
			if err != nil {
				return err
			}

			defer recurseIter.Close()

			for recurseIter.Next() {
				return nil
			}

			return GraphSubjectDrop(ctx, g, triple.Object.(rdf.BlankNode), opts)
		}()
		if err != nil {
			return fmt.Errorf("recursing: %w", err)
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}
