
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
</tr>
<tr align="right" style="border-top-style:none;border-bottom-style:none;">
	<td style="border-top-style:none;border-bottom-style:none;"></td>
	<td style="border-top-style:none;border-bottom-style:none;"></td>
	{{ range .Total }}
		<td style="border-top-style:none;border-bottom-style:none;">
		{{ Comma . }}
		</td>
	{{ end }}
	<td style="border-top-style:none;border-bottom-style:none;"></td>
	<td style="border-top-style:none;border-bottom-style:none;"></td>
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
	<td style="border-top-style:none;border-bottom-style:none;"></td>
</tr>
<tr align="center" style="border-top-style:none;">
	<td style="border-top-style:none; border-bottom-style:none;">現頁</td>
	<td style="border-top-style:none; border-bottom-style:none;">累計獲得</td>
	{{ $e := .Eventid }}
	{{ $u := .Userno }}
	{{ range .Ifrm }}
		<td style="border-top-style:hidden;border-bottom-style:none;">
			<a href="list-cntrbS?eventid={{ $e }}&userno={{ $u }}&ifrm={{ . }}">累計順</a>
		</td>
	{{ end }}
	<td style="border-top-style:none;border-bottom-style:none;">リスナー名</td>
	<td style="border-top-style:none;border-bottom-style:none;">Tlsnid</td>
</tr>
<tr align="center" style="border-top-style:none;">
	<td style="border-top-style:none; border-bottom-style:none;border-bottom-style:none;">順位</td>
	<td style="border-top-style:none; border-bottom-style:none;border-bottom-style:none;">ポイント</td>
	{{ $e := .Eventid }}
	{{ $u := .Userno }}
	{{ range .Ifrm }}
		<td style="border-top-style:hidden;border-bottom-style:none;">
			<a href="list-cntrbS?eventid={{ $e }}&userno={{ $u }}&ifrm={{ . }}&sort=D">増分順</a>
		</td>
	{{ end }}
	<td style="border-top-style:none;border-bottom-style:none;"></td>
	<td style="border-top-style:none;border-bottom-style:none;">(履歴)</td>
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
</tr>