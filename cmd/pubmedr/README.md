# pubmedr

MappCPD worker that fetches articles from [Pubmed](https://www.ncbi.nlm.nih.gov/pubmed/) 
and inserts them into the primary database table: `ol_resource`. From there they will 
be picked up by `mongr` -> `algr` and then visible in the Resource library search.

![resources](https://docs.google.com/drawings/d/1zJ4pQCb94syzpCvoqRBXwbMUvs8LhpFlFE2Gax6LTfM/pub?w=691&h=431)

## Configuration

**Env vars**

```bash
# Admin auth credentials
MAPPCPD_ADMIN_PASS="demo-user"
MAPPCPD_ADMIN_USER"demo-pass"

# API
MAPPCPD_API_URL="https://mappcpd-api.com"

# Pubmed return batch size - too big can time out
MAPPCPD_PUBMED_RETMAX=200

# Local or remote (http) JSON config file (see below)
BATCH_FILE="https://s3-ap-southeast-1.amazonaws.com/demo-mappcpd/public/pubmedr/pubmed.json"
```

The option for a remote batch config file was added so that the config 
did not have to be uploaded with the repo. It contains an array of one 
or more Pubmed fetch configurations, eg:

```json
[{
		"run": true,
		"category": "Cardiology",
		"searchTerm": "%28cardiology%5BMeSH%20Terms%5D%20OR%20cardiology%5BAll%20Fields%5D%29%20AND%20loattrfree%20full%20text%5BFilter%5D%20AND%20medline%5BFilter%5D%20AND%20jsubsetaim%5Btext%5D",
		"relDate": 1000,
		"attributes": {
			"category": "Cardiology",
			"free": true,
			"public": true,
			"source": "Pubmed"
		},
		"resourceTypeID": 80
	},
	{
		"run": true,
		"category": "Physiotherapy",
		"searchTerm": "%28physiotherapy%5BMeSH%20Terms%5D%20OR%20physiotherapy%5BAll%20Fields%5D%29%20AND%20loattrfree%20full%20text%5BFilter%5D%20AND%20medline%5BFilter%5D%20AND%20jsubsetaim%5Btext%5D",
		"relDate": 1000,
		"attributes": {
			"category": "Physiotherapy",
			"free": true,
			"public": true,
			"source": "Pubmed"
		},
		"resourceTypeID": 80
	}
]
```

The intention was to be able to change the configuration without having to commit 
a code update. However, it is a bit awkward to update so a local `pubmed.json` was 
added to the _root_ dir of the repo.

**Note: The `pubmed.json` is hard-coded and should be accessed via a flag. 
See issue #72.

As the time of writing it looks like this:

```json
[
  {
    "run": true,
    "category": "Cardiology",
    "searchTerm": "loattrfree%20full%20text%5BFilter%5D%20AND%20(%22Am%20Heart%20J%22%5Bjour%5D%20OR%20%22Am%20J%20Cardiol%22%5Bjour%5D%20OR%20%22Arterioscler%20Thromb%20Vasc%20Biol%22%5Bjour%5D%20OR%20%22Atherosclerosis%22%5Bjour%5D%20OR%20%22Basic%20Res%20Cardiol%22%5Bjour%5D%20OR%20%22Cardiovasc%20Res%22%5Bjour%5D%20OR%20%22Chest%22%5Bjour%5D%20OR%20%22Circulation%22%5Bjour%5D%20OR%20%22Circ%20Arrhythm%20Electrophysiol%22%5Bjour%5D%20OR%20%22Circ%20Cardiovasc%20Genet%22%5Bjour%5D%20OR%20%22Circ%20Cardiovasc%20Imaging%22%5Bjour%5D%20OR%20%22Circ%20Cardiovasc%20Qual%20Outcomes%22%5Bjour%5D%20OR%20%22Circ%20Cardiovasc%20Interv%22%5Bjour%5D%20OR%20%22Circ%20Heart%20Fail%22%5Bjour%5D%20OR%20%22Circ%20Res%22%5Bjour%5D%20OR%20%22ESC%20Heart%20Fail%22%5Bjour%5D%20OR%20%22Eur%20Heart%20J%22%5Bjour%5D%20OR%20%22Eur%20Heart%20J%20Cardiovasc%20Imaging%22%5Bjour%5D%20OR%20%22Eur%20Heart%20J%20Acute%20Cardiovasc%20Care%22%5Bjour%5D%20OR%20%22Eur%20Heart%20J%20Cardiovasc%20Pharmacother%22%5Bjour%5D%20OR%20%22Eur%20Heart%20J%20Qual%20Care%20Clin%20Outcomes%22%5Bjour%5D%20OR%20%22Eur%20J%20Heart%20Fail%22%5Bjour%5D%20OR%20%22Eur%20J%20Vasc%20Endovasc%20Surg%22%5Bjour%5D%20OR%20%22Europace%22%5Bjour%5D%20OR%20%22Heart%22%5Bjour%5D%20OR%20%22Heart%20Lung%20Circ%22%5Bjour%5D%20OR%20%22Heart%20Rhythm%22%5Bjour%5D%20OR%20%22JACC%20Cardiovasc%20Interv%22%5Bjour%5D%20OR%20%22JACC%20Cardiovasc%20Imaging%22%5Bjour%5D%20OR%20%22JACC%20Heart%20Fail%22%5Bjour%5D%20OR%20%22J%20Am%20Coll%20Cardiol%22%5Bjour%5D%20OR%20%22J%20Am%20Heart%20Assoc%22%5Bjour%5D%20OR%20%22J%20Am%20Soc%20Echocardiogr%22%5Bjour%5D%20OR%20%22J%20Card%20Fail%22%5Bjour%5D%20OR%20%22J%20Cardiovasc%20Electrophysiol%22%5Bjour%5D%20OR%20%22J%20Cardiovasc%20Magn%20Reson%22%5Bjour%5D%20OR%20%22J%20Heart%20Lung%20Transplant%22%5Bjour%5D%20OR%20%22J%20Hypertens%22%5Bjour%5D%20OR%20%22J%20Mol%20Cell%20Cardiol%22%5Bjour%5D%20OR%20%22J%20Thorac%20Cardiovasc%20Surg%22%5Bjour%5D%20OR%20%22J%20Vasc%20Surg%22%5Bjour%5D%20OR%20%22Nat%20Rev%20Cardiol%22%5Bjour%5D%20OR%20%22Prog%20Cardiovasc%20Dis%22%5Bjour%5D%20OR%20%22Resuscitation%22%5Bjour%5D%20OR%20%22Stroke%22%5Bjour%5D)",
    "relDate": 1000,
    "attributes": {
      "free": true,
      "public": true,
      "source": "Pubmed"
    },
    "resourceTypeID": 80
  }
]
```

Fields in the config:

`run` : `true/false` - switch on or off

`category` : Descriptive, does nothing yet

`searchTerm` : the url encoded search term, test [here](https://www.ncbi.nlm.nih.gov/pubmed/advanced)

`relDate` : include articles published up to this many days back

`attributes` : a json string used for faceting, optional however `category` 
should be included for multi-category resource libraries

`resourceTypeID` : id of the resource type from primary database 
`ol_resource_type` table, used to provide facet search for *video*, 
*audio*, *document* etc.

The search term in the json above extracts articles from a specific set of 
journals (see below).

The journals have been primarily selected from the highest ranking cardiology 
journals at <https://www.scimagojr.com/se>.

The search terms use the pubmed journal names in abbreviated form (journal codes). 
A complete list can be found at <https://www.ncbi.nlm.nih.gov/books/NBK3827/table/pubmedhelp.T.journal_lists/>

The current list:

- American heart journal (Am Heart J)
- The American journal of cardiology (Am J Cardiol)
- Arteriosclerosis, thrombosis, and vascular biology (Arterioscler Thromb Vasc Biol)
- Atherosclerosis (Atherosclerosis)
- Basic research in cardiology (Basic Res Cardiol)
- Cardiovascular research (Cardiovasc Res)
- Chest (Chest)
- Circulation (Circulation)
- Circulation. Arrhythmia and electrophysiology (Circ Arrhythm Electrophysiol)
- Circulation. Cardiovascular genetics (Circ Cardiovasc Genet)
- Circulation. Cardiovascular imaging (Circ Cardiovasc Imaging)
- Circulation. Cardiovascular quality and outcomes (Circ Cardiovasc Qual Outcomes)
- Circulation. Cardiovascular interventions (Circ Cardiovasc Interv)
- Circulation. Heart failure (Circ Heart Fail)
- Circulation research - (Circ Res)
- ESC heart failure (ESC Heart Fail)
- European heart journal (Eur Heart J)
- European heart journal cardiovascular Imaging (Eur Heart J Cardiovasc Imaging)
- European heart journal. Acute cardiovascular care (Eur Heart J Acute Cardiovasc Care)
- European heart journal. Cardiovascular pharmacotherapy (Eur Heart J Cardiovasc Pharmacother)
- European heart journal. Quality of care & clinical outcomes (Eur Heart J Qual Care Clin Outcomes)
- European journal of heart failure (Eur J Heart Fail)
- European journal of vascular and endovascular surgery : the official journal (Eur J Vasc Endovasc Surg)
- Europace : European pacing, arrhythmias, and cardiac electrophysiology (Europace)
- Heart - British Cardiac Society (Heart)
- Heart, lung & circulation (Heart Lung Circ)
- Heart rhythm (Heart Rhythm)
- JACC. Cardiovascular interventions (JACC Cardiovasc Interv)
- JACC. Cardiovascular imaging (JACC Cardiovasc Imaging)
- JACC. Heart failure (JACC Heart Fail)
- Journal of the American College of Cardiology (J Am Coll Cardiol)
- Journal of the American Heart Association (J Am Heart Assoc)
- Journal of the American Society of Echocardiography (J Am Soc Echocardiogr)
- Journal of cardiac failure (J Card Fail)
- Journal of cardiovascular electrophysiology (J Cardiovasc Electrophysiol)
- Journal of cardiovascular magnetic resonance (J Cardiovasc Magn Reson)
- The Journal of heart and lung transplantation (J Heart Lung Transplant)
- Journal of hypertension (J Hypertens)
- Journal of molecular and cellular cardiology (J Mol Cell Cardiol)
- The Journal of thoracic and cardiovascular surgery (J Thorac Cardiovasc Surg)
- Journal of vascular surgery (J Vasc Surg)
- Nature reviews. Cardiology (Nat Rev Cardiol)
- Progress in cardiovascular diseases (Prog Cardiovasc Dis)
- Resuscitation (Resuscitation)
- Stroke (Stroke)

## Pubmed Notes

The Pubmed query that fetches the article abstract (efetch) supports XML 
and *not* JSON. Go can access nested XML fields nicely so this works well.

Here's an example:
https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&id=17284678&retmode=xml&rettype=abstract

Pubmed search terms are complex. There is a lot of [documentation](https://www.ncbi.nlm.nih.gov/books/NBK25499/) and searches 
 can be tested [here](https://www.ncbi.nlm.nih.gov/pubmed/advanced).

Here is an example:
```
(((((((Cardiology) OR Cardiology[MeSH Terms]) AND Heart) OR Heart[MeSH Terms]))) AND "freetext"[Filter] AND jsubsetaim[text]) 
```

This first bit looks for articles related to "Cardiology" or "Heart"
```
(((((((Cardiology) OR Cardiology[MeSH Terms]) AND Heart) OR Heart[MeSH Terms])))
```

This bit filters results that include links to full text articles, free of charge. 
Although many listed as free are often not - yes I'm looking at you, Elsevier.
```
AND "freetext"[Filter]
```

... and the last section include articles from "Core Clinical" journals, or the [Abridged Index Medicus](https://www.nlm.nih.gov/bsd/aim.html).
 This filters also accommodates [Index Medicus](https://en.wikipedia.org/wiki/Index_Medicus) using the term `jsubsetim`  
```
AND jsubsetaim[text])
```

The other filter we need to apply is on `<MedlineCitation Status="Value">` node. Details on the `Summary` attribute 
can be found [here](https://www.nlm.nih.gov/bsd/licensee/elements_descriptions.html). 
  
From this information it seems as though we should include articles with the Status values of _MEDLINE_, 
_PubMed-not-Medline_, and _OLDMEDLINE_.    


(
    (Cardiology OR Heart)
    AND "freetext"[Filter] 
    AND "medline"[filter]
    AND jsubsetaim[text]  
    AND 
    (
    "random allocation"[MeSH Terms]
    OR
    "therapeutic use"[Subheading]
    )
)

Yields 18,000+

(
    (Cardiology OR Heart)
    AND "freetext"[Filter] 
    AND "medline"[filter]
    AND jsubsetaim[text]  
)

Yields 64,000+

(excluding time filters)

The above can be pasted directly into the search box here:

https://www.ncbi.nlm.nih.gov/pubmed/?term=((Cardiology+OR+Heart)+AND+%22freetext%22[Filter]+AND+%22medline%22[filter]+AND+jsubsetaim[text])

### Search Term Breakdown

This searches all text as well as mesh terms:
```
(Cardiology OR Heart)
```
 The above is translated as:
```
(("cardiology"[MeSH Terms] OR "cardiology"[All Fields]) OR ("heart"[MeSH Terms] OR "heart"[All Fields]))
```


This filters citations with links to free, full-text articles
```
AND "freetext"[Filter]
``` 

This filters MEDLINE citations: 
```
AND "medline"[filter]
```
Thus, we end up filtering articles with Medline Citation status = "MEDLINE", ie:
```xml
<MedlineCitation Status="MEDLINE" Owner="NLM">
   ...
</MedlineCitation>
```
This has two implications. Firstly, "MEDLINE" status means the articles have passed through  
quality processes. Secondly, it changes the fields that are returned containing useful keywords.

Without this filter the useful keywords are either contained in either `<Keywords>` or `<MeshHeading>`. 
  The reason is that the MeshHeading descriptors are only created once the article is elevated to MEDLINE status. 
  Given that we are applying this filter, the keywords *should* always be contained in the 
  `<MeshHeading>` node. For example:
  
```xml
<MeshHeading>
  <DescriptorName UI="D013318" MajorTopicYN="Y">Stroke Volume</DescriptorName>
</MeshHeading>
```


Note... these might be useful as well... strongest evidence:
"random allocation"[MeSH Terms]
"therapeutic use"[Subheading]

NOT medline[sb]) AND english[la] AND (systematic[sb] OR ((clinical[tiab] AND trial[tiab]) OR clinical trial[pt] OR random*[tiab] OR random allocation[mh] OR therapeutic use[sh])) AND free full text[sb]

## References

General: 
https://www.ncbi.nlm.nih.gov/books/NBK25499/

MedlineCitation Status: 
https://www.nlm.nih.gov/bsd/licensee/elements_descriptions.html

Differet databases & restricting to Medline only:
https://www.nlm.nih.gov/pubs/factsheets/dif_med_pub.html
 
Abridged Index Medicus:
https://www.nlm.nih.gov/bsd/aim.html

Subject Filters[sb]: 
https://www.nlm.nih.gov/bsd/pubmed_subsets.html

CareSearch example:
https://www.caresearch.com.au/caresearch/tabid/1743/Default.aspx

