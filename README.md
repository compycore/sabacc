# Sabacc

[![Build Status](https://travis-ci.org/jessemillar/sabacc.svg?branch=master)](https://travis-ci.org/jessemillar/sabacc)

In Disney's newly-opened [Star Wars: Galaxy's Edge](https://disneyparks.disney.go.com/star-wars-galaxys-edge/) lands, you can buy a deck of [Sabacc playing cards](https://starwars.fandom.com/wiki/Sabacc). Over Christmas break, [Jesse Millar](https://jessemillar.com) became obsessed with playing Sabacc with his brothers and was quite disappointed when he could no longer play the game in person at the end of the holiday. This repo represents an attempt at creating a digital version of Sabacc so that people can play with each other regardless of distance or time commitments.

There were a few goals for this project:
- Async play to allow for a lack of time constraint (e.g. play casually throughout the day)
- Not require money to play or maintain (use free hosting and deployment technologies)

Everything happens via URL links that are emailed to participants. There's nothing other than honesty stopping players from taking multiple turns or hijacking someone else's turn.

## Deploy

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

## Notes

### Repository Branches

#### `master` branch

The `master` branch contains the Golang microservice that runs on Heroku and orchestrates the game logic and the sending of notification emails.

#### `gh-pages` branch

The `gh-pages` branch has a simple, static-file UI hosted on GitHub Pages that allows players to take actions on their turn.

## Credits

- [Mailjet](https://github.com/mailjet/mailjet-apiv3-go) for sending emails
- [Echo](https://echo.labstack.com) for HTTP server
- [warpspeed.js](https://fdossena.com/?p=warpspeed/i.frag) for starfield effect
- [SweetAlert2](https://sweetalert2.github.io) for pretty prompts
- [Skeleton](http://getskeleton.com) for CSS boilerplate
- [Baraja-JS](https://github.com/nuxy/baraja-js) for card animations
