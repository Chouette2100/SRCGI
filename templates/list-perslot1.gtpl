<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
<html>
<body>
<p id="Top">
<br>
</p>
<button type="button" onclick="location.href='top'">top</button>　
<button type="button" onclick="location.href='currentevent'">開催中イベント一覧表</button>　
<button type="button" onclick="location.href='top?eventid={{.Event_ID}}'">このルームの表示項目選択</button><br>
<h2>配信枠毎の獲得ポイント</h2>
<p style="padding-left:2em;color:blue">
一回の配信でも配信中に10分間程度獲得ポイントの変化がない場合は
複数回の配信とみなされることがあります。<br><br>
イベント全期間に渡ってデータを取得していない場合は以下のような現象が起きることがあります。<br>
・配信開始・終了時刻が"n/a"と表示される。<br>
・複数の配信が一つの配信として扱われる（獲得ポイントが合算される）<br>
<span style="color:crimson">・（表示上の）最初の枠のデータが正しくない。</span>（この場合でも累積獲得ポイントは正しい）<br><br>
<span style="color:crimson">最終枠の獲得ポイントはイベント最終日翌日の13時30分に最終結果に更新されます。</span><br><br>
※　配信開始時刻、配信終了時刻は推定値です。<br>
　　配信の開始時刻は取得できるのですが、現在それはやっていません。したがって配信開始時刻は実際のものとは異なります。<br><br>

</p>
<p style="padding-left:2em">
<a href="https://www.showroom-live.com/event/{{.Event_ID}}">{{ .Event_name }}</a>（{{.Event_ID}}）<br>
{{ .Period }}<br>
</p>
