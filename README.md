# Gone Fishing - Tracking our travel spends

<p align="center">
  <img src="/doc/iphone-full.png" width=150px" />
  <img src="/doc/content-full-open.png" width=240px" />
</p>

## Context

<img align="right" width="100" height="100" src="/doc/on-the-road.jpg">

In 2018, my girlfriend and I decided to travel for 6 months in South America.
Travelling over an extended period like that requires
requires having a budget and sticking to it as much as possible.
To help us with that, and as an excuse to keep coding on the road, I decided to
create an application and learn Go

### Why another application?

Most tracking apps require you to manually enter all expenses. Thanks to new
banks like [Monzo](https://monzo.com/) providing an API, this task can now be
automated. We also wanted to use [Splitwise](https://secure.splitwise.com/) for
other expenses not paid out of our Monzo accounts.

6 countries means 6 currencies. Some applications can handle a fixed currency
rate but they’re never 100% accurate and don’t include any bank fees for paying
in a foreign currency. Plugging directly into my bank API means I can reflect
and track my actual spending more accurately and in my own currency.

## Requirements

Being able to see:

 1) our total spending (itemised) each day.
 2) he average and median of our total daily spending per week and overall since
 the start of our trip.

## Architecture and solution

This project is a Go application based on [Gin Gonic](https://gin-gonic.com/)
framework. Gin has a small footprint and is good to spin up a barebone web
application easily. This app fetches all expenses from:

+ our bank Monzo.
+ specifically-identified groups in the Splitwise app that we used to track
  other expenses.

We wanted to track our daily spending but there is an issue when simply relying on
these APIs results. How to reflect expenses paid at one point in time but relating to a different
period of time? For example, an advance hostel booking for 3 nights.

 2 possible solutions came up:
  + Store metadata on this application in a database.
  + Collocate metadata next to each expense entry.

While storing metadata requires a more complex architecture relying on
relational models, it turns out that both Monzo and Splitwise allow expenses to
be annotated with plain text. Therefore parsing this block of text was enough to
provide a way to manipulate our expenses as we wanted and avoid setting up
a database!

### Spreading expenses

So, say we paid for our 3-night hostel stay on the first day but wanted to
allocate the price per night to reflect a more accurate daily budget.
To do this, I chose to apply a special note format and
[parse it](/internal/expenses/parse.go#L13-L40) to generate new entries in our
list of expenses.

For example, we spent 2 nights at San Francisco Hotel in Salta (Argentina) and
paid at then end. The entry on Monzo can be annotated with `#nights 27/11/18 to
28/11/18`.

<p align="center">
  <img src="/doc/monzo-nights.png" width=220px" />
</p>

Here we wanted to track that we spent £26.79 (£53.58 / 2) on Tuesday 27 November
2018 and on Wednesday 28 November 2018.

<p align="center">
  <img src="/doc/nights-example.png" width=400px" />
</p>

Splitwise offers a similar note feature that I leveraged as well to track
non-Monzo expenses.

<p align="center">
  <img src="/doc/splitwise-spread.png" width=220px" />
</p>


### Cash entries

Tracking cash during a trip like this is a bit more tricky.  With so many
transactions made in cash in South America, our automated tracking with the
Monzo API only got us so far. Also, since ATM withdrawal fees sometimes apply,
It was difficult to get an accurate sense of our spending.

To cope with this, I decided to follow a [similar
approach](/internal/expenses/parse.go#L42-L75) adapted to cash
entries. An ATM withdrawal note starting with `#cash` will treat each line as
a separate entry.

<p align="center">
  <img src="/doc/monzo-atm.png" width=220px" />
  <img src="/doc/monzo-atm-note.png" width=220px" />
</p>

Each cash entry is converted in GBP [using the real ATM withdrawal
fees](/internal/expenses/parse.go#L46).

<p align="center">
  <img src="/doc/cash-example.png" width=400px" />
</p>


Finally, I also used this text annotation to ignore some expenses that might
have occurred while travelling but were not related to the trip like gifts to
friends, with the keyword `#ignore`.

## Iterations

The first iteration on this project was to fetch and display all expenses
from the Monzo account and Splitwise group in a simple list, ordered by date.
Having something in production straight away was really important for me and
allowed us to use the gist of this application quickly.

Later I added a layer of caching to prevent hitting API limits while
developing and also to make the app slightly faster in South America where the
internet speed is not always good. I decided to use memory cache as I wanted to
quickly read/write direct blobs of data from the app without having to create
models in DB.

## Timezones issue

Because we crossed multiple timezones during this trip, I decided to [shift all
expense entries](/internal/expenses/expenses.go#L109) to the local timezone in
which the spend occurred.

## Install

To run this project you'll need

  + Docker (or a Redis server)
  + Go 1.13.x

Then you'll need the following environment variables that you can set in
a `.env` with

```bash
cat << EOF > .env
export APP_URL="http://127.0.0.1:8080/"
export START_DATE="2018-09-26T00:00:00Z"
export MONZO_ACCOUNT_ID="<YOUR_MONZO_CRED>"
export MONZO_CLIENT_ID="<YOUR_MONZO_CLIENT_ID>"
export MONZO_CLIENT_SECRET="<YOUR_MONZO_CLIENT_SECRET>"
export SPLITWISE_CLIENT_ID="<YOUR_SPLITWISE_CLIENT_ID>"
export SPLITWISE_CLIENT_SECRET="<YOUR_SPLITWISE_CLIENT_SECRET>"
EOF
```

```bash
docker-compose up
make run
```

You can run tests with `make test` and build with `make build`.

