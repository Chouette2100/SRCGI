<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<body>
<button type="button" onclick="location.href='top'">「SHOWROOMイベント結果表示」画面に戻る</button>
<br><br>
<p style="padding-left:2em">
新規イベントと参加ルームの登録
<p style="padding-left:4em;color:blue">
{{.Msg}}<br>
<form>
<table>
<tr><td style="width:4em"></td><td>イベントのID</td><td><input type="text" name="eventid" value="{{.Eventid}}" readonly >（イベントページのURLの最後のフィールド）
</td></tr>
<tr><td style="width:4em"></td><td>イベント名</td><td>{{.Eventname}}</td></tr>
<tr><td style="width:4em"></td><td>イベント期間</td><td>{{.Period}}</td></tr>
<tr><td style="width:4em"></td><td>イベント参加ルーム数</td><td>{{.Noroom}}</td></tr>
<tr><td style="width:4em"></td><td><label>ＤＢに登録する順位の範囲</td><td><input type="number" name="breg" value="1" min="1" max="200" required pattern="[0-9]+"><label>位から
<input type="number" name="ereg" value="10" size="3" required min="1" max="200"><label>位まで</td></tr>
<tr><td></td><td>獲得ポイントデータ取得のタイミング</td><td>毎時<input type="text" name="modmin" value="{{.Stm}}" size="2" required pattern="[0-9]+">分
<input type="text" name="modsec" value="{{.Sts}}" size="2" required pattern="[0-9]+">秒から
<input type="text" name="intervalmin" value="5" size="2" required pattern="[5]">分おきに取得する。</td></tr>
<tr><td></td><td>日々の獲得ポイントのリセット時刻</td><td>毎日<input type="text" name="resethh" value="4" size="2" required pattern="[0-9]+">時<input type="text" name="resetmm" value="0" size="2" required pattern="[0-9]+">分</td></tr>
