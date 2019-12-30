// TODO Calculate your hand's value
// TODO Read URI params
// TODO Show game start/setup page if no query param info

var database;
var pickAction = "Pick your action!";

function init() {
	decodeParams();
  checkTurn();
}

function decodeParams() {
  database = decodeURIComponent(window.location.search.substr(1));
	console.log(database);
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
        swal(pickAction, {
          icon: "success",
        });
      } else {
        swal("Please wait for your turn. You'll receive another email when it's time.").then(() => {
          wipePage();
        });
      }
    });
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
