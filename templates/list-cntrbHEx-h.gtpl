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
            <td><button type="button" onclick="location.href='top?eventid={{.Eventid}}'">イベントトップ</button></td>
            <td></td>
            <td><button type="button"
                    onclick="location.href='graph-total?eventid={{.Eventid}}&maxpoint={{.Maxpoint}}&gscale={{.Gscale}}'">獲得ポイントグラフ</button>
            </td>
            <td></td>
        </tr>
        <tr>
            <td><button type="button" onclick="location.href='list-last?eventid={{.Eventid}}'">直近の獲得ポイント</button></td>
            <td><button type="button" onclick="location.href='list-cntrb?eventid={{.Eventid}}&userno={{.Userno}}&ie={{ .Ie }}'">枠別貢献ポイント</button></td>
            <td></td>
            <td></td>
        </tr>
        </table>
    </table>
    <br>
    <p style="color: red;">リスナーさんの過去のイベントでの貢献ポイントの履歴です<br>
    このデータは獲得ポイントデータ取得時に枠別貢献ポイントの取得を指示した結果得られたデータのみを対象にしています<br>
    いくら貢献ポイントが大きくても枠別貢献ポイントが取得されていない場合はここには表示されません。<br>
    ニーズがあるようなら、ちゃんと、つまり可能なかぎり過去イベントの貢献ポイントを取得してやってみたいとは思っています。</p>
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
            <td align="center">“{{ .Listener }}” （ Tlsnid = {{ .Tlsnid }} ）</td>
        </tr>
    </table>
    <br>
    <table>
        <tr>
            <td width="400" align="left">
                {{ if ne .Tlsnid_b -1 }}
                <button type="button"
                    onclick="location.href='list-cntrbHEx?eventid={{.Eventid}}&userno={{.Userno}}&tlsnid={{.Tlsnid_b}}'">{{
                    .Listener_b }}</button>
                {{ else }}
                -----------
                {{ end }}
            </td>
            <td width="400" align="right">
                {{ if ne .Tlsnid_f -1 }}
                <button type="button"
                    onclick="location.href='list-cntrbHEx?eventid={{.Eventid}}&userno={{.Userno}}&tlsnid={{.Tlsnid_f}}'">{{
                    .Listener_f }}</button>
                {{ else }}
                -----------
                {{ end }}
            </td>
        </tr>
    </table>