# Opinionated example of deployment via Salt Stack

## Assumptions:
  * user & group srv:srv exist, this is used as a generic service runner
  * you want to run the service on port 80, directly exposed to the interwebs (you really want to add a reverse SSL proxy in between)
  * You have, or are willing to deploy Salt Stack.
  * A bit of assembly is required since you need to move pillar.sls into your saltroot/pillar/ and the rest into saltroot/wttr/
  * You want metric-sm units. Just roll your own wegorc to change this

## Caveats:
  * Doesn't do enough to make a recent master checkout work, i.e. needs further improvement. Latest known working revision is 0d76ba4a3e112694665af6653040807835883b22
