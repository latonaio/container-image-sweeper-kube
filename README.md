# container-image-sweeper-kube

container-image-sweeper-kube は、不要になった Docker イメージを自動で削除するツールです。AION のマイクロサービスとして定期実行のほか、コマンドラインから一度きりの実行も可能です。


## 実行環境

* Linux
* Docker
* Go 1.17 以上


## ローカル上での実行

```sh
make run
```

## AION 上での実行

* Docker イメージをビルドします。

```sh
make docker-build
```

* 以下の定義を AION の `services.yml` に貼り付けてください。

```yaml
  container-image-sweeper-kube:
    startup: yes
    always: yes
    scale: 1
    volumeMountPathList:
      - /var/run/docker.sock:/var/run/docker.sock
    # 環境変数での設定のカスタマイズを行う場合、以下の例のように記載
    # env:
    #   INTERVAL: 3000
```


## 設定のカスタマイズ

環境変数により設定のカスタマイズを行うことができます。必要に応じて行ってください。

| 環境変数名          | 型   | デフォルト値 | 概要                                                                                |
| ------------------- | ---- | ------------ | ----------------------------------------------------------------------------------- |
| `DAEMONIZE`         | bool | `true`       | デーモンとして起動するかどうか設定します。                                          |
| `RETAIN`            | int  | `3`          | イメージのリポジトリ名ごとに残すイメージ数を設定します。                            |
| `KEEP_LATEST`       | bool | `true`       | `RETAIN` の数に関わらず、タグ名が `latest` のイメージを常に残すかどうか設定します。 |
| `PRUNE_IMAGES`      | bool | `true`       | dangling なイメージを削除するか設定します。                                         |
| `PRUNE_BUILD_CACHE` | bool | `false`      | dangling なビルドキャッシュを削除するか設定します。                                 |
| `INTERVAL`          | int  | `3000`       | デーモン実行時の実行間隔を秒単位で設定します。                                      |
