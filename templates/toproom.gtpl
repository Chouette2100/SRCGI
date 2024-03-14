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
    <br>
    <br>特定のジャンルでの順位やアイドルのジャンルは除いた順位がほしいと思われる方もいらっしゃると思いますが、現時点では取得しているデータに不備があり対応できません。
    <br>（何が問題かは下記の結果をご覧になればわかると思いますが、今後データの修復を行う予定はありますのでどうにかなるかもしれません）
    <br>なお、ルームのジャンルや名称はこの結果を表示した時点のものであり、イベント開催時のものではありません。
    <br>（ジャンルやルーム名をイベント開催時のものにすることは全く不可能ではなさそうなので検討はしています）
    <br>
    <br>
    <form>
    <input type="submit" fromaction="toproom" formmethod="get" value="再表示">
    イベント終了日時が <input type="date" name="from" value="{{ FormatTime .From "2006-01-02" }}"> 午前0時から
    <input type="date" name="to" value="{{ FormatTime .To "2006-01-02" }}"> 午前0時までの中から
    上位<input style="width: 4em" type="number" name="olim" value="{{ .Olim }}"> ルームを表示する。
    </form>
    ※　結果が表示されるまで十数秒要します。<br>
    ※　最終結果発表日時はイベント終了日翌日のお昼頃ですので終了日時の上限は一昨日としています。
    <br><br>
    <table>
        <tr style="text-align: center">
            <td>獲得ポイント</td>
            <td>ルーム</td>
            <td>ジャンル</td>
            <td>イベント（イベント順位）</td>
            <td>イベント終了日時</td>
        </tr>
        {{ range .TopRoomList }}
        <tr>
            <td align="right">{{ Comma .Point }}</td>
            <td><a href="https://www.showroom-live.com/room/profile?room_id={{ .Room_id }}" target="_blank" rel="noopener noreferrer">{{ .Room_name }}</a></td>
            <td>{{ .Genre }}</td>
            <td><a href="https://www.showroom-live.com/event/{{ .Event_id }}" target="_blank" rel="noopener noreferrer">{{ .Event_name }}</a>（{{ .Rank }}）</td>
            <td>〜{{ FormatTime .Event_endtime "2006-01-02" }}</td>
        </tr>
        {{ end }}
        <tr>
    </table>
    {{/*
    <p>
        {{ .ErrMsg }}
    </p>
    */}}
</body>

</html>