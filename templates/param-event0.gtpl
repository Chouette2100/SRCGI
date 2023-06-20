<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<html>
<body>
<button type="button" onclick="location.href='top?eventid={{.Event_ID}}'">「SHOWROOMイベント結果表示」画面に戻る</button>
<br><br>
<p style="padding-left:2em">
イベント設定の変更
<p style="padding-left:4em">
このイベントに関わる設定を変更するときは、設定値を書き換えた後「設定変更」ボタンを押してください。<br>
<form>
<table>
<tr><td style="width:4em"></td><td>イベントのID</td><td><input type="hidden" name="eventid" value="{{.Event_ID}}" >{{.Event_ID}}</td></tr>
<tr><td style="width:4em"></td><td>イベント名</td><td>{{.Event_name}}</td></tr>
<tr><td style="width:4em"></td><td>イベント期間</td><td>{{.Period}}</td></tr>
<tr><td style="width:4em"></td><td>イベント参加ルーム数</td><td>{{.NoEntry}}（最新のデータでない可能性あり）</td></tr>
<tr><td style="width:4em"></td><td><label>ＤＢに登録する順位の範囲</td><td><input type="number" name="fromorder" value="{{.Fromorder}}" min="1" max="200" required pattern="[0-9]+"><label>位から
<input type="number" name="toorder" value="{{.Toorder}}" size="3" required min="1" max="200"><label>位まで</td></tr>
<tr><td></td><td>獲得ポイントデータ取得のタイミング</td><td>毎時<input type="text" name="modmin" value="{{.Modmin}}" size="2" required pattern="[0-9]+">分
<input type="text" name="modsec" value="{{.Modsec}}" size="2" required pattern="[0-9]+">秒から
<input type="text" name="intervalmin" value="{{.Intervalmin}}" size="2" required pattern="[0-9]+">分おきに取得する。</td></tr>
<tr><td></td><td>日々の獲得ポイントのリセット時刻</td><td>毎日<input type="text" name="resethh" value="{{.Resethh}}" size="2" required pattern="[0-9]+">時<input type="text" name="resetmm" value="{{.Resetmm}}" size="2" required pattern="[0-9]+">分</td></tr>
