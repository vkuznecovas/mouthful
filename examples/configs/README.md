# Configuration options

This section covers the configuration options you can change in the config.json of mouthful. As the name suggests, mouthful uses JSON as its config format. This directory also contains examples of configuration files for all of the supported databases. A good place to start is the [sqlite config file](./sqlite/config.json).

## Sections

The config consists of the following sections: 

* Root
* Moderation
* Api
* Client
* Notification
* Database

### Root

The root section contains all the other sections as well as the honeypot variable.

| Variables     | Use           | Type | Required  | Default value | Recommended setting |
| ------------- |:-------------:| :---:| :-------: |  :----------: |  :----------------: |
| honeypot     | determines if honeypot functionality will be used | bool | false | false | true |
| moderation     | contains the moderation settings for mouthful  | object | true |  | [see below](#Moderation) |
| api     | changes the api behaviour | object | true |  | [see below](#Api) |
| client     | changes client behaviour | object | true |  | [see below](#Client) |
| notification     | changes notification behaviour  | object | true |  |  [see below](#Notification)|
| database     | allows for configuring the data store | object | true |  |  [see below](#Database)|


### Moderation

Moderation section is responsible for editing the moderation functionality.

| Variables     | Use           | Type | Required  | Default value | Recommended setting |
| ------------- |:-------------:| :---:| :-------: |  :----------: |  :----------------: |
| enabled     | determines if moderation functionality will be used. If moderation is turned on, variables below become required | bool | false | false | up to you |
| adminPassword     | sets the administration panel password. | string | true |  |  Please make sure to set it to something that's strong and not a couple of symbols long. |
| sessionDurationSeconds     | determines the length of an admin session or how long until you are forced to log in again. | int | true |  | 21600 |
| maxCommentLength     | determines the maximum comment length. Setting to a value of 0 or below allows for unlimited length | int | true | 0 | 1000 |
| maxAuthorLength     | determines the maximum author length. Setting to a value of 3 or below defaults to no limit | int | true | 50 | 35 |
| path     | the path you'll run the admin panel from | string | false | "/" | none |
| oauthCallbackOrigin | the base url of your API | string | true if using oauth | "" | fully fledged url of your admin panel |
| disablePasswordLogin | disables the passsword authentication for admin panel if set to true | bool | false | false | true if using oauth, false otherwise | 
| oauthProviders | determines which oauth providers will be used for mouthful admin panel, [see below](#oauth-providers)| array | false | none | your preference |
| periodicCleanup | determines if periodic cleanup is used and all its preferences, [see below](#periodic-cleanup)| object | false | none | your preference |

#### Oauth providers

The oauth providers is responsible for setting up your mouthful installation for oauth use. You can use as many providers as you like, or as few as you want. For an example config, head to [example oauth config file](./oauth/config.json)

| Variables     | Use           | Type | Required  | Default value | Recommended setting |
| ------------- |:-------------:| :---:| :-------: |  :----------: |  :----------------: |
| enabled     | determines if the provider is used or not | bool | false | false | up to you |
| name | Name of the oauth provider. The list is [available below](#supported-oauth-providers)| string | true | none | up to you |
| secret | Secret of the oauth provider. You'll have to head to the providers page to figure it out. | string | true | none | up to you |
| key | Key or id of the oauth provider. You'll have to head to the providers page to figure it out. | string | true | none | up to you |
| adminUserIds | Ids of the users that will be assigned admin status. | array of strings | true | none | up to you |

#### Periodic cleanup

Periodic cleanup enables mouthful to run background jobs cleaning  old deleted comments as well as unconfirmed comments that have not been confirmed for a predetermined amount of seconds. 
| Variables     | Use           | Type | Required  | Default value | Recommended setting |
| ------------- |:-------------:| :---:| :-------: |  :----------: |  :----------------: |
| enabled     | determines if the cleanup functionality is used or not | bool | true | false | up to you |
| removeDeleted | determines if the cleanup job will clean soft deleted comments | bool | true | false | up to you |
| removeUnconfirmed | determines if the cleanup job will clean unconfirmed comments | bool | true | false | up to you |
| unconfirmedTimeoutSeconds | the amount of seconds that it takes for an unconfirmed comment to be marked for deletion | int | true | none | up to you |
| deletedTimeoutSeconds | the amount of seconds that it takes for a deleted comment to be marked for deletion | int | true | none | up to you |
| removeDeletedPeriodSeconds | determines how often the deletion job for soft-deleted comments is run | int | false | 86400 | up to you |
| removeUnconfirmedPeriodSeconds | determines how often the deletion job for unconfirmed comments is run | int | false | 86400 | up to you |

If this all seems confusing, [see the example](./cleanup/config.json) and [its readme](./cleanup/README.md).


##### Supported Oauth providers

Currently mouthful supports 37 oauth providers:

* amazon
* battlenet
* bitbucket
* box
* dailymotion
* deezer
* digitalocean
* discord
* dropbox
* eveonline
* facebook
* fitbit
* github
* gitlab
* gplus
* heroku
* influxcloud
* instagram
* intercom
* lastfm
* linkedin
* meetup
* microsoftonline
* naver
* onedrive
* salesforce
* slack
* soundcloud
* spotify
* stripe
* twitch
* twitter
* uber
* vk
* wepay
* xero
* yahoo
* yammer

### Api

Api section changes the behaviour of the back-end API.

| Variables     | Use           | Type | Required  | Default value | Recommended setting |
| ------------- |:-------------:| :---:| :-------: |  :----------: |  :----------------: |
| port     | sets the port for API to run on | bool | false | 8080 | up to you |
| bindAddress | sets the address that the api will listen on | string | 0.0.0.0 | up to you |
| logging     | determines if gin logging will be enabled for the api | bool | false | true | true
| debug     | enables or disables the gin debug mode with more verbal logging. | bool | false | false | false |
| cache     | cache settings for the api | object | true |  | [see below](#api.cache) |
| cors     | cors settings for the api | object | true |  | [see below](#api.cors) |
| rateLimiting     | rate limiting settings for the api | object | true |  | [see below](#api.rateLimiting) |


#### api.cache

The cache section determines the API cache behaviour.

| Variables     | Use           | Type | Required  | Default value | Recommended setting |
| ------------- |:-------------:| :---:| :-------: |  :----------: |  :----------------: |
| enabled     | determines if cache functionality will be used. If cache is turned on, variables below become required | bool | false | false | up to you |
| expiryInSeconds     | determines the cache expiry time | int | true |  | 300 |
| intervalInSeconds     | determines how often we'll check for expired cache items | int | true |  | 10 |

#### api.cors

The cors section determines which origins will be allowed to access your backend. 

| Variables     | Use           | Type | Required  | Default value | Recommended setting |
| ------------- |:-------------:| :---:| :-------: |  :----------: |  :----------------: |
| enabled     | determines if cors functionality will be used. If cors is turned on, variables below become required | bool | false | false | true |
| allowedOrigins     | changes the list of allowed origins | array of strings | true |  | The address of your website |

#### api.rateLimiting

The rate limiting section is responsible for setting limitations on post rate for users.


| Variables     | Use           | Type | Required  | Default value | Recommended setting |
| ------------- |:-------------:| :---:| :-------: |  :----------: |  :----------------: |
| enabled     | determines if rateLimiting functionality will be used. If rateLimiting is turned on, variables below become required | bool | false | false | up to you |
| allowedOpostsHourrigins     | how many posts a single user is allowed to make per hour | int | true | | 100 |


### Client

The client section is responsible for setting the client side behaviour.

| Variables     | Use           | Type | Required  | Default value | Recommended setting |
| ------------- |:-------------:| :---:| :-------: |  :----------: |  :----------------: |
| useDefaultStyle     | determines if the default mouthful styling will be applied to the client | bool | false | false | false, if you'll override the styling |
| pageSize     | a limit on how many posts/replies to show | int | true | | 10 |

### Notification

The notification section is responsible for setting notification behaviour for new comments.

| Variables     | Use           | Type | Required  | Default value | Recommended setting |
| ------------- |:-------------:| :---:| :-------: |  :----------: |  :----------------: |
| enabled     | determines if notifications about new comments will be send | bool | false | false | up to you |
| url     | url to send an http post request to when a new comment is received  | string | true | | up to you |

### Database

The database section determines the data source mouthful will use. 

| Variables     | Use           | Type | Required  | Default value | Recommended setting |
| ------------- |:-------------:| :---:| :-------: |  :----------: |  :----------------: |
| dialect     | determines the type of the database we'll use. All the remaining fields depend on this value. | string | true |  | sqlite3, because comments are not big data |



Currently supported databases are:
#### sqlite

[Here's an example sqlite config file](./sqlite/config.json)

| Variables     | Use           | Type | Required  | Default value | Recommended setting |
| ------------- |:-------------:| :---:| :-------: |  :----------: |  :----------------: |
| database     | path to the database file sqlite will use | string | true |  | ./mouthful.db |


#### mysql

[Here's an example mysql config file](./mysql/config.json)

| Variables     | Use           | Type | Required  | Default value | Recommended setting |
| ------------- |:-------------:| :---:| :-------: |  :----------: |  :----------------: |
| database     | database name | string | true |  | mouthful |
| host     | database host | string | true |  | localhost |
| username     | username we'll use to connect to the database | string | true |  | mouthful_user |
| password     | password  we'll use to connect to the databse | string | true |  | a strong password |
| port     | port  we'll use to connect to the databse | string | true |  | 3306 |


#### postgres


[Here's an example postgres config file](./postgres/config.json)

| Variables     | Use           | Type | Required  | Default value | Recommended setting |
| ------------- |:-------------:| :---:| :-------: |  :----------: |  :----------------: |
| database     | database name | string | true |  | mouthful |
| host     | database host | string | true |  | localhost |
| username     | username we'll use to connect to the database | string | true |  | mouthful_user |
| password     | password  we'll use to connect to the databse | string | true |  | a strong password |
| port     | port  we'll use to connect to the databse | string | true |  | 3306 |
| sslEnabled     | determines if we'll use ssl to connect to the database | bool | false |  | Depends on your db setup |


#### dynamodb


[Here's an example dynamodb config file](./dynamodb/config.json)

| Variables     | Use           | Type | Required  | Default value | Recommended setting |
| ------------- |:-------------:| :---:| :-------: |  :----------: |  :----------------: |
| dynamoDBThreadReadUnits     | read units for the thread table | int | true |  | depends on your load |
| dynamoDBCommentReadUnits    | read units for the comment table | int | true |  | depends on your load |
| dynamoDBThreadWriteUnits     | write units for the thread table | int | true |  | depends on your load |
| dynamoDBCommentWriteUnits     | write units for the comment table | int | true |  | depends on your load |
| dynamoDBIndexWriteUnits     | write units for the comment index | int | true |  | depends on your load |
| dynamoDBIndexReadUnits    | write units for the comment index | int | true |  | depends on your load |
|awsRegion | determines the aws region we'll connect to | string | true |  | |
|awsAccessKeyID | your aws access key id. Can be overriden by env variable `AWS_ACCESS_KEY_ID` | string | true |  |  |
|awsSecretAccessKey | your secret aws access key. Can be overriden by env variable `AWS_SECRET_ACCESS_KEY` | string | true |  |  |
| dynamoDBEndpoint | endpoint of the aws dynamodb. Mostly used for testing. | string | false | none | none |



You can find suggested provision ratios for dynamodb in the [example dynamodb config file](./dynamodb/config.json). They are experimental atm. Will be updated once more data on usage is collected.


