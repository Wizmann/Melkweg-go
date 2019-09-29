cd `dirname "$0"` > /dev/null
BASEDIR=`pwd`
cd - > /dev/null

export GOPATH="$BASEDIR"

[[ ":$PATH:" != *"$GOPATH"* ]] && export PATH="$GOPATH:${PATH}"

