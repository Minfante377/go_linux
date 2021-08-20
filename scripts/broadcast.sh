if [ "$3" == "-help" ]; then
	echo "Help:"
	echo "Usage: /script broadcast.sh "msg""
  	exit
fi

PARAMS=("${@}")
USER_ID=$1
TELEGRAM_TOKEN=$2
MSG="${PARAMS[@]:2}"

SCRIPT=$(readlink -f $0)
SCRIPTPATH=`dirname $SCRIPT`
GOROOT=$(cd $SCRIPTPATH; cd ../; echo $PWD)
source "$GOROOT/.secrets"

if [ "$USER_ID" != "$ADMIN" ]; then
	echo "Only admin user is allowed to broadcast"
	exit
fi

DB=$(ls $GOROOT | grep --color=never .db)


USERS=(`python3 $GOROOT/scripts/get_users.py $GOROOT/$DB`)
FORMAT_MSG="Broadcast message from your admin $USER_ID: $MSG"
REPLACE="%20"
for user in "${USERS[@]}"
do
	echo "Sending "${FORMAT_MSG}" to $user"
	curl -X GET "https://api.telegram.org/bot$TELEGRAM_TOKEN/sendMessage?chat_id=$user&text="${FORMAT_MSG// /$REPLACE}""
	echo
done
