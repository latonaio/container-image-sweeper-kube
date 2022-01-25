package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/latonaio/golang-logging-library/logger"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
)

var log = logger.NewLogger()

type App struct {
	cmdArgs *CommandArgs
}

func Command() *cobra.Command {
	cmdArgs := NewCommandArgs()
	cmd := &cobra.Command{
		Use: os.Args[0],
		Run: func(cmd *cobra.Command, args []string) {
			app := &App{cmdArgs: cmdArgs}
			if err := app.main(); err != nil {
				log.Fatal(err)
			}
		},
	}
	cmdArgs.Set(cmd.Flags())

	return cmd
}

func (a *App) main() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if !a.cmdArgs.Daemonize {
		// デーモン化しない場合、そのまま実行で良い
		return a.mainImpl(ctx)
	}

	// 常時起動し、IntervalSec ごとに実行する
	a.daemonizedMain(ctx)

	// (到達しない)
	return nil
}

// 常時起動
func (a *App) daemonizedMain(ctx context.Context) {
	// Ctrl + C の検出
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	// 一定時間ごとに ticker.C チャンネルにシグナルを送信
	ticker := time.NewTicker(time.Second * time.Duration(a.cmdArgs.IntervalSec))
	defer ticker.Stop()

LOOP:
	for {
		// 初回は ticker を待たずに実行、その後は <-ticker.C を待つ
		if err := a.mainImpl(ctx); err != nil {
			log.Warn("errored while executing: %+v", err)
		}

		select {
		case <-signalCh:
			// Ctrl + C の検出時
			log.Info("signal received, shutting down...")
			break LOOP

		case <-ticker.C:
			// 一定時間経過後
			continue
		}
	}
}

func (a *App) mainImpl(ctx context.Context) error {
	// 開始と終了のログ
	log.Debug("started")
	defer func() { log.Debug("exited") }()

	m, err := NewDockerManager()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	return a.process(ctx, m)
}

func (a *App) process(ctx context.Context, m *DockerManager) error {
	// エラーオブジェクト
	// 無視可能なエラーが発生した場合はそのエラーを記録して続行し、最後にまとめて返す
	var errs error

	// dangling なイメージとビルドキャッシュの削除
	if a.cmdArgs.PruneImages {
		errs = multierr.Append(errs, m.PruneImages(ctx))
	}

	if a.cmdArgs.PruneBuildCache {
		errs = multierr.Append(errs, m.PruneBuildCache(ctx))
	}

	// 削除対象のイメージを取得
	imgMap, err := m.GetImages(ctx)
	if err != nil {
		return fmt.Errorf("failed to get image list: %w", err)
	}

	imgIDs := a.getImagesToBeDeleted(imgMap)
	if len(imgIDs) != 0 {
		log.Info("deleting images: %v", imgIDs)
	}

	// イメージ削除を実行
	for _, imgID := range imgIDs {
		errs = multierr.Append(err, m.DeleteImage(ctx, imgID))
	}

	return errs
}

func (a *App) getImagesToBeDeleted(imgMap map[string][]*TagInfo) []string {
	imgTags := make([]string, 0)
	for repoName, tagInfos := range imgMap {
		// 残すイメージ数の残りカウント
		remainCount := a.cmdArgs.RetainCount
		// リポジトリごとのタグ情報: ビルド日時が新しい順にソート済
		for _, tagInfo := range tagInfos {
			// latest を削除対象としない場合
			if a.cmdArgs.KeepLatest && tagInfo.Tag == "latest" {
				continue
			}

			// 残す個数分は削除対象としない
			if remainCount > 0 {
				remainCount--
				continue
			}

			// 削除対象に追加
			imgTags = append(imgTags, fmt.Sprintf("%v:%v", repoName, tagInfo.Tag))
		}
	}

	return imgTags
}
