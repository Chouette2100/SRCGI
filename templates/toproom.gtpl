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
    ジャンル変更の影響が大きすぎて手が回らないのでジャンルの選択機能は使用できなくしています。
    <br>
    <br>イベント終了日時が2023年10月1日00時以前のイベントは一部をのぞいてデータを取得していません。
    <br>イベント最終結果（獲得ポイント）が発表されるルーム（＝ランキングイベントの上位50位まで）が対象です。
    <br>ただし、すべての対象ルームのデータが取得できているかについては検証が不十分です。
    <br>もし"漏れ"に気づかれたら、□□イベントの○○さんがいない、みたいな指摘をいただければうれしいです。
    <br>また「unknown」はジャンルのデータが対象としているテーブルにないという意味です。履歴データから特定できるものもあるかもしれませんが。
    <br>試行錯誤しながら作ってきたものなのでこのあたりはご理解・ご勘弁のほどをお願いします。
    <br>
    <br>なお、ルームのジャンルや名称はこの結果を表示した時点（またはそれに近い時点）のものであり、イベント開催時のものではありません。
    <br>（ジャンルやルーム名をイベント開催時のものにすることが可能か、またできるとしてそうすることに意味があるか、は検討中です）
    <br>
    <br>
    <form>
    	<fieldset>
    <input type="submit" fromaction="toproom" formmethod="get" value="再表示">
    イベント終了日時が <input type="date" name="from" value="{{ FormatTime .From "2006-01-02" }}"> 午前0時から
    <input type="date" name="to" value="{{ FormatTime .To "2006-01-02" }}"> 午前0時までの中から
    上位<input style="width: 4em" type="number" name="olim" value="{{ .Olim }}"> ルームを表示する。
    {{/*
    <div>
    	<fieldset>
		<legend>ランキングの対象となるジャンルを選んでください</legend>

        {{ range .Genrelist }}
        <input type="checkbox" id="genre{{.Genre_id}}" name="genre{{.Genre_id}}" {{if .Checked }}checked{{end}} />
            <label for="genre{{.Genre_id}}">{{.Genre_name}}</label>
        {{ end }}
        </fieldset>
    </div>
    */}}
        </fieldset>
    </form>
    ※　結果が表示されるまで十数秒要します。<br>
    ※　最終結果発表日時はイベント終了日翌日のお昼頃ですので終了日時の上限は一昨日としています。
    <br><br>
    <table>
        <tr style="text-align: center">
            <td>獲得ポイント</td>
            <td>ルーム　リンク先は「終了イベント一覧」のルーム検索結果</td>
            <td>ジャンル</td>
            <td>イベント（イベント順位）　リンク先は「直近の獲得ポイント一覧」（表示に時間がかかるケースあり）</td>
            <td>イベント終了日時</td>
        </tr>
        {{ range .TopRoomList }}
        <tr>
            <td align="right">{{ Comma .Point }}</td>
            {{/*
            <td><a href="closedevents?userno={{ .Room_id }}&mode=0&path=5" target="_blank" rel="noopener noreferrer">{{ .Room_name }}</a></td>
            */}}
            <td><a href="closedevents?userno={{ .Room_id }}&mode=0&path=5">{{ .Room_name }}</a></td>
            <td>{{ .Genre }}</td>
            {{/*
            <td><a href="list-last?eventid={{ .Event_id }}&roomid={{.Room_id}}" target="_blank" rel="noopener noreferrer">{{ .Event_name }}</a>（{{ .Rank }}）</td>
            */}}
            <td><a href="list-last?eventid={{ .Event_id }}&roomid={{.Room_id}}">{{ .Event_name }}</a>（{{ .Rank }}）</td>
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