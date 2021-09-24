# go actions workflow webhook

A workflow that can be used to call a remote endpoint with github workflow environment variables as the payload written in golang. There are a bunch of other workflows that do the same thing, I recommend using them instead. I created this as I wanted to do it in golang.

## Setup

This action requires two secrets to be set, `WEBHOOK_URL` and `WEBHOOK_SECRET`. Also ensure that your webhook is using https. You can also add extra data to your webhook through `EXTRAS` as a json string.

## Usage

```
- name: Invoke deployment hook
  uses: mlioo/go-action-workflow-webhook@main
  with:
    webhook_url: ${{ secrets.WEBHOOK_URL }}
    webhook_secret: ${{ secrets.WEBHOOK_SECRET }}
```

If you need more data you can add extras to the request as a json string which will add it to the extras property

```
- name: Invoke deployment hook
  uses: mlioo/go-action-workflow-webhook@main
  with:
    webhook_url: ${{ secrets.WEBHOOK_URL }}
    webhook_secret: ${{ secrets.WEBHOOK_SECRET }}
    extras: '{"hello":"world","version":1}'
```
