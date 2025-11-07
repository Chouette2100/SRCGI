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
    {{/*}}
    <tr>
        <td><button type="button" onclick="location.href='eventtop?eventid={{.Eventid}}'">イベントトップ</button></td>
        <td><button type="button" onclick="location.href='list-last?eventid={{.Eventid}}'">直近の獲得ポイント</button></td>
        <td><button type="button"
                onclick="location.href='graph-total?eventid={{.Eventid}}&maxpoint={{.Maxpoint}}&gscale={{.Gscale}}'">獲得ポイントグラフ</button>
        </td>
        <td></td>
    </tr>
    {{*/}}
</table>
<br><br>
<p style="padding-left:2em">
新規イベントと参加ルームの登録
<p style="padding-left:4em;color:red">
現在開催前の取得対象ルームの登録は行えません<br>
取得対象ルームは自動的に設定されます。取得範囲を広げたり、取得対象ルームを追加する場合は<br>
開催後イベントのページ（例えば獲得ポイント一覧）に入ってから「イベントトップ」を選び、<br>
取得範囲を広げるには「イベントパラメータの設定」をお使いください。<br>
一覧にないルームを追加するには「(DB登録済み)イベント参加ルーム一覧（確認・編集）」をお使いください。<br>
お手数をおかけしますがよろしくお願いします。
</p>
<p style="padding-left:4em;color:blue">
{{.Msg}}<br>
<form>
<table>
<tr><td style="width:4em"></td><td>イベントのID</td><td><input type="text" name="eventid" value="{{.Eventid}}" readonly size="50">（イベントページのURLの最後のフィールド）
</td></tr>
<tr><td style="width:4em"></td><td>イベント名</td><td>{{.Eventname}}</td></tr>
<tr><td style="width:4em"></td><td>イベント期間</td><td>{{.Period}}</td></tr>
<tr><td style="width:4em"></td><td>イベント参加ルーム数</td><td>{{.Noroom}}</td></tr>
<tr><td style="width:4em"></td><td><label>ＤＢに登録する順位の範囲</td><td><input type="number" name="breg" value="1" size="3" min="1" max="20" required pattern="[0-9]+"><label>位から
<input type="number" name="ereg" value="20" size="3" required min="1" max="30"><label>位まで</td></tr>
<tr><td></td><td>獲得ポイントデータ取得のタイミング</td><td>毎時<input type="text" name="modmin" value="{{.Stm}}" size="2" required pattern="[0-9]+">分
<input type="text" name="modsec" value="{{.Sts}}" size="2" required pattern="[0-9]+">秒から
<input type="text" name="intervalmin" value="5" size="2" required pattern="[0-9]+">分おきに取得する。</td></tr>
<tr><td></td><td>日々の獲得ポイントのリセット時刻</td><td>毎日<input type="text" name="resethh" value="4" size="2" required pattern="[0-9]+">時<input type="text" name="resetmm" value="0" size="2" required pattern="[0-9]+">分</td></tr>
