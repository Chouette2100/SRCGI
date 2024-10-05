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
      <td><button type="button" onclick="location.href='listgs'">ギフトランキング</button></td>
      <td></td>
      <td></td>
      <td></td>
    </tr>
</table>
<br>
<p><a href="{{ .Url }}">{{ .Campaignname }}</a>（{{ .Campaignid }}）</p>
<p>　　{{ .Grname }}（{{ .Grid }}）</p>
<br>
<form>
{{/*
　　　　ギフト種別を選択する
{{ $grid :=  .Grid }}
<select name="giftid" type="text">
{{ range .GiftRanking }}
    {{ if eq .Grid $grid }}
        <option selected value="{{ .Grid }}">{{ .Grname }}</option>
    {{ else }}
        <option value="{{ .Grid }}">{{ .Grname }}</option>
    {{ end }}
{{ end }}
</select>
*/}}

<input id="giftid" name="giftid" type="hidden" value="{{ .Grid }}">
<input id="userno" name="userno" type="hidden" value="{{ .Userno }}">

　　横(時刻)表示数
<input value="{{ .Maxacq }}" name="maxacq" type="number" size="5" min="1" max="15">

　　縦(ルーム)表示数
<input value="{{ .Limit }}" name="limit" type="number" size="5" min="20" max="500">

 　　<input type="submit" value="この条件で再表示" formaction="listgsc" formmethod="GET">

</form>
<br>

<table>
    <tr>
        <td>
        {{ if ne .Nft -1 }}
            <button type="button" onclick="location.href='listgsc?giftid={{.Grid}}&userno={{ .Userno }}&ie={{.Nft}}'">先頭に戻る</button>
        {{ else }}
            -----------
        {{ end }}
        </td>
        <td>
        {{ if ne .Npb -1 }}
            <button type="button" onclick="location.href='listgsc?giftid={{.Grid}}&userno={{ .Userno }}&ie={{.Npb}}'">１ページ戻る</button>
        {{ else }}
            -----------
        {{ end }}
        </td>
        <td>
        {{ if ne .N1b -1 }}
            <button type="button" onclick="location.href='listgsc?giftid={{.Grid}}&userno={{ .Userno }}&ie={{.N1b}}'">一枠分戻る</button>
        {{ else }}
            -----------
        {{ end }}
        </td>
        <td>
        {{ if ne .N1f -1 }}
            <button type="button" onclick="location.href='listgsc?giftid={{.Grid}}&userno={{ .Userno }}&ie={{.N1f}}'">一枠分進む</button>
        {{ else }}
            -----------
        {{ end }}
        </td>
        <td>
        {{ if ne .Npf -1 }}
            <button type="button" onclick="location.href='listgsc?giftid={{.Grid}}&userno={{ .Userno }}&ie={{.Npf}}'">１ページ進む</button>
        {{ else }}
            -----------
        {{ end }}
        </td>
        <td>
        {{ if ne .Nlt -1 }}
            <button type="button" onclick="location.href='listgsc?giftid={{.Grid}}&userno={{ .Userno }}&ie={{.Nlt}}'">最後に進む</button>
        {{ else }}
            -----------
        {{ end }}
        </td>
    </tr>
</table>

<table border="1">
<tr align="center" style="border-bottom-style:none;">
	<td style="border-bottom-style:none;">年-月-日</td>
	{{ range .Stime }}
		<td style="border-bottom-style:none;">
		{{ t2s . "06-01-02" }}
		</td>
	{{ end }}
	<td style="border-bottom-style:none;"></td>
</tr>
<tr align="center" style="border-top-style:none;border-bottom-style:none;">
	<td style="border-top-style:none;border-bottom-style:none;">時:分</td>
	{{ range .Stime }}
		<td style="border-top-style:none;border-bottom-style:none;">
		{{ t2s . "15:04" }}
		</td>
	{{ end }}
	<td style="border-top-style:none;border-bottom-style:none;">リスナー（userid）</td>
</tr>
<tr align="center" style="border-top-style:none;">
	<td style="border-top-style:none;"></td>
	{{ range .Nof }}
		<td style="border-top-style:hidden;">
			（ {{ . }} ）
		</td>
	{{ end }}
	<td style="border-top-style:none;"></td>
</tr>
{{ $i := .Ncr }}
{{ range .Gsclist }}
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
			{{ else if eq . 0 }}
            ---
			{{ else }}
			{{ Comma . }}
			{{ end }}
		</td>
	{{end}}
	<td>{{ .Name }}（{{ .Viewerid }}）</td>
	</tr>
{{end}}
</table>
</body>
</html>
