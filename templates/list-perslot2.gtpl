
{{ range . }}
<p style="padding-left:4em">
<a href="#R{{.Roomid}}">『{{ .Roomname }}』の配信枠毎の獲得ポイントへ</a>（ページ内）
</p>
{{ end }}
<p>
<br><br>
</p>
{{ range . }}
<a name="R{{.Roomid}}" id="R{{.Roomid}}" href="https://www.showroom-live.com/room/profile?room_id={{.Roomid}}">{{ .Roomname }}</a>
（{{.Roomid}}）
</p>
<table border="1">
<tr align="center"><td>配信日</td><td>開始時刻</td><td>終了時刻</td><td>獲得ポイント</td><td>累積獲得ポイント</td></tr>
{{ range .Perslotlist }}
<tr align="center"><td>{{.Dstart}}</td><td>{{.Tstart}}</td><td>{{.Tend}}</td><td align="right">{{.Point}}</td><td align="right">{{.Tpoint}}</td></tr>
{{ end }}
</table>
<p style="padding-left:12em">
<a href="#Top">ページ先頭へ</a>
</p>
<br>
{{ end }}
</body>
</html>
