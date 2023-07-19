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
    {{ $tn := .TimeNow }}

    <button type="button" onclick="location.href='top'">Top</button>　
    <button type="button" onclick="location.href='currentevent'">開催中イベント一覧表</button>　
    {{/*}}
    <button type="button" onclick="location.href='scheduledevent'">開催予定イベント一覧表</button>　
    {{*/}}

    <br>
    <br>
    <p>開催予定イベント一覧表</p>
    <p>　　作成中。 現時点ではブロックイベントとイベントボックスの展開は行っていません。
        <br>これから、子イベントの展開、獲得ポイントデータ取得の登録や参加予定ルームの一覧あたりを作っていく予定です。</p>
    <p>　　 イベント数： {{ .Totalcount }} （{{ UnixTimeToStr $tn }}）</p>
    <br>
    <table>
        <tr style="text-align: center">
            <td>イベントID</td>
            <td>イベント名</td>
            <td>開始日時</td>
            <td>終了日時</td>
            <td>申込期限</td>
            <td>レベル<br>制限</td>
            <td>ランク<br>制限</td>
            <td style="width: 30%">備考</td>
        </tr>
        {{ range .Eventlist }}
        {{ if eq .Is_box_event true }}
        <tr style="background-color: darkgrey">
            {{ else }}
        <tr>
            {{ end }}
            <td>
                {{ .Event_id }}
            </td>
            <td>
                <a href="https://showroom-live.com/event/{{ .Event_url_key }}">{{ .Event_name }}</a>
            </td>
            <td>
                {{ UnixTimeToStr .Started_at }}
            </td>
            <td>
                {{ UnixTimeToStr .Ended_at }}
            </td>
            {{ $t := UnixTimeToStr .Offer_ended_at }}
            {{ if gt $tn .Offer_ended_at }}
            <td style="color: red">
                {{ $t }}
            </td>
            {{ else }}
            <td style="color: black">
                {{ $t }}
            </td>
            {{ end }}
            <td style="text-align: center">
                {{ if and ( gt .Required_level_max 0 ) ( lt .Required_level_max 99998 ) }}
                ～{{ .Required_level_max }}
                {{ end }}
                {{ if and ( gt .Required_level 0 ) ( lt .Required_level 99998 ) }}
                {{ .Required_level }}～
                {{ end }}
                {{ if eq .Is_entry_scope_inner true }}
                限定
                {{ end }}
            </td>
            <td style="text-align: center; font-size: 80%">
                {{ range .League_ids }}
                {{ if eq . 9 }}
                SS
                {{ end }}
                {{ if eq . 10 }}
                S
                {{ end }}
                {{ if eq . 20 }}
                A
                {{ end }}
                {{ if eq . 30 }}
                B
                {{ end }}
                {{ if eq . 40 }}
                C
                {{ end }}
                {{ end }}
            </td>
            <td style="width: 30%; font-size: 80%">
                {{ $sp := ""}}
                {{ range .Tag_list}}
                {{ $sp }}
                {{ .Tag_name}}
                {{ $sp = "/" }}
                {{ end }}

            </td>
        </tr>
        {{ end }}
        <tr>
    </table>
    <p>
        {{ .ErrMsg}}
    </p>
    <p><a href="https://zenn.dev/">Zenn</a> - <a
            href="https://zenn.dev/chouette2100/books/d8c28f8ff426b7">SHOWROOMのAPI、その使い方とサンプルソース</a></p>
</body>

</html>