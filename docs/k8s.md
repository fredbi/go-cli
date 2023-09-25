# Configuring a Kubernetes deployment

Typically, the configuration files are held in one or several `Configmap` resources, mounted by your deployed container.

Secret files can be mounted from `Secret` resources in the container, accessible as plain files.

Alternatively, Kubernetes may expose secrets as environment variables: `viper` takes care of loading them in the registry.

Normally, we don't want to expose secrets via CLI flags.

Example (e.g. volumes & container section of a k8s PodTemplateSpec):
```yaml
volumes:
  - name: config
    configMap:          # <- expose config file from ConfigMap resource to the pod's containers
      name: 'app-config'
  - name: secret-config # <- expose secrets file from Secret as file resource to the pod's containers
    secret:
      secret_name: 'app-secret-config'

containers:
  - name: app-container
    ...
    env:
      - name: CONFIG_DIR
        value: '/etc/app'
      - name: SECRET_URL
        valueFrom:
          secretKeyRef: # <- expose config value as an environment variable to the container
          name: 'app-secret-url'
          key: secretUrl

    volumeMounts:
      - mountPath: '/etc/app' # <- mount config file(s) as /etc/app/{key(s)} file(s)
        name: config
      - mountPath: '/etc/app/config.d'
```
