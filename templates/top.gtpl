<!DOCTYPE html>
<meta name="viewport" content="width=device-width, initial-scale=1.0" charset="UTF-8">
<html>

<head>
    <style>
        .p1 {
            white-space: pre-wrap;
            margin-left: 25;
        }
    </style>
</head>

<body>
    <H2>SHOWROOM イベント 獲得ポイント一覧</H2>

    {{/*
    <p style="color:red;">現在下記のURLで新しいバージョンのテスト中です。 </p>
    <p style="padding-left:2em">
        <a href="http://21.matrix.jp/TEST/CGI.cgi/top"> http://21.matrix.jp/TEST/CGI.cgi
            /top</a>
    </p>
    */}}
    {{/*
    <p style="color:red;">現在システムの更新作業中です。しばらくお待ちください。</p>
    <p style="color:red;">更新や障害の状況については<a
            href="https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/807ad9">だいじなお知らせ</a>をご覧ください！</p>
    */}}
    <p style="color:green;">更新や障害の状況については<a
            href="https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/807ad9">だいじなお知らせ</a>をご覧ください！
        {{/*　　サーバ移行のため<span style="color:red;">月曜日の早朝に1時間ほどサーバを停止</span>します。*/}}</p>
    <p style="color:green;">使用中にヘンな現象が発生したら<a
            href="https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/03187b">トラブルシューティング</a>をご覧ください！</p>
        {{/*
    <p style="color:red;">
        『STU48 × 「naive（ナイーブ）」 PRアンバサダー決定オーディション』は獲得ポイントの時間推移が取得できないイベントです。結果は毎日一回イベントページに発表されています。</p>
    <p style="color:red;">
        2023-04-20「CGIの更新を行いました。今回のバージョンではブロックイベントへの暫定的な対応がされています。詳細は「<a
            href="https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/807ad9">だいじなお知らせ</a>」をご覧ください。
    </p>
    */}}

    {{/*
    <p style="color:blue;">
        特定のルームを獲得ポイント取得の対象とするときは『(DB登録済み)イベント参加ルーム一覧（確認・編集）』の下の方にある
        『一覧にないルームの追加』の機能を使います。
        <br>ユーザIDはプロフィールやファンルームのURLの最後にある５桁あるいは６桁の整数です。
    </p>
    */}}

    {{/*
    <p style="color:blue;">
        『直近の獲得ポイント一覧』の画面の表にあるルームの配信画面へのリンク『LIVE』について正常に動作しないケースがありますがこれについてはじょじょに解決していく予定です。状況については随時<a
            href="https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/807ad9">だいじなお知らせ</a>でお知らせします。</p>
    <p style="color:blue;">午前11時30分頃獲得ポイントが減る（＝増分がマイナスになる）ことがあります。これは
        <br>　　『重複アカウントによる応援は禁止です。重複アカウントによる応援ポイント分は発覚次第、減算を行います。』
        <br>という告知に該当するものと思われます。不審に思われる方がいらっしゃるかもしれないので念のため書いておきます。
    </p>
    */}}

    <p style="color:green;">
        このWebサーバ(＝CGIとしても動作可能)とそれに関わるプログラムのースコードをGithubで公開しています。
        <br>詳しくは<a href="https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/4fccae">『SHOWROOM イベント
            獲得ポイント一覧』関連のソースの公開について(1)</a>をごらんください。
    </p>

    {{/*
    <p style="color:red;">ここはテスト中の新しいバージョンです。現在のバージョンは下記URLから。</p>
    <p style="padding-left:2em">
        <a href="http://21.matrix.jp/SHOWROOM/CGI.cgi/top"> http://21.matrix.jp/SHOWROOM
            /CGI.cgi
            /top</a>
    </p>
    */}}
    <br>
    <p style="padding-left:2em">
        <a href="currentevents">開催中イベント一覧</a> （1.獲得ポイントデータを見たいイベントを選択する。 2.獲得ポイントデータを取得したいイベントを選択する）
    </p>
    <p style="padding-left:2em">
        <a href="scheduledevents">開催予定イベント一覧</a> （1.獲得ポイントデータを取得の予約を行うイベントを選択する）
    </p>
    <p style="padding-left:2em">
        <a href="closedevents">終了イベント一覧</a> （1.過去のイベントを検索し結果を参照する<span style="color: red">、ルームによる検索は開催中イベントも検索対象</span>）
    </p>
    <br>
    <p style="padding-left:2em">
        <a href="currentdistrb">配信中ルーム一覧</a> （配信開始から間もないルームの一覧が表示されます（ジャンル問わず30ルーム表示されるようにしました））
    </p>
    <p style="padding-left:2em">
        <a href="scheduledeventssvr">開催予定イベント一覧（内容詳細）</a> （開催が予定されているイベントの詳細が表示されます）
    </p>
    <p style="padding-left:2em"> -----------------------------------<br>
    {{/*
    <br>実験中
    <br>
    <br><a href="toproom" style="padding-left:2em">最近のイベントの獲得ポイント上位のルーム</a>（結果が表示されるまで十数秒要します）
    */}}
    <p style="padding-left:2em">
    <a href="/listgs">「SHOWROOMライバー王決定戦」ギフトランキング </a>  （新規機能・テスト中）
    {{/*
    <br>
    <span style="color: red">（ジャンルランキング、決勝Sリーグ、決勝Rリーグについてはギフトの貢献ランキング200位までを追加）</span>
    */}}
    </p>
    <form action="/showrank" method="get" >
        <p style="padding-left:2em">
        <a href="showrank">SHOWランクが上位のルーム</a>
            　　　<label for="name">一覧に追加するルーム(ルーム番号をカンマ区切りで入力): </label>
            <input type="text" name="unlist" id="unlist" />
            <input type="submit" value="ルームを追加して一覧を表示する!" />
        </form>
    </p>
    <p style="padding-left:2em; background:yellow; width:54em;">
    <a href="m-cntrbrank-listener">月別イベント・リスナー貢献ポイントランキング</a>（新規機能・結果が表示されるまで数十秒かかります）
    </p>
    <p style="padding-left:2em; background:yellow; width:54em;">
    <a href="m-cntrbrank-Lg">月別貢献ポイントランキング（リスナー/ルーム）</a>（どのリスナーさんがどのルームを応援しているかわかるように作りました）
    </p>
    <p style="padding-left:2em">
    <a href="toproom">最近のイベントの獲得ポイント上位のルーム</a>（結果が表示されるまで30秒以上要します）
    </p>
    <p style="padding-left:2em">
    最近のイベントの貢献ポイント上位のリスナー
    </p>
    {{/*}}
    <p style="padding-left:4em">
        <a href="#newevent">獲得ポイントデータを取得するイベントのイベントID(Event_url_key)による</a>新規登録（ページ内）
    </p>
    <br>
    <p>イベント選択（最近のイベントから選ぶ...イベント名をクリックしてください）</p>

    <table border="1">
        <tr>
            <td>状態</td>
            <td style="border-right-style:none;">開始</td>
            <td style="border-right-style:none;border-left-style:none;"> - </td>
            <td style="border-left-style:none;">終了</td>
            <td style="border-right-style:none;"> 　　　　　　 イベント名をクリックをクリックすると「イベントトップ（このイベントの表示項目選択）」が表示されます。
                <br> 　　　　　　 "一覧"、"グラフ"のリンクは「直近の獲得ポイントリスト」、「獲得ポイントの推移グラフ」へのショートカットです。
            </td>
            <td style="border-left-style:none;"></td>
            <td>ベース</td>
            <td>Mm</td>
            <td>Ms</td>
        </tr>
        {{ range . }}
        <tr>
            <td>{{ .Status }}</td>
            <td style="border-right-style:none;">{{ .S_start }}</td>
            <td style="border-right-style:none;border-left-style:none;"> - </td>
            <td style="border-left-style:none;">{{ .S_end }}</td>
            <td style="border-right-style:none;">
                <a href="list-last?eventid={{ .EventID }}">一覧</a>
                <a href="graph-total?eventid={{ .EventID }}&maxpoint={{ .Maxpoint }}&gscale={{.Gscale}}">グラフ</a>
                <a href="top?eventid={{ .EventID }}">{{ .EventName }}</a>
            </td>
            <td style="border-left-style:none;"><a
                    href="https://www.showroom-live.com/event/{{ .EventID }}">イベントページへ</a></td>

            <td>
                {{ if ne .Pntbasis 0 }}
                {{ .Pbname }}
                {{ else }}
                ------
                {{ end }}
            </td>
            <td align="right">
                {{ .Modmin }}
            </td>
            <td align="right">
                {{ .Modsec }}
            </td>
        </tr>
        {{ end }}
    </table>
    {{*/}}