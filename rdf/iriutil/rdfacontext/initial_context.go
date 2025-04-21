package rdfacontext

import "github.com/dpb587/rdfkit-go/rdf/iriutil"

var (
	initialContext_as = iriutil.PrefixMapping{
		Prefix:   "as",
		Expanded: "https://www.w3.org/ns/activitystreams#",
	}
	initialContext_csvw = iriutil.PrefixMapping{
		Prefix:   "csvw",
		Expanded: "http://www.w3.org/ns/csvw#",
	}
	initialContext_dcat = iriutil.PrefixMapping{
		Prefix:   "dcat",
		Expanded: "http://www.w3.org/ns/dcat#",
	}
	initialContext_dqv = iriutil.PrefixMapping{
		Prefix:   "dqv",
		Expanded: "http://www.w3.org/ns/dqv#",
	}
	initialContext_duv = iriutil.PrefixMapping{
		Prefix:   "duv",
		Expanded: "https://www.w3.org/ns/duv#",
	}
	initialContext_grddl = iriutil.PrefixMapping{
		Prefix:   "grddl",
		Expanded: "http://www.w3.org/2003/g/data-view#",
	}
	initialContext_jsonld = iriutil.PrefixMapping{
		Prefix:   "jsonld",
		Expanded: "http://www.w3.org/ns/json-ld#",
	}
	initialContext_ldp = iriutil.PrefixMapping{
		Prefix:   "ldp",
		Expanded: "http://www.w3.org/ns/ldp#",
	}
	initialContext_ma = iriutil.PrefixMapping{
		Prefix:   "ma",
		Expanded: "http://www.w3.org/ns/ma-ont#",
	}
	initialContext_oa = iriutil.PrefixMapping{
		Prefix:   "oa",
		Expanded: "http://www.w3.org/ns/oa#",
	}
	initialContext_odrl = iriutil.PrefixMapping{
		Prefix:   "odrl",
		Expanded: "http://www.w3.org/ns/odrl/2/",
	}
	initialContext_org = iriutil.PrefixMapping{
		Prefix:   "org",
		Expanded: "http://www.w3.org/ns/org#",
	}
	initialContext_owl = iriutil.PrefixMapping{
		Prefix:   "owl",
		Expanded: "http://www.w3.org/2002/07/owl#",
	}
	initialContext_prov = iriutil.PrefixMapping{
		Prefix:   "prov",
		Expanded: "http://www.w3.org/ns/prov#",
	}
	initialContext_qb = iriutil.PrefixMapping{
		Prefix:   "qb",
		Expanded: "http://purl.org/linked-data/cube#",
	}
	initialContext_rdf = iriutil.PrefixMapping{
		Prefix:   "rdf",
		Expanded: "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
	}
	initialContext_rdfa = iriutil.PrefixMapping{
		Prefix:   "rdfa",
		Expanded: "http://www.w3.org/ns/rdfa#",
	}
	initialContext_rdfs = iriutil.PrefixMapping{
		Prefix:   "rdfs",
		Expanded: "http://www.w3.org/2000/01/rdf-schema#",
	}
	initialContext_rif = iriutil.PrefixMapping{
		Prefix:   "rif",
		Expanded: "http://www.w3.org/2007/rif#",
	}
	initialContext_rr = iriutil.PrefixMapping{
		Prefix:   "rr",
		Expanded: "http://www.w3.org/ns/r2rml#",
	}
	initialContext_sd = iriutil.PrefixMapping{
		Prefix:   "sd",
		Expanded: "http://www.w3.org/ns/sparql-service-description#",
	}
	initialContext_skos = iriutil.PrefixMapping{
		Prefix:   "skos",
		Expanded: "http://www.w3.org/2004/02/skos/core#",
	}
	initialContext_skosxl = iriutil.PrefixMapping{
		Prefix:   "skosxl",
		Expanded: "http://www.w3.org/2008/05/skos-xl#",
	}
	initialContext_ssn = iriutil.PrefixMapping{
		Prefix:   "ssn",
		Expanded: "http://www.w3.org/ns/ssn/",
	}
	initialContext_sosa = iriutil.PrefixMapping{
		Prefix:   "sosa",
		Expanded: "http://www.w3.org/ns/sosa/",
	}
	initialContext_time = iriutil.PrefixMapping{
		Prefix:   "time",
		Expanded: "http://www.w3.org/2006/time#",
	}
	initialContext_void = iriutil.PrefixMapping{
		Prefix:   "void",
		Expanded: "http://rdfs.org/ns/void#",
	}
	initialContext_wdr = iriutil.PrefixMapping{
		Prefix:   "wdr",
		Expanded: "http://www.w3.org/2007/05/powder#",
	}
	initialContext_wdrs = iriutil.PrefixMapping{
		Prefix:   "wdrs",
		Expanded: "http://www.w3.org/2007/05/powder-s#",
	}
	initialContext_xhv = iriutil.PrefixMapping{
		Prefix:   "xhv",
		Expanded: "http://www.w3.org/1999/xhtml/vocab#",
	}
	initialContext_xml = iriutil.PrefixMapping{
		Prefix:   "xml",
		Expanded: "http://www.w3.org/XML/1998/namespace",
	}
	initialContext_xsd = iriutil.PrefixMapping{
		Prefix:   "xsd",
		Expanded: "http://www.w3.org/2001/XMLSchema#",
	}
)

func InitialContext() iriutil.PrefixMappingList {
	return iriutil.PrefixMappingList{
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
}
