export GOROOT=/usr/local/go
export GOPATH=/opt/ibm/caley/
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
yum install unixODBC unixODBC-devel
cp cayley/graph/sql/odbc/odbcinst.ini /etc/odbcinst.ini
cp cayley/graph/sql/odbc/odbc.ini /etc/odbc.ini
yum install wget -y
wget https://dl.google.com/go/go1.12.6.linux-amd64.tar.gz
tar -xzf go1.12.6.linux-amd64.tar.gz
mv go /usr/local
/opt/ibm/caley/cayley init -c /opt/ibm/caley/config.json 
/opt/ibm/caley/cayley load -c /opt/ibm/caley/config.json -i /opt/ibm/caley/data/testdata.nq
