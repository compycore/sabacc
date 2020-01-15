// Configuration
var backendEndpoint = "https://jessemillar-sabacc.herokuapp.com/sabacc";

// Global variables
var turnTaken = false;
var pickAction = "Pick your action!";
var database = JSON.parse(decodeURIComponent(window.location.search.substr(1)));
var warp;
var rollInterval;

function init() {
  warp = new WarpSpeed("canvas", { speedAdjFactor: 0.02 });

  if (database && database.rematch && database.rematch.length > 0) {
    startRematch();
  } else if (database && database.players.length > 0) {
    // Play the game if there's a game going
    populatePage();
  } else {
    startNewGame();
  }
}

function startRematch() {
  allPlayers = [];
  database.players = [];

  // Start a rematch
  for (var i = 0; i < database.rematch.length; i++) {
    allPlayers.push(database.rematch[i].email);

    database.players.push({
      email: database.rematch[i].email
    });
  }

  database.rematch = null;

  Swal.fire({
    title: "Rematch?",
    text:
      "Do you want to start a rematch with " +
      arrayToSentence(allPlayers) +
      "?",
    icon: "question",
    showCancelButton: true,
    confirmButtonText: "Yes"
  }).then(result => {
    if (result.value) {
      endTurn(function() {
        Swal.fire({
          icon: "success",
          title: "Rematch started!",
          text:
            "A rematch has started with " +
            arrayToSentence(allPlayers) +
            ". The second player listed will now receive an email! You can close this browser window.",
        });
      });
    }
  });
}

function startNewGame() {
  Swal.fire({
    title: "Start a new game!",
    text:
      "Enter email addresses (separated by a comma) of the people you want to play with:",
    input: "text",
    inputAttributes: {
      autocapitalize: "off"
    }
  }).then(result => {
    if (
      result.value.split(",").length > 1 &&
      result.value.split(",").length <= 8
    ) {
      var emailAddresses = result.value.split(",");
      database = {
        players: []
      };

      for (var i = 0; i < emailAddresses.length; i++) {
        database.players.push({
          email: emailAddresses[i]
        });
      }

      endTurn(function() {
        Swal.fire(
          "A new game has started with " +
            result.value.split(",").join(", ") +
            ". The first player listed will now receive an email! You can now close this window."
        );
      });
    } else if (result.value.split(",").length == 1) {
      Swal.fire("Please enter more than one email address.").then(() => {
        location.reload(false);
      });
    } else if (result.value.split(",").length > 8) {
      Swal.fire("Please enter less than eight email addresses.").then(() => {
        location.reload(false);
      });
    }
  });
}

// Don't populate the page before we know we're dealing with the right person
function populatePage() {
  document.getElementById("container").style = "display: block;";

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

function populateDiscardPile() {
  addCardToHand(
    "discard-pile",
    database.discards[database.discards.length - 1],
    "promptSwap();"
  );
}

function addCardToHand(barajaDivId, card, onclick) {
  var cardCount = getLiCount(barajaDivId);

  // TODO Use the nice deck animations by uncommenting and fixing the below stuff
  var hand = document.getElementById(barajaDivId);
  if (cardCount > 1) {
    // hand = window.baraja(document.getElementById(barajaDivId));
  }

  var li = document.createElement("li");
  var image = document.createElement("img");
  image.className = "sabacc-card";
  image.src = getCardFilename(card);
  li.appendChild(image);

  if (card != "back") {
    if (!onclick) {
      li.setAttribute(
        "onClick",
        "javascript: swap(" + JSON.stringify(card) + ");"
      );
    } else {
      li.setAttribute("onClick", "javascript: " + onclick);
    }
  }

  if (cardCount > 1) {
    // hand.add(li.outerHTML);
    hand.appendChild(li);
  } else {
    hand.appendChild(li);
  }
}

// Populate card image divs
function populateYourHand() {
  var hand = database.players[database.turn].hand;

  for (var i = 0; i < hand.length; i++) {
    addCardToHand("your-hand-cards", hand[i]);
  }

  fanCards("your-hand-cards");
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
      document.getElementById("container").innerHTML +=
        '<div class="u-full-width">';
      document.getElementById("container").innerHTML +=
        '<ul id="enemyHand' + i + '" class="baraja-container"></ul>';
      document.getElementById("container").innerHTML += "</div>";

      var hand = database.players[i].hand;

      for (var j = 0; j < hand.length; j++) {
        addCardToHand("enemyHand" + i, "back");
      }

      fanCards("enemyHand" + i);
    }
  }
}

function fanCards(divId) {
  // TODO Figure out a better way to wait for the DOM to be ready
  setTimeout(function() {
    var baraja = window.baraja(document.getElementById(divId));
    baraja.fan();
  }, 500);
}

function gain() {
  Swal.fire({
    title: "Do you want to discard a card first?",
    text: "You can discard a card before drawing a new one if you'd like.",
    icon: "question",
    showCancelButton: true,
    focusConfirm: false,
    showCloseButton: true,
    confirmButtonText: "Just draw",
    cancelButtonText: "Discard then draw",
  }).then(result => {
    if (result.value) {
      Swal.fire({
        title: "You drew...",
        text: "the " + getCardString(database.draw) + "!",
        imageUrl: getCardFilename(database.draw)
      }).then(() => {
        database.players[database.turn].hand.push(database.draw);
        delete database.draw;
        endTurn();
      });
    } else if (result.dismiss === Swal.DismissReason.cancel) {
      Swal.fire("Please tap on the card in your hand you wish to discard.", {
        icon: "info"
      });
    }
  });
}

function swap(card) {
  if (!turnTaken) {
    var topDiscardCard = getCardString(
      database.discards[database.discards.length - 1]
    );

    Swal.fire({
      title: "Your " + getCardString(card) + "?",
      html:
        "Do you want to discard your " +
        getCardString(card) +
        " and blindly draw a new card from the deck?" +
        "<br><br>" +
        "Or do you want to swap your " +
        getCardString(card) +
        " with the " +
        topDiscardCard +
        " that's on top of the discard pile?",
      icon: "question",
      showCloseButton: true,
      showCancelButton: true,
      focusConfirm: false,
      confirmButtonText: "Discard and draw",
      cancelButtonText: "Swap with " + topDiscardCard,
    }).then(result => {
      // Find the object for the card in question in the player's hand
      var cardIndexInHand = database.players[database.turn].hand.findIndex(
        element => element.value == card.value && element.stave == card.stave
      );

      // We hit the "confirm" button (which is actually "Discard and draw")
      if (result.value) {
        // Remove the card in question from the player's hand
        database.players[database.turn].hand.splice(cardIndexInHand, 1);
        // Put the draw card in the player's hand
        database.players[database.turn].hand.push(database.draw);
        // Wipe the drawn card
        database.draw = "";
        endTurn();
      } else if (result.dismiss === Swal.DismissReason.cancel) {
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
      }
    });
  }
}

function stand() {
  Swal.fire({
    title: "You want to stand?",
    text: "You're sure you want to do nothing for your turn?",
    icon: "question",
    showCancelButton: true,
    focusConfirm: false,
    showCloseButton: true,
    confirmButtonText: "Yes",
    cancelButtonText: "No",
  }).then(result => {
    if (result.value) {
      endTurn();
    } else if (result.dismiss === Swal.DismissReason.cancel) {
      Swal.fire(pickAction);
    }
  });
}

function trash() {
  Swal.fire({
    title: "Are you sure you want to trash?",
    text: "If you trash, you drop out of the game permanently.",
    icon: "warning",
    showCancelButton: true,
    focusConfirm: false,
    showCloseButton: true,
    confirmButtonText: "Yes",
    cancelButtonText: "No",
  }).then(result => {
    if (result.value) {
      database.players.splice(database.turn, 1);
      endTurn(function() {
        Swal.fire("You've withdrawn from the game.", {
          icon: "success",
          confirmButtonText: "'Til the Spire."
        });
      });
    } else if (result.dismiss === Swal.DismissReason.cancel) {
      Swal.fire("You're still in the game!");
    }
  });
}

function endTurn(callback) {
  if (database.turn == database.dealer && database.round > 0) {
    promptDealerRoll();
  } else {
    saveData(callback);
  }
}

function saveData(callback) {
  punchItChewie();

  Swal.fire({
    title: "Saving data...",
    text: "Please don't close the page. This may take a moment.",
    showConfirmButton: false
  });

  // Make an API call to the backend with the updated database info
  $.ajax({
    url: backendEndpoint + "?" + encodeURIComponent(JSON.stringify(database)),
    crossDomain: true
  }).done(function(data) {
    if (callback == null) {
      Swal.fire({
        title: "Data saved!",
        text: "Please wait for the next email.",
        icon: "success",
        confirmButtonText: "Patience, young padawan"
      }).then(wipePage());
    } else {
      callback();
      wipePage();
    }
  });
}

function wipePage() {
  turnTaken = true;
  document.getElementById("actions-header").innerHTML = "";
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
  } else if (card == "back") {
    return "images/cards/back.png";
  } else {
    return "images/cards/" + getCardString(card, "-") + ".png";
  }
}

function getCardString(card, separator = " ") {
  var cardColor = getCardColor(card.value);

  // Humans don't care about the stave and they like uppercase letters
  if (separator == " ") {
    if (card.value == 0) {
      return "Black 0";
    } else {
      return (
        cardColor.charAt(0).toUpperCase() +
        cardColor.slice(1) +
        separator +
        Math.abs(card.value)
      );
    }
  } else {
    return (
      card.stave + separator + cardColor + separator + Math.abs(card.value)
    );
  }
}

function promptSwap() {
  Swal.fire(
    "Tap a card in your hand to swap it with the one in the discard pile."
  );
}

function getLiCount(ulId) {
  return document.getElementById(ulId).getElementsByTagName("li").length;
}

function punchItChewie() {
  warp.TARGET_SPEED = 50;
}

function rollD6() {
  return Math.floor(Math.random() * Math.floor(6)) + 1;
}

function animateRollDice() {
  $("#dice1").attr("class", "cube show" + rollD6());
  $("#dice2").attr("class", "cube show" + rollD6());
}

function promptDealerRoll() {
  Swal.fire({
    title: "Time to roll!",
    text:
      "You were the dealer this round which means you get to roll the dice!",
    icon: "info"
  }).then(() => {
    showDice();
  });
}

function showDice() {
  Swal.fire({
    confirmButtonText: "Roll!",
    html:
      "<div class='dice-container'>" +
      "<div id='dice1' class='cube'>" +
      "<div class='top'></div>" +
      "<div class='front'></div>" +
      "<div class='left'></div>" +
      "<div class='back'></div>" +
      "<div class='right'></div>" +
      "<div class='bottom'></div>" +
      "</div>" +
      "</div>" +
      "<div class='dice-container'>" +
      "<div id='dice2' class='cube'>" +
      "<div class='top'></div>" +
      "<div class='front'></div>" +
      "<div class='left'></div>" +
      "<div class='back'></div>" +
      "<div class='right'></div>" +
      "<div class='bottom'></div>" +
      "</div>" +
      "</div>"
  }).then(() => {
    rollDice();
  });
}

function rollDice() {
  rollInterval = setInterval("animateRollDice()", 200);

  Swal.fire({
    showConfirmButton: false,
    html:
      "<div class='dice-container'>" +
      "<div id='dice1' class='cube'>" +
      "<div class='top'></div>" +
      "<div class='front'></div>" +
      "<div class='left'></div>" +
      "<div class='back'></div>" +
      "<div class='right'></div>" +
      "<div class='bottom'></div>" +
      "</div>" +
      "</div>" +
      "<div class='dice-container'>" +
      "<div id='dice2' class='cube'>" +
      "<div class='top'></div>" +
      "<div class='front'></div>" +
      "<div class='left'></div>" +
      "<div class='back'></div>" +
      "<div class='right'></div>" +
      "<div class='bottom'></div>" +
      "</div>" +
      "</div>"
  });

  setTimeout(function() {
    rollDiceToSide(rollD6(), rollD6());
  }, 2000);
}

function rollDiceToSide(side1, side2) {
  clearInterval(rollInterval);
  $("#dice1").attr("class", "cube show" + side1);
  $("#dice2").attr("class", "cube show" + side2);

  setTimeout(function() {
    if (side1 == side2) {
      showDiceResultDiscard();
    } else {
      showDiceResultNoDiscard();
    }
  }, 2000);
}

function showDiceResultNoDiscard() {
  Swal.fire({
    title: "Dice didn't match!",
    text: "That means everyone keeps their current hands.",
    icon: "info"
  }).then(() => {
    database.rolled = true;
    saveData();
  });
}

function showDiceResultDiscard() {
  Swal.fire({
    title: "Dice matched!",
    text: "That means everyone has to discard their current hands!",
    icon: "warning"
  }).then(() => {
    database.rolled = true;

    for (var i = 0; i < database.players.length; i++) {
      database.discards.concat(database.players[i].hand);
      database.players[i].hand = [];
    }

    saveData();
  });
}

function arrayToSentence(arr) {
  return (
    arr.slice(0, -2).join(", ") +
    (arr.slice(0, -2).length ? ", " : "") +
    arr.slice(-2).join(" and ")
  );
}
