{{/*
Expand the name of the chart.
*/}}
{{- define "dsv.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "dsv.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create a DNS name i.e. a fully-qualified domain name (FQDN) for the webhook.
*/}}
{{- define "dsv.dnsname" -}}
{{- print (include "dsv.name" .) "." .Release.Namespace ".svc" -}}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "dsv.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "dsv.labels" -}}
helm.sh/chart: {{ include "dsv.chart" . }}
{{ include "dsv.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
dsv-filter-name: {{ .Chart.Name }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "dsv.selectorLabels" -}}
app.kubernetes.io/name: {{ include "dsv.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
dsv-filter-name: {{ .Chart.Name }}
{{- end }}
