<html>

<head>
	<style>
		.p1 {
			white-space: pre-wrap;
			margin-left: 25;
		}
	</style>
</head>

<body>
    {{if (HasAnnouncement)}}
    <div style="padding: 15px; margin: 0 0 20px 0; border-radius: 6px; background-color: {{(GetAnnouncement).BgColor}}; color: {{(GetAnnouncement).TextColor}}; font-size: 16px; font-weight: bold; text-align: center; border: 2px solid {{(GetAnnouncement).TextColor}}; word-wrap: break-word;">
      {{(GetAnnouncement).Message}}
    </div>
    {{end}}

	{{/*	メンテナンス */}}
    <button type="button" onclick="location.href='top'">Top</button>　
    <button type="button" onclick="location.href='currentevents'">開催中イベント一覧表</button>　
    <button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧表</button>　
    <button type="button" onclick="location.href='closedevents'">終了イベント一覧表</button>　
	{{/*	メンテナンス ここまで	*/}}
    <br>

