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
    <br>
    <br>
{{/* <h2>{{.function}}</h2> */}}
<p style="padding-left:2em">
{{.Msg001}}{{.Eventid}}{{.Msg002}}
</p>
</body>
</html>
