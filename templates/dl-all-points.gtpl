<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<body>
{{ $f := .Filename }}
  {{ with .Eventinf }}
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
<a href="{{ $f }}-1.csv" download="{{ $f }}-1.csv">1. 『{{ .Event_name }}』（{{.Eventid}} / {{.Ieventid}}）獲得ポイントデータのダウンロード(CSV)</a>
<br><br>
<a href="{{ $f }}-2.csv" download="{{ $f }}-2.csv">2. 『{{ .Event_name }}』（{{.Eventid}} / {{.Ieventid}}）獲得ポイントデータのダウンロード(UTF-8 BOMつきCSV)</a>
  {{ end }}
<br><br>
<br>
※ 1.はふつうのCSVファイルです。LibreOffice Calc、Googleスプレッドシートなどで開けます。<br>
※ 2.はUTF-8 BOM(注)付きのCSVファイルです。Microsoft Excelのために用意しましたが、LibreOffice Calc、Googleスプレッドシートなどでも開けます。<br>
<br>
注 UTF-8 BOM ... ファイル先頭に0xEF, 0xBB, 0xBFの3バイトを付加しています。これが問題になる場合は1.を使ってください。<br>
<br>
<br>
・獲得ポイントに変化がないところのデータはありませせん（いろんな事情でデータが出力されていることもあります）
<br>
・欠測の場合もデータがないわけですが、データがない部分の両側が同じ値かどうかでどちらか判断できます 。
<br>
・イベント開始時（ふつう開始日の18時00分）のデータ（=0ポイント）は除いてあります。
<br>
・イベント参加を途中でやめたルームも最後の獲得ポイントが大きいときは表示されることがあります。このことの可否については検討中です。
<br>
・要望があればJSONあるいやyamlなどでのデータ出力も可能です。
</body>
</html>
