# Configuration options

This section covers the configuration options you can change in the config.json of mouthful. As the name suggests, mouthful uses JSON as its config format. This directory also contains examples of configuration files for all of the supported databases. A good place to start is the [sqlite config file](./sqlite/config.json).

## Sections

The config consists of the following sections: 

* Root
* Moderation
* Api
* Client
* Database

### Root

The root section contains all the other sections as well as the honeypot variable.

| Variables     | Use           | Type | Required  | Default value | Recommended setting |
| ------------- |:-------------:| :---:| :-------: |  :----------: |  :----------------: |
| honeypot     | determines if honeypot functionality will be used | bool | false | false | true |
| moderation     | contains the moderation settings for mouthful  | object | true |  | [see below](#Moderation) |
| api     | changes the api behaviour | object | true |  | [see below](#Api) |
| client     | changes client behaviour | object | true |  | [see below](#Client) |
| database     | allows for configuring the data store | object | true |  |  [see below](#Database)|


### Moderation

Moderation section is responsible for editing the moderation functionality.

| Variables     | Use           | Type | Required  | Default value | Recommended setting |
| ------------- |:-------------:| :---:| :-------: |  :----------: |  :----------------: |
| enabled     | determines if moderation functionality will be used. If moderation is turned on, variables below become required | bool | false | false | up to you |
| adminPassword     | sets the administration panel password. | string | true |  |  Please make sure to set it to something that's strong and not a couple of symbols long. |
| sessionDurationSeconds     | determines the length of an admin session or how long until you are forced to log in again. | int | true |  | 21600 |
| maxCommentLength     | determines the maximum comment length. Setting to a value of 0 or below allows for unlimited length | int | true | 0 | 1000 |


### Api

Api section changes the behaviour of the back-end API.

| Variables     | Use           | Type | Required  | Default value | Recommended setting |
| ------------- |:-------------:| :---:| :-------: |  :----------: |  :----------------: |
| port     | sets the port for API to run on | bool | false | 8080 | up to you |
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

You can find suggested provision ratios for dynamodb in the [example dynamodb config file](./dynamodb/config.json). They are experimental atm. Will be updated once more data on usage is collected.