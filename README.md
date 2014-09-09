diegocanaryapp
==============

Simple canary app to test long-running Diego deployments

Usage
=====

Deploy 20 instances of the canary app to your Runtime/Diego cluster:

```
# e.g. app_name=diego-canary-app DATADOG_API_KEY=1234notgonnatellyou DEPLOYMENT_NAME=ketchup
cf api api.$DEPLOYMENT_NAME.cf-app.com
cf login
# ...
# find or create org/space named 'canaries'/'canaries', and target
cf push $app_name --no-start
cf set-env $app_name CF_DIEGO_BETA true
cf set-env $app_name CF_DIEGO_RUN_BETA true
cf set-env $app_name DATADOG_API_KEY $DATADOG_API_KEY
cf set-env $app_name DEPLOYMENT_NAME $DEPLOYMENT_NAME
cf start $app_name
cf scale $app_name -i 20
```

Create a graph in Datadog to track that 5 instances are up and running:

1. Use a time-series graph.
2. The metric should **sum** over `diego.canary.app.instance`, filtering by the tag `deployment:$DEPLOYMENT_NAME`.
3. (Optional) Add the appropriate metric to your dashboard configuration in the `datadog-config` repo and run the rake task to push the changes to your timeboard.

Set up up/down monitoring and email alerting via Pingdom:

1. Log in to Pingdom using the Diego account in LastPass.
2. Add an up/down alert with the URL `http://$app_name.$DEPLOYMENT_NAME.cf-app.com`