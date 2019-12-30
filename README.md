# Sabacc

There are two branches: `master` and `gh-pages`

## `master` branch

The `master` branch contains the Golang microservice that runs on Heroku and orchestrates the game logic and the sending of player turn emails.

## `gh-pages` branch

The `gh-pages` branch has a simple, static-file UI hosted on GitHub Pages that allows players to take actions on their turn. The UI only does a couple things:

1. Loosely verify that the player taking their turn is actually supposed to be (e.g. "Is [address] your email address?")
1. Display cards in your hand and the discard pile
1. Show buttons for taking actions
1. Sends an HTTP request to the Heroku microservice backend to actually take the selected action

## Design

Everything happens via URL links that are emailed to participants. There's nothing other than honesty stopping players from taking multiple turns.

Things to keep track of in the URL params:

- Round count
- Which player is taking their turn
- Player email addresses
- Player hands (try to "encrypt" so it's not easily readable; at least show a button instead of a link in the email)
- Discard pile

### Sample URL

```
https://jessemillar.github.io/sabacc?round=1&
```

### Notes

- See [this link](https://golang.org/pkg/net/url/#example_Values) for an example of how to encode an array for a URL in Golang
