FROM sebvaucher/sgx-base
WORKDIR /usr/src/app/
RUN /bin/bash -c 'git clone https://github.com/intel/intel-sgx-ssl.git' intel-sgx-ssl
RUN /bin/bash -c 'cd intel-sgx-ssl && git reset --hard e52f41972fd0963993feff7a447e37a57a66c6e5 && cd ..'
COPY openssl-1.1.1d.tar.gz ./intel-sgx-ssl/openssl_source/openssl-1.1.1d.tar.gz
WORKDIR /usr/src/app/intel-sgx-ssl/Linux
RUN /bin/bash -c 'source /opt/intel/sgxsdk/environment && make all DEBUG=0 SGX_MODE=SIM VERBOSE=1 && make install'
WORKDIR /usr/src/app
COPY . .
RUN /bin/bash -c 'source /opt/intel/sgxsdk/environment && make clean && make SGX_PRERELEASE=1 SGX_MODE=HW SGX_DEBUG=0'
