// TODO Implement dice rolling when the current player is the dealer and just ended their turn
// Configuration
var backendEndpoint = "https://jessemillar-sabacc.herokuapp.com/sabacc";

// Global variables
var turnTaken = false;
var pickAction = "Pick your action!";
var database = JSON.parse(decodeURIComponent(window.location.search.substr(1)));
console.log(database);

function init() {
  if (database && database.rematch && database.rematch.length > 0) {
    playerString = "";
    database.players = [];

    // Start a rematch
    for (var i = 0; i < database.rematch.length; i++) {
      playerString += database.rematch[i].email + ",";

      database.players.push({
        email: database.rematch[i].email
      });
    }

    database.rematch = null;

    endTurn(function() {
      swal(
        "A rematch has started with " +
          playerString.split(",").join(", ") +
          ". The first player listed will now receive an email! You can now close this window."
      );
    });
  } else if (database && database.players.length > 0) {
    // Play the game if there's a game going
    populatePage();
    swal(pickAction);
  } else {
    startNewGame();
  }
}

function startNewGame() {
  swal(
    "Enter email addresses (separated by a comma) of the people you want to play with:",
    {
      content: "input"
    }
  ).then(value => {
    if (value.split(",").length > 1 && value.split(",").length <= 8) {
      var emailAddresses = value.split(",");
      database = {
        players: []
      };

      for (var i = 0; i < emailAddresses.length; i++) {
        database.players.push({
          email: emailAddresses[i]
        });
      }

      endTurn(function() {
        swal(
          "A new game has started with " +
            value.split(",").join(", ") +
            ". The first player listed will now receive an email! You can now close this window."
        );
      });
    } else if (value.split(",").length == 1) {
      swal("Please enter more than one email address.").then(() => {
        location.reload(false);
      });
    } else if (value.split(",").length > 8) {
      swal("Please enter less than eight email addresses.").then(() => {
        location.reload(false);
      });
    }
  });
}

// Don't populate the page before we know we're dealing with the right person
function populatePage() {
  populateRound();
  populateScore();
  populateYourHand();
  populateDiscardPile();
  populateEnemyHands();
}

function populateRound() {
  document.getElementById("actions-header").innerHTML =
    "Actions (Round: " + database.round + ")";
}

function populateScore() {
  var player = database.players[database.turn];

  document.getElementById("your-hand-header").innerHTML =
    "Your Hand (Score: " + player.score + ")";
}

// Populate card image divs
function populateYourHand() {
  var hand = database.players[database.turn].hand;

  for (var i = 0; i < hand.length; i++) {
    var card = hand[i];

    var ul = document.getElementById("your-hand-cards");
    var li = document.createElement("li");
    var image = document.createElement("img");
    image.className = "sabacc-card";
    image.src = getCardFilename(card);
    image.onclick = "swap(" + JSON.stringify(card) + ")";
    li.appendChild(image);
    ul.appendChild(li);
  }

  fanCards("your-hand-cards");
}

function populateDiscardPile() {
  document.getElementById("discard-pile").innerHTML +=
    "<div class='two columns'><img src='" +
    getCardFilename(database.discards[database.discards.length - 1]) +
    "' class='u-max-full-width' onclick='promptSwap()' style='cursor: pointer;' /></div>";
}

function promptSwap() {
  swal("Tap a card in your hand to swap it with the one in the discard pile.");
}

function populateEnemyHands() {
  for (var i = 0; i < database.players.length; i++) {
    if (database.players[i].email != database.players[database.turn].email) {
      document.getElementById("container").innerHTML += '<div class="row">';
      document.getElementById("container").innerHTML +=
        '<div class="u-full-width">';
      document.getElementById("container").innerHTML +=
        "<h4>" + database.players[i].email + "'s Hand</h4>";
      document.getElementById("container").innerHTML += "</div>";
      document.getElementById("container").innerHTML += "</div>";
      document.getElementById("container").innerHTML += '<div class="row">';

      var hand = database.players[i].hand;

      for (var j = 0; j < hand.length; j++) {
        document.getElementById("container").innerHTML +=
          '<div class="two columns"><img src="images/cards/back.png" class="u-max-full-width" /></div>';
      }

      document.getElementById("container").innerHTML += "</div>";
    }
  }
}

function fanCards(divId) {
	// TODO Figure out a better way to wait for the DOM to be ready
  setTimeout(function() {
    var cardCount = document.getElementById(divId).getElementsByTagName("li")
      .length;
    var baraja = window.baraja(document.getElementById(divId));

    baraja.fan({
      direction: "right",
      easing: "ease-out",
      origin: { x: 50, y: 200 },
      speed: 500,
      range: cardCount * 10,
      center: true
    });
  }, 500);
}

function gain() {
  swal({
    title: "Do you want to discard a card first?",
    text: "You can discard a card before drawing a new one if you'd like.",
    icon: "warning",
    buttons: ["Nah.", "Yeah!"]
  }).then(willDiscard => {
    if (willDiscard) {
      swal("Please tap on the card in your hand you wish to discard.", {
        icon: "info"
      });
    } else {
      swal({
        title: "You drew...",
        text: "the " + getCardString(database.draw) + "!",
        icon: getCardFilename(database.draw)
      }).then(() => {
        database.players[database.turn].hand.push(database.draw);
        delete database.draw;
        endTurn();
      });
    }
  });
}

function swap(card) {
  if (!turnTaken) {
    swal({
      title: "What do you want to do with this card?",
      text:
        "Do you want to swap your " +
        getCardString(card) +
        " with the " +
        getCardString(database.discards[database.discards.length - 1]) +
        " that's on top of the discard pile? Or discard your " +
        getCardString(card) +
        " and blindly draw a new card from the deck?",
      icon: "warning",
      buttons: {
        cancel: {
          text: "Cancel",
          value: "cancel",
          visible: true,
          className: "",
          closeModal: true
        },
        gain: {
          text: "Discard and draw",
          value: "gain",
          visible: true,
          className: "",
          closeModal: true
        },
        swap: {
          text: "Swap with discard pile",
          value: "swap",
          visible: true,
          className: "",
          closeModal: true
        }
      }
    }).then(value => {
      // Find the object for the card in question in the player's hand
      var cardIndexInHand = database.players[database.turn].hand.findIndex(
        element => element.value == card.value && element.stave == card.stave
      );

      if (value == "swap") {
        // Remove the card in question from the player's hand
        database.players[database.turn].hand.splice(cardIndexInHand, 1);
        // Put the top of the discard pile in the player's hand
        database.players[database.turn].hand.push(
          database.discards[database.discards.length - 1]
        );
        // Remove the card that was just added to the player's hand from the discard pile
        database.discards.splice(database.discards.length - 1, 1);
        // Put the card in the discard pile
        database.discards.push(card);
        endTurn();
      } else if (value == "gain") {
        // Remove the card in question from the player's hand
        database.players[database.turn].hand.splice(cardIndexInHand, 1);
        // Put the draw card in the player's hand
        database.players[database.turn].hand.push(database.draw);
        // Wipe the drawn card
        database.draw = "";
        endTurn();
      } else {
        swal(pickAction);
      }
    });
  }
}

function stand() {
  swal({
    title: "You want to stand?",
    text: "You're sure you want to do nothing for your turn?",
    icon: "warning",
    buttons: ["Nah.", "Yeah!"]
  }).then(willStand => {
    if (willStand) {
      endTurn();
    } else {
      swal(pickAction);
    }
  });
}

function trash() {
  swal({
    title: "Are you sure you want to trash?",
    text: "If you trash, you drop out of the game permanently.",
    icon: "warning",
    buttons: ["Nope!", "Yes."],
    dangerMode: true
  }).then(willDelete => {
    if (willDelete) {
      database.players.splice(database.turn, 1);
      endTurn(function() {
        swal("You've withdrawn from the game.", {
          icon: "success",
          button: "'Til the Spire."
        });
      });
    } else {
      swal("You're still in the game!");
    }
  });
}

function endTurn(callback) {
  swal({
    title: "Saving data...",
    text: "Please don't close the page. This may take a moment.",
    buttons: false
  });

  // Make an API call to the backend with the updated database info
  $.ajax({
    url: backendEndpoint + "?" + encodeURIComponent(JSON.stringify(database)),
    crossDomain: true
  }).done(function(data) {
    if (callback == null) {
      swal({
        title: "Data saved!",
        text: "Please wait for the next email.",
        icon: "success",
        button: "Patience, young padawan."
      }).then(wipePage());
    } else {
      callback();
      wipePage();
    }
  });
}

function wipePage() {
  turnTaken = true;
  document.getElementById("action-buttons").innerHTML = "";
}

function getCardColor(cardValue) {
  if (cardValue > 0) {
    return "green";
  } else {
    return "red";
  }
}

function getCardFilename(card) {
  if (card.value == 0) {
    return "images/cards/zero.png";
  } else {
    return "images/cards/" + getCardString(card, "-") + ".png";
  }
}

function getCardString(card, separator = " ") {
  return (
    card.stave +
    separator +
    getCardColor(card.value) +
    separator +
    Math.abs(card.value)
  );
}
