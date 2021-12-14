# Copyright 2018 Intel Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# -----------------------------------------------------------------------------

FROM hyperledger/sawtooth-shell:chime

RUN apt-get update \
 && apt-get install -y -q \
    golang-1.13-go \
    git \
    libssl-dev \
    wget \
 && apt install ca-certificates libgnutls30 -y

ENV GOPATH=/go
ENV PATH=$PATH:/go/bin:/usr/lib/go-1.13/bin

RUN go get -u \
    github.com/moorada/sawtooth-auction-go/src/sawtooth_auction_client
# RUN wget https://raw.githubusercontent.com/moorada/sawtooth-auction-go/main/src/sawtooth_auction_client/test.sh
# RUN chmod +x test.sh
