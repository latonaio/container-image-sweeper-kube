package args

import (
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type CommandArgs struct {
	Daemonize       bool
	RetainCount     int
	KeepLatest      bool
	PruneImages     bool
	PruneBuildCache bool
	IntervalSec     int
}

const (
	defaultDaemonize   = true
	defaultRetainCount = 3
	defaultKeepLatest  = true
	defaultPruneImages = true
	defaultIntervalSec = 3000
)

func NewCommandArgs() *CommandArgs {
	return &CommandArgs{
		Daemonize:   defaultDaemonize,
		RetainCount: defaultRetainCount,
		KeepLatest:  defaultKeepLatest,
		PruneImages: defaultPruneImages,
		IntervalSec: defaultIntervalSec,
	}
}

func (a *CommandArgs) Set(flags *pflag.FlagSet) {
	flags.BoolP("daemonize", "d", a.Daemonize, "run as daemon mode")
	flags.IntP("retain", "r", a.RetainCount, "number of image tags to keep, per repository")
	flags.Bool("keep-latest", a.KeepLatest, "always retain image with \"latest\" tag")
	flags.Bool("prune-images", a.PruneImages, "prune dangling images")
	flags.Bool("prune-build-cache", a.PruneBuildCache, "prune dangling build cache")
	flags.IntP("interval", "i", a.IntervalSec, "(daemon mode only) deletion interval, in seconds")
	flags.Parse(os.Args[1:])

	// 環境変数と引数を利用する: 優先順位は 引数 → 環境変数 → デフォルト値
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	v.BindPFlags(flags)

	a.Daemonize = v.GetBool("daemonize")
	a.RetainCount = v.GetInt("retain")
	a.KeepLatest = v.GetBool("keep-latest")
	a.PruneImages = v.GetBool("prune-images")
	a.PruneBuildCache = v.GetBool("prune-build-cache")
	a.IntervalSec = v.GetInt("interval")
}
