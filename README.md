diegocanaryapp
==============

**Note**: This repository should be imported as `code.cloudfoundry.org/diegocanaryapp`.

Simple canary app to test long-running Diego deployments

Usage
=====

Deploy instances of the canary app to your Runtime/Diego cluster:

```
# target the CF org and space intended for the app
export DEPLOYMENT_NAME='cf-your-deployment-diego'
./scripts/deploy
```

- The app name defaults to `diego-canary-app`, and can be overridden with the `APP_NAME` environment variable.
- The instance count defaults to 20, and can be overridden with the `INSTANCE_COUNT` environment variable.
- The domain defaults to the default domain for the targeted space, and can be overriden with the `DOMAIN` environment variable.
- For a dry run, set the `CF_COMMAND` environment variable to `'echo cf'`.
- If using a syslog forwarding service (to Papertrail, for example), specify the name of the service to bind to the app as `SERVICE_TO_BIND`.

Datadog
=======

The `datadog-config` repo has config for a diego board that has a graph for
number of instances that are up.  Instance number `n` of the app will emit the
`diego.canary.app.instance` metric with the tags `deployment:$DEPLOYMENT_NAME` 
and `diego-canary-app:n`.
To enable sending a heartbeat to DataDog, set the following environment
variable:
export DATADOG_API_KEY='your-api-key'
Pingdom
=======

Set up up/down monitoring and email alerting via Pingdom:

1. Log in to Pingdom.
2. Add an uptime check for the URL of the canary app, and choose the alert policy for the appropriate team.
