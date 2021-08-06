if [ "$3" == "-help" ]; then
	echo "Help:"
	echo "Usage: /script get_file.sh <file_full_path>"
  	exit
fi

USER_ID=$1
TELEGRAM_TOKEN=$2
FILE_PATH=$3

echo "Uploading $FILE_PATH"
curl -F document=@$FILE_PATH https://api.telegram.org/bot$TELEGRAM_TOKEN/sendDocument?chat_id=$USER_ID
