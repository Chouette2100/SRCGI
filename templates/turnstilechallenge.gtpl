{{define "turnstilechallenge"}}
<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
    {{if .TurnstileSiteKey}}
    <script src="https://challenges.cloudflare.com/turnstile/v0/api.js" async defer></script>
    {{end}}
</head>
<body>
		<!-- Turnstileチャレンジ表示 -->
		<div style="border: 2px solid #4A90E2; padding: 20px; border-radius: 5px; max-width: 600px; background-color: #f9f9f9; margin: 20px auto;">
			<h3>セキュリティチェック</h3>
			{{if .TurnstileError}}
			<p style="color: red; font-weight: bold;">{{.TurnstileError}}</p>
			{{end}}
			<p>このページを表示するには、セキュリティチェックを完了してください。</p>
			<p>「確認して続行」ボタンを押すとクッキーが保存されます</p>
			<form method="POST" action="{{if .ChallengeAction}}{{.ChallengeAction}}{{else}}/{{end}}">
				{{if .RequestID}}
				<input type="hidden" name="requestid" value="{{.RequestID}}">
				{{end}}
				{{range $key, $value := .ChallengeHiddenFields}}
				<input type="hidden" name="{{$key}}" value="{{$value}}">
				{{end}}
				<div class="cf-turnstile" data-sitekey="{{.TurnstileSiteKey}}" data-theme="light"></div>
				<br>
				<button type="submit" style="padding: 10px 20px; background-color: #4A90E2; color: white; border: none; border-radius: 5px; cursor: pointer; font-size: 16px;">確認して続行</button>
			</form>
		</div>
		<br><br>
</body>
</html>
{{end}}
