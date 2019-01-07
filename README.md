# fame

fame is a platform which is used by "Rettungshunde Niederösterreich" (see https://rettungshunde.eu) to manage their members, dogs and their rescue-missions

## Developer installation
### Install go
requires at least go version 1.10 for go dep
See https://golang.org
It is recommended to add the $GOPATH/bin directory to your $PATH variable

### Install python
Install python via package manager or https://www.python.org/downloads/

### Clone repo in go path
```bash
cd $GOPATH/src
mkdir -p github.com/gerhardgruber
cd github.com/gerhardgruber
git clone git@github.com:gerhardgruber/fame.git
cd fame
```

### Install go dep
```bash
go get -u github.com/golang/dep/cmd/dep
```

### Install dependencies
```bash
dep ensure
```

### Compile (and install)
```bash
go install ./bin/fame_server
```

### Create database (e. g. for MySQL)
```mysql
CREATE DATABASE `fame` /*!40100 DEFAULT CHARACTER SET utf8 */
```

### Run GUI
```bash
fame_server
npm install
npm start
```

