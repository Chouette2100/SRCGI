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

    <p style="color:green;">更新や障害の状況については<a
            href="https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/807ad9">だいじなお知らせ</a>をご覧ください！
    <p style="color:green;">使用中にヘンな現象が発生したら<a
            href="https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/03187b">トラブルシューティング</a>をご覧ください！</p>



    <p style="color:green;">
        このWebサーバ(＝CGIとしても動作可能)とそれに関わるプログラムのースコードをGithubで公開しています。
        <br>詳しくは<a href="https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/4fccae">『SHOWROOM イベント
            獲得ポイント一覧』関連のソースの公開について(1)</a>をごらんください。
    </p>

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
    <p style="padding-left:2em">
    <a href="/listgs">「SHOWROOMライバー王決定戦」ギフトランキング </a>  （新規機能・テスト中）
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
    <p style="padding-left:4em">
    <a href="/experimental">保守・資料・実験</a>
    </p>