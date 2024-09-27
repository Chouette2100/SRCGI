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
  <td><button type="button" onclick="location.href='listvgs'">ファンランキング</button></td>
<td></td>
<td></td>
<td><button type="button" onclick="location.href='listgs?giftid={{.Grid}}'">ギフトランキングリスト</button></td>
  </tr>
</table>
  <h2>{{.Grname}}</h2>
<form>
<input type="submit" value="再描画" formaction="graphgs" formmethod="POST" style="background-color: khaki">
<label>　　縮尺　
<input type="radio" name="gscale" value="100" {{ if eq .Gscale 100 }} checked {{ end }}> 100%
<input type="radio" name="gscale" value="90" {{ if eq .Gscale 90 }} checked {{ end }}> 90%
<input type="radio" name="gscale" value="80" {{ if eq .Gscale 80 }} checked {{ end }}> 80%
<input type="radio" name="gscale" value="70" {{ if eq .Gscale 70 }} checked {{ end }}> 70%
<label>　　表示する最大ポイント　<input type="text" name="maxpoint" value="{{.Maxpoint}}" size="10" required pattern="[0-9]+">
<label>（表示範囲を制限しない場合は"0"とする）
<input type="hidden" name="campaignid" value="{{.Campaignid}}">　　
<input type="hidden" name="giftid" value="{{.Grid}}">　　
{{/*
<input
      type="checkbox"
      id="resetcolor"
      name="resetcolor"
      value="on" />
    <label for="resetcolor">グラフ線の配色を初期化する</label>
<br><input type="submit" value="表示するルームを選ぶ・グラフ線の色を変える"
            formaction="edit-user" formmethod="POST" style="background-color: khaki">
*/}}
</form>
<br><br>
<img src="{{.Filename}}" alt="" width="{{.Gscale}}%">
</body>
</html>
