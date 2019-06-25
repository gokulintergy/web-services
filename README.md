# MappCPD Web Services

[![Build Status](https://travis-ci.org/mappcpd/web-services.svg?branch=master)](https://travis-ci.org/mappcpd/web-services)

Web services supporting MappCPD / HeartOne.

- [cmd/](/cmd/README.md) - executable packages
  - [algr/](/cmd/algr/README.md) - worker to sync Algolia indexes
  - [backupdb](/cmd/backupdb/README.md) - worker to backup MySQL database to Dropbox
  - [fixr/](/cmd/fixr/README.md) - utility to check and fix data
  - [pubmedr/](/cmd/pubmedr/README.md) - worker to fetch pubmed articles
  - [syncr/](/cmd/syncr/README.md) - worker to sync data from MySQL to MongoDB
  - [webd/](/cmd/webd/README.md) - web services API
- [internal/](/internal/README.md) - internal packages

## Configuration

Sample `.env` shows required env vars. Most of these are required to be present
even if they have dummy values. For example, if the Amazon SES email service
will not be used the config values must still exist or the web services server
will not start.

Note these are in alphabetical order, rather than being grouped into related services.

```bash

# AWS S3 -----------------------------------------------------------------------

# The API currently plays a small role in facilitating direct uploads 
# from the client to S3, thus bypassing the server. To do this the API 
# issues a signed url to the client requiring the following credentials:
AWS_ACCESS_KEY_ID="ABC....RST"
AWS_REGION="ap-southeast-1"
AWS_SECRET_ACCESS_KEY="fghjkl...asdfg"

# AWS Simple Email Service --------------------------------------------------------
# Used for package https://github.com/8o8/email that wraps email services, these 
# are required to switch over to SES.
AWS_SES_ACCESS_KEY_ID="AK....GH"
AWS_SES_REGION="us-east-1"
AWS_SES_SECRET_ACCESS_KEY="OyL....8qW"

# Dropbox access token for backupdb service
DROPBOX_ACCESS_TOKEN="OAAA....rKwv"

# Mailgun email sending service
MAILGUN_API_KEY="key-bab...a96"
MAILGUN_DOMAIN="mx.csanz.edu.au"

# Admin creds to access the API endpoints, required by some services
MAPPCPD_ADMIN_USER="admin-user"
MAPPCPD_ADMIN_PASS="admin-pass"

# Algolia index creds
MAPPCPD_ALGOLIA_API_KEY="01f...c09"
MAPPCPD_ALGOLIA_APP_ID="IYL...BL7"

# Member titles to EXCLUDE when updating the directory index
MAPPCPD_ALGOLIA_DIRECTORY_EXCLUDE_TITLES="non-member,trainee,admin member,applicant,limited access"

# Algolia index names
MAPPCPD_ALGOLIA_DIRECTORY_INDEX="DIRECTORY"
MAPPCPD_ALGOLIA_MEMBERS_INDEX="MEMBERS"
MAPPCPD_ALGOLIA_MODULES_INDEX="MODULES"
MAPPCPD_ALGOLIA_ORGANISATIONS_INDEX="ORGANISATIONS"
MAPPCPD_ALGOLIA_QUALIFICATIONS_INDEX="QUALIFICATIONS"
MAPPCPD_ALGOLIA_RESOURCES_INDEX="RESOURCES"


# BASE URL of the web API
MAPPCPD_API_URL="https://mappcpd-api.io"

# Token stuff
MAPPCPD_JWT_SIGNING_KEY="anyTokenSigningKey"
MAPPCPD_JWT_TTL_HOURS=4

# MongoDB
MAPPCPD_MONGO_DBNAME="dbname"
MAPPCPD_MONGO_DESC="Descriptive MongoDB source - shows in responses"
MAPPCPD_MONGO_URL="mongodb://user:pass@host1:port,host2:port/databasename?replicaSet=??????"

# Specify MX service - mailgun,sendgrid, or ses
MAPPCPD_MX_SERVICE="mailgun"

# MySQL
MAPPCPD_MYSQL_DESC="Descriptive MySQL source - shows in responses"
MAPPCPD_MYSQL_URL="user:pass@tcp(host:3306)/dbname"

# Pubmed
# config for pubmed service
MAPPCPD_PUBMED_BATCH_FILE="https://url.to.jsonfile.com/file.json"
# max results to return per request
MAPPCPD_PUBMED_RETMAX=200

# Short link redirection
# prefix for shortlink IDs, eg 'r' in http://link.io/r123
MAPPCPD_SHORT_LINK_PREFIX="r"
# base RUL for short link redirector (linkr)
MAPPCPD_SHORT_LINK_URL="https://link.to"

# Sendgrid email service
SENDGRID_API_KEY="SG.fHT...Tga"
```

## Services architecture

How the services fit together:

![resources](https://docs.google.com/drawings/d/1zJ4pQCb94syzpCvoqRBXwbMUvs8LhpFlFE2Gax6LTfM/pub?w=691&h=431)

The `pubmedr` service request new articles from Pubmed (**1**), processes
responses (**2**) and accesses the API (**3**) to add batches of articles to the resources table in MySQL (**4**).

The `syncr` service is responsible for the synchronisation of required data from
the primary MySQL store, to the MongoDB database (**5**). Its main
responsibility is the copying of new or updated member and resource records.

Some API actions will trigger an update of the corresponding member record
immediately. For example, when a field in the in the `member` table is updated,
the corresponding MongoDB document is updated immediately afterwards. However,
if data in a related member table is updated, the corresponding MongoDB document
may be updated only when `syncr` is run.

The `algr` service fetches document records from MongoDB (**6**) and updates the
corresponding Algolia indexes (**7**). It can do partial or complete index rebuilds
depending on flags. Complete rebuilds are minimised on large indexes to reduce
the number of indexing operations.

[`linkr`](https://github.com/34South/linkr) is a (short) link redirection
service. Resource search results from the Algolia index are delivered directly
to users via a javascript client (**8**). Resource links go first to the `linkr`
service (**9**) which records stats in the MongoDB database (**10**) and then redirect to the resource URL (**11**).

**Note:** Initially these services were separate projects and each accessed the
databases via the API. They have since been integrated and updated to access the
database directly using the `internal` packages. The `pubmedr` service still
uses the API and will also be updated at some stage. The `linkr` service runs
as a separate, standalone web srever.

See [MappCPD Architecture](https://github.com/mappcpd/architecture/wiki) for more info.

## References

Project structure based on Bill Kennedy's [package oriented
design](https://www.goinggo.net/2017/02/package-oriented-design.html).
