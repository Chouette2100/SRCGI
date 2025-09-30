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
<p>枠別リスナー別貢献ポイント取得対象の決定　　<span style="color:red;">初めて使うときは表の後にある説明をよく読んでください！</span></p>
<table style="text-align: center">
<tr><td style="width:2em"></td><td><a href="https://www.showroom-live.com/event/{{.Eventid}}">{{.Eventname}}</a></td></tr>
<tr><td style="width:2em"></td><td>{{.Period}}</td></tr>
</table>

{{/*
{{ range . }}
<form acttion="edit-user?eventid={{.Eventid}}" method="POST" id="{{.Formid}}"></form>
{{ end }}
*/}}
<form acttion="edit-cntrbpoints" method="POST">
<table border="1">
<tr>
<th>ルーム名</th>
<th>Prof./FR/Cnt.</th>
<th>ジャンル</th>
<th align="center" style="border-right-style:none;">ランク</th>
<th align="right" style="border-left-style:none;"></th>
<th>レベル</th>
<th>フォロ数</th>
<th>獲得ポイント</th>
<th>貢献<br>取得</th>
</tr>

{{ $i := 0 }}
<input type="hidden" name="eventid" value="{{.Eventid}}">

{{ range .Roominfolist }}
<tr>
<td><a href="https://www.showroom-live.com/{{.Account}}">{{.Name}}
{{/*
</a><input type="hidden" name="userid" value="{{.Userno}}" form="{{.Formid}}">
</a><input type="hidden" name="eventid" value="{{.Eventid}}" form="{{.Formid}}">
*/}}
</td>
<td>
<a href="https://www.showroom-live.com/room/profile?room_id={{.Userno}}">Prof.</a>/
<a href="https://www.showroom-live.com/room/fan_club?room_id={{.Userno}}">FR</a>/
<a href="https://www.showroom-live.com/event/contribution/{{.Eventid}}?room_id={{.Userno}}">Cnt.</a>
</td>
<td>{{.Genre}}</td>
<td align="center" style="border-right-style:none;">{{.Rank}}</td>
<td align="right" style="border-left-style:none;">{{.Nrank}}</td>
<td align="right">{{.Slevel}}</td><td align="right">{{.Sfollowers}}</td><td align="right">{{.Spoint}}</td>

<td align="center">
  {{ if eq .Iscntrbpoint "Checked" }}
    取得対象
    <input type="hidden" name="usr{{$i}}" value="0" >
  {{ else }}
    <input type="checkbox" name="usr{{$i}}"    value="{{.Userno}}" ></td>
  {{ end }}
</td>

</tr>

{{ $i = ( add $i 1 )}}
{{ end }}


</table>
<br>
　　<input type="submit" value="決定">　（チェックの入ったルームを枠別リスナー別貢献ポイント取得の対象とします。「決定」は取り消すことはできません。
</form>
<br>
<p style="padding-left:4em">イベント参加者として登録されているルームの一覧です<br><br>
<span style="color:red">ただし、イベント開始前は以下の「一覧にないルームの追加」で追加したルームのみ表示され、<br>
イベント開始後、指定した順位の範囲にあるルームが自動的に追加されます。</span><br><br>
リスト・グラフのルーム名には「表示名」が使われます<br>
「グラフの色」の色見本は更新ボタンを押したあと入力値を反映します
（表示名は<a href="/edit-user?eventid={{.Eventid}}">「(DB登録済)イベント参加ルーム一覧（確認・編集）」</a>で変更できます）<br>
</p>

<p style="padding-left:4em">獲得ポイントは定期の獲得ポイント取得で得たものか、<br>
    最後に「イベント参加者ルーム情報の追加と更新」を実行したときのものです。<br>
    「データ取得」にチェックの入っていないルームの獲得ポイントは定期的に更新されませんが、<br>
    「イベント参加者ルーム情報の追加と更新」は
    <a href="/edit-user?eventid={{.Eventid}}">「(DB登録済)イベント参加ルーム一覧（確認・編集）」</a>からお願いします。
</p>

<font color="blue">
    <form action="new-user" method="GET">
        <p style="padding-left:4em"><span style="font-weight:bold;">一覧にないルームの追加</span>　　ユーザーID：
            <input type="hidden" name="eventid" value="{{.Eventid}}">
            <input type="hidden" name="func" value="newuser">
            <input type="text" name="roomid" value="999999" required pattern="[0-9]+">　
            <input type="submit" value="登録"><br>
            イベント参加ルームの追加は「イベント参加者ルーム情報の追加と更新」で行うのですが、<br>
            順位に関係なく特定のルームを追加したい、参加ルームが多すぎてルームのサムネがイベントページに表示されていない、<br>
            などのケースはここでユーザーIDを指定して追加することができます。<br>
            ユーザーIDというのはプロフィールやファンルームのページのURLの最後にある6桁（か5桁？）の数字のことです。<br>
            <br>
            <font color="red">
            今後、イベント参加者のリストやルーム名から追加するルームを選択する方法を作成する予定です。
            </font>
        </p>
    </form>
</font>

</body>

</html>