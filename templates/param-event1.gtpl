<tr><td></td><td>ポイント差の基準とする配信者</td><td>
<select name="nobasis">
		{{ range . }}
			<option value="{{.Userno}}" {{.Selected}}>{{.Userlongname}}　({{.Userno}})</option>
		{{ end }}
	</select>
</td></tr>
