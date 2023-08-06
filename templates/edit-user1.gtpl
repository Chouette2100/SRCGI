<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
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
        <td><button type="button"
                onclick="location.href='graph-total?eventid={{.Eventid}}&maxpoint={{.Maxpoint}}&gscale={{.Gscale}}'">獲得ポイントグラフ</button>
        </td>
        <td></td>
    </tr>
</table>
<br><br>
<p>(DB登録済)イベント参加ルーム一覧（確認・編集）　　<span style="color:red;">初めて使うときは表の後にある説明をよく読んでください！</span></p>
<table style="text-align: center">
<tr><td style="width:2em"></td><td><a href="https://www.showroom-live.com/event/{{.Eventid}}">{{.Eventname}}</a></td></tr>
<tr><td style="width:2em"></td><td>{{.Period}}</td></tr>
</table>
