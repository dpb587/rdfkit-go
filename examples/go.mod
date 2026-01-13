module examples

go 1.25.5

require github.com/dpb587/rdfkit-go v0.0.0-00010101000000-000000000000

require (
	github.com/apparentlymart/go-textseg/v16 v16.0.0 // indirect
	github.com/dpb587/cursorio-go v0.0.0-20250717044249-e1d8c928b30d // indirect
	github.com/dpb587/inspecthtml-go v0.0.0-20250906134739-0d404c86637b // indirect
	github.com/dpb587/inspectjson-go v0.0.0-20251203142639-90f1442149fb // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/tomnomnom/linkheader v0.0.0-20250811210735-e5fe3b51442e // indirect
	golang.org/x/net v0.47.0 // indirect
)

replace github.com/dpb587/rdfkit-go => ..

replace github.com/dpb587/rdfkit-go/cmd/rdfkit => ../cmd/rdfkit
