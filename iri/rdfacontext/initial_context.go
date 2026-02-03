package rdfacontext

import "github.com/dpb587/rdfkit-go/iri"

var (
	initialContext_as = iri.PrefixMapping{
		Prefix:   "as",
		Expanded: "https://www.w3.org/ns/activitystreams#",
	}
	initialContext_csvw = iri.PrefixMapping{
		Prefix:   "csvw",
		Expanded: "http://www.w3.org/ns/csvw#",
	}
	initialContext_dcat = iri.PrefixMapping{
		Prefix:   "dcat",
		Expanded: "http://www.w3.org/ns/dcat#",
	}
	initialContext_dqv = iri.PrefixMapping{
		Prefix:   "dqv",
		Expanded: "http://www.w3.org/ns/dqv#",
	}
	initialContext_duv = iri.PrefixMapping{
		Prefix:   "duv",
		Expanded: "https://www.w3.org/ns/duv#",
	}
	initialContext_grddl = iri.PrefixMapping{
		Prefix:   "grddl",
		Expanded: "http://www.w3.org/2003/g/data-view#",
	}
	initialContext_jsonld = iri.PrefixMapping{
		Prefix:   "jsonld",
		Expanded: "http://www.w3.org/ns/json-ld#",
	}
	initialContext_ldp = iri.PrefixMapping{
		Prefix:   "ldp",
		Expanded: "http://www.w3.org/ns/ldp#",
	}
	initialContext_ma = iri.PrefixMapping{
		Prefix:   "ma",
		Expanded: "http://www.w3.org/ns/ma-ont#",
	}
	initialContext_oa = iri.PrefixMapping{
		Prefix:   "oa",
		Expanded: "http://www.w3.org/ns/oa#",
	}
	initialContext_odrl = iri.PrefixMapping{
		Prefix:   "odrl",
		Expanded: "http://www.w3.org/ns/odrl/2/",
	}
	initialContext_org = iri.PrefixMapping{
		Prefix:   "org",
		Expanded: "http://www.w3.org/ns/org#",
	}
	initialContext_owl = iri.PrefixMapping{
		Prefix:   "owl",
		Expanded: "http://www.w3.org/2002/07/owl#",
	}
	initialContext_prov = iri.PrefixMapping{
		Prefix:   "prov",
		Expanded: "http://www.w3.org/ns/prov#",
	}
	initialContext_qb = iri.PrefixMapping{
		Prefix:   "qb",
		Expanded: "http://purl.org/linked-data/cube#",
	}
	initialContext_rdf = iri.PrefixMapping{
		Prefix:   "rdf",
		Expanded: "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
	}
	initialContext_rdfa = iri.PrefixMapping{
		Prefix:   "rdfa",
		Expanded: "http://www.w3.org/ns/rdfa#",
	}
	initialContext_rdfs = iri.PrefixMapping{
		Prefix:   "rdfs",
		Expanded: "http://www.w3.org/2000/01/rdf-schema#",
	}
	initialContext_rif = iri.PrefixMapping{
		Prefix:   "rif",
		Expanded: "http://www.w3.org/2007/rif#",
	}
	initialContext_rr = iri.PrefixMapping{
		Prefix:   "rr",
		Expanded: "http://www.w3.org/ns/r2rml#",
	}
	initialContext_sd = iri.PrefixMapping{
		Prefix:   "sd",
		Expanded: "http://www.w3.org/ns/sparql-service-description#",
	}
	initialContext_skos = iri.PrefixMapping{
		Prefix:   "skos",
		Expanded: "http://www.w3.org/2004/02/skos/core#",
	}
	initialContext_skosxl = iri.PrefixMapping{
		Prefix:   "skosxl",
		Expanded: "http://www.w3.org/2008/05/skos-xl#",
	}
	initialContext_ssn = iri.PrefixMapping{
		Prefix:   "ssn",
		Expanded: "http://www.w3.org/ns/ssn/",
	}
	initialContext_sosa = iri.PrefixMapping{
		Prefix:   "sosa",
		Expanded: "http://www.w3.org/ns/sosa/",
	}
	initialContext_time = iri.PrefixMapping{
		Prefix:   "time",
		Expanded: "http://www.w3.org/2006/time#",
	}
	initialContext_void = iri.PrefixMapping{
		Prefix:   "void",
		Expanded: "http://rdfs.org/ns/void#",
	}
	initialContext_wdr = iri.PrefixMapping{
		Prefix:   "wdr",
		Expanded: "http://www.w3.org/2007/05/powder#",
	}
	initialContext_wdrs = iri.PrefixMapping{
		Prefix:   "wdrs",
		Expanded: "http://www.w3.org/2007/05/powder-s#",
	}
	initialContext_xhv = iri.PrefixMapping{
		Prefix:   "xhv",
		Expanded: "http://www.w3.org/1999/xhtml/vocab#",
	}
	initialContext_xml = iri.PrefixMapping{
		Prefix:   "xml",
		Expanded: "http://www.w3.org/XML/1998/namespace",
	}
	initialContext_xsd = iri.PrefixMapping{
		Prefix:   "xsd",
		Expanded: "http://www.w3.org/2001/XMLSchema#",
	}
)

var initialContext = iri.PrefixMappingList{
	initialContext_as,
	initialContext_csvw,
	initialContext_dcat,
	initialContext_dqv,
	initialContext_duv,
	initialContext_grddl,
	initialContext_jsonld,
	initialContext_ldp,
	initialContext_ma,
	initialContext_oa,
	initialContext_odrl,
	initialContext_org,
	initialContext_owl,
	initialContext_prov,
	initialContext_qb,
	initialContext_rdf,
	initialContext_rdfa,
	initialContext_rdfs,
	initialContext_rif,
	initialContext_rr,
	initialContext_sd,
	initialContext_skos,
	initialContext_skosxl,
	initialContext_ssn,
	initialContext_sosa,
	initialContext_time,
	initialContext_void,
	initialContext_wdr,
	initialContext_wdrs,
	initialContext_xhv,
	initialContext_xml,
	initialContext_xsd,
}

func NewInitialContext() *iri.PrefixManager {
	return iri.NewPrefixManager(initialContext)
}

func AppendInitialContext(base iri.PrefixMappingList) iri.PrefixMappingList {
	return append(base, initialContext...)
}
