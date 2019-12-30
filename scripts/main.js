// TODO Calculate your hand's value

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
      text: "It's hellojessemillar@gmail.com's turn. Are you hellojessemillar@gmail.com?",
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
    var cardColor;

    if (card.value > 0) {
      cardColor = "green";
    } else {
      cardColor = "red";
    }

    document.getElementById("your-hand-cards").innerHTML += "<div class='two columns'><img src='images/cards/" + card.stave + "-" + cardColor + "-" + Math.abs(card.value) + ".jpg' class='u-max-full-width' onclick='discard()' /></div>";
  }
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
    })
    .then((willDiscard) => {
      if (willDiscard) {
        swal("Please tap on the card in your hand you wish to discard.", {
          icon: "info",
        });
      } else {
        swal(pickAction);
      }
    });
}

function discard(card) {
  swal({
      title: "Discard this card?",
      text: "You want to discard " + card.name + "?",
      icon: "warning",
      buttons: ["Nah.", "Yeah!"],
    })
    .then((willDiscard) => {
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
    })
    .then((willStand) => {
      if (willStand) {
        swal("You chose to stand.", {
          icon: "success",
        }).then(() => {
          endTurn()
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
    })
    .then((willDelete) => {
      if (willDelete) {
        swal("You've withdrawn from the game.", {
          icon: "success",
          button: "'Til the Spire.",
        }).then(() => {
          wipePage()
        });
      } else {
        swal("You're still in the game!");
      }
    });
}

function endTurn() {
  // TODO Make an API call
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
