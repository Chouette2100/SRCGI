
{{ $i := .Ncr }}
{{ range .Vgslist }}
	<tr>
	<td align="right">
		{{ if eq .Orderno 0 }}
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
	<td>{{ .Viewername }}<td>
	<td align="right">{{ .Viewerid }}</td>
	</tr>
{{end}}
</table>
<br>
<a href="/listgs?giftid=486">人気ライバー王</a> giftid=486<br>
<a href="/listgs?giftid=490">新人スタートダッシュ</a> giftid=490<br>
<a href="/listgs?giftid=494">アイドル王</a> giftid=494<br>
<a href="/listgs?giftid=495">俳優王</a> giftid=495<br>
<a href="/listgs?giftid=496">アナウンサー王</a> giftid=496<br>
<a href="/listgs?giftid=497">グローバル王</a> giftid=497<br>
<a href="/listgs?giftid=498">声優王</a> giftid=498<br>
<a href="/listgs?giftid=499">芸人王</a> giftid=499<br>
<a href="/listgs?giftid=500">タレント王</a> giftid=500<br>
<a href="/listgs?giftid=501">ライバー王</a> giftid=501<br>
<a href="/listgs?giftid=502">モデル王</a> giftid=502<br>
<a href="/listgs?giftid=503">バーチャル王</a> giftid=503<br>
<a href="/listgs?giftid=504">アーティスト王</a> giftid=504<br>
<br>
<a href="/listvgs?giftid=206">最強ファンランキング giftid=206<br>
</body>
</html>
