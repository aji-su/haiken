# haiken

Misskey のタイムラインから 575 を検出する bot です。

オリジナルである https://github.com/theoria24/FindHaiku4Mstdn を Misskey 向けに実装したものです。

## 仕様

- ホームタイムラインを監視し、 575 を見つけたらリプライします。
  - 公開範囲がパブリック、ホームの投稿のみが対象です。
  - 形態素解析の内容もリプライします。
- フォローされたらフォローバックします。
- フォロワーからのメンションに反応します。
  - 「俳句検出を停止してください」と送られたらフォロー解除します。
  - それ以外の場合は形態素解析の内容をリプライします。

## development

```bash
# Copy and edit .env file
cp .env.example .env

# Run
docker compose up -d
docker compose logs -f

# set platform:
docker buildx build --platform linux/amd64 -t haiken_app .
```

## testing

```bash
docker compose exec app go test -v ./...
```

## running

```bash
docker build -t aji-su/haiken:latest .
docker run --env-file .env -d aji-su/haiken:latest
```

## deploy

Deploy to Amazon ECS, see [aws.yml](.github/workflows/aws.yml)
