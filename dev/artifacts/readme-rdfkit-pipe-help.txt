Decode and re-encode using supported encoding formats

Usage:
  rdfkit pipe [flags]

Flags:
  -h, --help                       help for pipe
  -i, --in string                  path or IRI for reading (default stdin)
      --in-base string             override the base IRI of the resource
      --in-param stringArray       extra decode configuration parameters (syntax "KEY[=VALUE]")
      --in-param-io stringArray    extra read configuration parameters (syntax "KEY[=VALUE]")
      --in-type string             name or alias for the decoder (default detect)
  -o, --out string                 path or IRI for writing (default stdout)
      --out-base string            override the base IRI of the resource
      --out-param stringArray      extra encode configuration parameters (syntax "KEY[=VALUE]")
      --out-param-io stringArray   extra write configuration parameters (syntax "KEY[=VALUE]")
      --out-type string            name or alias for the encoder (default detect or nquads)

Encodings:

  org.json-ld.document (decode)

    Aliases: jsonld
    File Extensions: .jsonld
    Media Types: application/ld+json

    --in-param captureTextOffsets[=bool]
      Capture the line+column offsets for statement properties

    --in-param tokenizer.lax[=bool]
      Accept and recover common syntax errors

  org.w3.n-quads (decode, encode)

    Aliases: n-quads, nq, nquads
    File Extensions: .nq
    Media Types: application/n-quads

    --in-param captureTextOffsets[=bool]
      Capture the line+column offsets for statement properties

    --out-param ascii[=bool]
      Use escape sequences for non-ASCII characters

  org.w3.n-triples (decode, encode)

    Aliases: n-triples, nt, ntriples
    File Extensions: .nt
    Media Types: application/n-triples

    --in-param captureTextOffsets[=bool]
      Capture the line+column offsets for statement properties

    --out-param ascii[=bool]
      Use escape sequences for non-ASCII characters

  org.w3.rdf-json (decode, encode)

    Aliases: rdf-json, rdfjson, rj
    File Extensions: .rj
    Media Types: application/rdf+json

    --in-param captureTextOffsets[=bool]
      Capture the line+column offsets for statement properties

  org.w3.rdf-xml (decode)

    Aliases: rdf-xml, rdfxml, xml
    File Extensions: .rdf
    Media Types: application/rdf+xml

    --in-param captureTextOffsets[=bool]
      Capture the line+column offsets for statement properties

  org.w3.trig (decode)

    Aliases: trig
    File Extensions: .trig
    Media Types: application/trig

    --in-param captureTextOffsets[=bool]
      Capture the line+column offsets for statement properties

  org.w3.turtle (decode, encode)

    Aliases: ttl, turtle
    File Extensions: .ttl
    Media Types: text/turtle

    --in-param captureTextOffsets[=bool]
      Capture the line+column offsets for statement properties

    --out-param buffered[=bool]
      Load all statements into memory before writing any output

    --out-param iris.useBase[=bool]
      Prefer IRIs relative to the resource IRI

    --out-param iris.usePrefix=string...
      Prefer IRIs using a prefix. Use the syntax of "{prefix}:{iri}", "rdfa-context", or "none"

    --out-param resources[=bool]
      Write nested statements and resource descriptions (implies buffered=true)

  public.html (decode)

    Aliases: htm, html, xhtml
    File Extensions: .htm, .html, .xhtml
    Media Types: application/xhtml+xml, text/html, text/xhtml+xml

    --in-param captureTextOffsets[=bool]
      Capture the line+column offsets for statement properties
