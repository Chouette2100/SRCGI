
{{/*}}
<br>
========================================================================
<br>
<p>イベント選択（配信者から選ぶ）</p>

<p>1. 配信者の選択（イベント選択対象の絞り込み条件/ルーム情報の表示対象）<br>
　　ここでいう配信者とは各イベントの「ポイント差の基準とする配信者」のことを言います。<br>
　　配信者が設定されていないイベントは「ポイント差の基準とする配信者が設定されていない」を選択してください。

<form>

<p style="padding-left:2em">
配信者　
<select name="userno">
		{{ range . }}
			<option value="{{.Userno}}" {{.Selected}}>{{.Userlongname}}　({{.Userno}})</option>
		{{ end }}
	</select>

<span style="padding-left:2em"><input type="submit" value="決定" formaction="top" formmethod="GET"></span>

</p>

<span style="padding-left:2em">
ルーム情報（レベルとフォロワー数の推移、配信者を選択しておくこと）<br>
</span>
<span style="padding-left:4em">
<input type="submit" value="表示" formaction="list-level" formmethod="GET" >
　　
<input type="checkbox" name="levelonly" value="1" checked="checked">レベルが変化したところだけ表示する
</span>
</form>
</p>

{{*/}}