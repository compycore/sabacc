// TODO Show current round

// Configuration
var backendEndpoint = "https://jessemillar-sabacc.herokuapp.com/sabacc"

// Global variables
var database;
var pickAction = "Pick your action!";
var database = JSON.parse(decodeURIComponent(window.location.search.substr(1)));
console.log(database);

function init() {
  // Play the game if there's a game going
  if (database.players.length > 0) {
    checkTurn();
  } else {
    // TODO Start page
  }
}

function checkTurn() {
  swal({
      title: "Is it your turn?",
      text: "It's " + database.players[database.turn].email + "'s turn. Are you " + database.players[database.turn].email + "?",
      icon: "warning",
      buttons: ["Nope!", "I'm the droid you're looking for."],
    })
    .then((willDelete) => {
      if (willDelete) {
        swal(pickAction).then(() => {
          populatePage();
        });
      } else {
        swal("Please wait for your turn. You'll receive another email when it's time.").then(() => {
          wipePage();
        });
      }
    });
}

// Don't populate the page before we know we're dealing with the right person
function populatePage() {
  populateScore();
  populateYourHand();
  populateDiscardPile();
  populateEnemyHands();
}

function populateScore() {
  var player = database.players[database.turn];
  var hand = player.hand;
  var score = 0;

  for (var i = 0; i < hand.length; i++) {
    score += hand[i].value;
  }

  document.getElementById("your-hand-header").innerHTML = "Your Hand (Score: " + score + ") [" + player.email + "]";
}

// Populate card image divs
function populateYourHand() {
  var hand = database.players[database.turn].hand;

  for (var i = 0; i < hand.length; i++) {
    var card = hand[i];

    document.getElementById("your-hand-cards").innerHTML += "<div class='two columns'><img src='" + getCardFilename(card) + "' class='u-max-full-width' onclick='discard()' /></div>";
  }
}

function populateDiscardPile() {
  document.getElementById("discard-pile").innerHTML += "<div class='two columns'><img src='" + getCardFilename(database.draw) + "' class='u-max-full-width' onclick='discard()' /></div>";
}

function populateEnemyHands() {
  for (var i = 0; i < database.players.length; i++) {
    if (database.players[i].email != database.players[database.turn].email) {
      document.getElementById("container").innerHTML += '<div class="row">'
      document.getElementById("container").innerHTML += '<div class="u-full-width">'
      document.getElementById("container").innerHTML += '<h4>' + database.players[i].email + '\'s Hand</h4>'
      document.getElementById("container").innerHTML += '</div>'
      document.getElementById("container").innerHTML += '</div>'
      document.getElementById("container").innerHTML += '<div class="row">'

      var hand = database.players[database.turn].hand;

      for (var j = 0; j < hand.length; j++) {
        document.getElementById("container").innerHTML += '<div class="two columns"><img src="images/cards/back.jpg" class="u-max-full-width" /></div>'
      }

      document.getElementById("container").innerHTML += '</div>'
    }
  }
}

function gain() {
  swal({
    title: "Do you want to discard a card first?",
    text: "You can discard a card before drawing a new one if you'd like.",
    icon: "warning",
    buttons: ["Nah.", "Yeah!"],
  }).then((willDiscard) => {
    if (willDiscard) {
      swal("Please tap on the card in your hand you wish to discard.", {
        icon: "info",
      });
    } else {
      // TODO Draw the card chosen by the backend
      swal({
        title: "You drew...",
        text: "the " + database.draw.stave + " " + getCardColor(database.draw.value) + " " + Math.abs(database.draw.value) + "!",
        icon: getCardFilename(database.draw),
      }).then(() => {
        database.players[database.turn].hand += database.draw;
        delete database.draw;
        endTurn();
      });
    }
  });
}

function discard(card) {
  swal({
    title: "Discard this card?",
    text: "You want to discard " + card.name + "?",
    icon: "warning",
    buttons: ["Nah.", "Yeah!"],
  }).then((willDiscard) => {
    if (willDiscard) {
      swal("Card discarded!").then(() => {
        // TODO Discard the card and make an API call
        endTurn();
      });
    } else {
      swal(pickAction);
    }
  });
}

function stand() {
  swal({
    title: "You want to stand?",
    text: "You're sure you want to do nothing for your turn?",
    icon: "warning",
    buttons: ["Nah.", "Yeah!"],
  }).then((willStand) => {
    if (willStand) {
      endTurn()

      swal("You chose to stand.", {
        icon: "success",
      });
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
    dangerMode: true,
  }).then((willDelete) => {
    if (willDelete) {
      database.players.splice(database.turn, 1);
      endTurn();

      swal("You've withdrawn from the game.", {
        icon: "success",
        button: "'Til the Spire.",
      });
    } else {
      swal("You're still in the game!");
    }
  });
}

function endTurn() {
  // TODO Make the backend increase the round and player turn
  // Make an API call to the backend with the updated database info
  $.ajax({
      url: backendEndpoint + "?" + encodeURIComponent(JSON.stringify(database)),
      crossDomain: true
    })
    .done(function(data) {
      if (console && console.log) {
        console.log("Sample of data:", data.slice(0, 100));
      }
    });

  swal({
    title: "Turn over.",
    text: "Your turn is now over! Please wait for the next email.",
    icon: "success",
    button: "Patience, young padawan.",
  }).then(
    wipePage()
  );
}

function wipePage() {
  document.body.innerHTML = '';
  window.close();
}

function getCardColor(cardValue) {
  if (cardValue > 0) {
    return "green";
  } else {
    return "red";
  }
}

function getCardFilename(card) {
  return "images/cards/" + card.stave + "-" + getCardColor(card.value) + "-" + Math.abs(card.value) + ".jpg";
}
