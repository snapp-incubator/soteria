# To notify matrix about CI job

MATRIX_SERVER="https://matrix.snapp.cab"
MATRIX_MSQTYPE=m.text
MX_TXN="`date "+%s"`$(( RANDOM % 9999 ))"


BODY="[Gitlab🚀] CI run by *$GITLAB_USER_LOGIN* on the $CI_PROJECT_TITLE project: \n      Job name : $CI_JOB_NAME \n      Commit   : $CI_COMMIT_TITLE \n      Job url    : $CI_JOB_URL "

# Post into maint room
curl -X PUT --header 'Content-Type: application/json' --header 'Accept: application/json' -d "{
\"msgtype\": \"$MATRIX_MSQTYPE\",
\"body\": \"$BODY\"
      }" "$MATRIX_SERVER/_matrix/client/unstable/rooms/$MATRIX_ROOM_ID/send/m.room.message/$MX_TXN?access_token=$MATRIX_ACCESS_TOKEN"  >/dev/nul 2>&1
