# Copyright (c) 2019  InterDigital Communications, Inc
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

FROM debian:9.6-slim

COPY ./meep-gis-engine /templates/sandbox/meep-gis-engine
COPY ./meep-loc-serv /templates/sandbox/meep-loc-serv
COPY ./meep-metrics-engine /templates/sandbox/meep-metrics-engine
COPY ./meep-mg-manager /templates/sandbox/meep-mg-manager
COPY ./meep-rnis /templates/sandbox/meep-rnis
COPY ./meep-sandbox-ctrl /templates/sandbox/meep-sandbox-ctrl
COPY ./meep-tc-engine /templates/sandbox/meep-tc-engine
COPY ./meep-virt-chart-templates /templates/scenario/meep-virt-chart-templates
COPY ./entrypoint.sh /entrypoint.sh
COPY ./meep-virt-engine /meep-virt-engine

ENV HELM_VERSION="v2.16.1"
RUN mkdir -p /active \
    && apt-get update \
    && apt-get install -y wget \
    && wget -q https://storage.googleapis.com/kubernetes-helm/helm-${HELM_VERSION}-linux-amd64.tar.gz -O - | tar -xzO linux-amd64/helm > /usr/local/bin/helm \
    && chmod +x /usr/local/bin/helm \
    && chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
