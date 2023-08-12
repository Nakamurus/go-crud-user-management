# baseurlの定義
baseurl="http://localhost:8080"

# 複数のユーザーを登録
http POST $baseurl/user name=User1 email=user1@example.com password=password1 id=d1d3728f-77d7-4bd0-b90c-e479a5cb7241
http POST $baseurl/user name=User2 email=user2@example.com password=password2 id=d1d3728f-77d7-4bd0-b90c-e479a5cb7250

# ユーザー一覧の表示
http GET $baseurl/users


# ユーザーのログイン (jwt_tokenはレスポンスから取得)
http POST $baseurl/login email=user2@example.com password=password2

jwt_token=$(http POST $baseurl/login email=user1@example.com password=password1 | jq -r .token)
echo $jwt_token
# ユーザーの情報取得
http GET $baseurl/user/d1d3728f-77d7-4bd0-b90c-e479a5cb7241 'Authorization:Bearer '$jwt_token

# ユーザー情報の更新
jwt_token2=$(http PUT $baseurl/me/d1d3728f-77d7-4bd0-b90c-e479a5cb7241 name=NewUser1 email=newuser1@example.com 'Authorization:Bearer '$jwt_token | jq -r .token)
echo $jwt_token2

# パスワードの変更
http PUT $baseurl/me/d1d3728f-77d7-4bd0-b90c-e479a5cb7241/password old_password=password1 new_password=new_password1 'Authorization:Bearer '$jwt_token2
# ユーザーの削除
http DELETE $baseurl/me/d1d3728f-77d7-4bd0-b90c-e479a5cb7241 'Authorization:Bearer '$jwt_token2

# ログアウト
http POST $baseurl/login email=user2@example.com password=password2
jwt_token=$(http POST $baseurl/login email=user2@example.com password=password2 | jq -r .token)

http POST $baseurl/me/refresh-token "Authorization:Bearer $jwt_token"
new_jwt_token=$(http POST $baseurl/me/refresh-token "Authorization:Bearer $jwt_token" | jq -r .token)
echo $jwt_token
echo $new_jwt_token
http POST $baseurl/me/logout "Authorization:Bearer $new_jwt_token"

