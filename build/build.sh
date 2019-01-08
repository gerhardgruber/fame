set -e

echo Server build starting
echo Instanz:     ${INSTANCE}
echo Tag-Version: ${TAG_VERSION}
export FAME_SERVER=doctype.documatrix.com

# set folder DEV Prefix
if [ "${INSTANCE}" = "PROD" ]; then
  export DEV_FOLDER=
else
  export DEV_FOLDER=dev
fi
# resolve dependencies
echo resolve dependencies
dep ensure -v
export STATUS=$?
if [ $STATUS -gt 0 ] ; then
    echo Error while resolving dependencies: $STATUS
    exit $STATUS
fi

# compile source
echo check compiler version
go version

echo compile source
go install ./bin/fame_server
export STATUS=$?
if [ $STATUS -gt 0 ] ; then
    echo Error while compiling: $STATUS
    exit $STATUS
fi

# copy to remote server
echo copy to remote server
scp $GOPATH/bin/fame_server sophy@${FAME_SERVER}:/opt/${DEV_FOLDER}/fame_server/delivery/bin/
scp ./i18n/* sophy@${FAME_SERVER}:/opt/${DEV_FOLDER}/fameserver/delivery/i18n/
scp ./package.json sophy@${FAME_SERVER}:/opt/${DEV_FOLDER}/fame_server/delivery/
export STATUS=$?
if [ $STATUS -gt 0 ] ; then
    echo Error while copying: $STATUS
    exit $STATUS
fi

# start the service on the remote server
echo stop, copy and start the service on the remote server
ssh sophy@${FAME_SERVER} << FINISHED
    if [ "${INSTANCE}" = "TEST" ]; then
      sudo stop fame_dev
    else
      sudo stop fame
    fi
    mv /opt/${DEV_FOLDER}/fame_server/delivery/bin/* /opt/${DEV_FOLDER}/fame_server/bin/
    mv /opt/${DEV_FOLDER}/fame_server/delivery/i18n/* /opt/${DEV_FOLDER}/fame_server/i18n/
    mv /opt/${DEV_FOLDER}/fame_server/delivery/package.json /opt/${DEV_FOLDER}/fame_server/
    if [ "${INSTANCE}" = "TEST" ]; then
      sudo start fame_dev
    else
      sudo start fame
    fi
FINISHED

export STATUS=$?
if [ $STATUS -gt 0 ] ; then
    echo Error while starting the service: $STATUS
    exit $STATUS
fi

# finished
echo Server build finished


# start webapp build
echo "Starting webapp build..."

echo "Installing dependencies..."
npm install
export STATUS=$?
if [ $STATUS -gt 0 ] ; then
    echo Error while installing dependencies: $STATUS
    exit $STATUS
fi

echo "Modify API_ROOT"
if [ "${INSTANCE}" = "PROD" ]; then
  echo "export const API_ROOT = window.location.protocol + \"//\" + window.location.hostname + \"/api\";" > webapp/config/fame.ts
  echo "export const WS_API_ROOT = \"wss://\" + window.location.hostname + \"/wsapi\";" >> webapp/config/fame.ts
fi

echo "Starting build..."
npm run-script build
export STATUS=$?
if [ $STATUS -gt 0 ] ; then
    echo Error while building: $STATUS
    exit $STATUS
fi

sed -i 's/bundle.js/bundle-'"${VERSION}"'.js/g' webapp/dist.html
scp webapp/dist/* sophy@${FAME_SERVER}:/opt/${DEV_FOLDER}/fame_server/webapp/
scp webapp/dist/bundle.js sophy@${FAME_SERVER}:/opt/${DEV_FOLDER}/fame_server/webapp/bundle-${VERSION}.js
scp webapp/static/*.png sophy@${FAME_SERVER}:/opt/${DEV_FOLDER}/fame_server/webapp/
scp webapp/static/*.jpg sophy@${FAME_SERVER}:/opt/${DEV_FOLDER}/fame_server/webapp/
scp webapp/dist.html sophy@${FAME_SERVER}:/opt/${DEV_FOLDER}/fame_server/webapp/index.html
export STATUS=$?
if [ $STATUS -gt 0 ] ; then
    echo Error while copying: $STATUS
    exit $STATUS
fi

# finished
echo "Webapp build finished"
