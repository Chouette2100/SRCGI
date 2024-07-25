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
<td><button type="button" onclick="location.href='top?eventid={{.eventid}}'">イベントトップ</button></td>
<td><button type="button" onclick="location.href='list-last?eventid={{.eventid}}'">直近の獲得ポイント</button></td>
<td><button type="button" onclick="location.href='graph-total?eventid={{.eventid}}&maxpoint={{.maxpoint}}&gscale={{.gscale}}'">獲得ポイントグラフ</button></td>
<td></td>
  </tr>
</table>
  <h2>獲得ポイントグラフ</h2>
<form>
<input type="submit" value="再描画" formaction="graph-total" formmethod="POST" style="background-color: khaki">
<label>　　縮尺　
<input type="radio" name="gscale" value="100" {{ if eq .gscale "100" }} checked {{ end }}> 100%
<input type="radio" name="gscale" value="90" {{ if eq .gscale "90" }} checked {{ end }}> 90%
<input type="radio" name="gscale" value="80" {{ if eq .gscale "80" }} checked {{ end }}> 80%
<input type="radio" name="gscale" value="70" {{ if eq .gscale "70" }} checked {{ end }}> 70%
<label>　　表示する最大ポイント　<input type="text" name="maxpoint" value="{{.maxpoint}}" size="10" required pattern="[0-9]+">
<label>（表示範囲を制限しない場合は"0"とする）
<input type="hidden" name="eventid" value="{{.eventid}}">
　　<input type="submit" value="表示するルームを選ぶ" formaction="edit-user" formmethod="POST" style="background-color: khaki">
</form>
<br><br>
<img src="{{.filename}}" alt="" width="{{.gscale}}%">
</body>
</html>
