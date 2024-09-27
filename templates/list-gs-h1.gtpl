<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0"  charset="UTF-8">
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
     <td><button type="button" onclick="location.href='listvgs'">ファンランキング</button></td>
     <td></td>
     <td></td>
     <td><button type="button" onclick="location.href='graphgs?campaignid={{.Campaignid}}&giftid={{.Grid}}'">ギフトランキンググラフ</button></td>
    </tr>

</table>
<br>
<p><a href="{{ .Url }}">{{ .Campaignname }}</a>（{{ .Campaignid }}）</p>
<p>　　{{ .Grname }}（{{ .Grid }}）</p>
<br>
<form>
　　　　ギフト種別を選択する
{{ $grid :=  .Grid }}
<select name="giftid" type="text">
{{ range .GiftRanking }}
    {{ if eq .Grid $grid }}
        <option selected value="{{ .Grid }}">{{ .Grname }}</option>
    {{ else }}
        <option value="{{ .Grid }}">{{ .Grname }}</option>
    {{ end }}
{{ end }}
</select>

　　横(時刻)表示数
<input value="{{ .Maxacq }}" name="maxacq" type="number" size="5" min="1" max="15">

　　縦(ルーム)表示数
<input value="{{ .Limit }}" name="limit" type="number" size="5" min="20" max="500">

 　　<input type="submit" value="この条件で再表示" formaction="listgs" formmethod="GET">

</form>
<br>
<table>
    <tr>
        <td>
        {{ if ne .Nft -1 }}
            <button type="button" onclick="location.href='listgs?giftid={{.Grid}}&ie={{.Nft}}&limit={{.Limit}}&maxacq={{.Maxacq}}'">先頭に戻る</button>
        {{ else }}
            -----------
        {{ end }}
        </td>
        <td>
        {{ if ne .Npb -1 }}
            <button type="button" onclick="location.href='listgs?giftid={{.Grid}}&ie={{.Npb}}&limit={{.Limit}}&maxacq={{.Maxacq}}'">１ページ戻る</button>
        {{ else }}
            -----------
        {{ end }}
        </td>
        <td>
        {{ if ne .N1b -1 }}
            <button type="button" onclick="location.href='listgs?giftid={{.Grid}}&ie={{.N1b}}&limit={{.Limit}}&maxacq={{.Maxacq}}'">一枠分戻る</button>
        {{ else }}
            -----------
        {{ end }}
        </td>
        <td>
        {{ if ne .N1f -1 }}
            <button type="button" onclick="location.href='listgs?giftid={{.Grid}}&ie={{.N1f}}&limit={{.Limit}}&maxacq={{.Maxacq}}'">一枠分進む</button>
        {{ else }}
            -----------
        {{ end }}
        </td>
        <td>
        {{ if ne .Npf -1 }}
            <button type="button" onclick="location.href='listgs?giftid={{.Grid}}&ie={{.Npf}}&limit={{.Limit}}&maxacq={{.Maxacq}}'">１ページ進む</button>
        {{ else }}
            -----------
        {{ end }}
        </td>
        <td>
        {{ if ne .Nlt -1 }}
            <button type="button" onclick="location.href='listgs?giftid={{.Grid}}&ie={{.Nlt}}&limit={{.Limit}}&maxacq={{.Maxacq}}'">最後に進む</button>
        {{ else }}
            -----------
        {{ end }}
        </td>
    </tr>
</table>
