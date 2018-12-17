cd `dirname "$0"` > /dev/null
BASEDIR=`pwd`
cd - > /dev/null

export GOPATH="$BASEDIR"
export PATH=$GOPATH/bin:$PATH

