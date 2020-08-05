package main

var configTemplate string = `# Generated config

listening-port={{ .Config.ListeningPort }}

listening-ip={{ .LocalIP }}
external-ip={{ .RemoteIP }}

min-port={{ .Config.MinPort }}
max-port={{ .Config.MaxPort }}

use-auth-secret
static-auth-secret={{ .Config.StaticAuthSecret }}

realm={{ .Config.Realm }}
server-name={{ .Config.ServerName }}

user-quota={{ .Config.UserQuota }}
total-quota={{ .Config.TotalQuota }}

no-tcp-relay

log-file=stdout

{{- range $ip := .Config.DeniedPeerIP }}
denied-peer-ip={{ $ip }}
{{- end }}

{{- if .Config.AllowedPeerIP }}
allowed-peer-ip={{ .Config.AllowedPeerIP }}
{{- end }}

no-cli
`
