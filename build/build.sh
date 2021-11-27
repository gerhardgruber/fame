set -e

echo Server build starting
echo Instanz:     ${INSTANCE}
echo Tag-Version: ${TAG_VERSION}

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
cp $GOPATH/bin/fame_server /opt/fame_server/delivery/bin/
cp ./i18n/* /opt/fame_server/delivery/i18n/
cp ./package.json /opt/fame_server/delivery/
export STATUS=$?
if [ $STATUS -gt 0 ] ; then
    echo Error while copying: $STATUS
    exit $STATUS
fi

# start the service on the remote server
sudo service fame_server stop
mv /opt/fame_server/delivery/bin/* /opt/fame_server/bin/
mv /opt/fame_server/delivery/i18n/* /opt/fame_server/i18n/
mv /opt/fame_server/delivery/package.json /opt/fame_server/
sudo service start fame_server

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
npm ci
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
cp webapp/dist/* /opt/fame_server/webapp/
cp webapp/dist/bundle.js /opt/fame_server/webapp/bundle-${VERSION}.js
cp webapp/static/*.jpg /opt/fame_server/webapp/
cp webapp/dist.html /opt/fame_server/webapp/index.html
export STATUS=$?
if [ $STATUS -gt 0 ] ; then
    echo Error while copying: $STATUS
    exit $STATUS
fi

# finished
echo "Webapp build finished"
