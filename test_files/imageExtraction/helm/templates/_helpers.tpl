{{/*
Expand the name of the chart.
*/}}
{{- define "containers.name" -}}
{{- .Chart.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}

{{- define "get-fullname" -}}
{{- $context := index . 0 -}}
{{- $service := index . 1 -}}
{{- printf "%s-%s-%s" $context.Release.Name $context.Chart.Name $service | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "containers.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "containers.labels" -}}
helm.sh/chart: {{ include "containers.chart" . }}
{{ include "containers.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "containers.selectorLabels" -}}
app.kubernetes.io/name: {{ include "containers.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create platform header for SCA usage
*/}}
{{- define "sca.platform" -}}
{{- $tenantMode := ternary "ST" "MT"  (eq .Values.deploymentType "single_tenant") }}
{{- $env := default "N/A" .Values.scaServices.iamDomain | replace "https://" "" | replace ".cxast.net" "" | replace ".checkmarx.net" "" }}
{{- printf "AST|SCA_ver:%v|Env:%s|Tenant_mode:%s" .Chart.Version $env $tenantMode }}
{{- end }}

{{- define "sca.astDomain" -}}
{{- $env := default "N/A" .Values.networking.domain | replace "https://" "" | replace ".cxast.net" "" | replace ".checkmarx.net" "" }}
{{- printf "%s" $env }}
{{- end }}
