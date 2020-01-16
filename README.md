# Sabacc

[![Build Status](https://travis-ci.org/jessemillar/sabacc.svg?branch=master)](https://travis-ci.org/jessemillar/sabacc)

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
