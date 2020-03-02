# Subtitles Translator Tool
Simple web application which accepts files of [Subrip format](https://en.wikipedia.org/wiki/SubRip) and translates their content to a given language.

Nothing fancy, just created for a friend :D

## Setup
The whole project is a Google Project using various Google products:
1. [Google Translation](https://cloud.google.com/translate/docs) for the actual translation of the subtitles
2. [Google Storage](https://cloud.google.com/products/storage/) for hosting the web page
3. [Google Functions](https://cloud.google.com/functions) for reading the subtitles file, and translating into the given language using the Google Tranlation API

All the above products are located within the same Google Project.


### Google Translation API
Enable API as described in the documentation.

### Frontend
Simple `html` page to upload and send the file.
The page is hosted in a Google Storaga Bucket.

#### Create Bucket
```
gsutil mb gs://[BUCKET_NAME]/
```
Configure with respective flags regarding storage class, location, and persmissions.

#### Upload `index.html` to the bucket
```
gsutil cp index.html gs://[BUCKET_NAME]/index.html
```

### Backend
Google Function to accept the HTTP call, and to return back the translated file.

#### Deploy function
```
cd functions/subTranslate/
gcloud functions deploy HandleTranslate --runtime go113 --trigger-http --set-env-vars GCP_PROJECT=[GCP_PROJECT]
```

#### Call function
```
curl -v  "localhost:8081/HandleTranslate?language=de" -X POST --data-binary @subs/test.srt
```