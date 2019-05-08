# Gone Fishing Expenses

## Context

My girlfriend and I decided to travel for 6 months in South America.
Travelling over an extended period like that requires to provision a budget and
stick to it as much as possible. To help us with that and, as I
wanted to keep up with some coding while being on the road,
I decided to create an application and learn Go!

### Why another application?

Most tracking app requires to manually enter all expenses. Thanks to new banks
like Monzo providing solid APIs, this task can be automated. I also wanted to
keep using SplitWise for some expenses like the ones paid with a credit card.

6 countries means 6 currencies to deal with. Some applications can manage a fixed
currency exchange rate but cannot be fully accurate and include your bank fees.
Using my bank API means I can reflect a more accurate spending in my own
currency.

## Requirement

Being able to access our daily total spendings as well as an average (and median) per week
and overall.

## Architecture and solution

Initially, I started with a Go application based on Gin Gonic framework. Gin has
a small footprint and is good to spin up a barebone web application easily. Both
Monzo and Splitwise returned

// GRAPH Joint account + Splitwise group => normalise

We quickly realised that some expenses should be manipulated slightly. Here
2 possible solutions came up:

+ Store metadata on GF application in relational DB
+ Annotate expenses directly on external services (Monzo & Splitwise)

While storing metadata on GF required a more complex architecture relying on
relational model, it turns out that both Monzo and Splitwise allow expenses to be
annotated with plain text. Therefore parsing this chunck of text is
enough to provide a way to manipulate our expenses.

For example, when an expense need to be spreaded over 3 days it could be
annotated with `#spread 01/04/19 to 03/04/19` Following simple text formats,
3 features we required came up:

+ Spreading over days. For expenses like hostels where we spend more than one
  nights, or flights, we wanted to spread the cost over the appropriate period. We used `#spread 01/04/19 to 03/04/19` or `#nights 01/04/19 to 03/04/19`
+ Ignore for expenses not related to our trip like gifts to friends etc.
+ Handling the cash entries. This is the most challenging part as we wanted to
  be able to reflect all our spending in cash.

## Iterations

The first iteration on this project was to fetch and display all expenses
from the Monzo account and Splitwise'group in a simple list, ordered by date.
Having something in production straight away was really important for me and
allowed us to use the gist of this application quickly.

Later I added a layer of caching to prevent hitting API limits while
developing and also to make the app slightly faster in South America where the
internet speed is not always good. I decided to use memory cache as I wanted to
quickly read/write direct blob of data from the app without having to create
models in DB.

## Timezone issue

// Describe the timezone issue where we need to reflect an expense in its local
context.
