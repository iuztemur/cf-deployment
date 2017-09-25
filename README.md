# cf-deployment

**This repo is still a work in progress for certain use cases.
Take a look at <a href='#readiness'>this table</a>
to see if it's recommended that you use it.**

### Table of Contents
* <a href='#purpose'>Purpose</a>
* <a href='#readiness'>Is `cf-deployment` ready to use?</a>
* <a href='#deploying-cf'>Deploying CF</a>
* <a href='#contributing'>Contributing</a>
* <a href='#setup'>Setup and Prerequisites</a>
* <a href='#ops-files'>Ops Files</a>
* <a href='#ci'>CI</a>

## <a name='purpose'></a>Purpose
This repo contains a canonical manifest
for deploying Cloud Foundry without the use of `cf-release`,
relying instead on individual component releases.
It will replace the [manifest generation scripts in cf-release][cf-release-url]
when `cf-release` is deprecated.
It uses several newer features
of the BOSH director and CLI.
Older directors may need to be upgraded
and have their configurations extended
in order to support `cf-deployment`.

`cf-deployment` embodies several opinions
about Cloud Foundry deployment.
It:
- prioritizes readability and meaning to a human operator.
  For instance, only necessary configuration is included.
- emphasizes security and production-readiness by default.
  - bosh's `--vars-store` feature is used
  to generate strong passwords, certs, and keys.
  There are no default credentials, even in bosh-lite.
  - TLS/SSL features are enabled on every job which supports TLS.
- uses three AZs, of which two are used to provide redundancy for most instance groups.
The third is used only for instance groups
that should not have even instance counts,
such as etcd and consul.
- uses Diego natively,
does not support DEAs,
and enables diego-specific features
such as ssh access to apps by default.
- deploys jobs to handle platform data persistence
using the cf-mysql release for databases
and the CAPI release's WebDAV job for blob storage.
- assumes load-balancing will be handled by the IaaS
or an external deployment.

### <a name='readiness'></a> Is `cf-deployment` ready to use?

| Use Case | Is cf-deployment ready? | Blocked On |
| -------- | ----------------------- | ---------- |
| Test and development | Yes | |
| New production deployments | No | Downtime testing |
| Existing production deployments using cf-release | No | Migration tools |

We've been testing cf-deployment for some time,
and many of the development teams in the Cloud Foundry organization
are using it for development and testing.
If that describes your use case,
you can use cf-deployment as your manifest.

If you're hoping to use cf-deployment for a new _production_ deployment,
we still wouldn't suggest using cf-deployment.
We still need to be able to make some guarantees
about app availability during rolling deploys.
When we think cf-deployment is ready,
we'll update this section and make announcements on the cf-dev mailing list.

### Can I Transition from `cf-release`?
A migration will be possible.
It will be easier for some configurations
than others.

The Release Integration team
is working on a transition path from `cf-release`.
We don't advise anybody attempt the migration yet.
Our in-progress tooling and documentation can be found at
https://github.com/cloudfoundry/cf-deployment-transition

## <a name='deploying-cf'></a>Deploying CF
Deployment instructions have become verbose,
so we've moved them into a [dedicated deployment guide here](deployment-guide.md).

See the rest of this document for more on the new CLI, deployment vars, and configuring your BOSH director.

## <a name='contributing'></a>Contributing
Although the default branch for the repository is `master`,
we ask that all pull requests be made against
the `develop` branch.
Please also take a look at the ["style guide"](texts/style-guide.md),
which lays out some guidelines for adding properties or jobs
to the deployment manifest.

Before submitting a pull request,
please run `scripts/test`
which interpolates all of our ops files
with the `bosh` cli.

We ask that pull requests and other changes be successfully deployed,
and tested with the latest sha of CATs.

## <a name='setup'></a>Setup and Prerequisites
`cf-deployment` relies on newer BOSH features,
and requires a bosh director with a valid cloud-config that has been configured with a certificate authority.
It also requires the new `bosh` CLI,
which it relies on to generate and fill-in needed variables.

### BOSH CLI
`cf-deployment` requires the new [BOSH CLI](https://github.com/cloudfoundry/bosh-cli).

### BOSH `cloud-config`
`cf-deployment` assumes that
you've uploaded a compatible [cloud-config](http://bosh.io/docs/cloud-config.html) to the BOSH director.
The cloud-config produced by `bbl` is compatible by default,
which covers GCP and AWS.
For bosh-lite, you can use the cloud-config in the `bosh-lite` directory of this repo.
We have not yet tested cf-deployment against other IaaSes,
so you may need to do some engineering work to figure out the right cloud config (and possibly ops files)
to get it working for `cf-deployment`.

### Deployment variables and the var-store
`cf-deployment.yml` requires additional information
to provide environment-specific or sensitive configuration
such as the system domain and various credentials.
To do this in the default configuration,
we use the `--vars-store` flag in the new BOSH CLI.
This flag takes the name of a `yml` file that it will read and write to.
Where necessary credential values are not present,
it will generate new values
based on the type information stored in `cf-deployment.yml`.

Necessary variables that BOSH can't generate
need to be supplied as well.
Though in the default case
this is just the system domain,
some ops files introduce additional variables.
See the summary for the particular ops files you're using
for any additional necessary variables.

There are three ways to supply
such additional variables.

1. They can be provided by passing individual `-v` arguments.
   The syntax for `-v` arguments is
   `-v <variable-name>=<variable-value>`.
   This is the recommended method for supplying
   the system domain.
2. They can be provided in a yaml file
   accessed from the command line with the
   `-l` or `--vars-file` flag.
   This is the recommended method for configuring
   external persistence services.
3. They can be inserted directly in `--vars-store` file
   alongside BOSH-managed variables.
   This can confuse things,
   but you may find the simplicity worth it.

Variables passed with `-v` or `-l`
will override those already in the var store,
but will not be stored there.

## <a name='ops-files'></a>Ops Files
The configuration of CF represented by `cf-deployment.yml` is intended to be a workable, secure, fully-featured default.
When the need arises to make different configuration choices,
we accomplish this with the `-o`/`--ops-file` flags.
These flags read a single `.yml` file that details operations to be performed on the manifest
before variables are generated and filled.
We've supplied some common manifest modifications in the `operations` directory.

| Name | Purpose | Notes |
|:---  |:---     |:---   |
| [old-droplet-mitigation.yml](operations/legacy/old-droplet-mitigation.yml) | Mitigates against old droplets that may still have a legacy security vulnerability. | See comment in the ops file for more details. |
| [aws.yml](operations/aws.yml) | Overrides the loggregator ports to 4443. | It is required to have a separate port from the standard HTTPS port (443) for loggregator traffic in order to use the AWS load balancer. |
| [disable-router-tls-termination.yml](operations/disable-router-tls-termination.yml) | Eliminates keys related to performing tls/ssl termination within the gorouter job. | Useful for deployments where tls termination is performed prior to the gorouter - for instance, on AWS, such termination is commonly done at the ELB. This also eliminates the need to specify `((router_ssl.certificate))` and `((router_ssl.private_key))` in the var files. |
| [configure-default-router-group.yml](operations/configure-default-router-group.yml) | Allows deployer to configure reservable ports for default tcp router group by passing variable `default_router_group_reservable_ports`. |  |
| [enable-privileged-container-support.yml](operations/enable-privileged-container-support.yml) | Enables diego privileged container support on cc-bridge. | This opsfile might not be compatible with opsfiles that inline bridge functionality to cloud-controller. |
| [gcp.yml](operations/gcp.yml) | Intentionally left blank for backwards compatibility. | It previously overrode the static IP addresses assigned to some instance groups, as GCP networking features allow them to all co-exist on the same subnet despite being spread across multiple AZs. |
| [rename-deployment.yml](operations/rename-deployment.yml) | Allows a deployer to rename the deployment by passing a variable `deployment_name` |  |
| [rename-network.yml](operations/rename-network.yml) | Allows a deployer to rename the network by passing a variable `network_name` |  |
| [scale-to-one-az.yml](operations/scale-to-one-az.yml) | Scales cf-deployment down to a single instance per instance group, placing them all into a single AZ. | Effectively halves the deployment's footprint. Should be applied before other ops files. |
| [tcp-routing-gcp.yml](operations/tcp-routing-gcp.yml) | Intentionally left blank for backwards compatibility. | Previously added TCP routers for GCP. `cf-deployment.yml` now does this by default. |
| [use-blobstore-cdn.yml](operations/use-blobstore-cdn.yml) | Adds support for accessing the `droplets` and `resource_pool` blobstore resources via signed urls over a cdn. | This assumes that you are using the same keypair for both buckets. Introduces [new variables](operations/example-vars-files/vars-use-blobstore-cdn.yml) |
| [use-external-dbs.yml](operations/use-external-dbs.yml) | Removes the MySQL instance group, cf-mysql release, and all cf-mysql variables. **Warning**: this does not migrate data, and will delete existing database instance groups. | This requires an external data store.   Introduces [new variables](operations/example-vars-files/vars-use-external-dbs.yml) for DB connection details which will need to be provided at deploy time. This must be applied _before_ any ops files that removes jobs that use a database, such as the ops file to remove the routing API. |
| [use-postgres.yml](operations/use-postgres.yml) | Replaces the MySQL instance group with a postgres instance group. **Warning**: this will lead to total data loss if applied to an existing deployment with MySQL or removed from an existing deployment with postgres. |  |
| [use-s3-blobstore.yml](operations/use-s3-blobstore.yml) | Replaces local WebDAV blobstore with external s3 blobstore. | Introduces [new variables](operations/example-vars-files/vars-use-s3-blobstore.yml) for s3 credentials and bucket names. |
| [use-azure-storage-blobstore.yml](operations/use-azure-storage-blobstore.yml) | Replaces local WebDAV blobstore with external Azure Storage blobstore. | Introduces [new variables](operations/example-vars-files/vars-use-azure-storeage-blobstore.yml) for Azure credentials and container names. |
| [windows-cell.yml](operations/windows-cell.yml) | Deploys a windows diego cell, adds releases necessary for windows. |  |
| [bosh-lite.yml](operations/bosh-lite.yml) | Enables `cf-deployment` to be deployed on `bosh-lite`. | See [bosh-lite](iaas-support/bosh-lite/README.md) documentation. |
| [bypass-cc-bridge-privileged-containers.yml](operations/bypass-cc-bridge-privileged-containers.yml) | Use privileged containers for staging and running buildpack apps and tasks. | |
| [bypass-cc-bridge.yml](operations/bypass-cc-bridge.yml) | Bypass CC bridge | To enable privileged container support, also apply the bypass-cc-bridge-privileged-containers.yml ops file. |
| [cf-syslog-skip-cert-verify.yml](operations/cf-syslog-skip-cert-verify.yml) | This disables TLS verification when connecting to a HTTPS syslog drain. | |
| [enable-cc-rate-limiting.yml](operations/enable-cc-rate-limiting.yml) | Enable rate limiting for UAA-authenticated endpoints. | Introduces variables `cc_rate_limiter_general_limit` and `cc_rate_limiter_unauthenticated_limit` |
| [enable-uniq-consul-node-name.yml](operations/enable-uniq-consul-node-name.yml) | Configure Diego cell `consul_agent` jobs to have a unique id per instance. |  |
| [scale-down-etcd-for-cluster-changes.yml](operations/scale-down-etcd-for-cluster-changes.yml) | Scales `etcd` cluster to one node. |  |
| [use-compiled-releases.yml](operations/use-compiled-releases.yml) | Instead of having your BOSH Director compile each release, use this opsfile to use pre-compiled releases for a deployment speed improvement. |  |
| [use-latest-stemcell.yml](operations/use-latest-stemcell.yml) | Use the latest stemcell available on your BOSH director instead of the one in `cf-deployment.yml` |  |

### A note on `experimental` and `test` ops-files
The `operations` directory includes two subdirectories
for "experimental" and "test" ops-files.

#### Experimental
"Experimental" ops-files represent configurations
that we expect to promote to blessed configuration eventually,
meaning that,
once the configurations have been sufficiently validated,
they will become part of cf-deployment.yml
and the ops-files will be removed.
For details, see [experimental README](operations/experimental/README.md).

#### Test
"Test" ops-files are configurations
that we run in our testing pipeline
to enable certain features.
We include them in the public repository
(rather than in our private CI repositories)
for a few reasons,
depending on the particular ops-file.

Some files are included
because we suspect that the configurations will be commonly needed
but not easily generalized.
For example,
`add-persistent-isolation-segment.yml` shows how a deployer can add an isolated Diego cell,
but the ops-file is hard to apply repeatably.
In this case, the ops-file is an example.

Others,
like `cfr-to-cfd-transition.yml`,
will eventually be promoted to the `operations` directory,
but are still being modified regularly.
In this case, the ops-file is included for public visibility.

## <a name='ci'></a>CI
The [ci](https://release-integration.ci.cf-app.com/teams/main/pipelines/cf-deployment) for `cf-deployment`
automatically bumps to the latest versions of its component releases on the `develop` branch.
These bumps, along with any other changes made to `develop`, are deployed to a single long-running environment
and tested with CATs before being merged to master if CATs goes green.
There is not presently any versioning scheme,
or way to correlate which version of CATs is associated with which sha of cf-deployment,
other than the history in CI.
As `cf-deployment` matures, we'll address versioning.
The configuration for our pipeline can be found [here](https://github.com/cloudfoundry/runtime-ci/pipelines/cf-deployment.yml).

[cf-deployment-concourse-url]: https://release-integration.ci.cf-app.com/teams/main/pipelines/cf-deployment
[cf-release-url]: https://github.com/cloudfoundry/cf-release/tree/master/templates
