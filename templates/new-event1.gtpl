<tr><td></td><td>ポイント差の基準とする配信者</td><td>

<select name="nobasis">
		{{ range . }}
			<option value="{{.Userno}}" {{.Selected}}>{{.Userlongname}}　({{.Userno}})</option>
		{{ end }}
	</select>

　（基準とする配信者がリストにないときは「基準とする配信者を設定しない」を選びイベント登録後「イベントパラメータの設定」から設定してください）</td></tr> 
