<h3>掲示板　（直接の連絡は<a href="https://twitter.com/Seppina1/" target="_blank" rel="noopener noreferrer">こちら</a>へ）</h3>
<div>
	<form action='/write-bbs' method='POST' style="color: {{ .Manager }};">
		<input type="hidden" name=color value="{{ .Manager }}" />
		{{ range $i, $v := .Cntlist }}
			<input type="hidden" name="cnt{{ $i }}" value="{{ Add $i 1 }}" />
		{{ end }}
		<textarea name='title' cols='40' rows='1'>件名</textarea>
		<textarea name='name' cols='20' rows='1'>お名前</textarea><br>
		<div>
			{{/*	意図したとおりに動作しない！？
			{{ if or ( lt $i 4 ) ( ne .Manager "black" ) }}
			*/}}
			{{ range $i, $v := .Cntlist }}
				{{ if lt $i 4 }}
					<input type="radio" id="CntChoice{{ Add $i 1}}" name="cntw" value="{{ Add $i 1 }}" {{ if eq $i 0 }} checked {{ end }}}} />
					<label for="CntChoice{{ Add $i 1 }}">{{ CntToName $i }}</label>
				{{ end }}
			{{ end }}
			{{ if ne .Manager "black" }}
				<input type="radio" id="CntChoice5" name="cntw" value="5"  />
				<label for="CntChoice5">{{ CntToName 4 }}</label>
			{{ end }}
			{{/*
			{{ end }}
			*/}}
		</div>
		<textarea name='body' cols='60' rows='2' required maxlength="370" placeholder='不具合・要望・質問・その他なんでも（文字数は入力できる範囲（３２０文字程度）でお願いします）'></textarea><br>
		<input type='submit' value='書込'>
	</form>
</div>
<hr>
<p>投稿一覧</p>
<form action='?' method='POST'>
	<input type="hidden" name=offset value="{{ .Offset }}" />
	<input type="hidden" name=limit value="{{ .Limit }}" />
	<input type="hidden" name=from value="disp-bbs" />
	<input type="hidden" name=manager value="{{ .Manager }}" />
	<fieldset>
		<legend>表示したいジャンルを選んでください</legend>
		<input type="radio" id="CntDispA" name="cntr" value="9" checked />
		<label for="CntDispA">{{ CntToName 5 }}</label>
		{{ $cntr := .Cntr }}
		{{ range $i, $v := .Cntlist }}
		{{/*
		<input type="checkbox" id="cnt{{ $i }}" name="cnt{{ $i }}" value="{{ Add $i 1 }}" {{ if gt $v 0 }} checked {{
			end }}}} />
		<label for="cnt{{ $i }}">{{ CntToName $i }}</label>
		*/}}
		<input type="radio" id="CntDisp{{ Add $i 1}}" name="cntr" value="{{ Add $i 1 }}" {{ if eq ( Add $i 1) $cntr }} checked {{ end }}}} />
		<label for="CntDisp{{ Add $i 1 }}">{{ CntToName $i }}</label>
		{{ end }}
		<input type='submit' name="action" value='再表示(top)' formaction="/disp-bbs" />
		<input type='submit' name="action" value='next' formaction="/disp-bbs" />
		<input type='submit' name="action" value='prev.' formaction="/disp-bbs" />
	</fieldset>

{{ range .Loglist }}
<div style="color: {{ .Color }};">
	No.{{.ID}}【{{ CntToName ( Add .Cntw -1 ) }}】「{{ htmlEscapeString .Title}}」({{ htmlEscapeString .Name}})
	（{{FormatTime .CTime "2006-01-02 15:04" }}）
	<p class="p1">{{ htmlEscapeString .Body }}</p>
	<hr>
</div>
{{ end }}
		<input type='submit' name="action" value='再表示(top)' formaction="/disp-bbs" />
		<input type='submit' name="action" value='next' formaction="/disp-bbs" />
		<input type='submit' name="action" value='prev.' formaction="/disp-bbs" />
</form>
