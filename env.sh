BASEDIR=`readlink -f $(dirname "$0")`
export GOPATH="$BASEDIR"

export PATH=$GOPATH/bin:$PATH


