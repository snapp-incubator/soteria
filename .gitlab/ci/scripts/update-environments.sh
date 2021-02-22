#!/usr/bin/env bash
set -eu

APP_NAME=soteria
APP_USERNAME=soteria
APP_HOSTNAME="emqx-0${CI_NODE_INDEX}.app.afra.snapp.infra"
APP_HOME_PATH=/var/lib/${APP_USERNAME}
APP_CONFIG_PATH="/etc/soteria"
ENV_DEFAULT_FILE=.gitlab/ci/env/env.conf
ENV_HOSTNAME=$APP_HOSTNAME
ENV_SOTERIA_JWT_KEYS_PATH=${ENV_SOTERIA_JWT_KEYS_PATH}
ENV_THIRD_PARTY_JWT_PRIVATE_KEY_PRODUCTION=${THIRD_PARTY_JWT_PRIVATE_KEY_PRODUCTION}
ENV_PASSENGER_JWT_PUBLIC_KEY_PRODUCTION=${ENV_PASSENGER_JWT_PUBLIC_KEY_PRODUCTION}
ENV_DRIVER_JWT_PUBLIC_KEY_PRODUCTION=${DRIVER_JWT_PUBLIC_KEY_PRODUCTION}
ENV_THIRD_PARTY_JWT_PUBLIC_KEY_PRODUCTION=${THIRD_PARTY_JWT_PUBLIC_KEY_PRODUCTION}

cat << EOF > /etc/resolv.conf
nameserver 172.16.76.22
nameserver 172.21.49.230
EOF

mkdir -p  /root/.ssh && touch /root/.ssh/id_rsa
echo "$DEPLOYER_PRIVATE_KEY" > /root/.ssh/id_rsa
chmod 0600 /root/.ssh/id_rsa

sed -i "2s/\(.*\)/\#Last update\:\ $(date)/" $ENV_DEFAULT_FILE


while read LINE
do


    if echo $LINE | egrep -v "^#|^$" >/dev/null 2>&1
    then

        ENV="$(echo $LINE | sed 's/\(.*\)=\(.*\)/\1/g' )"
        ENV_NAME="ENV_$ENV"
        ENV_VALUE="$(echo $LINE | sed 's/\(.*\)=\(.*\)/\2/g' )"

        echo -e "\e[33m# Ckeck variable \e[1m$ENV_NAME\e[0m"

            if [ -v $ENV_NAME ]
            then

                if [ $ENV_VALUE = "'___'" ]
                then

                    sed -i "s/$ENV\(\s\)*=\(\s\)*'\(.*\)'/$ENV=\'${!ENV_NAME}\'/g"  $ENV_DEFAULT_FILE
                    [ $? -eq 0 ] && echo -e "\e[32m==> Changed variable \e[1m$ENV_NAME\e[0m \e[32mfrom \e[1m$ENV_VALUE\e[0m \e[32mto \e[1m😜\e[0m "

                else

                    sed -i "s/$ENV\(\s\)*=\(\s\)*'\(.*\)'/$ENV=\'${!ENV_NAME}\'/g"  $ENV_DEFAULT_FILE
                    [ $? -eq 0 ] && echo -e "\e[32m==> Changed variable \e[1m$ENV_NAME\e[0m \e[32mfrom \e[1m$ENV_VALUE\e[0m \e[32mto \e[1m${!ENV_NAME}\e[0m "
                fi

            fi

    fi

done < $ENV_DEFAULT_FILE



  rsync -e 'ssh -o "StrictHostKeyChecking=no"' -avz "$ENV_DEFAULT_FILE" "$APP_USERNAME@$APP_HOSTNAME:$APP_CONFIG_PATH/$APP_NAME.conf"
  [ $? -eq 0 ] && echo -e "\e[32m# Connecting to Server \e[1m$APP_HOSTNAME\e[0m \e[32m and update \e[1m$APP_CONFIG_PATH/$APP_NAME.conf\e[0m"
  echo

  if [ -v UPDATE_ENVIRONMENT_VARIABLE ]
  then
     ssh -o StrictHostKeyChecking=no "$APP_USERNAME@$APP_HOSTNAME" "

       echo -e '\e[33m# Restarting $APP_NAME ...\e[0m'
       sudo systemctl restart $APP_NAME.service

     "
  fi

## JWT keys

mkdir jwt 
echo "${ENV_DRIVER_JWT_PUBLIC_KEY_PRODUCTION}" > "jwt"/0.pem
echo "${ENV_PASSENGER_JWT_PUBLIC_KEY_PRODUCTION}" > "jwt"/1.pem
echo "${ENV_THIRD_PARTY_JWT_PUBLIC_KEY_PRODUCTION}" > "jwt"/100.pem
echo "${ENV_THIRD_PARTY_JWT_PRIVATE_KEY_PRODUCTION}" > "jwt"/100.private.pem
rsync -e 'ssh -o "StrictHostKeyChecking=no"' -avzr ./jwt "$APP_USERNAME@$APP_HOSTNAME:${ENV_SOTERIA_JWT_KEYS_PATH}"
rm -r jwt 

### Green service config file
# sed -i "1s/\(.*\)/\#  Application: $APP_NAME-green /" $ENV_DEFAULT_FILE
sed -i "2s/\(.*\)/\#  Last update: $(date) /" $ENV_DEFAULT_FILE
sed -i "s/SOTERIA_HTTP_PORT\(\s\)*=\(\s\)*'\(.*\)'/ENV_NAME=\'9998\'/g"  $ENV_DEFAULT_FILE
rsync -e 'ssh -o "StrictHostKeyChecking=no"' -avz "$ENV_DEFAULT_FILE" "$APP_USERNAME@$APP_HOSTNAME:$APP_CONFIG_PATH/$APP_NAME-green.conf"
