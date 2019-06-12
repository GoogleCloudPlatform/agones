# Copyright 2019 Google LLC All Rights Reserved.
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

resource "kubernetes_service_account" "tiller" {
  metadata {
    name      = "tiller"
    namespace = "kube-system"
  }
  automount_service_account_token = true
}

resource "kubernetes_cluster_role_binding" "tiller" {
  metadata {
    name = "tiller"
  }

  role_ref {
    kind      = "ClusterRole"
    name      = "cluster-admin"
    api_group = "rbac.authorization.k8s.io"
  }

  subject {
    kind = "ServiceAccount"
    name = "tiller"

    api_group = ""
    namespace = "kube-system"
  }

  depends_on = ["kubernetes_service_account.tiller"]
}

provider "kubernetes" {
  version = "~> 1.5"
  load_config_file       = false
  host                   = "${var.host}"
  token                  = "${var.token}"
  cluster_ca_certificate = "${var.cluster_ca_certificate}"
}

provider "helm" {
  version = "~> 0.7"

  debug           = true
  install_tiller  = true
  service_account = "${kubernetes_service_account.tiller.metadata.0.name}"
  tiller_image    = "gcr.io/kubernetes-helm/tiller:v2.12.3"

  kubernetes {
    load_config_file = false
    host                   = "${var.host}"
    token                  = "${var.token}"
    cluster_ca_certificate = "${var.cluster_ca_certificate}"
  }
}

# In Terraform version 0.12 Interpolation would only evaluate one branch of a condition
# https://github.com/hashicorp/terraform/issues/15605
# so we can remove this and change values in helm_release to:
# 
#  values = [
#    "${length(var.values_file) == 0 ? ""  : file("${var.values_file}"))}"
#  ]
data "null_data_source" "values_file" {
  count = "${length(var.values_file) == 0 ? 0 : 1}"
  inputs = {
    "values" = "${file("${var.values_file}")}"
  }
}


data "helm_repository" "agones" {
  name = "agones"
  url  = "https://agones.dev/chart/stable"

  depends_on = ["kubernetes_cluster_role_binding.tiller"]
}


# TODO: remove - not needed in Terraform 0.12
locals {
  values = {
    params = "${join("", data.null_data_source.values_file.*.outputs.values)}"
  } 
  # Skip image tag if it is not needed
  # for installing latest image it would use chart value
  tag_name = "${var.agones_version != "" ? "agones.image.tag" : "skip"}"
}

resource "helm_release" "agones" {
  name  = "agones"
  force_update = "true"
  repository = "${data.helm_repository.agones.metadata.0.name}"
  chart = "${var.chart}"
  timeout = 420

  values = [
    # Switch in terraform 0.12 to:
    # "${length(var.values_file) == 0 ? ""  : file("${var.values_file}"))}"
    "${length(var.values_file) == 0 ? ""  : local.values["params"]}"
  ]

  set {
    name  = "crds.CleanupOnDelete"
    value = "${var.crd_cleanup}"
  }
  set {
    name = "${local.tag_name}"
    value = "${var.agones_version}"
  }
  set {
    name  = "agones.image.registry"
    value = "${var.image_registry}"
  }
  set {
    name  = "agones.image.controller.pullPolicy"
    value = "${var.pull_policy}"
  }
  set {
    name  = "agones.image.sdk.alwaysPull"
    value = "${var.always_pull_sidecar}"
  }
  set {
    name  = "agones.image.controller.pullSecret"
    value = "${var.image_pull_secret}"
  }
  set {
    name  = " agones.ping.http.serviceType"
    value = "${var.ping_service_type}"
  }
  set {
    name  = "agones.ping.udp.serviceType"
    value = "${var.ping_service_type}"
  }
  version = "${var.agones_version}"
  namespace  = "agones-system"

  depends_on = ["null_resource.helm_init", "kubernetes_cluster_role_binding.tiller"]
}

provider "null" {
  version = "~> 2.1"
}

# Creates folder with repositories so that helm provider would not fail
resource "null_resource" "helm_init" {
  triggers = {
    always_run = "${timestamp()}"
  }
  provisioner "local-exec" {
    command = "helm init --client-only"
  }
}
