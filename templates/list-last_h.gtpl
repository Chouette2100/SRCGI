{{define "list-last_h"}}
<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<meta http-equiv="refresh" content="{{.SecondsToReload}}; URL=list-last?eventid={{.Eventid}}&userno={{.userno}}&detail={{.Detail}}">
{{ $detail := .Detail }}
<html>
<head>
    <style type="text/css">
        .bgct {
            background-color: paleturquoise;
        }
    </style>

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
      <td><button type="button" onclick="location.href='top?eventid={{.Eventid}}'">イベントトップ</button></td>
      <td><button type="button" onclick="location.href='list-last?eventid={{.Eventid}}'">直近の獲得ポイント</button></td>
      <td><button type="button" onclick="location.href='graph-total?eventid={{.Eventid}}&maxpoint={{.Maxpoint}}&gscale={{.Gscale}}'">獲得ポイントグラフ</button></td>
      <td></td>
        </tr>
      </table>
<p>直近の獲得ポイント一覧　　<span style="color:red;">初めて使うときは表の後にある注意事項を読んでください！</span>　このページはブックマーク可能です。</p>
<p style="padding-left:2em">
{{.UpdateTime}}
</p>
<p style="padding-left:4em">
{{.NextTime}}
<br>
{{ if ne .NextTime "イベントは終了しています。" }}
{{.ReloadTime}}
{{ else }}
<span style="color: red;">最終結果の反映はイベント終了日翌日の12時過ぎです。<span>
{{ end }}
</p>
<table>
<tr><td align="center"><a href="https://www.showroom-live.com/event/{{.Eventid}}">{{.EventName}}</a>（{{.Eventid}}）</td></tr>
<tr><td align="center">{{.Period}}</td></tr>
</table>
{{ if eq .Detail "1" }}
<button type="button" onclick="location.href='list-last?eventid={{.Eventid}}&userno={{.userno}}&detail=0'">ルーム詳細情報を表示しない</button>
{{ else }}
<button type="button" onclick="location.href='list-last?eventid={{.Eventid}}&userno={{.userno}}&detail=1'">ルーム詳細情報を表示する</button>
{{ end }}
{{end}}
