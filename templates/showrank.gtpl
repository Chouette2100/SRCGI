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
    <p>
    {{/*}}
    (外部リンク)<br>
        <a href="https://zenn.dev/">Zenn</a> - <a
            href="https://zenn.dev/chouette2100/books/d8c28f8ff426b7">SHOWROOMのAPI、その使い方とサンプルソース</a></p>
   <div style="text-indent: 2rem;"><a href="https://chouette2100.com/">記事/ソース/CGI一覧表</a>　（証明書の期限が切れてしまっていました。2022年8月29日、有効な証明書に切り替えました）</div>
    <p>-------------------------------------------------------------</p>
    (サイト内リンク)<br>
    <div style="text-indent: 2rem;"><a href="t008top">t008:開催中イベント一覧</a></div>
    <p>-------------------------------------------------------------</p>
    {{*/}}
    <button type="button" onclick="location.href='top'">Top</button>　
    <button type="button" onclick="location.href='currentevents'">開催中イベント一覧表</button>　
    <button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧表</button>　
    <button type="button" onclick="location.href='closedevents'">終了イベント一覧表</button>　
    <br>
    <br>
    <p>最近のイベントの獲得ポイント上位のルーム</p>
    イベント終了日時が2023年10月1日00時以前のイベントは一部をのぞいてデータを取得していません。
    <br>イベント最終結果（獲得ポイント）が発表されるルーム（＝ランキングイベントの上位30位まで）が対象です。
    <br>ただし、すべての対象ルームのデータが取得できているかについては検証が不十分です。
    <br>もし"漏れ"に気づかれたら、□□イベントの○○さんがいない、みたいな指摘をいただければうれしいです。
    <br>また「unknown」はジャンルのデータが対象としているテーブルにないという意味です。履歴データから特定できるものもあるかもしれませんが。
    <br>試行錯誤しながら作ってきたものなのでこのあたりはご理解・ご勘弁のほどをお願いします。
    <br>
    <br>なお、ルームのジャンルや名称はこの結果を表示した時点（またはそれに近い時点）のものであり、イベント開催時のものではありません。
    <br>（ジャンルやルーム名をイベント開催時のものにすることが可能か、またできるとしてそうすることに意味があるか、は検討中です）
    <br>
    <br>
    {{ $l := "SS-5" }}
    {{ $n := "SS-5" }}
    {{ $c := "lightyellow" }}
    {{ $i := 1 }}
    <table>
        <tr style="text-align: center">
            <td>No.</td>
            <td>SHOWランク</td>
            <td>ルーム名（イベント履歴へのリンク）</td>
            <td>ルームレベル</td>
            <td>フォロワー数</td>
            <td>ファン数<br>（今月）</td>
            <td>ファンパワー<br>（〃）</td>
            <td>ファン数<br>（前月）</td>
            <td>ファンパワー<br>（〃）</td>
            <td>next_score</td>
            <td>prev_score</td>
            <td>データ取得日時</td>
        </tr>
        {{ range .Userlist }}
        {{ $n = Showrank .Rank }}
        {{ if ne $n $l }}
            {{ $l = $n }}
            {{ if eq $c "lightyellow" }}
                {{ $c = "lightcyan" }}
            {{ else }}
                {{ $c = "lightyellow" }}
            {{ end }}
        {{ end }}
        <tr bgcolor="{{$c}}" >
            <td align="right">{{ $i }}</td>
            <td align="center">{{ $n }}</td>
            <td><a href="closedevents?userno={{ .Userno }}&mode=0&path=5">{{ .User_name }}</a></td>
            <td align="right">{{ Comma .Level }}</td>
            <td align="right">{{ Comma .Followers }}</td>
            <td align="right">{{ Comma .Fans }}</td>
            <td align="right">{{ Comma .FanPower }}</td>
            <td align="right">{{ Comma .Fans_lst }}</td>
            <td align="right">{{ Comma .FanPower_lst }}</td>
            <td align="right">{{ Comma .Inrank }}</td>
            <td align="right">{{ Comma .Iprank }}</td>
            <td>{{ FormatTime .Ts "2006-01-02 15:04" }}</td>
        </tr>
        {{ $i = Add $i 1 }}
 
        {{ end }}
    </table>
    {{/*
    <p>
        {{ .ErrMsg }}
    </p>
    */}}
</body>

</html>