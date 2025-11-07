<html>

<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
    <style type="text/css">
        {{/*
        th,
        td {
            border: solid 1px;
        }

        table {
            border-collapse: collapse;
            /.
            width: 100%;
            ./
        }
        */}}

        .bgct {
            background-color: paleturquoise;
        }
    </style>
</head>


<body>
    <table style="border-width: 0;">
            <tr>
    <td><button type="button" onclick="location.href='top'">Top</button></td>
    <td><button type="button" onclick="location.href='currentevents'">開催中イベント一覧表</button></td>
    <td><button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧表</button></td>
    <td><button type="button" onclick="location.href='closedevents'">終了イベント一覧表</button></td>
            </tr>
            <tr>
      <td><button type="button" onclick="location.href='eventtop?eventid={{.Eventurl}}'">イベントトップ</button></td>
      <td></td>
      <td><button type="button" onclick="location.href='list-last?eventid={{.Eventurl}}'">直近の獲得ポイント</button></td>
      {{/*
      <td><button type="button" onclick="location.href='graph-total?eventid={{.Eventurl}}&maxpoint={{.Maxpoint}}&gscale={{.Gscale}}'">獲得ポイントグラフ</button></td>
      */}}
      <td><button type="button" onclick="location.href='graph-total?eventid={{.Eventurl}}'">獲得ポイントグラフ</button></td>
        </tr>
    </table>
    <br>
    <br>

    <p>終了イベント・最終結果</p>
    <br>

    {{ if eq .Eventid 0 }}


    <p style="padding-left:2em;">
        イベントを指定して、参加ルーム一覧を表示する。
    </p>


    <form>
        <p style="padding-left:2em;">
            参加ルーム一覧を表示したいイベントを選択し、「表示する」ボタンをクリックしてください。
        </p>
        <p style="padding-left:2em">
            <select name="eventid">
                {{ range .Eventlist }}
                <option value="{{ .Eventurl }}" zgotmplz="">{{ .Eventname }}</option>
                {{ end }}
            </select>
            <span style="padding: 2em;"></span>
            <input type="number" name="ib" min="1" max="9999" value="1"> 番目から
            <input type="number" name="ie" min="1" max="9999" value="10"> 番目まで
            <input type="submit" value="表示する" formaction="apieventroomlist" formmethod="GET"></span>
        </p>
        <p style="padding-left:2em;">
            表示ルーム数が多ければ多いほど表示に時間がかかります。10番目くらいまでにしてレスポンスをチェックしてください。
            <br>「9999番目まで」など大きな値にすればすべての参加ルームが表示されます（参加ルーム数が多いときは数十秒から数分レスポンスがなくなります）
        </p>
    </form>

    <br>
    <hr>
    <br>
    <p style="padding-left: 2em;">
        【解説】
    </p>

    <p style="padding-left: 2em;">
        イベントのEvent_id（eventid）を指定して実行することによりイベントに参加するルームの一覧を表示します。
        <br>
        <br>
    <p style="padding-left: 4em;">ソースをビルドしてローカルでWebサーバーとして起動した場合。
    <pre style="font-size: 1.2em;"><code>
            http://localhost:8080/apieventroomlist?eventid=29553&ib=1&ie=10
        </code></pre>
    </p>
    <p style="padding-left: 4em;">VPS(Ubuntu) Apache2のWebサーバーでCGIとして動いているサンプル
    <pre style="font-size: 1.2em;"><code>
            https://chouette2100.com:8443/cgi-bin/test/apieventroomlist?eventid=29553&ib=1&ie=10
        </code></pre>
    </p>
    <p style="padding-left: 2em;">
        <br>参加ルームリストのib番目からie番目までのルームが表示されます。ibあるいはieは省略可能で、省略すると ib=1, ie=10 として扱われます。
        <br>
        <br>eventidはイベントを識別する５桁程度の整数です。イベントのURLの最後にある英数字記号からなる文字列 Event_url_key とは異なるものです。
        <br>通常eventidを目にすることはなくeventidがわからない場合もあると思われるのでeventidを指定せずに実行すると、開催中イベント選択画面を表示するようにしました。
        <br>この画面がそれです。
    </p>
    <pre style="font-size: 1.2em;"><code>
            http://localhost:8080/apieventroomlist
            https://chouette2100.com:8443/cgi-bin/test/apieventroomlist
        </code></pre>
    <p style="padding-left: 2em;">あるいは「<a
            href="t008top">開催中イベント一覧表を表示する</a>」の「参加ルーム一覧へのリンク」からイベント参加ルーム一覧が表示できるようにしてあります。
        <br>

        <br>参加ルーム一覧へのリンクからは開催中のイベントのルームリストしか表示できませんが、eventidを指定すれば開催前、開催後のリストも見ることができます。
        <br>（ただし、表示内容には制約があります）
    </p>

    {{ else }}

    <p style="padding-left: 2em;">
        イベント： <a href="https://www.showroom-live.com/event/{{ .Eventurl }}">{{ .Eventname }}</a>
        （ eventid = {{ .Eventid }} ）
    </p>
    {{ if ne .Msg "" }}
    <p>
        Msg = {{ .Msg }}
    </p>
    {{ end }}

    <form>
        <p style="padding-left:4em">
            このページはブックマーク可能です（上のURLに対して同一のイベントに対する結果が得られます）
            {{/*}}
            表示ルーム数を変更するときは、表示数を入力し「再表示する」ボタンをクリックしてください。
            <br>
            <br>
            <input type="number" name="ib" min="1" max="9999" value="{{ .Ib }}"> 番目から
            <input type="number" name="ie" min="1" max="9999" value="{{ .Ie }}"> 番目まで
            <input type="hidden" name="eventid" value="{{ .Eventid }}">
            <input type="hidden" name="eventurlkey" value="{{ .Eventurl }}">
            <span style="padding-left:2em"><input type="submit" value="再表示する" formaction="eventroomlist"
                    formmethod="GET"></span>
                {{*/}}
        </p>
    </form>

    <p style="padding-left: 2em;">
    <table style="border: solid 1px; border-collapse: collapse;">
        <tr style="border-width: 1;">
            <td style="border: solid 1px;">順位</td>
            <td style="border: solid 1px;">ルーム名（ルーム状況へのリンク）</td>
            <td style="border: solid 1px;">獲得ポイント</td>
            {{/*}}
            <td></td>
            {{*/}}
            <td style="border: solid 1px;">Prof. / Live / FC / Cnt.</td>
            <td style="border: solid 1px;">配信中？</td>
            <td style="border: solid 1px;">次回配信時刻</td>
            <td style="border: solid 1px;">枠別リスナー別貢献pt</td>
        </tr>

        {{ $e := .Eventurl }}
        {{ $r := .Roomid }}

        {{ with .Roomlistinf }}
        {{ range .RoomList }}
        <tr {{ if eq .Room_id $r }} class=bgct {{ end }}>
            <td style="text-align: right; border: solid 1px;">
                {{ if ne .Rank -1 }}
                {{ .Rank }}
                {{ end }}
            </td>
            <td style="border: solid 1px;"><a href="/apiroomstatus?room_url_key={{ .Room_url_key }}" target="_blank" rel="noopener noreferrer">{{
                    .Room_name }}</a></td>
            <td style="text-align: right; border: solid 1px;">
                {{ if ne .Point -1 }}
                {{ Comma .Point }}
                {{ end }}
            </td>
                {{/*}}
            <td style="text-align: right; border: solid 1px;">
                {{ if ne .Rank 1 }}
                {{ Comma .Gap }}
                {{ end }}
            </td>
                {{*/}}
            <td style="border: solid 1px;"><a href="https://www.showroom-live.com/room/profile?room_id={{ .Room_id }}" target="_blank"
                    rel="noopener noreferrer">Prof.</a> /
                <a href="https://www.showroom-live.com/{{ .Room_url_key }}" target="_blank"
                    rel="noopener noreferrer">Live</a> /
                <a href="https://www.showroom-live.com/room/fan_club?room_id={{ .Room_id }}" target="_blank"
                    rel="noopener noreferrer">FC</a> /
                <a href="https://www.showroom-live.com/event/contribution/{{ $e }}?room_id={{ .Room_id }}"
                    target="_blank" rel="noopener noreferrer">Cnt.</a>
            </td>
            <td style="text-align: center; border: solid 1px;">
                {{ if eq .Islive true}}
                {{ UnixtimeToTime .Startedat "15:04"}} 〜
                {{ end }}
            </td>
            <td style="text-align: center; border: solid 1px;">
                {{ if ne .Nextlive 0 }}
                {{ UnixtimeToTime .Nextlive "01/02 15:04"}}
                {{ end }}
            </td>
            <td style="text-align: center; border: solid 1px;">
                {{ if eq .Isofficial true }}
                <a href="/cgi-bin/SC1/SC1/list-cntrb?eventid={{$e}}&userno={{.Room_id}}">表示</a>
                {{ end }}
            </td>
        </tr>
        {{ end }}
        {{ end }}
    </table>
    </p>

    {{ end }}
    <br>
    <hr>
    <br>
    {{/*
    {{ template "footer" }}
    */}}
</body>
</html>