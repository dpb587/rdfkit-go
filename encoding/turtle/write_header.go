package turtle

import (
	"fmt"
	"io"

	"github.com/dpb587/rdfkit-go/rdf/iriutil"
)

func WriteDocumentHeader(w io.Writer, base string, prefixes iriutil.PrefixMappingList) (int, error) {
	var written int

	if len(base) > 0 {
		n, err := fmt.Fprintf(w, "@base <%s> .\n", base)
		if err != nil {
			return written + n, fmt.Errorf("writing header: %v", err)
		}

		written += n
	}

	for _, t := range prefixes {
		n, err := fmt.Fprintf(w, "@prefix %s: <%s> .\n", t.Prefix, t.Expanded)
		if err != nil {
			return written + n, fmt.Errorf("writing header: %v", err)
		}

		written += n
	}

	return written, nil
}
