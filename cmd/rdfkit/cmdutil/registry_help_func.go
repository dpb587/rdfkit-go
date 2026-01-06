package cmdutil

import (
	"maps"
	"slices"
	"strings"

	"github.com/dpb587/kvstrings-go/kvstrings"
	"github.com/dpb587/rdfkit-go/encoding"
	"github.com/dpb587/rdfkit-go/rdfio/rdfiotypes"
	"github.com/spf13/cobra"
)

func RegistryHelpFunc(app *App) func(*cobra.Command, []string) {
	return func(c *cobra.Command, args []string) {
		usage := c.Long
		if usage == "" {
			usage = c.Short
		}

		usage = strings.TrimSpace(usage)

		if len(usage) > 0 {
			c.Println(usage)
			c.Println()
		}

		if c.Runnable() || c.HasSubCommands() {
			c.Print(c.UsageString())
		}

		c.Println()
		c.Println("Encodings:")

		ctiUniq := map[encoding.ContentTypeIdentifier]struct{}{}

		for cti := range app.Registry.EncoderManagers {
			if !strings.HasPrefix(string(cti), "internal.") {
				ctiUniq[cti] = struct{}{}
			}
		}

		for cti := range app.Registry.DecoderManagers {
			if !strings.HasPrefix(string(cti), "internal.") {
				ctiUniq[cti] = struct{}{}
			}
		}

		ctiList := slices.Collect(maps.Keys(ctiUniq))
		slices.SortFunc(ctiList, func(a, b encoding.ContentTypeIdentifier) int {
			return strings.Compare(string(a), string(b))
		})

		for _, cti := range ctiList {
			decoder := app.Registry.DecoderManagers[cti]
			encoder := app.Registry.EncoderManagers[cti]

			capabilities := []string{}

			if decoder != nil {
				capabilities = append(capabilities, "decode")
			}

			if encoder != nil {
				capabilities = append(capabilities, "encode")
			}

			c.Println()
			c.Println("  " + string(cti) + " (" + strings.Join(capabilities, ", ") + ")")
			c.Println()

			{
				var aliases []string

				for alias, ctiAlias := range app.Registry.Aliases {
					if cti == ctiAlias {
						aliases = append(aliases, alias)
					}
				}

				slices.SortFunc(aliases, strings.Compare)

				if len(aliases) > 0 {
					c.Println("    Aliases: " + strings.Join(aliases, ", "))
				}
			}

			{
				var fileExts []string

				for fileExt, ctiFileExt := range app.Registry.FileExts {
					if cti == ctiFileExt {
						fileExts = append(fileExts, fileExt)
					}
				}

				slices.SortFunc(fileExts, strings.Compare)

				if len(fileExts) > 0 {
					c.Println("    File Extensions: " + strings.Join(fileExts, ", "))
				}
			}

			{
				var mediaTypes []string

				for mediaType, ctiMediaType := range app.Registry.MediaTypes {
					if cti == ctiMediaType {
						mediaTypes = append(mediaTypes, mediaType)
					}
				}

				slices.SortFunc(mediaTypes, strings.Compare)

				if len(mediaTypes) > 0 {
					c.Println("    Media Types: " + strings.Join(mediaTypes, ", "))
				}
			}

			if decoder != nil {
				encodingHelpParams(c, "--in-param", decoder.NewDecoderParams())
			}

			if encoder != nil {
				encodingHelpParams(c, "--out-param", encoder.NewEncoderParams())
			}
		}
	}
}

func encodingHelpParams(c *cobra.Command, flagName string, d rdfiotypes.Params) {
	if d == nil {
		return
	}

	kvl := d.NewParamsCollection().GetDescriptors()
	if len(kvl) == 0 {
		return
	}

	slices.SortFunc(kvl, func(a, b kvstrings.KeyValueDescriptor[rdfiotypes.ParamMeta]) int {
		return strings.Compare(string(a.Key), string(b.Key))
	})

	for _, kv := range kvl {
		if kv.Meta.Hidden {
			continue
		}

		c.Println()
		c.Print("    " + flagName + " " + string(kv.Key))

		suffix := "=" + kv.Value.GetValueName()

		if kvv, ok := kv.Value.(kvstrings.ImpliedValueHandler); ok {
			if _, err := kvv.GetImpliedValue(); err == nil {
				suffix = "[" + suffix + "]"
			}
		}

		if kv.Value.GetImportValueBehavior() == kvstrings.ValueImportBehaviorPatch {
			suffix = suffix + "..."
		}

		c.Println(suffix)

		usage := &strings.Builder{}

		if len(kv.Meta.Usage) > 0 {
			usage.WriteString(kv.Meta.Usage)
		}

		if len(kv.Meta.ValueEnums) > 0 {
			usage.WriteString(" (enum ")
			usage.WriteString(strings.Join(kv.Meta.ValueEnums, ", "))
			usage.WriteString(")")
		}

		if usage.Len() > 0 {
			c.Println("      " + usage.String())
		}
	}
}
