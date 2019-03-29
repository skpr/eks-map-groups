EKS: Map Groups
===============

Translates IAM groups into EKS `mapUsers` for aws-iam-authenticator.

## Example

**group.yml**

```yaml
- name: group-a
  username: user-a
  groups:
    - system:masters
- name: group-b
  username: user-b
  groups:
    - system:masters
```

Run the command:

`eks-map-groups --file=group.yml --namespace=kube-system --configmap=aws-auth`

Synchronises the `mapUsers` field on the `kube-system/aws-auth` ConfigMap.

```yaml
data:
  mapUsers: |
    - userarn: arn:aws:iam::xxxxxxxxxxx:user/bob
      username: user-a
      groups:
      - system:masters
    - userarn: arn:aws:iam::xxxxxxxxxxx:user/tom
      username: user-b
      groups:
      - system:masters
```
