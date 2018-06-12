# Introducing mouthful-spoon

Spoon is the cli helper for mouthful, allowing for migration from other commenting engines, database backups and migrations.

# Here are a few examples of working with spoon

To import comments from disqus to mouthful:
`spoon migrate disqus --dump ./disqus.xml`

> special thanks to Reddit user [doenietzomoeilijk](https://www.reddit.com/user/doenietzomoeilijk) for providing the dump to make this possible

To import comments from sqlite to dynamodb:
`spoon migrate dynamodb --sqlite ./mouthful.db ./config.json`

To import comments from isso to mouthful:
`spoon migrate isso --isso ./isso.db`

To export comments from mouthful:
`spoon export --c ./config.json`

To restore a previous dump to mouthful:
`spoon export --c ./config.json --dump ./mouthful.dmp`
