package fluentbit

const valuesTemplate = `
---
image:
  repository: public.ecr.aws/aws-observability/aws-for-fluent-bit
  tag: {{ .Version }}
serviceAccount:
  annotations:
    {{ .IrsaAnnotation }}
  create: true
  name: {{ .ServiceAccount }}
env:
  - name: HOST_NAME
    valueFrom:
      fieldRef:
        fieldPath: spec.nodeName

#! https://docs.fluentbit.io/manual/administration/configuring-fluent-bit/configuration-file
config:
  service: |
    [SERVICE]
        Flush                     5
        Log_Level                 {{"{{ .Values.logLevel }}"}}
        Daemon                    off
        Parsers_File              parsers.conf
        Parsers_File              custom_parsers.conf
        HTTP_Server               On
        HTTP_Listen               0.0.0.0
        HTTP_Port                 {{"{{ .Values.service.port }}"}}
        storage.path              /var/fluent-bit/state/flb-storage/
        storage.sync              normal
        storage.checksum          off
        storage.backlog.mem_limit 5M

  #! https://docs.fluentbit.io/manual/pipeline/inputs
  #! TODO -- Look into new Fluent-bit 1.8, multiline.parser  docker, cri
  #! TODO -- Understand Docker_Mode_Parser container_firstline
  inputs: |
    [INPUT]
        Name                tail
        Tag                 application.*
        Exclude_Path        /var/log/containers/cloudwatch-agent*, /var/log/containers/fluent-bit*, /var/log/containers/aws-node*, /var/log/containers/kube-proxy*
        Path                /var/log/containers/*.log
        Docker_Mode         On
        Docker_Mode_Flush   5
        Docker_Mode_Parser  container_firstline
        Parser              docker_cwci
        DB                  /var/fluent-bit/state/flb_container.db
        Mem_Buf_Limit       50MB
        Skip_Long_Lines     On
        Refresh_Interval    10
        Rotate_Wait         30
        storage.type        filesystem
        Read_from_Head      {{ .TailReadFromHead }}

    [INPUT]
        Name                tail
        Tag                 application.*
        Path                /var/log/containers/fluent-bit*
        Parser              docker_cwci
        DB                  /var/fluent-bit/state/flb_log.db
        Mem_Buf_Limit       5MB
        Skip_Long_Lines     On
        Refresh_Interval    10
        Read_from_Head      {{ .TailReadFromHead }}

    [INPUT]
        Name                tail
        Tag                 application.*
        Path                /var/log/containers/cloudwatch-agent*
        Docker_Mode         On
        Docker_Mode_Flush   5
        Docker_Mode_Parser  cwagent_firstline
        Parser              docker_cwci
        DB                  /var/fluent-bit/state/flb_cwagent.db
        Mem_Buf_Limit       5MB
        Skip_Long_Lines     On
        Refresh_Interval    10
        Read_from_Head      {{ .TailReadFromHead }}

    [INPUT]
        Name                systemd
        Tag                 dataplane.systemd.*
        Systemd_Filter      _SYSTEMD_UNIT=docker.service
        Systemd_Filter      _SYSTEMD_UNIT=kubelet.service
        DB                  /var/fluent-bit/state/systemd.db
        Path                /var/log/journal

    [INPUT]
        Name                tail
        Tag                 dataplane.tail.*
        Path                /var/log/containers/aws-node*, /var/log/containers/kube-proxy*
        Docker_Mode         On
        Docker_Mode_Flush   5
        Docker_Mode_Parser  container_firstline
        Parser              docker_cwci
        DB                  /var/fluent-bit/state/flb_dataplane_tail.db
        Mem_Buf_Limit       50MB
        Skip_Long_Lines     On
        Refresh_Interval    10
        Rotate_Wait         30
        storage.type        filesystem
        Read_from_Head      {{ .TailReadFromHead }}


  #! https://docs.fluentbit.io/manual/pipeline/filters
  #! TODO -- Add option for Use_Kubelet
  #! TODO -- native docker and syslog-rfc3164 parsers, remove time?
  filters: |
    [FILTER]
        Name                kubernetes
        Match               application.*
        Kube_URL            https://kubernetes.default.svc:443
        Kube_Tag_Prefix     application.var.log.containers.
        Merge_Log           On
        Merge_Log_Key       log_processed
        K8S-Logging.Parser  On
        K8S-Logging.Exclude Off
        Labels              Off
        Annotations         Off

    [FILTER]
        Name                modify
        Match               dataplane.systemd.*
        Rename              MESSAGE                     message
        Rename              _HOSTNAME                   hostname
        Rename              _SYSTEMD_UNIT               systemd_unit
        Remove_regex        ^((?!hostname|systemd_unit|message).)*$

    [FILTER]
        Name                aws
        Match               dataplane.*
        imds_version        v1

    [FILTER]
        Name                rewrite_tag
        Match               application.*
        Rule                $log .* od.$kubernetes['namespace_name'].$kubernetes['pod_name'] false

    [FILTER]
        Name                kubernetes
        Match               dataplane.tail.*
        Kube_URL            https://kubernetes.default.svc:443
        Kube_Tag_Prefix     dataplane.tail.var.log.containers.
        Merge_Log           On
        Merge_Log_Key       log_processed
        K8S-Logging.Parser  On
        K8S-Logging.Exclude Off
        Labels              Off
        Annotations         Off

    [FILTER]
        Name                rewrite_tag
        Match               dataplane.tail.*
        Rule                $log .* d.$kubernetes['namespace_name'].$kubernetes['pod_name'].$kubernetes['host'] false   

    [FILTER]
        Name                rewrite_tag
        Match               dataplane.systemd.*
        Rule                $message .* ystemd.$systemd_unit[0].$hostname false

  #! https://docs.fluentbit.io/manual/pipeline/outputs
  outputs: |
    [OUTPUT]
        Name                cloudwatch_logs
        Match               od.*
        region              {{ .Region }}
        log_group_name      /aws/eks/{{ .ClusterName }}/application
        log_stream_prefix   p
        auto_create_group   true

    [OUTPUT]
        Name                cloudwatch_logs
        Match               ystemd.*
        region              {{ .Region }}
        log_group_name      /aws/eks/{{ .ClusterName }}/dataplane
        log_stream_prefix   s
        auto_create_group   true

    [OUTPUT]
        Name                cloudwatch_logs
        Match               d.*
        region              {{ .Region }}
        log_group_name      /aws/eks/{{ .ClusterName }}/dataplane
        log_stream_prefix   po
        auto_create_group   true

  #! https://docs.fluentbit.io/manual/pipeline/parsers
  customParsers: |
    [PARSER]
        Name                docker_cwci
        Format              json
        Time_Key            time
        Time_Format         %Y-%m-%dT%H:%M:%S.%LZ

    [PARSER]
        Name                syslog
        Format              regex
        Regex               ^(?<time>[^ ]* {1,2}[^ ]* [^ ]*) (?<host>[^ ]*) (?<ident>[a-zA-Z0-9_\/\.\-]*)(?:\[(?<pid>[0-9]+)\])?(?:[^\:]*\:)? *(?<message>.*)$
        Time_Key            time
        Time_Format         %b %d %H:%M:%S

    [PARSER]
        Name                container_firstline
        Format              regex
        Regex               (?<log>(?<="log":")\S(?!\.).*?)(?<!\\)".*(?<stream>(?<="stream":").*?)".*(?<time>\d{4}-\d{1,2}-\d{1,2}T\d{2}:\d{2}:\d{2}\.\w*).*(?=})
        Time_Key            time
        Time_Format         %Y-%m-%dT%H:%M:%S.%LZ

    [PARSER]
        Name                cwagent_firstline
        Format              regex
        Regex               (?<log>(?<="log":")\d{4}[\/-]\d{1,2}[\/-]\d{1,2}[ T]\d{2}:\d{2}:\d{2}(?!\.).*?)(?<!\\)".*(?<stream>(?<="stream":").*?)".*(?<time>\d{4}-\d{1,2}-\d{1,2}T\d{2}:\d{2}:\d{2}\.\w*).*(?=})
        Time_Key            time
        Time_Format         %Y-%m-%dT%H:%M:%S.%LZ
`
