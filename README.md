# Desktop TOTP Application

このアプリケーションは、TOTP（Time-based One-Time Password）コードを生成し、システムトレイに表示するGo言語で書かれたプログラムです。TOTPコードをクリップボードにコピーすることもできます。

## 特徴
- システムトレイにTOTPコードを表示
- TOTPコードをクリップボードにコピー
- PC起動時に自動的に実行（任意設定）

## 必要条件
- Go 1.23以降

## インストール

1. リポジトリをクローンします。

    ```sh
    git clone https://github.com/dqx0/desktop-totp.git
    cd desktop-totp
    ```

2. 依存関係をインストールします。

    ```sh
    go mod tidy
    ```

3. （任意）プログラムをビルドします。

    ```sh
    go build -ldflags "-H=windowsgui" -o C:\path\to\your\myapp.exe main.go
    ```

4. （任意）PC起動時に自動実行するために、スタートアップフォルダにショートカットを作成します。以下のPowerShellスクリプトを実行してください（Windowsの場合）。

    ```powershell
    $sourcePath = "C:\path\to\your\myapp.exe"
    $shortcutPath = "$env:APPDATA\Microsoft\Windows\Start Menu\Programs\Startup\myapp.lnk"

    $WScriptShell = New-Object -ComObject WScript.Shell
    $shortcut = $WScriptShell.CreateShortcut($shortcutPath)
    $shortcut.TargetPath = $sourcePath
    $shortcut.Save()
    ```

## 環境変数の設定

以下の環境変数を設定してください。

- `TOTP_SECRET`: TOTPのシークレットキー
- `TOTP_ICON_PATH`: アイコンファイル（.ico）のパス

## 使用方法

### ビルドした実行ファイルを使用する場合

1. プログラムを実行します。

    ```sh
    C:\path\to\your\myapp.exe
    ```

2. システムトレイにアイコンが表示されます。
3. 「Copy」をクリックすると、TOTPコードがクリップボードにコピーされます。

### `go run`コマンドを使用する場合

1. プログラムを実行します。

    ```sh
    go run main.go
    ```

2. システムトレイにアイコンが表示されます。
3. 「Copy」をクリックすると、TOTPコードがクリップボードにコピーされます。

### 自動実行の設定をした場合

1. 初回は再起動が必要です。
2. システムトレイにアイコンが表示されます。
3. 「Copy」をクリックすると、TOTPコードがクリップボードにコピーされます。

## ライセンス

このプロジェクトはApache 2.0ライセンスの下で公開されています。詳細については、LICENSEファイルを参照してください。

このプロジェクトは以下のライブラリを使用しています：

- `atotto/clipboard` (BSD 3-Clause License)
- `getlantern/systray` (Apache 2.0 License)
- `pquerna/otp` (Apache 2.0 License)