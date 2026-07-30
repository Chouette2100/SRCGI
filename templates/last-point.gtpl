<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<body>
    {{if (HasAnnouncement)}}
    <div style="padding: 15px; margin: 0 0 20px 0; border-radius: 6px; background-color: {{(GetAnnouncement).BgColor}}; color: {{(GetAnnouncement).TextColor}}; font-size: 16px; font-weight: bold; text-align: center; border: 2px solid {{(GetAnnouncement).TextColor}}; word-wrap: break-word;">
      {{(GetAnnouncement).Message}}
    </div>
    {{end}}

<p>{{.function}}</p>
<p style="padding-left:2em">
{{.comment}}
</p>
<p style="padding-left:2em">
<button type="button" onclick="history.back()">結果表示選択画面に戻る</button><br>
</p>
</body>
</html>
