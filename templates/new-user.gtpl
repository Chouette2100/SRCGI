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
        <td><button type="button" onclick="location.href='top?eventid={{.Event_ID}}'">イベントトップ</button></td>
        <td><button type="button" onclick="location.href='list-last?eventid={{.Event_ID}}'">直近の獲得ポイント</button></td>
        <td><button type="button"
                onclick="location.href='graph-total?eventid={{.Event_ID}}&maxpoint={{.Maxpoint}}&gscale={{.Gscale}}'">獲得ポイントグラフ</button>
        </td>
        <td></td>
    </tr>
</table>

<br><br>
<p>個別に指定した配信者の登録</p>

イベント

<table>
<tr><td style="width:4em"><td style="width:8em">イベント名</td><td><a href="https://www.showroom-live.com/event{{.Event_ID}}">{{.Event_name}}</a></td></tr>
<tr><td></td><td>期間</td><td>{{.Period}}</td></tr>
</table>
<p>
{{.Msg1}}
</p>
<form action="edit-user" method="GET">
<input type="hidden" name="func" value="newuser">
<input type="hidden" name="eventid" value="{{.Event_ID}}">
<input type="hidden" name="userid" value="{{.Roomid}}">
<table>
<tr><td style="width:4em"></td><td>ルーム名</td><td>{{.Roomname}}
　<a href="https://www.showroom-live.com/{{.Roomurlkey}}">配信ルーム</a>　<a href="https://www.showroom-live.com/room/profile?room_id={{.Roomid}}">プロフィール</a></td></tr>
<tr><td></td><td>ジャンル</td><td>{{.Genre}}</td></tr>
<tr><td></td><td>ランク</td><td>{{.Rank}}</td></tr>
<tr><td></td><td>レベル</td><td>{{.Level}}</td></tr>
<tr><td></td><td>フォロワー数</td><td>{{.Followers}}</td></tr>
<tr><td></td><td>表示名</td><td><input type="text" name="longname" value="{{.Longname}}" required >　グラフやリストに表示される配信者名です。</td></tr>
<tr><td></td><td>短縮表示名</td><td><input type="text" name="shortname" value="{{.Shortname}}" required >　1文字かせいぜい2文字まで</td></tr>
</table>
<p style="color:{{.Msg2color}}">
{{.Msg2}}
</p>
<table>
<tr><td style="width:4em"></td>
<td style="width:8em">
<input type="{{.Submit}}" value="登録する"><br>
</form>
</td><td>
<form action="edit-user" method="GET">
<input type="hidden" name="eventid" value="{{.Event_ID}}">
<input type="submit" value="{{.Label}}"><br>
</form>
</td></tr>
</body>
</html>
