package game_2048

const gameManifestTemplate = `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: {{ .Namespace }}
  name: deployment-2048
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: app-2048
  replicas: {{ .Replicas }}
  template:
    metadata:
      labels:
        app.kubernetes.io/name: app-2048
    spec:
      containers:
      - image: public.ecr.aws/l6m2t8p7/docker-2048:{{ or .Version "latest" }}
        imagePullPolicy: Always
        name: app-2048
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  namespace: {{ .Namespace }}
  name: service-2048
spec:
  ports:
    - port: 80
      targetPort: 80
      protocol: TCP
{{- if not .IngressHost }}
  type: LoadBalancer
{{- end }}
  selector:
    app.kubernetes.io/name: app-2048
{{- if .IngressHost }}
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  namespace: {{ .Namespace }}
  name: ingress-2048
  annotations:
    alb.ingress.kubernetes.io/scheme: internet-facing
    alb.ingress.kubernetes.io/target-type: ip
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTPS":443}]'
spec:
  ingressClassName: alb
  rules:
    - http:
        paths:
        - path: /
          pathType: Prefix
          backend:
            service:
              name: service-2048
              port:
                number: 80
  tls:
  - hosts:
    - {{ .IngressHost }}
{{- end }}
...`
