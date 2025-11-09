<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<head>
{{if .TurnstileSiteKey}}
<script src="https://challenges.cloudflare.com/turnstile/v0/api.js" async defer></script>
{{end}}
</head>
<body>
<table>
  <tr>
<td><button type="button" onclick="location.href='top'">トップ</button>　</td>
<td><button type="button" onclick="location.href='currentevents'">開催中イベント一覧</button></td>
<td><button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧</button></td>
<td><button type="button" onclick="location.href='closedevents'">終了イベント一覧</button></td>
  </tr>
  <tr>
<td><button type="button" onclick="location.href='eventtop?eventid={{.Eventid}}'">イベントトップ</button></td>
<td><button type="button" onclick="location.href='list-last?eventid={{.Eventid}}'">直近の獲得ポイント</button></td>
<td><button type="button" onclick="location.href='graph-total?eventid={{.Eventid}}'">獲得ポイントグラフ</button></td>
<td></td>
  </tr>
</table>
<br><br>
{{if and .TurnstileSiteKey (not .Filename)}}
<!-- Turnstileチャレンジ表示 -->
<div style="border: 2px solid #4A90E2; padding: 20px; border-radius: 5px; max-width: 600px; background-color: #f9f9f9;">
<h3>セキュリティチェック</h3>
{{if .TurnstileError}}
<p style="color: red; font-weight: bold;">{{.TurnstileError}}</p>
{{end}}
<p>{{ .Event_name }}（{{.Roomid}}）の貢献ランキングを表示するには、セキュリティチェックを完了してください。</p>
<form method="POST" action="contributors">
<input type="hidden" name="ieventid" value="{{.Ieventid}}">
<input type="hidden" name="roomid" value="{{.Roomid}}">
<div class="cf-turnstile" data-sitekey="{{.TurnstileSiteKey}}" data-theme="light"></div>
<br>
<button type="submit" style="padding: 10px 20px; background-color: #4A90E2; color: white; border: none; border-radius: 5px; cursor: pointer; font-size: 16px;">確認して続行</button>
</form>
</div>
<br><br>
{{else}}
<!-- データ表示 -->
{{if .Filename}}
<a href="{{ .Filename }}-1.csv" download="{{ .Filename }}-1.csv">1. {{ .Event_name }}（{{.Roomid}}）貢献ランキングのダウンロード(CSV)</a>
<br>
<a href="{{ .Filename }}-2.csv" download="{{ .Filename }}-2.csv">2. {{ .Event_name }}（{{.Roomid}}）貢献ランキングのダウンロード(UTF-8 BOMつきCSV)</a>
<br>
<br>
※ 1.はふつうのCSVファイルです。LibreOffice Calc、Googleスプレッドシートなどで開けます。<br>
※ 2.はUTF-8 BOM(注)付きのCSVファイルです。Microsoft Excelのために用意しましたが、LibreOffice Calc、Googleスプレッドシートなどでも開けます。<br>
※ どちらも文字化けする場合はお手数ですが下記の内容をコピーしてアプリに貼り付けてください。<br>
<br>
注 UTF-8 BOM ... ファイル先頭に0xEF, 0xBB, 0xBFの3バイトを付加しています。これが問題になる場合は1.を使ってください。<br>
<br><br>
※ 要望があればjsonあるいはyamlでのデータ出力も可能です。
<br><br>
        <table border="1" style="border-collapse: collapse; border: 1px solid black">
            <tr bgcolor="gainsboro" style="text-align: center">
                <td>順位</td>
                <td style="border-right: none;">視聴者名</td>
                <td style="border-left: none;">ユーザID</td>
                <td>貢献ポイント</td>
            </tr>
            {{range .Result}}
            <tr style="text-align: right;">
            <td>{{.Irank}}</td>
            <td style="text-align: left; border-right: none;">{{.Name}}</td>
            <td style="border-left: none;">{{.Viewerid}}</td>
            <td>{{.Point}}</td>
            </tr>
            {{end}}
        </table>
{{end}}
{{end}}
</body>
</html>
