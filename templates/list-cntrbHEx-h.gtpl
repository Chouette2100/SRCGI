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
    <p style="color: blue;">リスナーさんの過去のイベントでの貢献ポイントの履歴です<br>
    データは貢献100位までが発表されるイベント貢献ランキングから取得しています。<br>
    つまり貢献ポイントが1ptでもこのリストにあることもありますが、貢献ポイント10,000ptでもないことがあります。<br>
    </p>
    <p style="color: blue;">
    終了したイベントのこのデータは特定の日時（注１）以後では、イベント獲得ポイントの取得を指定したすべてのイベント・ユーザーついて表示されます。<br>
    通常獲得ポイントが上位・中位のルームは自動的にイベント獲得ポイント取得の対象となりますので、このデータも得られます。<br>
    この場合リスナー別/枠別貢献ポイントの取得を指定したか否かは無関係です。<br>
    開催中のイベントについてはこの機能は現時点ではリスナー別/枠別貢献ポイントの取得を指定しておく必要がありますが<br>
    今大急ぎで開催中イベントでも使えるように改修中です。<br>
    <BR>
    特定の日時（注１）以前では、リスナー別/枠別貢献ポイントの取得を指定したイベント・ルームのみが表示の対象となります。<br>
    注１　「特定の日時」とは現時点では2024年9月29日です。<br>
    この日時はなんらかの事情あるいはご要望により今後数ヶ月から数年さかのぼった日時に変更する可能性はあります。
    </p>
    <p style="color: blue;">
    なお、5月15日以後については最終枠の枠別貢献ポイントが取得できなくなっていましたが、これについては6月9日18時までにすべて修復しました。
    </p>
    <br>
    <div style="background: #ffe4e1; border: #ffe4e1 solid 2px; font-size: 100%; padding: 20px; width: 50em;">
    ※　イベント・ルームに対する貢献ポイントの月単位のランキングを作ってみました。<br><br>
    　　<a href="m-cntrbrank-listener">月別イベント・リスナー貢献ポイントランキング</a>（結果が表示されるまで数十秒かかります）<br><br>
    ※　これも貢献ランキングですがどのリスナーさんがどのルームを応援しているかわかるように作ったものです。<br><br>
    　　<a href="m-cntrbrank-Lg">月別貢献ポイントランキング（リスナー/ルーム）</a>（結果が表示されるまで数十秒かかります）
    </div>
    <br>
    <table>
    {{/*
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
    */}}
        <tr>
            <td align="center">“{{ .Listener }}” （ Tlsnid = {{ .Tlsnid }} ）</td>
        </tr>
    </table>
    <br>
    {{/*
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
    */}}

{{/*
    {{ $n := .Stnow }}
*/}}
   
