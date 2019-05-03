# syncr

A data syncronisation utility that ensures the primary data records are sync'd to the document database. 

This will replace `mongr` and will _not_ require API access.

## Usage

```bash
# sync all data updated within the last 14 days
syncr -b 14 -c all
```