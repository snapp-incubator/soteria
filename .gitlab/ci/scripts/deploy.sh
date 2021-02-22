#!/usr/bin/env bash
set -eu

APP_NAME=soteria
APP_USERNAME=soteria
APP_HOME_PATH=/var/lib/${APP_USERNAME}
APP_PATH=/usr/local/bin/${APP_NAME}/
APP_GREEN_PATH=/usr/local/bin/${APP_NAME}/${APP_NAME}-green
APP_RELEASES_PATH=${APP_HOME_PATH}/releases
APP_LATEST_RELEASE_PATH=${APP_RELEASES_PATH}/"$CI_COMMIT_REF_SLUG-$CI_COMMIT_SHORT_SHA"
APP_HOSTNAME=${APP_HOSTNAME:?Please specify "APP_HOSTNAME", e.g. "APP_HOSTNAME=emqx-01.app.afra.snapp.infra"};

#TOTAL_APP_NODES=5
UNHEALTHY_NODE_EXISTS=0

cat << EOF > /etc/resolv.conf
nameserver 172.16.76.22
nameserver 172.21.49.230
EOF

mkdir -p  /root/.ssh && touch /root/.ssh/id_rsa
echo "$DEPLOYER_PRIVATE_KEY" > /root/.ssh/id_rsa
chmod 0600 /root/.ssh/id_rsa


echo -e "\e[33m# Deployment on \e[1m$APP_HOSTNAME ...\e[0m"

if ssh -o StrictHostKeyChecking=no "$APP_USERNAME@$APP_HOSTNAME" "stat $APP_RELEASES_PATH/$CI_COMMIT_REF_SLUG-$CI_COMMIT_SHORT_SHA/$APP_NAME > /dev/null 2>&1"
then
  echo -e "\e[33m# The release has alredy exist.\e[0m"
else
    echo -e "\e[33m# Sending Artifact to Server\e[0m"
    rsync -e 'ssh -o "StrictHostKeyChecking=no"' -avz "./artifacts-$CI_COMMIT_REF_SLUG-$CI_COMMIT_SHORT_SHA.tar.gz" "$APP_USERNAME@$APP_HOSTNAME:$APP_HOME_PATH"
fi


ssh -o StrictHostKeyChecking=no "$APP_USERNAME@$APP_HOSTNAME" "
  mkdir -p ${APP_RELEASES_PATH} ${APP_LATEST_RELEASE_PATH}

  if [ -f "${APP_HOME_PATH}/artifacts-$CI_COMMIT_REF_SLUG-$CI_COMMIT_SHORT_SHA.tar.gz" ]; then
    echo -e '\e[33m# Extract Release\e[0m'
    tar -xzf ${APP_HOME_PATH}/artifacts-$CI_COMMIT_REF_SLUG-$CI_COMMIT_SHORT_SHA.tar.gz -C ${APP_LATEST_RELEASE_PATH}
    rm -f ${APP_HOME_PATH}/artifacts-$CI_COMMIT_REF_SLUG-$CI_COMMIT_SHORT_SHA.tar.gz
  fi

 "

echo -e "\e[33m# Trying to connect to 'http://$APP_HOSTNAME' ...\e[0m"
timeout=60

   
#done

echo
echo

#for i in $(seq 1 ${TOTAL_APP_NODES}); do
#  APP_HOSTNAME=$(printf "${APP_NAME}-%02d.app.afra.snapp.infra" ${i})
echo -e "\e[33m# Connecting to Server\e[0m"
  echo -e "\e[33m# Connecting to Server \e[1m$APP_HOSTNAME\e[0m"
  ssh -o StrictHostKeyChecking=no "$APP_USERNAME@$APP_HOSTNAME" "
    echo -e '\e[33m# Activating Latest Realese\e[0m'
    ln -sf ${APP_LATEST_RELEASE_PATH}/$APP_USERNAME ${APP_PATH}

    echo -e '\e[33m# Restarting $APP_NAME ...\e[0m'
    sudo systemctl restart $APP_NAME.service

  "
  echo
#done
