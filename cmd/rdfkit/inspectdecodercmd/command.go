package inspectdecodercmd

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/dpb587/rdfkit-go/cmd/rdfkit/cmdflags"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/nquads"
	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/spf13/cobra"
)

var tmpl = template.Must(template.New("inspect").Parse(`<!DOCTYPE html>
<html>
	<head>
		<title>inspect</title>
		<meta http-equiv="Content-Type" content="text/html;charset=utf-8" />
  	<script src="https://cdn.tailwindcss.com"></script>
		<style type="text/tailwindcss">
			#output {
				@apply cursor-default;
			}

			#output span[title]:hover {
				@apply cursor-pointer bg-amber-200;
			}
		</style>
	</head>
	<body>
		<div class="fixed inset-0">
			<div id="input" class="absolute inset-x-0 top-0 bottom-1/2 border-b border-stone-200"></div>
			<div id="output" class="text-sm font-mono px-2 py-1.5 leading-5 absolute -z-10 inset-x-0 bottom-0 top-1/2 overflow-y-auto bg-stone-100 text-stone-900">
				{{- range .ParseRows -}}
					<div>{{/*
						*/}}<span{{ with .SubjectRange }} title="{{ . }}"{{ end }}>{{ .Subject }}</span> {{/*
						*/}}<span{{ with .PredicateRange }} title="{{ . }}"{{ end }}>{{ .Predicate }}</span> {{/*
						*/}}<span{{ with .ObjectRange }} title="{{ . }}"{{ end }}>{{ .Object }}</span> {{/*
						*/}}{{ if .GraphName }}<span{{ with .GraphNameRange }} title="{{ . }}"{{ end }}>{{ .GraphName }}</span> {{ end }}{{/*
						*/}}.{{/*
					*/}}</div>
{{ end -}}
			</div>
		</div>

		<script src="https://unpkg.com/monaco-editor@latest/min/vs/loader.js"></script>
		<script>
			require.config({ paths: { vs: 'https://unpkg.com/monaco-editor@latest/min/vs' } });

			require(['vs/editor/editor.main'], function () {
				window.rdfkitEditor = monaco.editor.create(
					document.getElementById('input'),
					{
						value: {{ .Source }},
						language: {{ .Language }},
					},
				);

				document.getElementById('output').addEventListener('click', function (e) {
					let spanRange = null;

					for (let el = e.target; el; el = el.parentElement) {
						if (el.tagName === 'SPAN' && el.hasAttribute('title')) {
							spanRange = el;

							break;
						} else if (el.tagName === 'PRE') {
							break;
						}
					}

					if (!spanRange) {
						return;
					}
					
					const textPosition = spanRange.getAttribute('title').split(';', 2)[0].split(':', 2);
			
					const textPositionStart = textPosition[0].match(/^L(\d+)C(\d+)$/);
					if (!textPositionStart) {
						return;
					}
			
					const textPositionStartModel = new monaco.Position(parseInt(textPositionStart[1]), parseInt(textPositionStart[2]));
			
					const textPositionEnd = textPosition[1] ? textPosition[1].match(/^L(\d+)C(\d+)$/) : textPositionStart;
					if (!textPositionEnd) {
						return;
					}
			
					const textPositionEndModel = new monaco.Position(parseInt(textPositionEnd[1]), parseInt(textPositionEnd[2]));
			

					if (window.rdfkitCurrent) {
						window.rdfkitCurrent.className = window.rdfkitCurrent.className.replace(/ bg-amber-400/, '');
						window.rdfkitCurrentDecorations.clear();
					}
			
					window.rdfkitCurrent = spanRange;
					window.rdfkitCurrent.className += ' bg-amber-400';
					window.rdfkitCurrentDecorations = window.rdfkitEditor.createDecorationsCollection([{
						range: new monaco.Range(
								textPositionStartModel.lineNumber,
								textPositionStartModel.column,
								textPositionEndModel.lineNumber,
								textPositionEndModel.column,
						),
						options: {
							isWholeLine: false,
							className: 'bg-amber-400',
						},
					}]);
			
					window.rdfkitEditor.revealPositionNearTop(textPositionStartModel, monaco.editor.ScrollType.Smooth);
				});
			});
		</script>
	</body>
</html>`))

type parseRow struct {
	Subject      string
	SubjectRange string

	Predicate      string
	PredicateRange string

	Object      string
	ObjectRange string

	GraphName      string
	GraphNameRange string
}

func New() *cobra.Command {
	fIn := &cmdflags.EncodingInput{
		Path: "-",
		Type: "",
	}

	cmd := &cobra.Command{
		Use: "inspectdecoder",
		RunE: func(cmd *cobra.Command, args []string) error {
			var sourceBuffer = &bytes.Buffer{}

			bfIn, err := fIn.OpenTee(sourceBuffer)
			if err != nil {
				return fmt.Errorf("input: %v", err)
			}

			defer bfIn.Close()

			var prl []parseRow

			blankNodeIdentifierFactory := blanknodeutil.NewStringerInt64()

			b := &bytes.Buffer{}

			decoderTextOffsets, _ := bfIn.Decoder.(encoding.StatementTextOffsetsProvider)

			for bfIn.Decoder.Next() {
				quad := bfIn.Decoder.Quad()

				pr := parseRow{}

				switch s := quad.Triple.Subject.(type) {
				case rdf.BlankNode:
					pr.Subject = fmt.Sprintf("_:%s", blankNodeIdentifierFactory.GetBlankNodeIdentifier(s))
				case rdf.IRI:
					b.Reset()

					nquads.WriteIRI(b, s, false)

					pr.Subject = b.String()
				default:
					return fmt.Errorf("subject: invalid type: %T", s)
				}

				switch p := quad.Triple.Predicate.(type) {
				case rdf.IRI:
					b.Reset()

					nquads.WriteIRI(b, p, false)

					pr.Predicate = b.String()
				default:
					return fmt.Errorf("predicate: invalid type: %T", p)
				}

				switch o := quad.Triple.Object.(type) {
				case rdf.BlankNode:
					pr.Object = fmt.Sprintf("_:%s", blankNodeIdentifierFactory.GetBlankNodeIdentifier(o))
				case rdf.IRI:
					b.Reset()

					nquads.WriteIRI(b, o, false)

					pr.Object = b.String()
				case rdf.Literal:
					b.Reset()

					nquads.WriteLiteral(b, o, false)

					pr.Object = b.String()
				default:
					return fmt.Errorf("object: invalid type: %T", o)
				}

				if quad.GraphName != nil {
					switch g := quad.GraphName.(type) {
					case rdf.BlankNode:
						pr.GraphName = fmt.Sprintf("_:%s", blankNodeIdentifierFactory.GetBlankNodeIdentifier(g))
					case rdf.IRI:
						b.Reset()

						nquads.WriteIRI(b, g, false)

						pr.GraphName = b.String()
					default:
						return fmt.Errorf("graphName: invalid type: %T", g)
					}
				}

				if decoderTextOffsets != nil {
					if sourceRanges := decoderTextOffsets.StatementTextOffsets(); len(sourceRanges) > 0 {
						if v, ok := sourceRanges[encoding.SubjectStatementOffsets]; ok {
							pr.SubjectRange = v.OffsetRangeString()
						}

						if v, ok := sourceRanges[encoding.PredicateStatementOffsets]; ok {
							pr.PredicateRange = v.OffsetRangeString()
						}

						if v, ok := sourceRanges[encoding.ObjectStatementOffsets]; ok {
							pr.ObjectRange = v.OffsetRangeString()
						}

						if v, ok := sourceRanges[encoding.GraphNameStatementOffsets]; ok {
							pr.GraphNameRange = v.OffsetRangeString()
						}
					}
				}

				prl = append(prl, pr)
			}

			if err := bfIn.Decoder.Err(); err != nil {
				return fmt.Errorf("read: %s: %v", bfIn.Format, err)
			}

			tmplArgs := map[string]any{
				"Source":    sourceBuffer.String(),
				"Language":  "text",
				"ParseRows": prl,
			}

			switch fIn.Type {
			case "html":
				tmplArgs["Language"] = "html"
			case "jsonld", "rdfjson":
				tmplArgs["Language"] = "json"
			case "rdfxml":
				tmplArgs["Language"] = "xml"
			}

			err = tmpl.Execute(os.Stdout, tmplArgs)
			if err != nil {
				return fmt.Errorf("template: %v", err)
			}

			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&fIn.Path, "in", "i", fIn.Path, "")
	f.StringVar(&fIn.Type, "in-type", fIn.Type, "")

	return cmd
}
