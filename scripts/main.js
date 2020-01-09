// TODO Implement dice rolling when the current player is the dealer and just ended their turn
// Configuration
var backendEndpoint = "https://jessemillar-sabacc.herokuapp.com/sabacc";

// Global variables
var turnTaken = false;
var database;
var pickAction = "Pick your action!";
var database = JSON.parse(decodeURIComponent(window.location.search.substr(1)));
console.log(database);

function init() {
  if (database && database.rematch.length > 0) {
    for (var i = 0; i < database.rematch.length; i++) {
      database.players.push({
        email: database.rematch[i]
      });
    }

    swal(
      "A new game has started with " +
        value.split(",").join(", ") +
        ". The first player listed will now receive an email! You can now close this window."
    );
    endTurn(false);
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

      swal(
        "A new game has started with " +
          value.split(",").join(", ") +
          ". The first player listed will now receive an email! You can now close this window."
      );
      endTurn(false);
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
    "Your Hand (Score: " + player.score + ") [" + player.email + "]";
}

// Populate card image divs
function populateYourHand() {
  var hand = database.players[database.turn].hand;

  for (var i = 0; i < hand.length; i++) {
    var card = hand[i];

    document.getElementById("your-hand-cards").innerHTML +=
      "<div class='two columns'><img src='" +
      getCardFilename(card) +
      "' class='u-max-full-width' onclick='swap(" +
      JSON.stringify(card) +
      ")' style='cursor: pointer;' /></div>";
  }
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
          '<div class="two columns"><img src="images/cards/back.jpg" class="u-max-full-width" /></div>';
      }

      document.getElementById("container").innerHTML += "</div>";
    }
  }
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
          value: "swa",
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

      if ((value = "swap")) {
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
      } else if ((value = "gain")) {
        // Remove the card in question from the player's hand
        database.players[database.turn].hand.splice(cardIndexInHand, 1);
        // Put the draw card in the player's hand
        database.players[database.turn].hand.push(dabase.draw);
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
      endTurn(false);

      swal("You've withdrawn from the game.", {
        icon: "success",
        button: "'Til the Spire."
      });
    } else {
      swal("You're still in the game!");
    }
  });
}

function endTurn(showTurnOver = true) {
  // Make an API call to the backend with the updated database info
  $.ajax({
    url: backendEndpoint + "?" + encodeURIComponent(JSON.stringify(database)),
    crossDomain: true
  }).done(function(data) {
    if (showTurnOver) {
      swal({
        title: "Turn over.",
        text: "Your turn is now over! Please wait for the next email.",
        icon: "success",
        button: "Patience, young padawan."
      }).then(wipePage());
    } else {
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
    return "images/cards/zero.jpg";
  } else {
    return "images/cards/" + getCardString(card, "-") + ".jpg";
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
