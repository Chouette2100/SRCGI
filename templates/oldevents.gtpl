<html>

<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
    {{/*}}
    <style type="text/css">
        th,
        td {
            border: solid 1px;
        }

        table {
            border-collapse: collapse;
        }

        .noborder {
            border: 0px none;
        }
    </style>
    {{*/}}
</head>

<body>

    <br>
    {{/*
    (外部リンク)<br>
    <a href="https://zenn.dev/">Zenn</a> - <a
        href="https://zenn.dev/chouette2100/books/d8c28f8ff426b7">SHOWROOMのAPI、その使い方とサンプルソース</a>
    <div style="text-indent: 2rem;"><a
            href="https://chouette2100.com/">記事/ソース/CGI一覧表</a>　（証明書の期限が切れてしまっていました。2022年8月29日、有効な証明書に切り替えました）</div>
    <p>-------------------------------------------------------------</p>
    (サイト内リンク)<br>
    <div style="text-indent: 2rem;"><a href="t009top">t009:配信中ルーム一覧</a></div>
    <p>-------------------------------------------------------------</p>
    */}}

    
    <table>
        <tr>
            <td><button type="button" onclick="location.href='top'">トップ</button>　</td>
            <td><button type="button" onclick="location.href='currentevents'">開催中イベント一覧</button></td>
            <td><button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧</button></td>
            <td>終了イベント一覧</td>
        </tr>
    </table>
    <br>

        {{ $u := .User.Userno }}
        <br>{{ .User.User_name }} （ {{ .User.Userno }} ）
        <br>
        <br>
        <span style="margin-left: 4em; font-size: 1.2em;">処理は終了しました！</span><br><br>
        <button type="button" onclick="location.href='closedevents?userno={{.User.Userno}}&mode=0&path=5'"
                   style="margin-left: 6em; font-size: 1.2em; background-color: rgba(255, 0, 0, 0.5);" >確認(終了イベント一覧へ)</button>
        <br>
        <br>
        <hr>
        <br>
        <p>
            {{ .ErrMsg}}
        </p>

</body>

</html>