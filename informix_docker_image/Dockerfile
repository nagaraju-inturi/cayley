#FROM centos:7.4.1708
FROM centos
MAINTAINER nagaraju.inturi@hcl.com
ENV INFORMIXDIR /opt/ibm/informix/
ENV LD_LIBRARY_PATH /lib64
ENV GOROOT /usr/local/go
ENV GOPATH /opt/ibm/godb/
ENV ODBCINI /etc/odbc.ini
RUN mkdir /opt/ibm/  &&  mkdir /opt/ibm/informix
#VOLUME ["/opt/ibm/informix/", "/opt/ibm/data/"]
VOLUME ["/opt/ibm/data/"]
RUN yum -y update && yum -y install libaio sysstat libgcc libstdc compat-libstdc++-33 gcc kernel-headers kernel-devel net-tools.x86_64 xmvn.noarch epel-release git wget unixODBC unixODBC-devel
RUN /usr/sbin/groupadd -g 102 -f informix
RUN ["/usr/sbin/useradd", "-m", "-g 102", "-u 200", "informix"]
RUN /bin/bash -c 'echo "informix:changeme" | /usr/sbin/chpasswd'
ADD ibm.ids.14.10.FC1.LNX.tar /opt/ibm/informix/
ADD clientsdk.4.50.tar /opt/ibm/informix/
ADD csdk_response.txt /tmp/csdk_response.txt
ADD bundle.properties /tmp/bundle.properties
RUN /opt/ibm/informix/ids_install -i silent -f /tmp/bundle.properties; exit 0
RUN /opt/ibm/informix/installclientsdk  -i SILENT -f /tmp/csdk_response.txt
RUN /bin/bash -c 'chown informix:informix /opt/ibm/data/'
RUN /bin/bash -c '/bin/rm -rf /opt/ibm/informix/ids_install /opt/ibm/informix/IBM_Data_Server_Driver_Package /opt/ibm/informix/CSDK /opt/ibm/informix/ICONNECT /opt/ibm/informix/JDBC /opt/ibm/informix/SDK /opt/ibm/informix/SERVER'

RUN wget https://dl.google.com/go/go1.12.6.linux-amd64.tar.gz; exit 0
RUN tar -xzf go1.12.6.linux-amd64.tar.gz
RUN mv go /usr/local/

USER informix
RUN /bin/bash -c 'cp /opt/ibm/informix//etc/onconfig.std /opt/ibm/informix//etc/onconfig'
#RUN /bin/bash -c 'cp /opt/ibm/informix//etc/sqlhosts.std /opt/ibm/informix//etc/sqlhosts'
ADD graphdb.sql /opt/ibm/
ADD odbcinst.ini /etc/
ADD odbc.ini /etc/
ADD config.json /etc/

USER root

WORKDIR /opt/ibm
ENTRYPOINT ["/opt/ibm/boot.sh"]
CMD ["--start"]
EXPOSE 60000
EXPOSE 64210
ADD boot.sh /opt/ibm/
RUN git clone https://github.com/nagaraju-inturi/cayley
