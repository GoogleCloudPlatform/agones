#!/usr/bin/env bash

# Copyright 2018 Google Inc. All Rights Reserved.
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
set -e

export SHELL="/bin/bash"
export KUBECONFIG="/root/.kube/config"
mkdir -p /go/src/agones.dev/agones/ /root/.kube/
cp -r /workspace/. /go/src/agones.dev/agones/
cd /go/src/agones.dev/agones/build
if [ "$1" = 'local' ]
then
        gcloud auth login
fi
gcloud container clusters get-credentials e2e-test-cluster \
        --zone=us-west1-c --project=agones-images
kubectl port-forward statefulset/consul 8500:8500 &
echo "Waiting consul port-forward to launch on 8500..."
timeout 60 bash -c 'until printf "" 2>>/dev/null >>/dev/tcp/$0/$1; do sleep 1; done' 127.0.0.1 8500
echo "consul port-forward launched. Starting e2e tests..."
consul lock -child-exit-code=true -timeout 30m -try 5m -verbose LockE2E /root/e2e.sh
killall -q kubectl

