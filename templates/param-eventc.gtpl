<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<head>
<style type="text/css">
.fblue { color:blue;font-style:italic;font-size:1.2em;font-weight:bold;margin: 10px; }
</style>
</head>
<body>
<table>
    <tr>
        <td><button type="button" onclick="location.href='top'">トップ</button>　</td>
        <td><button type="button" onclick="location.href='currentevents'">開催中イベント一覧</button></td>
        <td><button type="button" onclick="location.href='scheduledevents'">開催予定イベント一覧</button></td>
        <td><button type="button" onclick="location.href='closedevents'">終了イベント一覧</button></td>
    </tr>
    <tr>
        <td><button type="button" onclick="location.href='top?eventid={{.Event_ID}}'">イベントトップ</button></td>
        <td><button type="button" onclick="location.href='list-last?eventid={{.Event_ID}}'">直近の獲得ポイント</button></td>
        <td><button type="button"
                onclick="location.href='graph-total?eventid={{.Event_ID}}&maxpoint={{.Maxpoint}}&gscale={{.Gscale}}'">獲得ポイントグラフ</button>
        </td>
        <td></td>
    </tr>
</table>
<br><br>
<p style="padding-left:2em">
イベント設定の変更（確認）
<p style="padding-left:4em">
イベントに関する設定値を変更しました。もう一度設定を変更するときは「設定変更に戻る」ボタンを押してください。<br>
<form>
<table>
<tr><td style="width:4em"></td><td>イベントのID</td><td><input type="hidden" name="eventid" value="{{.Event_ID}}" >{{.Event_ID}}</td></tr>
<tr><td style="width:4em"></td><td>イベント名</td><td>{{.Event_name}}</td></tr>
<tr><td style="width:4em"></td><td>イベント期間</td><td>{{.Period}}</td></tr>
<tr><td style="width:4em"></td><td>イベント参加ルーム数</td><td>{{.NoEntry}}（最新のデータでない可能性あり）</td></tr>
<tr><td style="width:4em"></td><td><label>ＤＢに登録する順位の範囲</td><td><span class="fblue">{{.Fromorder}}</span>位から<span class="fblue">{{.Toorder}}</span>位まで</td></tr>
<tr><td></td><td>獲得ポイントデータ取得のタイミング</td><td>毎時<span class="fblue">{{.Modmin}}</span>分<span class="fblue">{{.Modsec}}</span>秒から<span class="fblue">{{.Intervalmin}}</span>分おきに取得する。（不適切な設定値は修正されています）</td></tr>
<tr><td></td><td>日々の獲得ポイントのリセット時刻</td><td>毎日<span class="fblue">{{.Resethh}}</span>時<span class="fblue">{{.Resetmm}}</span>分</td></tr>
<tr><td></td><td>ポイント差の基準とする配信者</td><td><span class="fblue">{{.Nobasis}}</span></td></tr>
<tr><td></td><td>目標ポイント</td><td><span class="fblue">{{.Target}}</span></td></tr>
<tr><td></td><td>最大表示数</td><td><span class="fblue">{{.Maxdsp}}</span></td></tr>
<tr><td></td><td>カラーマップ</td><td><span class="fblue">{{.Cmap}}</span></td></tr>
<tr><td></td><td></td><td align="right"><input type="submit" value="設定変更に戻る" formaction="param-event" formmethod="POST" style="background-color: khaki"></td></tr>
<tr><td></td><td></td><td align="right"><input type="submit" value="終了" formaction="top?eventid={{.Event_ID}}" formmethod="POST" style="background-color: khaki"></td></tr>
</table>
</form>
<p style="padding-left:4em">
「ＤＢに登録する順位の範囲」というのは、「登録」を実行した時点で指定した範囲にある順位のルームを登録する、ことを意味します。一度登録されたルームは削除されませんので、順位の入れ替わりがあったあと「登録」を行うとDBに登録されているルーム数は増えていきます。DBに登録されるルーム数には制限はありません。
</p>

</body>
</html>
