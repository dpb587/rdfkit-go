package encodingref_test

import (
	"testing"

	"github.com/dpb587/rdfkit-go/rdf"
	"github.com/dpb587/rdfkit-go/x/encodingref"
)

type htmlOption struct {
	CaptureTextOffsets *bool
	Rdfa               *htmlOptionRdfa
}

type htmlOptionRdfa struct {
	DefaultVocab *rdf.IRI
}

func TestUnmarshalFlags(t *testing.T) {
	tests := []struct {
		name    string
		flags   []string
		want    func(*testing.T, *htmlOption)
		wantErr bool
	}{
		{
			name:  "bool flag without value",
			flags: []string{"captureTextOffsets"},
			want: func(t *testing.T, got *htmlOption) {
				if got.CaptureTextOffsets == nil {
					t.Error("CaptureTextOffsets should not be nil")
					return
				}
				if !*got.CaptureTextOffsets {
					t.Error("CaptureTextOffsets should be true")
				}
			},
		},
		{
			name:  "nested field with value",
			flags: []string{`rdfa.defaultVocab=https://github.com/`},
			want: func(t *testing.T, got *htmlOption) {
				if got.Rdfa == nil {
					t.Error("Rdfa should not be nil")
					return
				}
				if got.Rdfa.DefaultVocab == nil {
					t.Error("DefaultVocab should not be nil")
					return
				}
				expected := "https://github.com/"
				if string(*got.Rdfa.DefaultVocab) != expected {
					t.Errorf("DefaultVocab = %q, want %q", *got.Rdfa.DefaultVocab, expected)
				}
			},
		},
		{
			name:  "multiple flags",
			flags: []string{"captureTextOffsets", `rdfa.defaultVocab=https://github.com/`},
			want: func(t *testing.T, got *htmlOption) {
				if got.CaptureTextOffsets == nil {
					t.Error("CaptureTextOffsets should not be nil")
					return
				}
				if !*got.CaptureTextOffsets {
					t.Error("CaptureTextOffsets should be true")
				}
				if got.Rdfa == nil {
					t.Error("Rdfa should not be nil")
					return
				}
				if got.Rdfa.DefaultVocab == nil {
					t.Error("DefaultVocab should not be nil")
					return
				}
				expected := "https://github.com/"
				if string(*got.Rdfa.DefaultVocab) != expected {
					t.Errorf("DefaultVocab = %q, want %q", *got.Rdfa.DefaultVocab, expected)
				}
			},
		},
		{
			name:    "nil pointer",
			flags:   []string{"captureTextOffsets"},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opt *htmlOption
			if tt.want != nil {
				opt = &htmlOption{}
			}

			err := encodingref.UnmarshalFlags(opt, tt.flags)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				tt.want(t, opt)
			}
		})
	}
}

func TestUnmarshalFlags_NonPointer(t *testing.T) {
	opt := htmlOption{}
	err := encodingref.UnmarshalFlags(opt, []string{"captureTextOffsets"})
	if err == nil {
		t.Error("expected error for non-pointer argument")
	}
}

func TestUnmarshalFlags_TypeParsing(t *testing.T) {
	type testStruct struct {
		StringVal  string
		IntVal     int
		Int64Val   int64
		UintVal    uint
		Uint64Val  uint64
		BoolVal    bool
		PtrIntVal  *int
		PtrBoolVal *bool
	}

	tests := []struct {
		name    string
		flags   []string
		want    func(*testing.T, *testStruct)
		wantErr bool
	}{
		{
			name:  "string value",
			flags: []string{"stringVal=hello world"},
			want: func(t *testing.T, got *testStruct) {
				if got.StringVal != "hello world" {
					t.Errorf("StringVal = %q, want %q", got.StringVal, "hello world")
				}
			},
		},
		{
			name:  "int value",
			flags: []string{"intVal=42"},
			want: func(t *testing.T, got *testStruct) {
				if got.IntVal != 42 {
					t.Errorf("IntVal = %d, want %d", got.IntVal, 42)
				}
			},
		},
		{
			name:  "int64 value",
			flags: []string{"int64Val=9223372036854775807"},
			want: func(t *testing.T, got *testStruct) {
				if got.Int64Val != 9223372036854775807 {
					t.Errorf("Int64Val = %d, want %d", got.Int64Val, 9223372036854775807)
				}
			},
		},
		{
			name:  "uint value",
			flags: []string{"uintVal=100"},
			want: func(t *testing.T, got *testStruct) {
				if got.UintVal != 100 {
					t.Errorf("UintVal = %d, want %d", got.UintVal, 100)
				}
			},
		},
		{
			name:  "bool value true",
			flags: []string{"boolVal=true"},
			want: func(t *testing.T, got *testStruct) {
				if !got.BoolVal {
					t.Error("BoolVal should be true")
				}
			},
		},
		{
			name:  "bool value false",
			flags: []string{"boolVal=false"},
			want: func(t *testing.T, got *testStruct) {
				if got.BoolVal {
					t.Error("BoolVal should be false")
				}
			},
		},
		{
			name:  "pointer int value",
			flags: []string{"ptrIntVal=123"},
			want: func(t *testing.T, got *testStruct) {
				if got.PtrIntVal == nil {
					t.Error("PtrIntVal should not be nil")
					return
				}
				if *got.PtrIntVal != 123 {
					t.Errorf("PtrIntVal = %d, want %d", *got.PtrIntVal, 123)
				}
			},
		},
		{
			name:  "pointer bool no value",
			flags: []string{"ptrBoolVal"},
			want: func(t *testing.T, got *testStruct) {
				if got.PtrBoolVal == nil {
					t.Error("PtrBoolVal should not be nil")
					return
				}
				if !*got.PtrBoolVal {
					t.Error("PtrBoolVal should be true")
				}
			},
		},
		{
			name:    "invalid int",
			flags:   []string{"intVal=notanumber"},
			wantErr: true,
		},
		{
			name:    "invalid bool",
			flags:   []string{"boolVal=notabool"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := &testStruct{}
			err := encodingref.UnmarshalFlags(opt, tt.flags)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				tt.want(t, opt)
			}
		})
	}
}
