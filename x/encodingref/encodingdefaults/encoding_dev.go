package encodingdefaults

import (
	"html/template"

	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/encoding/encodingtest"
	"github.com/dpb587/rdfkit-go/encoding/turtle"
	"github.com/dpb587/rdfkit-go/rdf/blanknodeutil"
	"github.com/dpb587/rdfkit-go/x/encodingref"
)

const ctiDevHtmlInspector encoding.ContentTypeIdentifier = "internal.dev.html-inspector"

type encodingDev struct{}

var _ encodingref.RegistryEncoding = &encodingDev{}

func (e encodingDev) NewDecoder(cti encoding.ContentTypeIdentifier, rr encodingref.ResourceReader, opts encodingref.DecoderOptions) (*encodingref.DecoderHandle, error) {
	return nil, encodingref.ErrEncodingNotSupported
}

func (e encodingDev) NewEncoder(cti encoding.ContentTypeIdentifier, ww encodingref.ResourceWriter, opts encodingref.EncoderOptions) (*encodingref.EncoderHandle, error) {
	switch cti {
	case encodingtest.DiscardEncoderContentTypeIdentifier:
		return &encodingref.EncoderHandle{
			Writer:  ww,
			Encoder: encodingtest.DiscardEncoder,
		}, nil
	case encodingtest.TriplesEncoderContentTypeIdentifier:
		return &encodingref.EncoderHandle{
			Writer:  ww,
			Encoder: encodingtest.NewTriplesEncoder(wrapWriter(ww, opts), encodingtest.TriplesEncoderOptions{}),
		}, nil
	case encodingtest.QuadsEncoderContentTypeIdentifier:
		return &encodingref.EncoderHandle{
			Writer:  ww,
			Encoder: encodingtest.NewQuadsEncoder(wrapWriter(ww, opts), encodingtest.QuadsEncoderOptions{}),
		}, nil
	case ctiDevHtmlInspector:
		return &encodingref.EncoderHandle{
			Encoder: encodingtest.NewBufferedQuadsTemplate(encodingtest.BufferedQuadsTemplateOptions{
				Output: wrapWriter(ww, opts),
				OutputTemplate: template.Must(template.New("inspect").Parse(`<!DOCTYPE html>
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
				{{- range .Quads -}}
					<div>{{/*
						*/}}<span{{ with .TextOffsets.Subject }} title="{{ .OffsetRangeString }}"{{ end }}>{{ .Encoded.Subject }}</span> {{/*
						*/}}<span{{ with .TextOffsets.Predicate }} title="{{ .OffsetRangeString }}"{{ end }}>{{ .Encoded.Predicate }}</span> {{/*
						*/}}<span{{ with .TextOffsets.Object }} title="{{ .OffsetRangeString }}"{{ end }}>{{ .Encoded.Object }}</span> {{/*
						*/}}{{ if .Encoded.GraphName }}<span{{ with .TextOffsets.GraphName }} title="{{ .OffsetRangeString }}"{{ end }}>{{ .Encoded.GraphName }}</span> {{ end }}{{/*
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
						language: {{ with .SourceType }}{{ . }}{{ else }}"text"{{ end }},
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
</html>`)),
				Formatter: turtle.NewTermFormatter(turtle.TermFormatterOptions{
					BlankNodeStringer: blanknodeutil.NewStringerInt64(),
				}),
			}),
		}, nil
	}

	return nil, encodingref.ErrEncodingNotSupported
}
