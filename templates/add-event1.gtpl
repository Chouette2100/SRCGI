<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
{{/*
<meta http-equiv="refresh" content="{{.SecondsToReload}}; URL=list-last?eventid={{.Eventid}}"> */}}
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
            <td><button type="button" onclick="location.href='eventtop?eventid={{.Event_ID}}'">イベントトップ</button></td>
            <td><button type="button" onclick="location.href='list-last?eventid={{.Event_ID}}'">直近の獲得ポイント</button></td>
            <td><button type="button"
                    onclick="location.href='graph-total?eventid={{.Event_ID}}&maxpoint={{.Maxpoint}}&gscale={{.Gscale}}'">獲得ポイントグラフ</button>
            </td>
            <td></td>
        </tr>
    </table>
    <br>
    (DB登録済)イベント参加ルーム一覧
    <br>
    <br>
    <table>
        <tr>
            <td>イベント</td>
            <td><a href="https://www.showroom-live.com/event/{{.Event_ID}}">{{ .Event_name }}</a>　（{{.Event_ID}}）</td>
        </tr>
        <tr>
            <td>期間</td>
            <td>{{.Period}}</td>
        </tr>
        <tr>
            <td>開始日時</td>
            <td>{{.Start_time}}</td>
        </tr>
        <tr>
            <td>終了日時</td>
            <td>{{.End_time}}</td>
        </tr>
    </table>