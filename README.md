# Pod Exec Defender

This project is a sample of use-case when pod/exec subresources needs to be matched by labels.

It's just webhook server that allows you to labels your pods with `exec-defender.sleshche.com: "activated"` 
and forbid pod/exec despite of user's RBAC permissions.

`make install` will install this sample on your cluster.

After it's done you can try:
```
  kubectl apply -f ./sample/protected-pod.yaml
```

and once it's started you won't be able to do:
```
  kubectl exec protected-pod -- echo hello
```
due error:
```
Error from server (You can't connect to pods which are labeled with `prevent-exec.defender.test.com`): admission webhook "webhook-server.podexec-defender.svc" denied the request: You can't connect to pods which are labeled with `prevent-exec.defender.test.com`
```

The motivation behind feature request with labels selector, when your webhook server is not available,
pod/exec must be blocked into pods with `exec-defender.sleshche.com` label but not to every pod on the cluster.

But when you use:
```
  objectSelector:
    matchExpressions:
      - key: exec-defender.sleshche.com
        operator: Exists
  rules:
  - apiGroups:
    - ""
    apiVersions:
    - v1
    operations:
    - CONNECT
    resources:
    - pods/exec
    scope: '*'
```
it won't match any pod/execs since [PodExecOption](https://github.com/kubernetes/kubernetes/blob/9cdd673a8b784f22709b133dde1da16332b14889/staging/src/k8s.io/api/core/v1/types.go#L5259) does not have labels fields
so, any object selector will lead to matching nothing.


