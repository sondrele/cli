# debuk

An Aiven Kafka Debug Cli to apply a `aivenApplication` and extract its credentials. The tool will apply
an [Protected](https://doc.nais.io/persistence/kafka/#accessing-topics-from-an-application-on-legacy-infrastructure) `aivenApplication`
in your specified namespace=`team`. This will give access to personal but time limited credential. This credential can
be used to debug an Aiven hosted kafka topic. **Debuk** extracts the fresh `aivenApplication` credentials and puts them
in `current` folder. The applied `aivenApplication` has a default `expireAt` and is set to 1 day.

## Prerequisite

* Authentication & Authorization
    * Connect to [naisdevice](https://doc.nais.io/device/)
    * Tool is used in GCP? please be sure to log in:

```
gcloud auth login --update-adc
```

* Update your [topic](https://doc.nais.io/persistence/kafka/#creating-topics-and-defining-access) resource & ACLs
    * Add `username` to your `Topic`.yaml ACLs and apply to your namespace.


* Install [kcat](https://github.com/edenhill/kcat) (kcat is the project formerly known as kafkacat) in preferred way.

## Use

First;

```
brew tap nais/tap
```

then;

```
brew install debuk  
```

check;

```
debuk version
```

You should be able to use command;

```
debuk [commands] [flags]
```

Available commands:

- apply
- version
- get

For help on individual commands, add `--help` short: `-h`.

### Flags

Flags provide modifiers to control how the action command operates.

#### Required

* apply
    * `--username`
        * short `-u`: Prefix before `@nav.no`.

    * `--team`
        * short `-t`: teamnamespace (default not supported).
* get
    * `--secret-name`
        * short `-s`: Secretname for your aiven application.

### Optional

* apply
    * `--pool`
        * short `-p` default: `nav-dev`: Preferred kafka pool.

    * `--expire`
        * short `-e` default: `1`: Time in days the created secret should be valid.

    * `--dest`
        * short `-d` default: `current`: Path to directory where secrets will be dropped of. For `current` with
          subfolder folder, e.g: `/.config`

    * `--secret-name`
        * short `-s` default: `namespace-username-(random-id)`: Preferred secret-name instead of the generated.

* version
    * `--commit`
        * short `-i` default: `false` : Get detailed information about this debuk version

* get
    * `--dest`
        * short `-d` default: `current`: Path to directory where secrets will be dropped of. For `current` with
          subfolder folder, e.g: `/.config`
    * `--config`
        * short `-c`: default: `all`: Config type, `all || kcat || .env`. `all` generates both .env and kcat config
          files.

## Available configuration files

After Successful `debuk` command a set of files will be available in `current` folder.

### All

- client-keystore.p12
- client-truststore.jks
- `username`.yaml (the yaml applied to specified namespace)
- kafka-ca.cert
- kafka-certificate.crt
- kafka-private-key.pem
- kafka-secret.env
- kcat.conf

#### Example

##### `username`.yaml

```yaml
apiVersion: aiven.nais.io/v1
kind: AivenApplication
metadata:
  name: my-user
  namespace: my-namespace
spec:
  secretName: default-my-user-d1037d
  protected: true
  kafka:
    pool: nav-dev
  userSpec:
    timeToLive: 1
```

### .env

- client-keystore.p12
- client-truststore.jks
- kafka-ca.cert
- kafka-certificate.crt
- kafka-private-key.pem
- kafka-secret.env

#### kafka-secret.env

```Properties
KAFKA_BROKERS:brokerurl.aivencloud.com:26484
KAFKA_PRIVATE_KEY=/path/to/kafka-private-key.pem
client.keystore.p12=/path/to/client-keystore.p12
client.truststore.jks=/path/to/.envs/client-truststore.jks
KAFKA_CA=/path/to/.envs/kafka-ca.cert
KAFKA_CERTIFICATE=/path/to/.envs/kafka-certificate.crt
KAFKA_CREDSTORE_PASSWORD:password
KAFKA_SCHEMA_REGISTRY:https://registry-url.aivencloud.com:26487
KAFKA_SCHEMA_REGISTRY_PASSWORD:password
KAFKA_SCHEMA_REGISTRY_USER:my-user
```

#### Example

### kcat

- kafka-ca.cert
- kafka-certificate.crt
- kafka-private-key.pem
- kcat.conf

#### Example

##### kcat.conf

```Properties
# Debuked 2021-09-01 15:26:00
# kcat -F kcat.conf
ssl.key.location=/path/to/current/folder/creds/my-username-kafka-private-key.pem
ssl.certificate.location=/path/to/current/folder/creds/my-username-kafka-certificate.crt
ssl.key.password=changeme
bootstrap.servers=https://boostrap-server.aivencloud.com:26484
ssl.ca.location=/path/to/current/folder/creds/my-user-kafka-ca.cert
security.protocol=ssl
```

`kcat.conf` can be used with [aiven-kcat](https://help.aiven.io/en/articles/2607674-using-kafkacat) to authenticate
against the Aiven hosted topics in GCP.

You can refer to this config with -F flag:

```
kcat -F kcat.conf
```

Alternatively, you can specify the same settings directly on the command line:

```
kcat \
    -b https://boostrap-server.aivencloud.com:26484 \
    -X security.protocol=ssl \
    -X ssl.key.location=service.key \
    -X ssl.certificate.location=service.cert \
    -X ssl.ca.location=ca.pem
```

For more details [aiven-kcat-help](https://help.aiven.io/en/articles/2607674-using-kafkacat)

## Flow

![Debuk under the hood](doc/debuk.png)

## Local Development

* Be sure to run your local cluster, recommend: [minkube](https://minikube.sigs.k8s.io/docs/start/).

Start minikube with a version < 1.22,
reason: [Feature removals](https://kubernetes.io/blog/2021/07/14/upcoming-changes-in-kubernetes-1-22/).

```
minikube start --kubernetes-version=v1.21.4
```

* Apply liberator CRDs.

```
kubectl apply -f path/to/liberator/crd/bases
```

* Create a `test` cluster.

```
kubectl create namespace test
```

* Create a [secret](https://doc.nais.io/persistence/kafka/#application-config) containing this data.

```yaml
apiVersion: v1
data:
  KAFKA_BROKERS: ...
  KAFKA_CA: ...
  KAFKA_CERTIFICATE: ...
  KAFKA_CREDSTORE_PASSWORD: ...
  KAFKA_PRIVATE_KEY: ...
  KAFKA_SCHEMA_REGISTRY: ...
  KAFKA_SCHEMA_REGISTRY_PASSWORD: ...
  KAFKA_SCHEMA_REGISTRY_USER: ...
  client.keystore.p12: ...
  client.truststore.jks: ...
kind: Secret
metadata:
  annotations:
    aivenator.aiven.nais.io/protected: "true"
    aivenator.aiven.nais.io/with-time-limit: "true"
    kafka.aiven.nais.io/pool: nav-test
    kafka.aiven.nais.io/serviceUser: service-user
  finalizers:
    - aivenator.aiven.nais.io/finalizer
  labels:
    team: test
    type: aivenator.aiven.nais.io
  name: test-user
  namespace: test
type: Opaque
```

```
kubectl apply -f path/to/secret
```

* Generate executable program and test your changes.

```
make debuk
```