<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
<html>

<body>
    <table>
        <tr>
            <td><button type="button" onclick="location.href='top'">トップ</button>　</td>
            <td><button type="button" onclick="location.href='currentevents'">開催中イベント一覧</button></td>
            <td><button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧</button></td>
            <td><button type="button" onclick="location.href='closedevents'">終了イベント一覧</button></td>
        </tr>
        <tr>
            <td><button type="button" onclick="location.href='eventtop?eventid={{.Eventid}}'">イベントトップ</button></td>
            <td></td>
            <td><button type="button"
                    onclick="location.href='graph-total?eventid={{.Eventid}}&maxpoint={{.Maxpoint}}&gscale={{.Gscale}}'">獲得ポイントグラフ</button>
            </td>
            <td></td>
        </tr>
        <tr>
            <td><button type="button" onclick="location.href='list-last?eventid={{.Eventid}}'">直近の獲得ポイント</button></td>
            <td><button type="button"
                    onclick="location.href='list-cntrb?eventid={{.Eventid}}&userno={{.Userno}}&ie={{ .Ie }}'">枠別貢献ポイント</button></td>
            <td></td>
            <td></td>
        </tr>
    </table>
    <p>枠別貢献ポイント一覧
        {{ if eq .Srt 1 }}
        （増分順）
        {{ else }}
        （累計順）
        {{ end }}
        （審査ポイント等は除く）
    </p>
    <table>
        <tr>
            <td align="center"><a
                    href="https://www.showroom-live.com/event/{{.Eventid}}">{{.Eventname}}</a>（{{.Eventid}}）</td>
        </tr>
        <tr>
            <td align="center">{{.Period}}</td>
        </tr>
        <br>
        <tr>
            <td align="center"><a
                    href="https://www.showroom-live.com/room/profile?room_id={{.Userno}}">{{.Username}}</a>（{{.Userno}}）
            </td>
        </tr>
        <br>
        <tr>
            <td align="center">{{ .S_stime }} ----- {{ .S_etime }} （ {{ .Ifrm1 }} ）</td>
        </tr>
    </table>
    <br>
    {{ if eq .Srt 1 }}
    <form>
        目標値（推定）　
        <input value="{{ .Target }}" name="target" type="number">　
            <input type="submit" value="目標値を変更する" formaction="list-cntrbS" formmethod="GET">
            <input type="hidden" value="{{ .Eventid }}" name="eventid">
            <input type="hidden" value="{{ .Userno }}" name="userno">
            <input type="hidden" value="{{ .Ifrm }}" name="ifrm">
            <input type="hidden" value="D" name="sort">
    </form>
    {{ end }}
    <br>
    <table>
        <tr>
            <td width="400" align="left">
                {{ if ne .Ifrm_b -1 }}
                {{ if eq .Srt 1 }}
                <button type="button"
                    onclick="location.href='list-cntrbS?eventid={{.Eventid}}&userno={{.Userno}}&ifrm={{.Ifrm_b}}&sort=D'">一つ前の配信枠へ</button>
                {{ else }}
                <button type="button"
                    onclick="location.href='list-cntrbS?eventid={{.Eventid}}&userno={{.Userno}}&ifrm={{.Ifrm_b}}'">一つ前の配信枠へ</button>
                {{ end }}
                {{ else }}
                -----------
                {{ end }}
            </td>
            <td width="400" align="right">
                {{ if ne .Ifrm_f -1 }}
                {{ if eq .Srt 1 }}
                <button type="button"
                    onclick="location.href='list-cntrbS?eventid={{.Eventid}}&userno={{.Userno}}&ifrm={{.Ifrm_f}}&sort=D'">次の配信枠へ</button>
                {{ else }}
                <button type="button"
                    onclick="location.href='list-cntrbS?eventid={{.Eventid}}&userno={{.Userno}}&ifrm={{.Ifrm_f}}'">次の配信枠へ</button>
                {{ end }}
                {{ else }}
                -----------
                {{ end }}
            </td>
        </tr>
    </table>