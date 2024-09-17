
<table border="1">
<tr align="center" style="border-bottom-style:none;">
	<td style="border-bottom-style:none;">年-月-日</td>
	{{ range .Stime }}
		<td style="border-bottom-style:none;">
		{{ t2s . "06-01-02" }}
		</td>
	{{ end }}
	<td style="border-bottom-style:none;"></td>
	<td style="border-bottom-style:none;"></td>
</tr>
<tr align="center" style="border-top-style:none;border-bottom-style:none;">
	<td style="border-top-style:none;border-bottom-style:none;">時:分</td>
	{{ range .Stime }}
		<td style="border-top-style:none;border-bottom-style:none;">
		{{ t2s . "15:04" }}
		</td>
	{{ end }}
	<td style="border-top-style:none;border-bottom-style:none;"></td>
	<td style="border-top-style:none;border-bottom-style:none;"></td>
</tr>
<tr align="center" style="border-top-style:none;">
	<td style="border-top-style:none;"></td>
	{{ range .Nof }}
		<td style="border-top-style:hidden;">
			（ {{ . }} ）
		</td>
	{{ end }}
	<td style="border-top-style:none;"></td>
	<td style="border-top-style:none;"></td>
</tr>