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