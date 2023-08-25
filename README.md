# linebot-otp



## SignUpの流れ

![image](https://github.com/shunsuke-tamura/linebot-otp/assets/74412997/61f383cf-eeb5-435b-8d9e-cb22c6b50167)

1. userが入力したid, name, passwordをサーバに送る
2. sessionIdとしてuuidを生成し、それをキーに送られてきた情報をredisに保存
3. frontendへsessionIdをレスポンスとして送る(cookieに保存される)
4. userが自分の端末でOTP送信用のLINEBotに”登録”と送る
5. サーバは受け取ったメッセージからLINEのUIDを取得し、生成したOTPをキーとしてLineUidをredisに保存する
6. 保存したOTPをLINEBotの返信として送信する
7. userはLINEBotから送信されたOTPをfrontendへ入力する
8. frontendからOTPをserverに送信する
9. サーバは受け取ったOTPをキーとするデータをredisから取得し、それに紐づいたLineUidを取得する
10. サーバはcookieからsesionIdを読み取り、sessionIdをキーとするデータをredisから取得し、それに紐づいたユーザー情報を取得する
11. サーバは取得したユーザー情報とLineUidをマージし、DBへUserをInsertする
12. 認証用のsessionIdを再び生成し、sessionIdをキーとしてUserをredisへ保存する
13. キーに使用したsessionIdをfrontendへレスポンスとして送る(cookieに保存される)

## SignInの流れ
 
![image](https://github.com/shunsuke-tamura/linebot-otp/assets/74412997/9e4563a8-84d3-42fb-b5d5-fc4fd81a6f3c)

 
1. userが入力したid, passwordをサーバに送る
2. サーバはDBからidを使ってUserを取得する
3. 取得したUserのpasswordと送られてきたパスワードが同じか検証する
4. sessionIdとOTPを生成し、sessionIdをキーにUserとOTPをredisに保存する
5. 保存したOTPをUserのLineUidを使用してプッシュメッセージで送信する
6. キーに使用したsessionIdをfrontendへレスポンスとして送る(cookieに保存される)
7. userはLINEBotから送信されたOTPをfrontendへ入力する
7. frontendからOTPをserverに送信する
8. サーバはcookieからsesionIdを読み取り、sessionIdをキーとするデータをredisから取得し、それに紐づいたデータを取得する
9. 取得したsessionに保存されたOTPと送られてきたOTPが一致するか検証する
10. 認証用のsessionIdを再び生成し、sessionIdをキーとしてUserをredisへ保存する
11. キーに使用したsessionIdをfrontendへレスポンスとして送る(cookieに保存される)
