{{/*
Expand the name of the chart.
*/}}
{{- define "kubeserial.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "kubeserial.fullname" -}}
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
Create chart name and version as used by the chart label.
*/}}
{{- define "kubeserial.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "kubeserial.labels" -}}
helm.sh/chart: {{ include "kubeserial.chart" . }}
{{ include "kubeserial.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "kubeserial.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kubeserial.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "kubeserial.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "kubeserial.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Webhook fullname
*/}}
{{- define "kubeserial.injectorFullname" -}}
{{ include "kubeserial.fullname" .}}-sidecar-injector
{{- end }}


{{/*
Webhook common labels
*/}}
{{- define "kubeserial.injectorLabels" -}}
helm.sh/chart: {{ include "kubeserial.chart" . }}
{{ include "kubeserial.injectorSelectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Webhook selector labels
*/}}
{{- define "kubeserial.injectorSelectorLabels" -}}
app.kubernetes.io/name: {{ include "kubeserial.name" . }}-sidecar-injector
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Device Monitor fullname
*/}}
{{- define "kubeserial.monitorFullname" -}}
{{ include "kubeserial.fullname" .}}-monitor
{{- end }}

{{/*
Device Monitor common labels
*/}}
{{- define "kubeserial.monitorLabels" -}}
helm.sh/chart: {{ include "kubeserial.chart" . }}
{{ include "kubeserial.monitorSelectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Device Monitor selector labels
*/}}
{{- define "kubeserial.monitorSelectorLabels" -}}
app.kubernetes.io/name: {{ include "kubeserial.monitorFullname" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
