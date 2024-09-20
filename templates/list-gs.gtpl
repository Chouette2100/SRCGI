
{{ $i := .Ncr }}
{{ range .Gslist }}
	<tr>
	<td align="right">
		{{ if eq .Orderno 999 }}
			---
		{{ else if ne .Orderno -1 }}
			{{ .Orderno }}
		{{ end }}
	</td>
	{{ range .Score }}
		<td align="right">
			{{ if eq . -1 }}
			n/a
			{{ else }}
			{{ Comma . }}
			{{ end }}
		</td>
	{{end}}
	<td><a href="https://www.showroom-live.com/room/profile?room_id={{.Userno}}">{{.User_name}}</a>（{{ .Rank }}）</td>
	<td align="right">
		<a href="https://www.showroom-live.com/r/{{.Url}}">Live</a> <a href="/closedevents?userno={{.Userno}}&mode=0&path=5">イベント履歴</a>
	</td>
	</tr>
{{end}}
</table>
<br>
<a href="/listgs?giftid=486&limit=100&maxacq={{.Maxacq}}">人気ライバー王</a> giftid=486<br>
<a href="/listgs?giftid=490&limit=50&maxacq={{.Maxacq}}">新人スタートダッシュ</a> giftid=490<br>
<a href="/listgs?giftid=494&limit=50&maxacq={{.Maxacq}}">アイドル王</a> giftid=494<br>
<a href="/listgs?giftid=495&limit=50&maxacq={{.Maxacq}}">俳優王</a> giftid=495<br>
<a href="/listgs?giftid=496&limit=50&maxacq={{.Maxacq}}">アナウンサー王</a> giftid=496<br>
<a href="/listgs?giftid=497&limit=50&maxacq={{.Maxacq}}">グローバル王</a> giftid=497<br>
<a href="/listgs?giftid=498&limit=50&maxacq={{.Maxacq}}">声優王</a> giftid=498<br>
<a href="/listgs?giftid=499&limit=50&maxacq={{.Maxacq}}">芸人王</a> giftid=499<br>
<a href="/listgs?giftid=500&limit=50&maxacq={{.Maxacq}}">タレント王</a> giftid=500<br>
<a href="/listgs?giftid=501&limit=50&maxacq={{.Maxacq}}">ライバー王</a> giftid=501<br>
<a href="/listgs?giftid=502&limit=50&maxacq={{.Maxacq}}">モデル王</a> giftid=502<br>
<a href="/listgs?giftid=503&limit=50&maxacq={{.Maxacq}}">バーチャル王</a> giftid=503<br>
<a href="/listgs?giftid=504&limit=50&maxacq={{.Maxacq}}">アーティスト王</a> giftid=504<br>
<br>
<a href="/listvgs?giftid=206&limit=50">最強ファンランキング giftid=206<br>
</body>
</html>
