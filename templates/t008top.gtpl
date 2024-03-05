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

    <br>
    (外部リンク)<br>
        <a href="https://zenn.dev/">Zenn</a> - <a
            href="https://zenn.dev/chouette2100/books/d8c28f8ff426b7">SHOWROOMのAPI、その使い方とサンプルソース</a></p>
   <div style="text-indent: 2rem;"><a href="https://chouette2100.com/">記事/ソース/CGI一覧表</a>　（証明書の期限が切れてしまっていました。2022年8月29日、有効な証明書に切り替えました）</div>
    <p>-------------------------------------------------------------</p>
    (サイト内リンク)<br>
    <div style="text-indent: 2rem;"><a href="t009top">t009:配信中ルーム一覧</a></div>
    <p>-------------------------------------------------------------</p>
    <p>開催中イベント一覧表</p>
    <div style="text-indent: 2rem;">イベント数： {{ .Totalcount }} （{{ UnixTimeToStr $tn }}）</div>
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
            <td  style="text-align: center; font-size: 80%">
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
</body>

</html>