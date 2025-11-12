		{{define "turnstilechallenge"}}
			<!-- Turnstileチャレンジ表示 -->
			<div style="border: 2px solid #4A90E2; padding: 20px; border-radius: 5px; max-width: 600px; background-color: #f9f9f9;">
				<h3>セキュリティチェック</h3>
				{{if .TurnstileError}}
				<p style="color: red; font-weight: bold;">{{.TurnstileError}}</p>
				{{end}}
				<p>{{ .Event_name }}（{{.Roomid}}）の貢献ランキングを表示するには、セキュリティチェックを完了してください。</p>
				<p>「確認して続行」ボタンを押すとクッキーが保存されます</p>
				<form method="POST" action="contributors">
					<input type="hidden" name="ieventid" value="{{.Ieventid}}">
					<input type="hidden" name="roomid" value="{{.Roomid}}">

					<input type="hidden" name="requestid" value="{{.RequestID}}">

					<div class="cf-turnstile" data-sitekey="{{.TurnstileSiteKey}}" data-theme="light"></div>
					<br>
					<button type="submit" style="padding: 10px 20px; background-color: #4A90E2; color: white; border: none; border-radius: 5px; cursor: pointer; font-size: 16px;">確認して続行</button>
				</form>
			</div>
			<br><br>
		{{ end }}
