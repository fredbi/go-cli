package cli

import (
	"net"
	"testing"
	"time"

	"github.com/fredbi/gflag"
	"github.com/fredbi/gflag/extensions"
	"github.com/fredbi/go-cli/cli/injectable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func testRootCommand(t *testing.T, asserter func(*testing.T, *viper.Viper)) *Command {
	// builds a root command with all sorts of flags with viper bindings
	return NewCommand(
		&cobra.Command{
			Use:   "test",
			Short: "explores cobra command with flags and config bindings",
			RunE: func(c *cobra.Command, _ []string) error {
				cfg := injectable.ConfigFromContext(c.Context(), viper.New)
				// assertions on flags bound to config
				asserter(t, cfg)

				return nil
			},
		},
		// single value flags
		WithFlag("bool-flag", false, "A bool flag", BindFlagToConfig("flags.bool")),
		WithFlag("true-bool-flag", true, "A bool flag defaulting to true", BindFlagToConfig("flags.true-bool")),
		WithFlag("int-flag", int(1), "An integer flag", BindFlagToConfig("flags.int")),
		WithFlag("int64-flag", int64(2), "An int64 flag", BindFlagToConfig("flags.int64")),
		WithFlag("float32-flag", float32(1.04), "A float32 flag", BindFlagToConfig("flags.float32")),
		WithFlag("float64-flag", 1.05, "A float64 flag", BindFlagToConfig("flags.float64")),
		WithFlag("duration-flag", time.Second, "A duration flag", BindFlagToConfig("flags.duration")),
		// TODO: uint, uint64, IP, IPMask, IPNet, []byte
		// extensions flags
		WithFlag("byte-size-flag", extensions.ByteSizeValue(1024*1024), "A byte-size extension flag", BindFlagToConfig("flags.byte-size")),
		// TODO: count
		// slice flags
		WithSliceFlag("string-slice-flag", []string{"a", "b"}, "A string-slice flag", BindFlagToConfig("flags.string-slice")),
		WithSliceFlag("ip-slice-flag", []net.IP{net.ParseIP("8.8.8.8"), net.ParseIP("127.0.0.1")}, "An IP-slice flag", BindFlagToConfig("flags.ip-slice")),
		// TODO: bools, durations, ints, ...
		// persistent flags
		WithFlag("persistent-bool-flag", false, "A persistent bool flag", BindFlagToConfig("flags.persistent.bool"), FlagIsPersistent()),
		WithFlag("persistent-int-flag", 2, "A persistent int flag", BindFlagToConfig("flags.persistent.int"), FlagIsPersistent()),
		WithFlag("persistent-float64-flag", 2.05, "A persistent float64 flag", BindFlagToConfig("flags.persistent.float64"), FlagIsPersistent()),
		WithFlag("persistent-duration-flag", 2*time.Second, "A persistent duration flag", BindFlagToConfig("flags.persistent.duration"), FlagIsPersistent()),
		/* TODO
		WithPersistentFlagFunc(func(flags *pflag.FlagSet) string {
			const userFlag = "extension-count"
			flags.StringVar(&globalFlags.User, userFlag, globalFlags.Defaults().User, "Originating user")
			return userFlag
		},
			FlagIsRequired(), BindFlagToConfig(keyUser),
		),
		*/
		// apply config to the command tree
		WithConfig(Config()),
	)
}
func TestCommandWithFlags(t *testing.T) {
	t.Run("should parse flags of all supported types", func(t *testing.T) {
		asserter := func(t *testing.T, cfg *viper.Viper) {
			require.EqualValues(t, true, cfg.GetBool("flags.bool"))
			require.EqualValues(t, false, cfg.GetBool("flags.true-bool"))
			require.EqualValues(t, 4, cfg.GetInt("flags.int"))
			require.EqualValues(t, int64(8), cfg.GetInt64("flags.int64"))
			require.EqualValues(t, 5.15, cfg.GetFloat64("flags.float32"))
			require.EqualValues(t, 2.1, cfg.GetFloat64("flags.float64"))
			require.EqualValues(t, 4*time.Minute, cfg.GetDuration("flags.duration"))

			t.Run("should retrieve extension flag from viper as a foreign type", func(t *testing.T) {
				sizeConfig := cfg.Get("flags.byte-size")
				require.NotNil(t, sizeConfig)

				// viper retrieves foreign types as a string
				sizeAsString, ok := sizeConfig.(string)
				require.Truef(t, ok, "expected a string representation of ByteSizeValue, got %T", sizeConfig)

				// NOTE: we should extend viper to retrieve the true value.
				// viper.Get() from flag values is currently hardcoding the retrieval of only a few native types.
				size := extensions.NewByteSizeValue(new(uint64), 0)
				require.NoError(t, size.UnmarshalFlag(sizeAsString))
				require.EqualValues(t, uint64(5*1000*1000), uint64(*size))

				marshalledValue, _ := size.MarshalFlag()
				require.Equal(t, sizeAsString, marshalledValue)
			})

			t.Run("should retrieve string slice flag from viper", func(t *testing.T) {
				require.EqualValues(t, []string{"x", "y", "u", "v", "w"}, cfg.GetStringSlice("flags.string-slice"))
			})

			t.Run("should retrieve slice not supported by viper", func(t *testing.T) {
				ipConfig := cfg.Get("flags.ip-slice")
				require.NotNil(t, ipConfig)
				ipsAsString, ok := ipConfig.(string)
				require.True(t, ok)

				// NOTE: we should extend viper to retrieve the true value.
				ips := make([]net.IP, 0, 2)
				v := gflag.NewFlagSliceValue(&ips, nil)
				require.NoError(t, v.UnmarshalFlag(ipsAsString))
				require.EqualValues(t, []net.IP{net.ParseIP("1.1.1.1"), net.ParseIP("255.255.255.254")}, ips)

				marshalledValue, _ := v.MarshalFlag()
				require.Equal(t, ipsAsString, marshalledValue)
			})

			t.Run("should retrieve persistent flags from config", func(t *testing.T) {
				require.EqualValues(t, true, cfg.GetBool("flags.persistent.bool"))
				require.EqualValues(t, 15, cfg.GetInt("flags.persistent.int"))
				require.EqualValues(t, 2.1, cfg.GetFloat64("flags.persistent.float64"))
				require.EqualValues(t, 5*time.Hour, cfg.GetDuration("flags.persistent.duration"))
			})
		}

		require.NoError(t,
			testRootCommand(t, asserter).ExecuteWithArgs(
				"--persistent-bool-flag",
				"--persistent-int-flag", "15",
				"--persistent-float64-flag", "2.10",
				"--persistent-duration-flag", "5h",
				"--bool-flag",
				"--true-bool-flag=false",
				"--int-flag", "4",
				"--int64-flag", "8",
				"--float32-flag", "5.15",
				"--float64-flag", "2.10",
				"--duration-flag", "4m",
				"--byte-size-flag", "5MB",
				// slice flags
				"--string-slice-flag", "x,y",
				"--string-slice-flag", "u,v,w", // append semantics
				"--ip-slice-flag", "1.1.1.1, 255.255.255.254",
			),
		)
	})
	// TODO: with subcommand
	// TODO: check defaults
	// TODO: parsing errors
	// TODO: support for map flags from pflag
}

// TODO: ExecuteContext()
// TODO: use Config()
