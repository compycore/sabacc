<p align="center">
  <img src="https://raw.githubusercontent.com/compycore/sabacc/gh-pages/images/logo.png">
</p>

# Sabacc

[![Build Status](https://travis-ci.org/compycore/sabacc.svg?branch=master)](https://travis-ci.org/compycore/sabacc) [![Go Report Card](https://goreportcard.com/badge/github.com/compycore/sabacc)](https://goreportcard.com/report/github.com/compycore/sabacc) [![Man Hours](https://img.shields.io/endpoint?url=https%3A%2F%2Fmh.%3Eessemillar.com%2Fhours%3Frepo%3Dhttps%3A%2F%2Fgithub.com%2Fcompycore%2Fsabacc.git)](https://jessemillar.com/r/man-hours)

In Disney's newly-opened [Star Wars: Galaxy's Edge](https://disneyparks.disney.go.com/star-wars-galaxys-edge/) lands, you can buy a deck of [Sabacc playing cards](https://starwars.fandom.com/wiki/Corellian_Spike) that comes with rules for the Corellian Spike version of play. Over Christmas break, [Jesse Millar](https://jessemillar.com) became obsessed with playing Sabacc with his brothers and was quite disappointed when he could no longer play the game in person at the end of the holiday. This repo represents an attempt at creating a digital version of Sabacc so that people can play with each other regardless of distance or time commitments.

There were a few goals for this project:
- Asynchronous play to allow for a lack of time constraint (e.g. play casually throughout the day)
- Not require money to play or maintain (use free hosting and deployment technologies)

## Deploy

There are a few easy steps to follow if you want to deploy your own copy of Sabacc:

1. Sign up for one or more free email API accounts from the list below and note down your API key(s)
	- [Mailjet](https://www.mailjet.com) (200 emails/day for free)
	- [SendGrid](https://sendgrid.com) (100 emails/day for free)
1. Deploy to Heroku using the button below (fill in the environment variables section on the Heroku site with the API keys you obtained above)

	[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

## How to Play

1. Navigate to [the static UI](https://compycore.com/sabacc) in your browser
1. Enter email addresses for the players you want to play with
1. Wait for an email notification that it's your turn

## Develop

If you want to deploy a copy of Sabacc without using Heroku or if you want to develop locally, you'll need to set the following environment variables:

```
// The port the Sabacc HTTP server listens on
export PORT="8080"
// The location of the HTML UI (either leave as the default or set to your own value)
export SABACC_UI_HREF="https://compycore.com/sabacc/"
```

You'll need to choose at least one supported email API provider and provide API key details:

```
// The public API key you got from Mailjet
export MAILJET_API_KEY_PUBLIC="..."
// The private API key you got from Mailjet
export MAILJET_API_KEY_PRIVATE="..."

// The API key you got from SendGrid
export SENDGRID_API_KEY="..."
```

Optional environment variables are defined below:

```
// Used to disable the sending of emails (only use when testing)
export SABACC_DEBUG=true
```

## Notes

### General Notes

- Everything happens via URL links that are emailed to participants. There are no accounts, databases, or "unnecessary" complexity. The game state is encoded as a URI parameter in each email that is sent which means there's nothing other than honesty stopping players from taking multiple turns, hijacking someone else's turn, or peeking at the cards in other players hands if they know how to use a browser debug console. We here at CompyCore are not concerned by this since it's quite easy to cheat in physical card games as well.

### Repository Branches

#### `master` branch

The [`master` branch](https://github.com/compycore/sabacc/tree/master) contains the Golang microservice that runs on Heroku and orchestrates the game logic and the sending of notification emails.

#### `gh-pages` branch

The [`gh-pages` branch](https://github.com/compycore/sabacc/tree/gh-pages) has a simple, static-file UI hosted on GitHub Pages that allows players to take actions on their turn.

## Credits

### Art

All art in Sabacc is custom and drawn and animated with love by [Jesse Millar](https://jessemillar.com).

### Libraries

- [Mailjet](https://github.com/mailjet/mailjet-apiv3-go) for sending emails
- [Echo](https://echo.labstack.com) for HTTP server
- [warpspeed.js](https://fdossena.com/?p=warpspeed/i.frag) for starfield effect
- [SweetAlert2](https://sweetalert2.github.io) for pretty prompts
- [Skeleton](http://getskeleton.com) for CSS boilerplate
- [Baraja-JS](https://github.com/nuxy/baraja-js) for card animations
