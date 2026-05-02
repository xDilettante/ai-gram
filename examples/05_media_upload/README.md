# 05 Media Upload

A small media example for sending photos and documents.

## What it does

- reads `AIGRAM_BOT_TOKEN` and `AIGRAM_CHAT_ID`;
- sends a photo by `file_id` when `AIGRAM_PHOTO_FILE_ID` is set;
- uploads a local photo when `AIGRAM_PHOTO_PATH` is set;
- sends a document by `file_id` when `AIGRAM_DOCUMENT_FILE_ID` is set;
- uploads a local document when `AIGRAM_DOCUMENT_PATH` is set;
- otherwise generates a small temporary text file and uploads it as a document.

No binary files are stored in the repository.

## Run

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
export AIGRAM_CHAT_ID='123456789'
go run ./examples/05_media_upload
```

Optional photo upload:

```bash
export AIGRAM_PHOTO_PATH='./photo.jpg'
go run ./examples/05_media_upload
```

Optional file IDs:

```bash
export AIGRAM_PHOTO_FILE_ID='existing_photo_file_id'
export AIGRAM_DOCUMENT_FILE_ID='existing_document_file_id'
go run ./examples/05_media_upload
```

Use a private test chat. Do not commit real tokens or private file IDs.
