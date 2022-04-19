COUNT=0
MAX_RETRIES=20
CURRENT_DATE=$(date +'%Y_%m_%d')

if [[ $# -ne 2 ]]; then
  echo "invalid number of arguments, provided $#, expected 2"
  exit 1
fi

PATH_TO_MANAGER=$1
PATH_TO_CONFIG=$2

if [ ! -f "${PATH_TO_MANAGER}" ]; then
    echo "cannot find account manager binary, provided path ${PATH_TO_CONFIG}"
    exit 1
fi

if [ ! -f "${PATH_TO_CONFIG}" ]; then
    echo "cannot find config file, provided path: ${PATH_TO_CONFIG}"
    exit 1
fi


while [ ${COUNT} -lt ${MAX_RETRIES} ]
do
  CURRENT_LOGS_FILE="logs_${CURRENT_DATE}_$(( COUNT+1 )).txt"

  ${PATH_TO_MANAGER} -config "${PATH_TO_CONFIG}" -type "reindex" >(tee -a "${CURRENT_LOGS_FILE}") 2>&1

  ERROR_OUTPUT=$(grep ERROR "${CURRENT_LOGS_FILE}" )

  if [ -z "${ERROR_OUTPUT}" ]
  then
    break
  else
    echo "Something went wrong error: ${ERROR_OUTPUT}"
    echo "Will retry "$(( COUNT+1 ))""
  fi

  COUNT=$(( COUNT+1 ))
done

if [ ${COUNT} -eq $(( MAX_RETRIES )) ]
then
  echo "Reindex process failed, check logs file ${LOGS_FILE}"
else
  echo "Reindex accounts with stake success"
fi


