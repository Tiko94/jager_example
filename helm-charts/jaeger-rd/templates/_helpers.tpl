{{/* vim: set filetype=helm: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "jaeger-rd.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "jaeger-rd.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "jaeger-rd.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "jaeger-rd.commonLabels" -}}
helm.sh/chart: {{ include "jaeger-rd.chart" . }}
{{ include "jaeger-rd.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/part-of: "jaeger-rd"
{{- end -}}

{{/*
Component labels
*/}}
{{- define "jaeger-rd.dataservice.labels" -}}
{{ include "jaeger-rd.commonLabels" . }}
app.kubernetes.io/component: "dataservice"
{{- end -}}

{{- define "jaeger-rd.loadbalancer.labels" -}}
{{ include "jaeger-rd.commonLabels" . }}
app.kubernetes.io/component: "loadbalancer"
{{- end -}}

{{- define "jaeger-rd.logservice.labels" -}}
{{ include "jaeger-rd.commonLabels" . }}
app.kubernetes.io/component: "logservice"
{{- end -}}

{{- define "jaeger-rd.web.labels" -}}
{{ include "jaeger-rd.commonLabels" . }}
app.kubernetes.io/component: "web"
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "jaeger-rd.selectorLabels" -}}
app.kubernetes.io/name: {{ include "jaeger-rd.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "jaeger-rd.dataservice.selectorLabels" -}}
app.kubernetes.io/component: "dataservice"
{{ include "jaeger-rd.selectorLabels" . }}
{{- end -}}

{{- define "jaeger-rd.loadbalancer.selectorLabels" -}}
app.kubernetes.io/component: "loadbalancer"
{{ include "jaeger-rd.selectorLabels" . }}
{{- end -}}

{{- define "jaeger-rd.logservice.selectorLabels" -}}
app.kubernetes.io/component: "logservice"
{{ include "jaeger-rd.selectorLabels" . }}
{{- end -}}

{{- define "jaeger-rd.web.selectorLabels" -}}
app.kubernetes.io/component: "web"
{{ include "jaeger-rd.selectorLabels" . }}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "jaeger-rd.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "jaeger-rd.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}
