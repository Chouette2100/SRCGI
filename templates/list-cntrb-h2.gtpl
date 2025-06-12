
<table border="1">
<tr align="center" style="border-bottom-style:none;">
	<td style="border-bottom-style:none;"></td>
	<td style="border-bottom-style:none;"></td>
	{{ range .S_stime }}
		<td style="border-bottom-style:none;">
		{{ . }}
		</td>
	{{ end }}
	<td style="border-bottom-style:none;"></td>
	<td style="border-bottom-style:none;"></td>
	<td style="border-bottom-style:none;"></td>
</tr>
<tr align="center" style="border-top-style:none;border-bottom-style:none;">
	<td style="border-top-style:none;border-bottom-style:none;"></td>
	<td style="border-top-style:none;border-bottom-style:none;"></td>
	{{ range .S_etime }}
		<td style="border-top-style:none;border-bottom-style:none;">
		{{ . }}
		</td>
	{{ end }}
	<td style="border-top-style:none;border-bottom-style:none;"></td>
	<td style="border-top-style:none;border-bottom-style:none;"></td>
	<td align="center" style="border-top-style:none;border-bottom-style:none;color=red;">過去イベントの</td>
</tr>
<tr align="right" style="border-top-style:none;border-bottom-style:none;">
	<td style="border-top-style:none;border-bottom-style:none;"></td>
	<td style="border-top-style:none;border-bottom-style:none;"></td>
	{{ range .Total }}
		<td style="border-top-style:none;border-bottom-style:none;">
		{{ Comma . }}
		</td>
	{{ end }}
	<td align="center" style="border-top-style:none;border-bottom-style:none;">リスナー名</td>
	<td align="center" style="border-top-style:none;border-bottom-style:none;">Tlsnid</td>
	<td align="center" style="border-top-style:none;border-bottom-style:none;">貢献ポイント</td>
</tr>
<tr align="right" style="border-top-style:none;border-bottom-style:none;">
	<td style="border-top-style:none;border-bottom-style:none;"></td>
	<td style="border-top-style:none;border-bottom-style:none;"></td>
	{{ range .Earned }}
		<td style="border-top-style:none;border-bottom-style:none;">
		/{{ Comma . }}
		</td>
	{{ end }}
	<td style="border-top-style:none;border-bottom-style:none;"></td>
	<td align="center" style="border-top-style:none;border-bottom-style:none;">枠別貢献pt一覧</td>
	<td align="center" style="border-top-style:none;border-bottom-style:none;">履歴</td>
</tr>
<tr align="center" style="border-top-style:none;">
	<td style="border-top-style:none; border-bottom-style:none;">現頁</td>
	<td style="border-top-style:none; border-bottom-style:none;">累計獲得</td>
	{{ $e := .Eventid }}
	{{ $u := .Userno }}
	{{ $i := .Ncr }}
	{{ range .Ifrm }}
		<td style="border-top-style:hidden;border-bottom-style:none;">
			<a href="list-cntrbS?eventid={{ $e }}&userno={{ $u }}&ifrm={{ . }}&ie={{ $i }}">累計順</a>
		</td>
	{{ end }}
	<td style="border-top-style:none;border-bottom-style:none;"></td>
	<td style="border-top-style:none;border-bottom-style:none;"></td>
	<td align="center" style="border-top-style:none;border-bottom-style:none;">（暫定）</td>
</tr>
<tr align="center" style="border-top-style:none;">
	<td style="border-top-style:none; border-bottom-style:none;border-bottom-style:none;">順位</td>
	<td style="border-top-style:none; border-bottom-style:none;border-bottom-style:none;">ポイント</td>
	{{ $e := .Eventid }}
	{{ $u := .Userno }}
	{{ $i := .Ncr }}
	{{ range .Ifrm }}
		<td style="border-top-style:hidden;border-bottom-style:none;">
			<a href="list-cntrbS?eventid={{ $e }}&userno={{ $u }}&ifrm={{ . }}&sort=D&ie={{ $i }}">増分順</a>
		</td>
	{{ end }}
	<td style="border-top-style:none;border-bottom-style:none;"></td>
	<td style="border-top-style:none;border-bottom-style:none;"></td>
	<td style="border-top-style:none;border-bottom-style:none;"></td>
</tr>
<tr align="center" style="border-top-style:none;">
	<td style="border-top-style:none;"></td>
	<td style="border-top-style:none;"></td>
	{{ range .Nof }}
		<td style="border-top-style:hidden;">
			（ {{ . }} ）
		</td>
	{{ end }}
	<td style="border-top-style:none;"></td>
	<td style="border-top-style:none;"></td>
	<td style="border-top-style:none;"></td>
</tr>