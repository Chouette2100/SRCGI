<html>

<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
    <style type="text/css">
        th,
        td {
            border: solid 1px;
        }

        table {
            border-collapse: collapse;
            /*
            width: 100%;
            */
        }
    </style>
</head>

<body>
    {{if (HasAnnouncement)}}
    <div style="padding: 15px; margin: 0 0 20px 0; border-radius: 6px; background-color: {{(GetAnnouncement).BgColor}}; color: {{(GetAnnouncement).TextColor}}; font-size: 16px; font-weight: bold; text-align: center; border: 2px solid {{(GetAnnouncement).TextColor}}; word-wrap: break-word;">
      {{(GetAnnouncement).Message}}
    </div>
    {{end}}

    <p>ルーム状況（ /api/room/status?roomid=..... ）</p>
    <p style="padding: 2em;">
        {{ range .Roomstatus }}
        {{ . }}<br>
        {{ end }}
    </p>
    <br>
    <hr>
    <br>
    {{/*}}
    {{ template "footer" }}
    {{*/}}
</body>

</html>