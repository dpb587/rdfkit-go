package rdfacontext

import "github.com/dpb587/rdfkit-go/rdf/iriutil"

var (
	widelyUsedInitialContext_cc = iriutil.PrefixMapping{
		Prefix:   "cc",
		Expanded: "http://creativecommons.org/ns#",
	}
	widelyUsedInitialContext_ctag = iriutil.PrefixMapping{
		Prefix:   "ctag",
		Expanded: "http://commontag.org/ns#",
	}
	widelyUsedInitialContext_dc = iriutil.PrefixMapping{
		Prefix:   "dc",
		Expanded: "http://purl.org/dc/terms/",
	}
	widelyUsedInitialContext_dcterms = iriutil.PrefixMapping{
		Prefix:   "dcterms",
		Expanded: "http://purl.org/dc/terms/",
	}
	widelyUsedInitialContext_dc11 = iriutil.PrefixMapping{
		Prefix:   "dc11",
		Expanded: "http://purl.org/dc/elements/1.1/",
	}
	widelyUsedInitialContext_foaf = iriutil.PrefixMapping{
		Prefix:   "foaf",
		Expanded: "http://xmlns.com/foaf/0.1/",
	}
	widelyUsedInitialContext_gr = iriutil.PrefixMapping{
		Prefix:   "gr",
		Expanded: "http://purl.org/goodrelations/v1#",
	}
	widelyUsedInitialContext_ical = iriutil.PrefixMapping{
		Prefix:   "ical",
		Expanded: "http://www.w3.org/2002/12/cal/icaltzd#",
	}
	widelyUsedInitialContext_og = iriutil.PrefixMapping{
		Prefix:   "og",
		Expanded: "http://ogp.me/ns#",
	}
	widelyUsedInitialContext_rev = iriutil.PrefixMapping{
		Prefix:   "rev",
		Expanded: "http://purl.org/stuff/rev#",
	}
	widelyUsedInitialContext_sioc = iriutil.PrefixMapping{
		Prefix:   "sioc",
		Expanded: "http://rdfs.org/sioc/ns#",
	}
	widelyUsedInitialContext_v = iriutil.PrefixMapping{
		Prefix:   "v",
		Expanded: "http://rdf.data-vocabulary.org/#",
	}
	widelyUsedInitialContext_vcard = iriutil.PrefixMapping{
		Prefix:   "vcard",
		Expanded: "http://www.w3.org/2006/vcard/ns#",
	}
	widelyUsedInitialContext_schema = iriutil.PrefixMapping{
		Prefix:   "schema",
		Expanded: "http://schema.org/",
	}
)

func WidelyUsedInitialContext() iriutil.PrefixMappingList {
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
		widelyUsedInitialContext_cc,
		widelyUsedInitialContext_ctag,
		widelyUsedInitialContext_dc,
		widelyUsedInitialContext_dcterms,
		widelyUsedInitialContext_dc11,
		widelyUsedInitialContext_foaf,
		widelyUsedInitialContext_gr,
		widelyUsedInitialContext_ical,
		widelyUsedInitialContext_og,
		widelyUsedInitialContext_rev,
		widelyUsedInitialContext_sioc,
		widelyUsedInitialContext_v,
		widelyUsedInitialContext_vcard,
		widelyUsedInitialContext_schema,
	}
}
